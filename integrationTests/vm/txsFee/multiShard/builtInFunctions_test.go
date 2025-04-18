package multiShard

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/scheduled"
	"github.com/multiversx/mx-chain-go/config"
	"github.com/multiversx/mx-chain-go/integrationTests"
	"github.com/multiversx/mx-chain-go/integrationTests/vm"
	"github.com/multiversx/mx-chain-go/integrationTests/vm/txsFee/utils"
	"github.com/multiversx/mx-chain-go/state"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/stretchr/testify/require"
)

func getZeroGasAndFees() scheduled.GasAndFees {
	return scheduled.GasAndFees{
		AccumulatedFees: big.NewInt(0),
		DeveloperFees:   big.NewInt(0),
		GasProvided:     0,
		GasPenalized:    0,
		GasRefunded:     0,
	}
}

// Test scenario
// 1. Do a SC deploy on shard 1
// 2. Do a ChangeOwnerAddress (owner of the deployed contract will be moved in shard 0)
// 3. Do a ClaimDeveloperReward (cross shard call , the transaction will be executed on the source shard and the destination shard)
// 4. Execute SCR from context destination on context source ( the new owner will receive the developer rewards)
func TestBuiltInFunctionExecuteOnSourceAndDestinationShouldWork(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	testContextSource, err := vm.CreatePreparedTxProcessorWithVMsMultiShard(
		0,
		config.EnableEpochs{
			PenalizedTooMuchGasEnableEpoch:                  integrationTests.UnreachableEpoch,
			DynamicGasCostForDataTrieStorageLoadEnableEpoch: integrationTests.UnreachableEpoch,
		},
		1)
	require.Nil(t, err)
	defer testContextSource.Close()

	testContextDst, err := vm.CreatePreparedTxProcessorWithVMsMultiShard(
		1,
		config.EnableEpochs{
			PenalizedTooMuchGasEnableEpoch:                  integrationTests.UnreachableEpoch,
			DynamicGasCostForDataTrieStorageLoadEnableEpoch: integrationTests.UnreachableEpoch,
		},
		1)
	require.Nil(t, err)
	defer testContextDst.Close()

	pathToContract := "../../wasm/testdata/counter/output/counter_old.wasm"
	scAddr, owner := utils.DoDeployOldCounter(t, testContextDst, pathToContract)
	require.Equal(t, uint32(1), testContextDst.ShardCoordinator.ComputeId(scAddr))
	require.Equal(t, uint32(1), testContextDst.ShardCoordinator.ComputeId(owner))
	gasAndFees := getZeroGasAndFees()
	testContextDst.TxFeeHandler.CreateBlockStarted(gasAndFees)
	utils.CleanAccumulatedIntermediateTransactions(t, testContextDst)

	newOwner := []byte("12345678901234567890123456789110")
	require.Equal(t, uint32(0), testContextDst.ShardCoordinator.ComputeId(newOwner))

	gasPrice := uint64(10)
	gasLimit := uint64(100000)

	txData := []byte(core.BuiltInFunctionChangeOwnerAddress + "@" + hex.EncodeToString(newOwner))
	tx := vm.CreateTransaction(1, big.NewInt(0), owner, scAddr, gasPrice, gasLimit, txData)
	returnCode, err := testContextDst.TxProcessor.ProcessTransaction(tx)
	require.Nil(t, err)
	require.Equal(t, vmcommon.Ok, returnCode)

	_, err = testContextDst.Accounts.Commit()
	require.Nil(t, err)

	utils.CheckOwnerAddr(t, testContextDst, scAddr, newOwner)

	accumulatedFees := testContextDst.TxFeeHandler.GetAccumulatedFees()
	require.Equal(t, big.NewInt(1000000), accumulatedFees)

	utils.CleanAccumulatedIntermediateTransactions(t, testContextDst)

	// do a sc call intra shard
	sndAddr := []byte("12345678901234567890123456789111")
	shardID := testContextDst.ShardCoordinator.ComputeId(sndAddr)
	require.Equal(t, uint32(1), shardID)

	scStateAcc, _ := testContextDst.Accounts.GetExistingAccount(scAddr)
	scUserAcc := scStateAcc.(state.UserAccountHandler)
	currentSCDevBalance := scUserAcc.GetDeveloperReward()

	gasLimit = uint64(50000)
	_, _ = vm.CreateAccount(testContextDst.Accounts, sndAddr, 0, big.NewInt(100000000))
	tx = vm.CreateTransaction(0, big.NewInt(0), sndAddr, scAddr, gasPrice, gasLimit, []byte("increment"))

	retCode, err := testContextDst.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	_, err = testContextDst.Accounts.Commit()
	require.Nil(t, err)

	expectedBalance := big.NewInt(99843270)
	vm.TestAccount(t, testContextDst.Accounts, sndAddr, 1, expectedBalance)

	accumulatedFees = testContextDst.TxFeeHandler.GetAccumulatedFees()
	require.Equal(t, big.NewInt(1156730), accumulatedFees)

	developerFees := testContextDst.TxFeeHandler.GetDeveloperFees()
	require.Equal(t, big.NewInt(115578), developerFees)

	// call get developer rewards
	gasLimit = 500000
	_, _ = vm.CreateAccount(testContextSource.Accounts, newOwner, 0, big.NewInt(10000000))
	txData = []byte(core.BuiltInFunctionClaimDeveloperRewards)
	tx = vm.CreateTransaction(0, big.NewInt(0), newOwner, scAddr, gasPrice, gasLimit, txData)

	// execute claim on source shard
	retCode, err = testContextSource.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	expectedBalance = big.NewInt(5000000)
	utils.TestAccount(t, testContextSource.Accounts, newOwner, 1, expectedBalance)

	accumulatedFees = testContextSource.TxFeeHandler.GetAccumulatedFees()
	require.Equal(t, big.NewInt(5000000), accumulatedFees)

	developerFees = testContextSource.TxFeeHandler.GetDeveloperFees()
	require.Equal(t, big.NewInt(499977), developerFees)

	utils.CleanAccumulatedIntermediateTransactions(t, testContextDst)

	// execute claim on destination shard
	retCode, err = testContextDst.TxProcessor.ProcessTransaction(tx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	txs := testContextDst.GetIntermediateTransactions(t)
	scr := txs[0]

	utils.ProcessSCRResult(t, testContextSource, scr, vmcommon.Ok, nil)

	expectedBalance = big.NewInt(499977 + 4515686 + currentSCDevBalance.Int64())
	utils.TestAccount(t, testContextSource.Accounts, newOwner, 1, expectedBalance)

}
