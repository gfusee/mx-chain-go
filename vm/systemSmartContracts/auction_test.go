package systemSmartContracts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"github.com/ElrondNetwork/elrond-go/vm"
	"github.com/ElrondNetwork/elrond-go/vm/mock"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/assert"
)

func createMockArgumentsForAuction() ArgsStakingAuctionSmartContract {
	args := ArgsStakingAuctionSmartContract{
		MinStakeValue: big.NewInt(10000000),
		MinStepValue:  big.NewInt(100000),
		TotalSupply:   big.NewInt(100000000000),
		UnBondPeriod:  uint64(100000),
		NumNodes:      uint32(10),
		Eei:           &mock.SystemEIStub{},
		SigVerifier:   &mock.MessageSignVerifierMock{},
	}

	return args
}

func createABid(totalStakeValue uint64, numBlsKeys uint32, maxStakePerNode uint64) AuctionData {
	data := AuctionData{
		RewardAddress:   []byte("addr"),
		RegisterNonce:   0,
		Epoch:           0,
		BlsPubKeys:      nil,
		TotalStakeValue: big.NewInt(0).SetUint64(totalStakeValue),
		LockedStake:     big.NewInt(0).SetUint64(totalStakeValue),
		MaxStakePerNode: big.NewInt(0).SetUint64(maxStakePerNode),
	}

	keys := make([][]byte, 0)
	for i := uint32(0); i < numBlsKeys; i++ {
		keys = append(keys, []byte(fmt.Sprintf("%d", rand.Uint32())))
	}
	data.BlsPubKeys = keys

	return data
}

func TestAuctionSC_calculateNodePrice_Case1(t *testing.T) {
	t.Parallel()

	expectedNodePrice := big.NewInt(20000000)
	stakingAuctionSC, _ := NewStakingAuctionSmartContract(createMockArgumentsForAuction())

	bids := []AuctionData{
		createABid(100000000, 100, 30000000),
		createABid(20000000, 100, 20000000),
		createABid(20000000, 100, 20000000),
		createABid(20000000, 100, 20000000),
		createABid(20000000, 100, 20000000),
		createABid(20000000, 100, 20000000),
	}

	nodePrice, err := stakingAuctionSC.calculateNodePrice(bids)
	assert.Equal(t, expectedNodePrice, nodePrice)
	assert.Nil(t, err)
}

func TestAuctionSC_calculateNodePrice_Case2(t *testing.T) {
	t.Parallel()

	expectedNodePrice := big.NewInt(20000000)
	args := createMockArgumentsForAuction()
	args.NumNodes = 5
	stakingAuctionSC, _ := NewStakingAuctionSmartContract(args)

	bids := []AuctionData{
		createABid(100000000, 1, 30000000),
		createABid(50000000, 100, 25000000),
		createABid(30000000, 100, 15000000),
		createABid(40000000, 100, 20000000),
	}

	nodePrice, err := stakingAuctionSC.calculateNodePrice(bids)
	assert.Equal(t, expectedNodePrice, nodePrice)
	assert.Nil(t, err)
}

func TestAuctionSC_calculateNodePrice_Case3(t *testing.T) {
	t.Parallel()

	expectedNodePrice := big.NewInt(12500000)
	args := createMockArgumentsForAuction()
	args.NumNodes = 5
	stakingAuctionSC, _ := NewStakingAuctionSmartContract(args)

	bids := []AuctionData{
		createABid(25000000, 2, 12500000),
		createABid(30000000, 3, 10000000),
		createABid(40000000, 2, 20000000),
		createABid(50000000, 2, 25000000),
	}

	nodePrice, err := stakingAuctionSC.calculateNodePrice(bids)
	assert.Equal(t, expectedNodePrice, nodePrice)
	assert.Nil(t, err)
}

func TestAuctionSC_calculateNodePrice_Case4ShouldErr(t *testing.T) {
	t.Parallel()

	stakingAuctionSC, _ := NewStakingAuctionSmartContract(createMockArgumentsForAuction())

	bid1 := createABid(25000000, 2, 12500000)
	bid2 := createABid(30000000, 3, 10000000)
	bid3 := createABid(40000000, 2, 20000000)
	bid4 := createABid(50000000, 2, 25000000)

	bids := []AuctionData{
		bid1, bid2, bid3, bid4,
	}

	nodePrice, err := stakingAuctionSC.calculateNodePrice(bids)
	assert.Nil(t, nodePrice)
	assert.Equal(t, vm.ErrNotEnoughQualifiedNodes, err)
}

func TestAuctionSC_selection_StakeGetAllocatedSeats(t *testing.T) {
	t.Parallel()

	args := createMockArgumentsForAuction()
	args.NumNodes = 5
	stakingAuctionSC, _ := NewStakingAuctionSmartContract(args)

	bid1 := createABid(25000000, 2, 12500000)
	bid2 := createABid(30000000, 3, 10000000)
	bid3 := createABid(40000000, 2, 20000000)
	bid4 := createABid(50000000, 2, 25000000)

	bids := []AuctionData{
		bid1, bid2, bid3, bid4,
	}

	// verify at least one is qualified from everybody
	expectedKeys := [][]byte{bid1.BlsPubKeys[0], bid3.BlsPubKeys[0], bid4.BlsPubKeys[0]}

	data := stakingAuctionSC.selection(bids)
	checkExpectedKeys(t, expectedKeys, data, len(expectedKeys))
}

func TestAuctionSC_selection_FirstBidderShouldTake50Percents(t *testing.T) {
	t.Parallel()

	args := createMockArgumentsForAuction()
	args.MinStepValue = big.NewInt(100000)
	args.MinStakeValue = big.NewInt(1)
	stakingAuctionSC, _ := NewStakingAuctionSmartContract(args)

	bids := []AuctionData{
		createABid(10000000, 10, 10000000),
		createABid(1000000, 1, 1000000),
		createABid(1000000, 1, 1000000),
		createABid(1000000, 1, 1000000),
		createABid(1000000, 1, 1000000),
		createABid(1000000, 1, 1000000),
		createABid(1000000, 1, 1000000),
		createABid(1000000, 1, 1000000),
		createABid(1000000, 1, 1000000),
		createABid(1000000, 1, 1000000),
		createABid(1000000, 1, 1000000),
	}

	data := stakingAuctionSC.selection(bids)
	//check that 50% keys belong to the first bidder
	checkExpectedKeys(t, bids[0].BlsPubKeys, data, 5)
}

func checkExpectedKeys(t *testing.T, expectedKeys [][]byte, data [][]byte, expectedNum int) {
	count := 0
	for _, key := range data {
		for _, expectedKey := range expectedKeys {
			if bytes.Equal(key, expectedKey) {
				count++
				break
			}
		}
	}
	assert.Equal(t, expectedNum, count)
}

func TestAuctionSC_selection_FirstBidderTakesAll(t *testing.T) {
	t.Parallel()

	stakingAuctionSC, _ := NewStakingAuctionSmartContract(createMockArgumentsForAuction())

	bids := []AuctionData{
		createABid(100000000, 10, 10000000),
		createABid(1000000, 1, 1000000),
		createABid(1000000, 1, 1000000),
		createABid(1000000, 1, 1000000),
		createABid(1000000, 1, 1000000),
		createABid(1000000, 1, 1000000),
		createABid(1000000, 1, 1000000),
		createABid(1000000, 1, 1000000),
		createABid(1000000, 1, 1000000),
		createABid(1000000, 1, 1000000),
		createABid(1000000, 1, 1000000),
	}

	data := stakingAuctionSC.selection(bids)
	//check that 100% keys belong to the first bidder
	checkExpectedKeys(t, bids[0].BlsPubKeys, data, 10)
}

func TestStakingAuctionSC_ExecuteStakeWithoutArgumentsShouldWork(t *testing.T) {
	t.Parallel()

	arguments := CreateVmContractCallInput()
	auctionData := createABid(25000000, 2, 12500000)
	auctionDataBytes, _ := json.Marshal(&auctionData)

	eei := &mock.SystemEIStub{}
	eei.GetStorageCalled = func(key []byte) []byte {
		if bytes.Equal(key, arguments.CallerAddr) {
			return auctionDataBytes
		}
		return nil
	}
	eei.SetStorageCalled = func(key []byte, value []byte) {
		if bytes.Equal(key, arguments.CallerAddr) {
			var auctionData AuctionData
			_ = json.Unmarshal(value, &auctionData)
			assert.Equal(t, big.NewInt(26000000), auctionData.TotalStakeValue)
		}
	}
	args := createMockArgumentsForAuction()
	args.Eei = eei
	args.NumNodes = uint32(5)

	stakingAuctionSC, _ := NewStakingAuctionSmartContract(args)

	arguments.Function = "stake"
	arguments.CallValue = big.NewInt(1000000)

	errCode := stakingAuctionSC.Execute(arguments)
	assert.Equal(t, vmcommon.Ok, errCode)
}

func TestStakingAuctionSC_ExecuteStakeAddedNewPubKeysShouldWork(t *testing.T) {
	t.Parallel()

	arguments := CreateVmContractCallInput()
	auctionData := createABid(25000000, 2, 12500000)
	auctionDataBytes, _ := json.Marshal(&auctionData)

	key1 := []byte("Key1")
	key2 := []byte("Key2")
	rewardAddr := []byte("tralala2")
	maxStakePerNoce := big.NewInt(500)

	eei := &mock.SystemEIStub{}
	eei.GetStorageCalled = func(key []byte) []byte {
		if bytes.Equal(key, arguments.CallerAddr) {
			return auctionDataBytes
		}
		return nil
	}
	eei.SetStorageCalled = func(key []byte, value []byte) {
		if bytes.Equal(key, arguments.CallerAddr) {
			var auctionData AuctionData
			_ = json.Unmarshal(value, &auctionData)
			assert.Equal(t, big.NewInt(26000000), auctionData.TotalStakeValue)
			assert.True(t, bytes.Equal(auctionData.BlsPubKeys[2], key1))
			assert.True(t, bytes.Equal(auctionData.BlsPubKeys[3], key2))
			assert.True(t, bytes.Equal(auctionData.RewardAddress, rewardAddr))
			assert.Equal(t, maxStakePerNoce, auctionData.MaxStakePerNode)
		}
	}

	args := createMockArgumentsForAuction()
	args.Eei = eei
	args.NumNodes = uint32(5)

	stakingAuctionSC, _ := NewStakingAuctionSmartContract(args)

	arguments.Function = "stake"
	arguments.CallValue = big.NewInt(1000000)
	arguments.Arguments = [][]byte{big.NewInt(2).Bytes(), key1, []byte("msg1"), key2, []byte("msg2"), maxStakePerNoce.Bytes(), rewardAddr}

	errCode := stakingAuctionSC.Execute(arguments)
	assert.Equal(t, vmcommon.Ok, errCode)
}

func TestStakingAuctionSC_ExecuteStakeUnStakeOneBlsPubKey(t *testing.T) {
	t.Parallel()

	arguments := CreateVmContractCallInput()
	auctionData := createABid(25000000, 2, 12500000)
	auctionDataBytes, _ := json.Marshal(&auctionData)

	stakedData := StakedData{
		RegisterNonce: 0,
		Staked:        true,
		UnStakedNonce: 1,
		UnStakedEpoch: 0,
		RewardAddress: []byte("tralala1"),
		StakeValue:    big.NewInt(0),
	}
	stakedDataBytes, _ := json.Marshal(&stakedData)

	eei := &mock.SystemEIStub{}
	eei.GetStorageCalled = func(key []byte) []byte {
		if bytes.Equal(key, arguments.CallerAddr) {
			return auctionDataBytes
		}
		if bytes.Equal(key, auctionData.BlsPubKeys[0]) {
			return stakedDataBytes
		}
		return nil
	}
	eei.SetStorageCalled = func(key []byte, value []byte) {
		var stakedData StakedData
		_ = json.Unmarshal(value, &stakedData)

		assert.Equal(t, false, stakedData.Staked)
	}

	args := createMockArgumentsForAuction()
	args.Eei = eei
	args.NumNodes = uint32(5)

	stakingAuctionSC, _ := NewStakingAuctionSmartContract(args)

	arguments.Function = "unStake"
	arguments.Arguments = [][]byte{auctionData.BlsPubKeys[0]}
	errCode := stakingAuctionSC.Execute(arguments)
	assert.Equal(t, vmcommon.Ok, errCode)
}

func TestStakingAuctionSC_ExecuteUnBound(t *testing.T) {
	t.Parallel()

	minStakeValue := big.NewInt(10000000)
	arguments := CreateVmContractCallInput()
	totalStake := uint64(25000000)

	auctionData := createABid(totalStake, 2, 12500000)
	auctionDataBytes, _ := json.Marshal(&auctionData)

	stakedData := StakedData{
		RegisterNonce: 0,
		Staked:        false,
		UnStakedNonce: 1,
		UnStakedEpoch: 0,
		RewardAddress: []byte("tralala1"),
		StakeValue:    big.NewInt(12500000),
	}
	stakedDataBytes, _ := json.Marshal(&stakedData)

	eei := &mock.SystemEIStub{}
	eei.GetStorageCalled = func(key []byte) []byte {
		if bytes.Equal(arguments.CallerAddr, key) {
			return auctionDataBytes
		}
		if bytes.Equal(key, auctionData.BlsPubKeys[0]) {
			return stakedDataBytes
		}

		return nil
	}

	eei.SetStorageCalled = func(key []byte, value []byte) {
		if bytes.Contains(key, arguments.CallerAddr) {
			var regData AuctionData
			_ = json.Unmarshal(value, &regData)

			assert.Equal(t, big.NewInt(0).Sub(big.NewInt(0).SetUint64(totalStake), minStakeValue), regData.LockedStake)
			assert.Equal(t, big.NewInt(0).Sub(big.NewInt(0).SetUint64(totalStake), minStakeValue), regData.TotalStakeValue)
		}
	}

	args := createMockArgumentsForAuction()
	args.Eei = eei
	args.NumNodes = uint32(5)

	stakingAuctionSC, _ := NewStakingAuctionSmartContract(args)

	arguments.Function = "unBond"
	arguments.Arguments = [][]byte{auctionData.BlsPubKeys[0]}
	errCode := stakingAuctionSC.Execute(arguments)
	assert.Equal(t, vmcommon.Ok, errCode)
}
