package processing

import (
	"sync"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go/consensus"
	"github.com/ElrondNetwork/elrond-go/dataRetriever"
	"github.com/ElrondNetwork/elrond-go/dblookupext"
	"github.com/ElrondNetwork/elrond-go/epochStart"
	"github.com/ElrondNetwork/elrond-go/errors"
	"github.com/ElrondNetwork/elrond-go/factory"
	"github.com/ElrondNetwork/elrond-go/genesis"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/ElrondNetwork/elrond-go/sharding/nodesCoordinator"
	"github.com/ElrondNetwork/elrond-go/update"
)

var _ factory.ComponentHandler = (*managedProcessComponents)(nil)
var _ factory.ProcessComponentsHolder = (*managedProcessComponents)(nil)
var _ factory.ProcessComponentsHandler = (*managedProcessComponents)(nil)

type managedProcessComponents struct {
	*processComponents
	factory              *processComponentsFactory
	mutProcessComponents sync.RWMutex
}

// NewManagedProcessComponents returns a news instance of managedProcessComponents
func NewManagedProcessComponents(pcf *processComponentsFactory) (*managedProcessComponents, error) {
	if pcf == nil {
		return nil, errors.ErrNilProcessComponentsFactory
	}

	return &managedProcessComponents{
		processComponents: nil,
		factory:           pcf,
	}, nil
}

// Create will create the managed components
func (mpc *managedProcessComponents) Create() error {
	pc, err := mpc.factory.Create()
	if err != nil {
		return err
	}

	mpc.mutProcessComponents.Lock()
	mpc.processComponents = pc
	mpc.mutProcessComponents.Unlock()

	return nil
}

// Close will close all underlying sub-components
func (mpc *managedProcessComponents) Close() error {
	mpc.mutProcessComponents.Lock()
	defer mpc.mutProcessComponents.Unlock()

	if mpc.processComponents == nil {
		return nil
	}

	err := mpc.processComponents.Close()
	if err != nil {
		return err
	}
	mpc.processComponents = nil

	return nil
}

// CheckSubcomponents verifies all subcomponents
func (mpc *managedProcessComponents) CheckSubcomponents() error {
	mpc.mutProcessComponents.Lock()
	defer mpc.mutProcessComponents.Unlock()

	if mpc.processComponents == nil {
		return errors.ErrNilProcessComponents
	}
	if check.IfNil(mpc.processComponents.nodesCoordinator) {
		return errors.ErrNilNodesCoordinator
	}
	if check.IfNil(mpc.processComponents.shardCoordinator) {
		return errors.ErrNilShardCoordinator
	}
	if check.IfNil(mpc.processComponents.interceptorsContainer) {
		return errors.ErrNilInterceptorsContainer
	}
	if check.IfNil(mpc.processComponents.resolversFinder) {
		return errors.ErrNilResolversFinder
	}
	if check.IfNil(mpc.processComponents.roundHandler) {
		return errors.ErrNilRoundHandler
	}
	if check.IfNil(mpc.processComponents.epochStartTrigger) {
		return errors.ErrNilEpochStartTrigger
	}
	if check.IfNil(mpc.processComponents.epochStartNotifier) {
		return errors.ErrNilEpochStartNotifier
	}
	if check.IfNil(mpc.processComponents.forkDetector) {
		return errors.ErrNilForkDetector
	}
	if check.IfNil(mpc.processComponents.blockProcessor) {
		return errors.ErrNilBlockProcessor
	}
	if check.IfNil(mpc.processComponents.blackListHandler) {
		return errors.ErrNilBlackListHandler
	}
	if check.IfNil(mpc.processComponents.bootStorer) {
		return errors.ErrNilBootStorer
	}
	if check.IfNil(mpc.processComponents.headerSigVerifier) {
		return errors.ErrNilHeaderSigVerifier
	}
	if check.IfNil(mpc.processComponents.headerIntegrityVerifier) {
		return errors.ErrNilHeaderIntegrityVerifier
	}
	if check.IfNil(mpc.processComponents.validatorsStatistics) {
		return errors.ErrNilValidatorsStatistics
	}
	if check.IfNil(mpc.processComponents.validatorsProvider) {
		return errors.ErrNilValidatorsProvider
	}
	if check.IfNil(mpc.processComponents.blockTracker) {
		return errors.ErrNilBlockTracker
	}
	if check.IfNil(mpc.processComponents.pendingMiniBlocksHandler) {
		return errors.ErrNilPendingMiniBlocksHandler
	}
	if check.IfNil(mpc.processComponents.requestHandler) {
		return errors.ErrNilRequestHandler
	}
	if check.IfNil(mpc.processComponents.txLogsProcessor) {
		return errors.ErrNilTxLogsProcessor
	}
	if check.IfNil(mpc.processComponents.headerConstructionValidator) {
		return errors.ErrNilHeaderConstructionValidator
	}
	if check.IfNil(mpc.processComponents.peerShardMapper) {
		return errors.ErrNilPeerShardMapper
	}
	if check.IfNil(mpc.processComponents.fallbackHeaderValidator) {
		return errors.ErrNilFallbackHeaderValidator
	}
	if check.IfNil(mpc.processComponents.nodeRedundancyHandler) {
		return errors.ErrNilNodeRedundancyHandler
	}
	if check.IfNil(mpc.processComponents.currentEpochProvider) {
		return errors.ErrNilCurrentEpochProvider
	}
	if check.IfNil(mpc.processComponents.scheduledTxsExecutionHandler) {
		return errors.ErrNilScheduledTxsExecutionHandler
	}
	if check.IfNil(mpc.processComponents.txsSender) {
		return errors.ErrNilTxsSender
	}
	if check.IfNil(mpc.processComponents.processedMiniBlocksTracker) {
		return process.ErrNilProcessedMiniBlocksTracker
	}
	return nil
}

// NodesCoordinator returns the nodes coordinator
func (mpc *managedProcessComponents) NodesCoordinator() nodesCoordinator.NodesCoordinator {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.nodesCoordinator
}

// ShardCoordinator returns the shard coordinator
func (mpc *managedProcessComponents) ShardCoordinator() sharding.Coordinator {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.shardCoordinator
}

// InterceptorsContainer returns the interceptors container
func (mpc *managedProcessComponents) InterceptorsContainer() process.InterceptorsContainer {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.interceptorsContainer
}

// ResolversFinder returns the resolvers finder
func (mpc *managedProcessComponents) ResolversFinder() dataRetriever.ResolversFinder {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.resolversFinder
}

// RoundHandler returns the roundHandler
func (mpc *managedProcessComponents) RoundHandler() consensus.RoundHandler {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.roundHandler
}

// EpochStartTrigger returns the epoch start trigger handler
func (mpc *managedProcessComponents) EpochStartTrigger() epochStart.TriggerHandler {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.epochStartTrigger
}

// EpochStartNotifier returns the epoch start notifier
func (mpc *managedProcessComponents) EpochStartNotifier() factory.EpochStartNotifier {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.epochStartNotifier
}

// ForkDetector returns the fork detector
func (mpc *managedProcessComponents) ForkDetector() process.ForkDetector {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.forkDetector
}

// BlockProcessor returns the block processor
func (mpc *managedProcessComponents) BlockProcessor() process.BlockProcessor {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.blockProcessor
}

// BlackListHandler returns the black list handler
func (mpc *managedProcessComponents) BlackListHandler() process.TimeCacher {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.blackListHandler
}

// BootStorer returns the boot storer
func (mpc *managedProcessComponents) BootStorer() process.BootStorer {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.bootStorer
}

// HeaderSigVerifier returns the header signature verification
func (mpc *managedProcessComponents) HeaderSigVerifier() process.InterceptedHeaderSigVerifier {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.headerSigVerifier
}

// HeaderIntegrityVerifier returns the header integrity verifier
func (mpc *managedProcessComponents) HeaderIntegrityVerifier() process.HeaderIntegrityVerifier {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.headerIntegrityVerifier
}

// ValidatorsStatistics returns the validator statistics processor
func (mpc *managedProcessComponents) ValidatorsStatistics() process.ValidatorStatisticsProcessor {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.validatorsStatistics
}

// ValidatorsProvider returns the validator provider
func (mpc *managedProcessComponents) ValidatorsProvider() process.ValidatorsProvider {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.validatorsProvider
}

// BlockTracker returns the block tracker
func (mpc *managedProcessComponents) BlockTracker() process.BlockTracker {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.blockTracker
}

// PendingMiniBlocksHandler returns the pending mini blocks handler
func (mpc *managedProcessComponents) PendingMiniBlocksHandler() process.PendingMiniBlocksHandler {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.pendingMiniBlocksHandler
}

// RequestHandler returns the request handler
func (mpc *managedProcessComponents) RequestHandler() process.RequestHandler {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.requestHandler
}

// TxLogsProcessor returns the tx logs processor
func (mpc *managedProcessComponents) TxLogsProcessor() process.TransactionLogProcessorDatabase {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.txLogsProcessor
}

// HeaderConstructionValidator returns the validator for header construction
func (mpc *managedProcessComponents) HeaderConstructionValidator() process.HeaderConstructionValidator {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.headerConstructionValidator
}

// PeerShardMapper returns the peer to shard mapper
func (mpc *managedProcessComponents) PeerShardMapper() process.NetworkShardingCollector {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.peerShardMapper
}

// FallbackHeaderValidator returns the fallback header validator
func (mpc *managedProcessComponents) FallbackHeaderValidator() process.FallbackHeaderValidator {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.fallbackHeaderValidator
}

// TransactionSimulatorProcessor returns the transaction simulator processor
func (mpc *managedProcessComponents) TransactionSimulatorProcessor() factory.TransactionSimulatorProcessor {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.txSimulatorProcessor
}

// WhiteListHandler returns the white list handler
func (mpc *managedProcessComponents) WhiteListHandler() process.WhiteListHandler {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.whiteListHandler
}

// WhiteListerVerifiedTxs returns the white lister verified txs
func (mpc *managedProcessComponents) WhiteListerVerifiedTxs() process.WhiteListHandler {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.whiteListerVerifiedTxs
}

// HistoryRepository returns the history repository
func (mpc *managedProcessComponents) HistoryRepository() dblookupext.HistoryRepository {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.historyRepository
}

// ImportStartHandler returns the import status handler
func (mpc *managedProcessComponents) ImportStartHandler() update.ImportStartHandler {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.importStartHandler
}

// RequestedItemsHandler returns the items handler for the requests
func (mpc *managedProcessComponents) RequestedItemsHandler() dataRetriever.RequestedItemsHandler {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.requestedItemsHandler
}

// NodeRedundancyHandler returns the node redundancy handler
func (mpc *managedProcessComponents) NodeRedundancyHandler() consensus.NodeRedundancyHandler {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.nodeRedundancyHandler
}

// AccountsParser returns the genesis accounts parser
func (mpc *managedProcessComponents) AccountsParser() genesis.AccountsParser {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.accountsParser
}

// CurrentEpochProvider returns the current epoch provider that can decide if an epoch is active or not on the network
func (mpc *managedProcessComponents) CurrentEpochProvider() process.CurrentNetworkEpochProviderHandler {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.currentEpochProvider
}

// ScheduledTxsExecutionHandler returns the scheduled transactions execution handler
func (mpc *managedProcessComponents) ScheduledTxsExecutionHandler() process.ScheduledTxsExecutionHandler {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.scheduledTxsExecutionHandler
}

// TxsSenderHandler returns the transactions sender handler
func (mpc *managedProcessComponents) TxsSenderHandler() process.TxsSenderHandler {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.txsSender
}

// HardforkTrigger returns the hardfork trigger
func (mpc *managedProcessComponents) HardforkTrigger() factory.HardforkTrigger {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.hardforkTrigger
}

// ProcessedMiniBlocksTracker returns the processed mini blocks tracker
func (mpc *managedProcessComponents) ProcessedMiniBlocksTracker() process.ProcessedMiniBlocksTracker {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.processComponents == nil {
		return nil
	}

	return mpc.processComponents.processedMiniBlocksTracker
}

// ReceiptsRepository returns the receipts repository
func (mpc *managedProcessComponents) ReceiptsRepository() factory.ReceiptsRepository {
	mpc.mutProcessComponents.RLock()
	defer mpc.mutProcessComponents.RUnlock()

	if mpc.receiptsRepository == nil {
		return nil
	}

	return mpc.processComponents.receiptsRepository
}

// IsInterfaceNil returns true if the interface is nil
func (mpc *managedProcessComponents) IsInterfaceNil() bool {
	return mpc == nil
}

// String returns the name of the component
func (mpc *managedProcessComponents) String() string {
	return factory.ProcessComponentsName
}
