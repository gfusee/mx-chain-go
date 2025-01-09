package leavesRetriever_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-go/common"
	"github.com/multiversx/mx-chain-go/testscommon"
	"github.com/multiversx/mx-chain-go/testscommon/enableEpochsHandlerMock"
	"github.com/multiversx/mx-chain-go/testscommon/hashingMocks"
	"github.com/multiversx/mx-chain-go/testscommon/marshallerMock"
	trieTest "github.com/multiversx/mx-chain-go/testscommon/state"
	"github.com/multiversx/mx-chain-go/trie"
	"github.com/multiversx/mx-chain-go/trie/leavesRetriever"
	"github.com/stretchr/testify/assert"
)

func TestNewLeavesRetriever(t *testing.T) {
	t.Parallel()

	t.Run("nil db", func(t *testing.T) {
		t.Parallel()

		lr, err := leavesRetriever.NewLeavesRetriever(nil, &marshallerMock.MarshalizerMock{}, &hashingMocks.HasherMock{}, 100)
		assert.Nil(t, lr)
		assert.Equal(t, leavesRetriever.ErrNilDB, err)
	})
	t.Run("nil marshaller", func(t *testing.T) {
		t.Parallel()

		lr, err := leavesRetriever.NewLeavesRetriever(testscommon.NewMemDbMock(), nil, &hashingMocks.HasherMock{}, 100)
		assert.Nil(t, lr)
		assert.Equal(t, leavesRetriever.ErrNilMarshaller, err)
	})
	t.Run("nil hasher", func(t *testing.T) {
		t.Parallel()

		lr, err := leavesRetriever.NewLeavesRetriever(testscommon.NewMemDbMock(), &marshallerMock.MarshalizerMock{}, nil, 100)
		assert.Nil(t, lr)
		assert.Equal(t, leavesRetriever.ErrNilHasher, err)
	})
	t.Run("new leaves retriever", func(t *testing.T) {
		t.Parallel()

		var lr common.TrieLeavesRetriever
		assert.True(t, check.IfNil(lr))

		lr, err := leavesRetriever.NewLeavesRetriever(testscommon.NewMemDbMock(), &marshallerMock.MarshalizerMock{}, &hashingMocks.HasherMock{}, 100)
		assert.Nil(t, err)
		assert.False(t, check.IfNil(lr))
	})
}

func TestLeavesRetriever_GetLeaves(t *testing.T) {
	t.Parallel()

	t.Run("get leaves from new instance", func(t *testing.T) {
		t.Parallel()

		tr := trieTest.GetNewTrie()
		trieTest.AddDataToTrie(tr, 25)
		rootHash, _ := tr.RootHash()

		lr, _ := leavesRetriever.NewLeavesRetriever(tr.GetStorageManager(), &marshallerMock.MarshalizerMock{}, &hashingMocks.HasherMock{}, 100000)
		leaves, iteratorId, err := lr.GetLeaves(10, rootHash, []byte(""), context.Background())
		assert.Nil(t, err)
		assert.Equal(t, 10, len(leaves))
		assert.Equal(t, 32, len(iteratorId))
	})
	t.Run("get leaves from existing instance", func(t *testing.T) {
		t.Parallel()

		tr := trieTest.GetNewTrie()
		trieTest.AddDataToTrie(tr, 25)
		rootHash, _ := tr.RootHash()

		lr, _ := leavesRetriever.NewLeavesRetriever(tr.GetStorageManager(), &marshallerMock.MarshalizerMock{}, &hashingMocks.HasherMock{}, 10000000)
		leaves1, iteratorId1, err := lr.GetLeaves(10, rootHash, []byte(""), context.Background())
		assert.Nil(t, err)
		assert.Equal(t, 10, len(leaves1))
		assert.Equal(t, 32, len(iteratorId1))
		assert.Equal(t, 1, len(lr.GetIterators()))
		assert.Equal(t, 1, len(lr.GetLruIteratorIDs()))

		leaves2, iteratorId2, err := lr.GetLeaves(10, rootHash, iteratorId1, context.Background())
		assert.Nil(t, err)
		assert.Equal(t, 10, len(leaves2))
		assert.Equal(t, 32, len(iteratorId2))
		assert.Equal(t, 2, len(lr.GetIterators()))
		assert.Equal(t, 2, len(lr.GetLruIteratorIDs()))

		assert.NotEqual(t, leaves1, leaves2)
		assert.NotEqual(t, iteratorId1, iteratorId2)
	})
	t.Run("traversing a trie saves all iterator instances", func(t *testing.T) {
		t.Parallel()

		tr := trieTest.GetNewTrie()
		trieTest.AddDataToTrie(tr, 25)
		rootHash, _ := tr.RootHash()

		lr, _ := leavesRetriever.NewLeavesRetriever(tr.GetStorageManager(), &marshallerMock.MarshalizerMock{}, &hashingMocks.HasherMock{}, 10000000)
		leaves1, iteratorId1, err := lr.GetLeaves(10, rootHash, []byte(""), context.Background())
		assert.Nil(t, err)
		assert.Equal(t, 10, len(leaves1))
		assert.Equal(t, 32, len(iteratorId1))
		assert.Equal(t, 1, len(lr.GetIterators()))
		assert.Equal(t, 1, len(lr.GetLruIteratorIDs()))

		leaves2, iteratorId2, err := lr.GetLeaves(10, rootHash, iteratorId1, context.Background())
		assert.Nil(t, err)
		assert.Equal(t, 10, len(leaves2))
		assert.Equal(t, 32, len(iteratorId2))
		assert.Equal(t, 2, len(lr.GetIterators()))
		assert.Equal(t, 2, len(lr.GetLruIteratorIDs()))

		leaves3, iteratorId3, err := lr.GetLeaves(10, rootHash, iteratorId2, context.Background())
		assert.Nil(t, err)
		assert.Equal(t, 5, len(leaves3))
		assert.Equal(t, 0, len(iteratorId3))
		assert.Equal(t, 2, len(lr.GetIterators()))
		assert.Equal(t, 2, len(lr.GetLruIteratorIDs()))
	})
	t.Run("iterator instances are evicted in a lru manner", func(t *testing.T) {
		t.Parallel()

		tr := trieTest.GetNewTrie()
		trieTest.AddDataToTrie(tr, 25)
		rootHash, _ := tr.RootHash()
		maxSize := uint64(1000)

		lr, _ := leavesRetriever.NewLeavesRetriever(tr.GetStorageManager(), &marshallerMock.MarshalizerMock{}, &hashingMocks.HasherMock{}, maxSize)
		iterators := make([][]byte, 0)
		_, id1, _ := lr.GetLeaves(5, rootHash, []byte(""), context.Background())
		iterators = append(iterators, id1)
		_, id2, _ := lr.GetLeaves(5, rootHash, id1, context.Background())
		iterators = append(iterators, id2)
		_, id3, _ := lr.GetLeaves(5, rootHash, id2, context.Background())
		iterators = append(iterators, id3)

		assert.Equal(t, 3, len(lr.GetIterators()))
		for i, id := range lr.GetLruIteratorIDs() {
			assert.Equal(t, iterators[i], id)
		}

		_, id4, _ := lr.GetLeaves(5, rootHash, id3, context.Background())
		iterators = append(iterators, id4)
		assert.Equal(t, 3, len(lr.GetIterators()))
		for i, id := range lr.GetLruIteratorIDs() {
			assert.Equal(t, iterators[i+1], id)
		}
	})
	t.Run("when an iterator instance is used, it is moved in the front of the eviction queue", func(t *testing.T) {
		t.Parallel()

		tr := trieTest.GetNewTrie()
		trieTest.AddDataToTrie(tr, 25)
		rootHash, _ := tr.RootHash()
		maxSize := uint64(100000)

		lr, _ := leavesRetriever.NewLeavesRetriever(tr.GetStorageManager(), &marshallerMock.MarshalizerMock{}, &hashingMocks.HasherMock{}, maxSize)
		iterators := make([][]byte, 0)
		_, id1, _ := lr.GetLeaves(5, rootHash, []byte(""), context.Background())
		iterators = append(iterators, id1)
		leaves1, id2, _ := lr.GetLeaves(5, rootHash, id1, context.Background())
		iterators = append(iterators, id2)
		_, id3, _ := lr.GetLeaves(5, rootHash, id2, context.Background())
		iterators = append(iterators, id3)

		assert.Equal(t, 3, len(lr.GetIterators()))
		for i, id := range lr.GetLruIteratorIDs() {
			assert.Equal(t, iterators[i], id)
		}

		leaves2, id4, _ := lr.GetLeaves(5, rootHash, id1, context.Background())
		assert.Equal(t, leaves1, leaves2)
		assert.Equal(t, id2, id4)

		assert.Equal(t, 3, len(lr.GetIterators()))
		retrievedIterators := lr.GetLruIteratorIDs()
		assert.Equal(t, 3, len(retrievedIterators))
		assert.Equal(t, id2, retrievedIterators[0])
		assert.Equal(t, id3, retrievedIterators[1])
		assert.Equal(t, id1, retrievedIterators[2])
	})
	t.Run("iterator not found", func(t *testing.T) {
		t.Parallel()

		tr := trieTest.GetNewTrie()
		trieTest.AddDataToTrie(tr, 25)
		rootHash, _ := tr.RootHash()
		maxSize := uint64(100000)

		lr, _ := leavesRetriever.NewLeavesRetriever(tr.GetStorageManager(), &marshallerMock.MarshalizerMock{}, &hashingMocks.HasherMock{}, maxSize)
		leaves, id, err := lr.GetLeaves(5, rootHash, []byte("invalid iterator"), context.Background())
		assert.Nil(t, leaves)
		assert.Equal(t, 0, len(id))
		assert.Equal(t, leavesRetriever.ErrIteratorNotFound, err)
	})
	t.Run("max size reached on the first iteration", func(t *testing.T) {
		t.Parallel()

		tr := trieTest.GetTrieWithData()
		rootHash, _ := tr.RootHash()
		maxSize := uint64(100)

		lr, _ := leavesRetriever.NewLeavesRetriever(tr.GetStorageManager(), &marshallerMock.MarshalizerMock{}, &hashingMocks.HasherMock{}, maxSize)
		leaves, id1, err := lr.GetLeaves(10, rootHash, []byte(""), context.Background())
		assert.Nil(t, err)
		assert.Equal(t, 2, len(leaves))
		assert.Equal(t, 0, len(id1))
		assert.Equal(t, 0, len(lr.GetIterators()))
	})
}

func TestLeavesRetriever_Concurrency(t *testing.T) {
	t.Parallel()

	numTries := 10
	numLeaves := 1000
	tries := buildTries(numTries, numLeaves)

	rootHashes := make([][]byte, 0)
	for _, tr := range tries {
		rootHash, _ := tr.RootHash()
		rootHashes = append(rootHashes, rootHash)

	}

	maxSize := uint64(1000000)
	lr, _ := leavesRetriever.NewLeavesRetriever(tries[0].GetStorageManager(), &marshallerMock.MarshalizerMock{}, &hashingMocks.HasherMock{}, maxSize)

	wg := &sync.WaitGroup{}
	wg.Add(numTries)
	for i := 0; i < numTries; i++ {
		go retrieveTrieLeaves(t, lr, rootHashes[i], numLeaves, wg)
	}
	wg.Wait()
}

func retrieveTrieLeaves(t *testing.T, lr common.TrieLeavesRetriever, rootHash []byte, numLeaves int, wg *sync.WaitGroup) {
	iteratorId := []byte("")
	numRetrievedLeaves := 0
	for {
		leaves, newId, err := lr.GetLeaves(100, rootHash, iteratorId, context.Background())
		assert.Nil(t, err)
		iteratorId = newId
		numRetrievedLeaves += len(leaves)
		fmt.Println("Retrieved leaves: ", numRetrievedLeaves, " for root hash: ", hex.EncodeToString(rootHash))
		if len(iteratorId) == 0 {
			break
		}
	}
	assert.Equal(t, numLeaves, numRetrievedLeaves)
	wg.Done()
}

func buildTries(numTries int, numLeaves int) []common.Trie {
	tries := make([]common.Trie, 0)
	tsm, marshaller, hasher := trieTest.GetDefaultTrieParameters()
	for i := 0; i < numTries; i++ {
		tr, _ := trie.NewTrie(tsm, marshaller, hasher, &enableEpochsHandlerMock.EnableEpochsHandlerStub{}, 5)
		addDataToTrie(tr, numLeaves)
		tries = append(tries, tr)
	}
	return tries
}

func addDataToTrie(tr common.Trie, numLeaves int) {
	for i := 0; i < numLeaves; i++ {
		_ = tr.Update(generateRandomByteArray(32), generateRandomByteArray(32))
	}
	_ = tr.Commit()
}

func generateRandomByteArray(size int) []byte {
	r := make([]byte, size)
	_, _ = rand.Read(r)
	return r
}
