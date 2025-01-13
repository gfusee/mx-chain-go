package leavesRetriever

import (
	"context"
	"sync"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/hashing"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-go/common"
	"github.com/multiversx/mx-chain-go/trie/leavesRetriever/dfsTrieIterator"
	logger "github.com/multiversx/mx-chain-logger-go"
)

var log = logger.GetOrCreate("trie/leavesRetriever")

type leavesRetriever struct {
	iterators      map[string]common.DfsIterator
	lruIteratorIDs [][]byte
	db             common.TrieStorageInteractor
	marshaller     marshal.Marshalizer
	hasher         hashing.Hasher
	size           uint64
	maxSize        uint64
	mutex          sync.RWMutex
}

// NewLeavesRetriever creates a new leaves retriever
func NewLeavesRetriever(db common.TrieStorageInteractor, marshaller marshal.Marshalizer, hasher hashing.Hasher, maxSize uint64) (*leavesRetriever, error) {
	if check.IfNil(db) {
		return nil, ErrNilDB
	}
	if check.IfNil(marshaller) {
		return nil, ErrNilMarshaller
	}
	if check.IfNil(hasher) {
		return nil, ErrNilHasher
	}

	return &leavesRetriever{
		iterators:      make(map[string]common.DfsIterator),
		lruIteratorIDs: make([][]byte, 0),
		db:             db,
		marshaller:     marshaller,
		hasher:         hasher,
		size:           0,
		maxSize:        maxSize,
	}, nil
}

// GetLeaves retrieves the leaves from the trie. If there is a saved checkpoint for the iterator id, it will continue to iterate from the checkpoint.
func (lr *leavesRetriever) GetLeaves(numLeaves int, rootHash []byte, iteratorID []byte, ctx context.Context) (map[string]string, []byte, error) {
	defer log.Trace("leaves retriever stats", "size", lr.size, "numIterators", len(lr.iterators))
	if len(iteratorID) == 0 {
		return lr.getLeavesFromNewInstance(numLeaves, rootHash, ctx)
	}

	lr.mutex.RLock()
	iterator, ok := lr.iterators[string(iteratorID)]
	lr.mutex.RUnlock()
	if !ok {
		return nil, nil, ErrIteratorNotFound
	}

	return lr.getLeavesFromCheckpoint(numLeaves, iterator, iteratorID, ctx)
}

func (lr *leavesRetriever) getLeavesFromNewInstance(numLeaves int, rootHash []byte, ctx context.Context) (map[string]string, []byte, error) {
	log.Trace("get leaves from new instance", "numLeaves", numLeaves, "rootHash", rootHash)
	iterator, err := dfsTrieIterator.NewIterator(rootHash, lr.db, lr.marshaller, lr.hasher)
	if err != nil {
		return nil, nil, err
	}

	return lr.getLeavesFromIterator(iterator, numLeaves, ctx)
}

func (lr *leavesRetriever) getLeavesFromCheckpoint(numLeaves int, iterator common.DfsIterator, iteratorID []byte, ctx context.Context) (map[string]string, []byte, error) {
	log.Trace("get leaves from checkpoint", "numLeaves", numLeaves, "iteratorID", iteratorID)
	lr.markIteratorAsRecentlyUsed(iteratorID)
	clonedIterator := iterator.Clone()

	return lr.getLeavesFromIterator(clonedIterator, numLeaves, ctx)
}

func (lr *leavesRetriever) getLeavesFromIterator(iterator common.DfsIterator, numLeaves int, ctx context.Context) (map[string]string, []byte, error) {
	leaves, err := iterator.GetLeaves(numLeaves, lr.maxSize, ctx)
	if err != nil {
		return nil, nil, err
	}

	if iterator.FinishedIteration() {
		return leaves, nil, nil
	}

	iteratorId := iterator.GetIteratorId()
	if len(iteratorId) == 0 {
		return leaves, nil, nil
	}

	shouldReturnId := lr.manageIterators(iteratorId, iterator)
	if !shouldReturnId {
		return leaves, nil, nil
	}
	return leaves, iteratorId, nil
}

func (lr *leavesRetriever) manageIterators(iteratorId []byte, iterator common.DfsIterator) bool {
	lr.mutex.Lock()
	defer lr.mutex.Unlock()

	newIteratorPresent := lr.saveIterator(iteratorId, iterator)
	if !newIteratorPresent {
		return false
	}
	lr.removeIteratorsIfMaxSizeIsExceeded()
	return true
}

func (lr *leavesRetriever) saveIterator(iteratorId []byte, iterator common.DfsIterator) bool {
	_, isPresent := lr.iterators[string(iteratorId)]
	if isPresent {
		return true
	}

	if iterator.Size() >= lr.maxSize {
		return false
	}

	lr.lruIteratorIDs = append(lr.lruIteratorIDs, iteratorId)
	lr.iterators[string(iteratorId)] = iterator
	lr.size += iterator.Size() + uint64(len(iteratorId))
	return true
}

func (lr *leavesRetriever) markIteratorAsRecentlyUsed(iteratorId []byte) {
	lr.mutex.Lock()
	defer lr.mutex.Unlock()

	for i, id := range lr.lruIteratorIDs {
		if string(id) == string(iteratorId) {
			lr.lruIteratorIDs = append(lr.lruIteratorIDs[:i], lr.lruIteratorIDs[i+1:]...)
			lr.lruIteratorIDs = append(lr.lruIteratorIDs, id)
			return
		}
	}
}

func (lr *leavesRetriever) removeIteratorsIfMaxSizeIsExceeded() {
	if lr.size <= lr.maxSize {
		return
	}

	idsToRemove := make([][]byte, 0)
	sizeOfRemoved := uint64(0)
	numOfRemoved := 0

	for i := 0; i < len(lr.lruIteratorIDs); i++ {
		id := lr.lruIteratorIDs[i]
		idsToRemove = append(idsToRemove, id)
		iterator := lr.iterators[string(id)]
		sizeOfRemoved += iterator.Size() + uint64(len(id))
		numOfRemoved++

		if lr.size-sizeOfRemoved <= lr.maxSize {
			break
		}
	}

	for _, id := range idsToRemove {
		delete(lr.iterators, string(id))
	}
	lr.lruIteratorIDs = lr.lruIteratorIDs[numOfRemoved:]
	lr.size -= sizeOfRemoved
}

// IsInterfaceNil returns true if there is no value under the interface
func (lr *leavesRetriever) IsInterfaceNil() bool {
	return lr == nil
}
