package incomingHeader

import (
	"encoding/hex"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-go/process"
)

const (
	minTopicsInEvent  = 4
	numTransferTopics = 3
)

type scrInfo struct {
	scr  *smartContractResult.SmartContractResult
	hash []byte
}

func (ihp *incomingHeaderProcessor) createIncomingSCRs(events []data.EventHandler) ([]*scrInfo, error) {
	scrs := make([]*scrInfo, 0, len(events))

	for _, event := range events {
		topics := event.GetTopics()
		if len(topics) < minTopicsInEvent || len(topics[1:])%numTransferTopics != 0 {
			log.Error("incomingHeaderHandler.createIncomingSCRs",
				"error", errInvalidNumTopicsIncomingEvent,
				"num topics", len(topics),
				"topics", topics)
			continue
		}

		scr := &smartContractResult.SmartContractResult{
			RcvAddr: topics[0],
			SndAddr: core.ESDTSCAddress,
			Data:    createSCRData(topics),
		}

		hash, err := core.CalculateHash(ihp.marshaller, ihp.hasher, scr)
		if err != nil {
			return nil, err
		}

		scrs = append(scrs, &scrInfo{
			scr:  scr,
			hash: hash,
		})
	}

	return scrs, nil
}

func createSCRData(topics [][]byte) []byte {
	numTokensToTransfer := len(topics[1:]) / numTransferTopics
	numTokensToTransferBytes := big.NewInt(int64(numTokensToTransfer)).Bytes()

	ret := []byte(core.BuiltInFunctionMultiESDTNFTTransfer +
		"@" + hex.EncodeToString(topics[0]) + // topics[0] = address
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

func (ihp *incomingHeaderProcessor) addSCRsToPool(scrs []*scrInfo) {
	cacheID := process.ShardCacherIdentifier(core.MainChainShardId, core.SovereignChainShardId)

	for _, scrData := range scrs {
		ihp.txPool.AddData(scrData.hash, scrData.scr, scrData.scr.Size(), cacheID)
	}
}
