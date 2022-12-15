package executionOrder

import (
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	"github.com/ElrondNetwork/elrond-go/storage"
)

type miniblockGetter struct {
	storer     storage.Storer
	marshaller marshal.Marshalizer
}

func newMiniblocksGetter(storer storage.Storer, marshaller marshal.Marshalizer) *miniblockGetter {
	return &miniblockGetter{
		storer:     storer,
		marshaller: marshaller,
	}
}

// GetScheduledMBs will return the scheduled miniblocks
func (bg *miniblockGetter) GetScheduledMBs(currentHeader, prevHeader data.HeaderHandler) ([]*block.MiniBlock, error) {
	scheduledMbs := make([]*block.MiniBlock, 0)
	if !shouldFetchFromStorageMbs(currentHeader) {
		return scheduledMbs, nil
	}

	for _, mbHeader := range prevHeader.GetMiniBlockHeaderHandlers() {
		isScheduled := mbHeader.GetProcessingType() == int32(block.Scheduled)
		if !isScheduled {
			continue
		}

		mb, errGet := bg.getMBByHash(mbHeader.GetHash())
		if errGet != nil {
			return nil, errGet
		}

		scheduledMbs = append(scheduledMbs, mb)
	}

	return scheduledMbs, nil
}

func shouldFetchFromStorageMbs(currentHeader data.HeaderHandler) bool {
	for _, mb := range currentHeader.GetMiniBlockHeaderHandlers() {
		mbType := mb.GetTypeInt32()
		if mbType == int32(block.InvalidBlock) || mbType == int32(block.SmartContractResultBlock) {
			return true
		}
	}
	return false
}

func (bg *miniblockGetter) getMBByHash(mbHash []byte) (*block.MiniBlock, error) {
	mbBytes, err := bg.storer.Get(mbHash)
	if err != nil {
		return nil, err
	}

	mb := &block.MiniBlock{}
	err = bg.marshaller.Unmarshal(mb, mbBytes)

	return mb, err
}
