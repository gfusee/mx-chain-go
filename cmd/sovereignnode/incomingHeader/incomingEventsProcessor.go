package incomingHeader

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/hashing"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-go/process"
	"github.com/multiversx/mx-chain-go/process/block"
)

const (
	minTopicsInEvent          = 4
	numTransferTopics         = 3
	numExecutedBridgeOpTopics = 2
)

const (
	topicIDExecutedBridgeOp = "executedBridgeOp"
	topicIDDeposit          = "deposit"
)

type scrInfo struct {
	scr  *smartContractResult.SmartContractResult
	hash []byte
}

type incomingEventsProcessor struct {
	txPool     TransactionPool
	marshaller marshal.Marshalizer
	hasher     hashing.Hasher
	pool       block.OutGoingOperationsPool
}

func (iep *incomingEventsProcessor) processIncomingEvents(events []data.EventHandler) ([]*scrInfo, error) {
	scrs := make([]*scrInfo, 0, len(events))

	for idx, event := range events {
		topics := event.GetTopics()

		var scr *scrInfo
		var err error
		switch string(event.GetIdentifier()) {
		case topicIDDeposit:
			scr, err = iep.createSCRInfo(topics, event)
			scrs = append(scrs, scr)
		case topicIDExecutedBridgeOp:
			err = iep.confirmBridgeOperation(topics)
		default:
			return nil, errInvalidIncomingEventIdentifier
		}

		if err != nil {
			return nil, fmt.Errorf("%w, event idx = %d", err, idx)
		}
	}

	return scrs, nil
}

func (iep *incomingEventsProcessor) createSCRInfo(topics [][]byte, event data.EventHandler) (*scrInfo, error) {
	// TODO: Check each param validity (e.g. check that topic[0] == valid address)
	if len(topics) < minTopicsInEvent || len(topics[1:])%numTransferTopics != 0 {
		log.Error("incomingHeaderHandler.createIncomingSCRs",
			"error", errInvalidNumTopicsIncomingEvent,
			"num topics", len(topics),
			"topics", topics)
		return nil, fmt.Errorf("%w; num topics = %d",
			errInvalidNumTopicsIncomingEvent, len(topics))
	}

	eventNonce := big.NewInt(0).SetBytes(event.GetData())
	scr := &smartContractResult.SmartContractResult{
		Nonce:          eventNonce.Uint64(),
		OriginalTxHash: nil, // TODO:  Implement this in MX-14321 task
		RcvAddr:        topics[0],
		SndAddr:        core.ESDTSCAddress,
		Data:           createSCRData(topics),
		Value:          big.NewInt(0),
	}

	hash, err := core.CalculateHash(iep.marshaller, iep.hasher, scr)
	if err != nil {
		return nil, err
	}

	return &scrInfo{
		scr:  scr,
		hash: hash,
	}, nil
}

func createSCRData(topics [][]byte) []byte {
	numTokensToTransfer := len(topics[1:]) / numTransferTopics
	numTokensToTransferBytes := big.NewInt(int64(numTokensToTransfer)).Bytes()

	ret := []byte(core.BuiltInFunctionMultiESDTNFTTransfer +
		"@" + hex.EncodeToString(numTokensToTransferBytes))

	for idx := 1; idx < len(topics[1:]); idx += 3 {
		transfer := []byte("@" +
			hex.EncodeToString(topics[idx]) + // tokenID
			"@" + hex.EncodeToString(topics[idx+1]) + //nonce
			"@" + hex.EncodeToString(topics[idx+2])) //value

		ret = append(ret, transfer...)
	}

	return ret
}

func (iep *incomingEventsProcessor) addSCRsToPool(scrs []*scrInfo) {
	cacheID := process.ShardCacherIdentifier(core.MainChainShardId, core.SovereignChainShardId)

	for _, scrData := range scrs {
		iep.txPool.AddData(scrData.hash, scrData.scr, scrData.scr.Size(), cacheID)
	}
}

func (iep *incomingEventsProcessor) confirmBridgeOperation(topics [][]byte) error {
	if len(topics) != numExecutedBridgeOpTopics {
		return fmt.Errorf("%w for %s", errInvalidNumTopicsIncomingEvent, topicIDExecutedBridgeOp)
	}

	hashOfHashes := topics[0]
	hash := topics[1]
	return iep.pool.ConfirmOperation(hashOfHashes, hash)
}
