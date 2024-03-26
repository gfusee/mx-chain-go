package block

import (
	"github.com/multiversx/mx-chain-go/common"
	factorySovereign "github.com/multiversx/mx-chain-go/factory/sovereign"
	"github.com/multiversx/mx-chain-go/process"
	"github.com/multiversx/mx-chain-go/process/block/bootstrapStorage"
	"github.com/multiversx/mx-chain-go/state"

	"github.com/multiversx/mx-chain-core-go/data"
	sovereignCore "github.com/multiversx/mx-chain-core-go/data/sovereign"
)

type blockProcessor interface {
	removeStartOfEpochBlockDataFromPools(headerHandler data.HeaderHandler, bodyHandler data.BodyHandler) error
}

type gasConsumedProvider interface {
	TotalGasProvided() uint64
	TotalGasProvidedWithScheduled() uint64
	TotalGasRefunded() uint64
	TotalGasPenalized() uint64
	IsInterfaceNil() bool
}

type peerAccountsDBHandler interface {
	MarkSnapshotDone()
}

type receiptsRepository interface {
	SaveReceipts(holder common.ReceiptsHolder, header data.HeaderHandler, headerHash []byte) error
	IsInterfaceNil() bool
}

type validatorStatsRootHashGetter interface {
	GetValidatorStatsRootHash() []byte
}

type sovereignChainHeader interface {
	GetExtendedShardHeaderHashes() [][]byte
	GetOutGoingMiniBlockHeaderHandler() data.OutGoingMiniBlockHeaderHandler
}

type crossNotarizer interface {
	getLastCrossNotarizedHeaders() []bootstrapStorage.BootstrapHeaderInfo
}

// OutGoingOperationsPool defines the behavior of a timed cache for outgoing operations
type OutGoingOperationsPool interface {
	Add(data *sovereignCore.BridgeOutGoingData)
	Get(hash []byte) *sovereignCore.BridgeOutGoingData
	Delete(hash []byte)
	GetUnconfirmedOperations() []*sovereignCore.BridgeOutGoingData
	ConfirmOperation(hashOfHashes []byte, hash []byte) error
	IsInterfaceNil() bool
}

// OutGoingOperationsPoolCreator defines the outgoing operations pool factory handler
type OutGoingOperationsPoolCreator interface {
	CreateOutGoingOperationPool() OutGoingOperationsPool
	IsInterfaceNil() bool
}

// BlockProcessorCreator defines the block processor factory handler
type BlockProcessorCreator interface {
	CreateBlockProcessor(argumentsBaseProcessor ArgBaseProcessor) (process.DebuggerBlockProcessor, error)
	IsInterfaceNil() bool
}

// HeaderValidatorCreator is an interface for creating header validators
type HeaderValidatorCreator interface {
	CreateHeaderValidator(args ArgsHeaderValidator) (process.HeaderConstructionValidator, error)
	IsInterfaceNil() bool
}

type runTypeComponentsHolder interface {
	AccountsCreator() state.AccountFactory
	DataCodecCreator() factorySovereign.DataDecoderCreator
	TopicsCheckerCreator() factorySovereign.TopicsCheckerCreator
	IsInterfaceNil() bool
}
