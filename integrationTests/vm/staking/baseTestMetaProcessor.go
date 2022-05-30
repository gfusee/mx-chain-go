package staking

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	arwenConfig "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/config"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go-core/display"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	"github.com/ElrondNetwork/elrond-go/common"
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/dataRetriever"
	"github.com/ElrondNetwork/elrond-go/epochStart"
	"github.com/ElrondNetwork/elrond-go/epochStart/metachain"
	"github.com/ElrondNetwork/elrond-go/factory"
	"github.com/ElrondNetwork/elrond-go/integrationTests"
	"github.com/ElrondNetwork/elrond-go/process"
	vmFactory "github.com/ElrondNetwork/elrond-go/process/factory"
	"github.com/ElrondNetwork/elrond-go/process/mock"
	"github.com/ElrondNetwork/elrond-go/sharding/nodesCoordinator"
	"github.com/ElrondNetwork/elrond-go/state"
	"github.com/ElrondNetwork/elrond-go/testscommon/stakingcommon"
	"github.com/ElrondNetwork/elrond-go/vm"
	"github.com/ElrondNetwork/elrond-go/vm/systemSmartContracts"
	"github.com/ElrondNetwork/elrond-go/vm/systemSmartContracts/defaults"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

const (
	stakingV4InitEpoch                       = 1
	stakingV4EnableEpoch                     = 2
	stakingV4DistributeAuctionToWaitingEpoch = 3
	addressLength                            = 15
	nodePrice                                = 1000
)

func haveTime() bool { return true }
func noTime() bool   { return false }

type nodesConfig struct {
	eligible    map[uint32][][]byte
	waiting     map[uint32][][]byte
	leaving     map[uint32][][]byte
	shuffledOut map[uint32][][]byte
	queue       [][]byte
	auction     [][]byte
}

// TestMetaProcessor -
type TestMetaProcessor struct {
	MetaBlockProcessor  process.BlockProcessor
	NodesCoordinator    nodesCoordinator.NodesCoordinator
	ValidatorStatistics process.ValidatorStatisticsProcessor
	EpochStartTrigger   integrationTests.TestEpochStartTrigger
	BlockChainHandler   data.ChainHandler
	NodesConfig         nodesConfig
	AccountsAdapter     state.AccountsAdapter
	Marshaller          marshal.Marshalizer
	TxCacher            dataRetriever.TransactionCacher
	TxCoordinator       process.TransactionCoordinator
	SystemVM            vmcommon.VMExecutionHandler
	BlockChainHook      process.BlockChainHookHandler
	StakingDataProvider epochStart.StakingDataProvider

	currentRound uint64
}

func newTestMetaProcessor(
	coreComponents factory.CoreComponentsHolder,
	dataComponents factory.DataComponentsHolder,
	bootstrapComponents factory.BootstrapComponentsHolder,
	statusComponents factory.StatusComponentsHolder,
	stateComponents factory.StateComponentsHandler,
	nc nodesCoordinator.NodesCoordinator,
	maxNodesConfig []config.MaxNodesChangeConfig,
	queue [][]byte,
) *TestMetaProcessor {
	saveNodesConfig(
		stateComponents.AccountsAdapter(),
		coreComponents.InternalMarshalizer(),
		nc,
		maxNodesConfig,
	)

	createDelegationManagementConfig(
		stateComponents.AccountsAdapter(),
		coreComponents.InternalMarshalizer(),
	)

	gasScheduleNotifier := createGasScheduleNotifier()
	blockChainHook := createBlockChainHook(
		dataComponents,
		coreComponents,
		stateComponents.AccountsAdapter(),
		bootstrapComponents.ShardCoordinator(),
		gasScheduleNotifier,
	)

	metaVmFactory := createVMContainerFactory(
		coreComponents,
		gasScheduleNotifier,
		blockChainHook,
		stateComponents.PeerAccounts(),
		bootstrapComponents.ShardCoordinator(),
		nc,
		maxNodesConfig[0].MaxNumNodes,
	)
	vmContainer, _ := metaVmFactory.Create()
	systemVM, _ := vmContainer.Get(vmFactory.SystemVirtualMachine)

	validatorStatisticsProcessor := createValidatorStatisticsProcessor(
		dataComponents,
		coreComponents,
		nc,
		bootstrapComponents.ShardCoordinator(),
		stateComponents.PeerAccounts(),
	)
	stakingDataProvider := createStakingDataProvider(
		coreComponents.EpochNotifier(),
		systemVM,
	)
	scp := createSystemSCProcessor(
		nc,
		coreComponents,
		stateComponents,
		bootstrapComponents.ShardCoordinator(),
		maxNodesConfig,
		validatorStatisticsProcessor,
		systemVM,
		stakingDataProvider,
	)

	txCoordinator := &mock.TransactionCoordinatorMock{}
	epochStartTrigger := createEpochStartTrigger(coreComponents, dataComponents.StorageService())

	eligible, _ := nc.GetAllEligibleValidatorsPublicKeys(0)
	waiting, _ := nc.GetAllWaitingValidatorsPublicKeys(0)
	shuffledOut, _ := nc.GetAllShuffledOutValidatorsPublicKeys(0)

	return &TestMetaProcessor{
		AccountsAdapter: stateComponents.AccountsAdapter(),
		Marshaller:      coreComponents.InternalMarshalizer(),
		NodesConfig: nodesConfig{
			eligible:    eligible,
			waiting:     waiting,
			shuffledOut: shuffledOut,
			queue:       queue,
			auction:     make([][]byte, 0),
		},
		MetaBlockProcessor: createMetaBlockProcessor(
			nc,
			scp,
			coreComponents,
			dataComponents,
			bootstrapComponents,
			statusComponents,
			stateComponents,
			validatorStatisticsProcessor,
			blockChainHook,
			metaVmFactory,
			epochStartTrigger,
			vmContainer,
			txCoordinator,
		),
		currentRound:        1,
		NodesCoordinator:    nc,
		ValidatorStatistics: validatorStatisticsProcessor,
		EpochStartTrigger:   epochStartTrigger,
		BlockChainHandler:   dataComponents.Blockchain(),
		TxCacher:            dataComponents.Datapool().CurrentBlockTxs(),
		TxCoordinator:       txCoordinator,
		SystemVM:            systemVM,
		BlockChainHook:      blockChainHook,
		StakingDataProvider: stakingDataProvider,
	}
}

func saveNodesConfig(
	accountsDB state.AccountsAdapter,
	marshaller marshal.Marshalizer,
	nc nodesCoordinator.NodesCoordinator,
	maxNodesConfig []config.MaxNodesChangeConfig,
) {
	eligibleMap, _ := nc.GetAllEligibleValidatorsPublicKeys(0)
	waitingMap, _ := nc.GetAllWaitingValidatorsPublicKeys(0)
	allStakedNodes := int64(len(getAllPubKeys(eligibleMap)) + len(getAllPubKeys(waitingMap)))

	maxNumNodes := allStakedNodes
	if len(maxNodesConfig) > 0 {
		maxNumNodes = int64(maxNodesConfig[0].MaxNumNodes)
	}

	stakingcommon.SaveNodesConfig(
		accountsDB,
		marshaller,
		allStakedNodes,
		1,
		maxNumNodes,
	)
}

func createDelegationManagementConfig(accountsDB state.AccountsAdapter, marshaller marshal.Marshalizer) {
	delegationCfg := &systemSmartContracts.DelegationManagement{
		MinDelegationAmount: big.NewInt(10),
	}
	marshalledData, _ := marshaller.Marshal(delegationCfg)

	delegationAcc := stakingcommon.LoadUserAccount(accountsDB, vm.DelegationManagerSCAddress)
	_ = delegationAcc.DataTrieTracker().SaveKeyValue([]byte("delegationManagement"), marshalledData)
	_ = accountsDB.SaveAccount(delegationAcc)
	_, _ = accountsDB.Commit()
}

func createGasScheduleNotifier() core.GasScheduleNotifier {
	gasSchedule := arwenConfig.MakeGasMapForTests()
	defaults.FillGasMapInternal(gasSchedule, 1)
	return mock.NewGasScheduleNotifierMock(gasSchedule)
}

func createEpochStartTrigger(
	coreComponents factory.CoreComponentsHolder,
	storageService dataRetriever.StorageService,
) integrationTests.TestEpochStartTrigger {
	argsEpochStart := &metachain.ArgsNewMetaEpochStartTrigger{
		Settings: &config.EpochStartConfig{
			MinRoundsBetweenEpochs: 10,
			RoundsPerEpoch:         10,
		},
		Epoch:              0,
		EpochStartNotifier: coreComponents.EpochStartNotifierWithConfirm(),
		Storage:            storageService,
		Marshalizer:        coreComponents.InternalMarshalizer(),
		Hasher:             coreComponents.Hasher(),
		AppStatusHandler:   coreComponents.StatusHandler(),
	}

	epochStartTrigger, _ := metachain.NewEpochStartTrigger(argsEpochStart)
	testTrigger := &metachain.TestTrigger{}
	testTrigger.SetTrigger(epochStartTrigger)

	return testTrigger
}

// Process -
func (tmp *TestMetaProcessor) Process(t *testing.T, numOfRounds uint64) {
	for r := tmp.currentRound; r < tmp.currentRound+numOfRounds; r++ {
		header := tmp.createNewHeader(t, r)
		tmp.createAndCommitBlock(t, header, haveTime)
	}

	tmp.currentRound += numOfRounds
}

func (tmp *TestMetaProcessor) createNewHeader(t *testing.T, round uint64) *block.MetaBlock {
	_, err := tmp.MetaBlockProcessor.CreateNewHeader(round, round)
	require.Nil(t, err)

	epoch := tmp.EpochStartTrigger.Epoch()
	printNewHeaderRoundEpoch(round, epoch)

	currentHeader, currentHash := tmp.getCurrentHeaderInfo()
	header := createMetaBlockToCommit(
		epoch,
		round,
		currentHash,
		currentHeader.GetRandSeed(),
		tmp.NodesCoordinator.ConsensusGroupSize(core.MetachainShardId),
	)

	return header
}

func (tmp *TestMetaProcessor) createAndCommitBlock(t *testing.T, header data.HeaderHandler, haveTime func() bool) {
	newHeader, blockBody, err := tmp.MetaBlockProcessor.CreateBlock(header, haveTime)
	require.Nil(t, err)

	err = tmp.MetaBlockProcessor.CommitBlock(newHeader, blockBody)
	require.Nil(t, err)

	time.Sleep(time.Millisecond * 50)
	tmp.updateNodesConfig(header.GetEpoch())
	tmp.displayConfig(tmp.NodesConfig)
}

func printNewHeaderRoundEpoch(round uint64, epoch uint32) {
	headline := display.Headline(
		fmt.Sprintf("Commiting header in epoch %v round %v", epoch, round),
		"",
		delimiter,
	)
	fmt.Println(headline)
}

func (tmp *TestMetaProcessor) getCurrentHeaderInfo() (data.HeaderHandler, []byte) {
	currentHeader := tmp.BlockChainHandler.GetCurrentBlockHeader()
	currentHash := tmp.BlockChainHandler.GetCurrentBlockHeaderHash()
	if currentHeader == nil {
		currentHeader = tmp.BlockChainHandler.GetGenesisHeader()
		currentHash = tmp.BlockChainHandler.GetGenesisHeaderHash()
	}

	return currentHeader, currentHash
}

func createMetaBlockToCommit(
	epoch uint32,
	round uint64,
	prevHash []byte,
	prevRandSeed []byte,
	consensusSize int,
) *block.MetaBlock {
	roundStr := strconv.Itoa(int(round))
	hdr := block.MetaBlock{
		Epoch:                  epoch,
		Nonce:                  round,
		Round:                  round,
		PrevHash:               prevHash,
		Signature:              []byte("signature"),
		PubKeysBitmap:          []byte(strings.Repeat("f", consensusSize)),
		RootHash:               []byte("roothash" + roundStr),
		ShardInfo:              make([]block.ShardData, 0),
		TxCount:                1,
		PrevRandSeed:           prevRandSeed,
		RandSeed:               []byte("randseed" + roundStr),
		AccumulatedFeesInEpoch: big.NewInt(0),
		AccumulatedFees:        big.NewInt(0),
		DevFeesInEpoch:         big.NewInt(0),
		DeveloperFees:          big.NewInt(0),
	}

	shardMiniBlockHeaders := make([]block.MiniBlockHeader, 0)
	shardMiniBlockHeader := block.MiniBlockHeader{
		Hash:            []byte("mb_hash" + roundStr),
		ReceiverShardID: 0,
		SenderShardID:   0,
		TxCount:         1,
	}
	shardMiniBlockHeaders = append(shardMiniBlockHeaders, shardMiniBlockHeader)
	shardData := block.ShardData{
		Nonce:                 round,
		ShardID:               0,
		HeaderHash:            []byte("hdr_hash" + roundStr),
		TxCount:               1,
		ShardMiniBlockHeaders: shardMiniBlockHeaders,
		DeveloperFees:         big.NewInt(0),
		AccumulatedFees:       big.NewInt(0),
	}
	hdr.ShardInfo = append(hdr.ShardInfo, shardData)

	return &hdr
}

func (tmp *TestMetaProcessor) updateNodesConfig(epoch uint32) {
	eligible, _ := tmp.NodesCoordinator.GetAllEligibleValidatorsPublicKeys(epoch)
	waiting, _ := tmp.NodesCoordinator.GetAllWaitingValidatorsPublicKeys(epoch)
	leaving, _ := tmp.NodesCoordinator.GetAllLeavingValidatorsPublicKeys(epoch)
	shuffledOut, _ := tmp.NodesCoordinator.GetAllShuffledOutValidatorsPublicKeys(epoch)

	rootHash, _ := tmp.ValidatorStatistics.RootHash()
	validatorsInfoMap, _ := tmp.ValidatorStatistics.GetValidatorInfoForRootHash(rootHash)

	auction := make([][]byte, 0)
	for _, validator := range validatorsInfoMap.GetAllValidatorsInfo() {
		if validator.GetList() == string(common.AuctionList) {
			auction = append(auction, validator.GetPublicKey())
		}
	}

	tmp.NodesConfig.eligible = eligible
	tmp.NodesConfig.waiting = waiting
	tmp.NodesConfig.shuffledOut = shuffledOut
	tmp.NodesConfig.leaving = leaving
	tmp.NodesConfig.auction = auction
	tmp.NodesConfig.queue = tmp.getWaitingListKeys()
}

func generateAddresses(startIdx, n uint32) [][]byte {
	ret := make([][]byte, 0, n)

	for i := startIdx; i < n+startIdx; i++ {
		ret = append(ret, generateAddress(i))
	}

	return ret
}

func generateAddress(identifier uint32) []byte {
	uniqueIdentifier := fmt.Sprintf("address-%d", identifier)
	return []byte(strings.Repeat("0", addressLength-len(uniqueIdentifier)) + uniqueIdentifier)
}
