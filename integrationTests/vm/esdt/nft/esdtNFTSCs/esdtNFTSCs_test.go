package esdtNFTSCs

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/stretchr/testify/require"

	"github.com/multiversx/mx-chain-go/integrationTests"
	"github.com/multiversx/mx-chain-go/integrationTests/vm/esdt"
	"github.com/multiversx/mx-chain-go/integrationTests/vm/esdt/nft"
)

func TestESDTNFTIssueCreateBurnSendViaAsyncViaExecuteOnSC(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}
	nodes, leaders := esdt.CreateNodesAndPrepareBalances(1)

	defer func() {
		for _, n := range nodes {
			n.Close()
		}
	}()

	initialVal := big.NewInt(10000000000)
	integrationTests.MintAllNodes(nodes, initialVal)

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	scAddress, tokenIdentifier := deployAndIssueNFTSFTThroughSC(t, nodes, leaders, &nonce, &round, "nftIssue", "@03@05")

	txData := []byte("nftCreate" + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + hex.EncodeToString(big.NewInt(1).Bytes()) + "@" + hex.EncodeToString([]byte("name")) +
		"@" + hex.EncodeToString(big.NewInt(10).Bytes()) + "@" + hex.EncodeToString(scAddress) +
		"@" + hex.EncodeToString([]byte("abc")) + "@" + hex.EncodeToString([]byte("NFT")))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		scAddress,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		scAddress,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		scAddress,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, 3, nonce, round)
	time.Sleep(time.Second)

	checkAddressHasNft(t, scAddress, scAddress, nodes, []byte(tokenIdentifier), 3, big.NewInt(1))
	checkAddressHasNft(t, scAddress, scAddress, nodes, []byte(tokenIdentifier), 2, big.NewInt(1))
	checkAddressHasNft(t, scAddress, scAddress, nodes, []byte(tokenIdentifier), 1, big.NewInt(1))

	txData = []byte("nftBurn" + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + hex.EncodeToString(big.NewInt(1).Bytes()) + "@" + hex.EncodeToString(big.NewInt(1).Bytes()))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		scAddress,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		scAddress,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, 3, nonce, round)
	time.Sleep(time.Second)

	checkAddressHasNft(t, scAddress, scAddress, nodes, []byte(tokenIdentifier), 2, big.NewInt(1))
	checkAddressHasNft(t, scAddress, scAddress, nodes, []byte(tokenIdentifier), 1, big.NewInt(0))

	destinationAddress := nodes[0].OwnAccount.Address
	txData = []byte("transferNftViaAsyncCall" + "@" + hex.EncodeToString(destinationAddress) +
		"@" + hex.EncodeToString([]byte(tokenIdentifier)) + "@" + hex.EncodeToString(big.NewInt(2).Bytes()) +
		"@" + hex.EncodeToString(big.NewInt(1).Bytes()) + "@" + hex.EncodeToString([]byte("NFT")))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		scAddress,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	txData = []byte("transfer_nft_and_execute" + "@" + hex.EncodeToString(destinationAddress) +
		"@" + hex.EncodeToString([]byte(tokenIdentifier)) + "@" + hex.EncodeToString(big.NewInt(3).Bytes()) +
		"@" + hex.EncodeToString(big.NewInt(1).Bytes()) + "@" + hex.EncodeToString([]byte("NFT")))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		scAddress,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)
	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, 3, nonce, round)
	time.Sleep(time.Second)

	checkAddressHasNft(t, scAddress, destinationAddress, nodes, []byte(tokenIdentifier), 2, big.NewInt(1))
	checkAddressHasNft(t, scAddress, destinationAddress, nodes, []byte(tokenIdentifier), 3, big.NewInt(1))
	checkAddressHasNft(t, scAddress, scAddress, nodes, []byte(tokenIdentifier), 2, big.NewInt(0))
	checkAddressHasNft(t, scAddress, scAddress, nodes, []byte(tokenIdentifier), 3, big.NewInt(0))
}

func TestESDTSemiFTIssueCreateBurnSendViaAsyncViaExecuteOnSC(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}
	nodes, leaders := esdt.CreateNodesAndPrepareBalances(1)

	defer func() {
		for _, n := range nodes {
			n.Close()
		}
	}()

	initialVal := big.NewInt(10000000000)
	integrationTests.MintAllNodes(nodes, initialVal)

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	scAddress, tokenIdentifier := deployAndIssueNFTSFTThroughSC(t, nodes, leaders, &nonce, &round, "sftIssue", "@03@04@05")

	txData := []byte("nftCreate" + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + hex.EncodeToString(big.NewInt(1).Bytes()) + "@" + hex.EncodeToString([]byte("name")) +
		"@" + hex.EncodeToString(big.NewInt(10).Bytes()) + "@" + hex.EncodeToString(scAddress) +
		"@" + hex.EncodeToString([]byte("abc")) + "@" + hex.EncodeToString([]byte("NFT")))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		scAddress,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	txData = []byte("nftAddQuantity" + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + hex.EncodeToString(big.NewInt(1).Bytes()) + "@" + hex.EncodeToString(big.NewInt(10).Bytes()))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		scAddress,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, 2, nonce, round)
	time.Sleep(time.Second)

	checkAddressHasNft(t, scAddress, scAddress, nodes, []byte(tokenIdentifier), 1, big.NewInt(11))

	txData = []byte("nftBurn" + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + hex.EncodeToString(big.NewInt(1).Bytes()) + "@" + hex.EncodeToString(big.NewInt(1).Bytes()))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		scAddress,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		scAddress,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, 2, nonce, round)
	time.Sleep(time.Second)

	checkAddressHasNft(t, scAddress, scAddress, nodes, []byte(tokenIdentifier), 1, big.NewInt(9))

	destinationAddress := nodes[0].OwnAccount.Address
	txData = []byte("transferNftViaAsyncCall" + "@" + hex.EncodeToString(destinationAddress) +
		"@" + hex.EncodeToString([]byte(tokenIdentifier)) + "@" + hex.EncodeToString(big.NewInt(1).Bytes()) +
		"@" + hex.EncodeToString(big.NewInt(5).Bytes()) + "@" + hex.EncodeToString([]byte("NFT")))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		scAddress,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	txData = []byte("transfer_nft_and_execute" + "@" + hex.EncodeToString(destinationAddress) +
		"@" + hex.EncodeToString([]byte(tokenIdentifier)) + "@" + hex.EncodeToString(big.NewInt(1).Bytes()) +
		"@" + hex.EncodeToString(big.NewInt(4).Bytes()) + "@" + hex.EncodeToString([]byte("NFT")))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		scAddress,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)
	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, 2, nonce, round)
	time.Sleep(time.Second)

	checkAddressHasNft(t, scAddress, destinationAddress, nodes, []byte(tokenIdentifier), 1, big.NewInt(9))
	checkAddressHasNft(t, scAddress, scAddress, nodes, []byte(tokenIdentifier), 1, big.NewInt(0))
}

func TestESDTTransferNFTBetweenContractsAcceptAndNotAcceptWithRevert(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}
	nodes, leaders := esdt.CreateNodesAndPrepareBalances(1)

	defer func() {
		for _, n := range nodes {
			n.Close()
		}
	}()

	initialVal := big.NewInt(10000000000)
	integrationTests.MintAllNodes(nodes, initialVal)

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	scAddress, tokenIdentifier := deployAndIssueNFTSFTThroughSC(t, nodes, leaders, &nonce, &round, "nftIssue", "@03@05")

	txData := []byte("nftCreate" + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + hex.EncodeToString(big.NewInt(1).Bytes()) + "@" + hex.EncodeToString([]byte("name")) +
		"@" + hex.EncodeToString(big.NewInt(10).Bytes()) + "@" + hex.EncodeToString(scAddress) +
		"@" + hex.EncodeToString([]byte("abc")) + "@" + hex.EncodeToString([]byte("NFT")))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		scAddress,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		scAddress,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, 2, nonce, round)
	time.Sleep(time.Second)

	checkAddressHasNft(t, scAddress, scAddress, nodes, []byte(tokenIdentifier), 2, big.NewInt(1))
	checkAddressHasNft(t, scAddress, scAddress, nodes, []byte(tokenIdentifier), 1, big.NewInt(1))

	destinationSCAddress := esdt.DeployNonPayableSmartContract(t, nodes, leaders, &nonce, &round, "../../testdata/nft-receiver.wasm")
	txData = []byte("transferNftViaAsyncCall" + "@" + hex.EncodeToString(destinationSCAddress) +
		"@" + hex.EncodeToString([]byte(tokenIdentifier)) + "@" + hex.EncodeToString(big.NewInt(1).Bytes()) +
		"@" + hex.EncodeToString(big.NewInt(1).Bytes()) + "@" + hex.EncodeToString([]byte("wrongFunctionToCall")))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		scAddress,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	txData = []byte("transfer_nft_and_execute" + "@" + hex.EncodeToString(destinationSCAddress) +
		"@" + hex.EncodeToString([]byte(tokenIdentifier)) + "@" + hex.EncodeToString(big.NewInt(2).Bytes()) +
		"@" + hex.EncodeToString(big.NewInt(1).Bytes()) + "@" + hex.EncodeToString([]byte("wrongFunctionToCall")))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		scAddress,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)
	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, 2, nonce, round)
	time.Sleep(time.Second)

	checkAddressHasNft(t, scAddress, destinationSCAddress, nodes, []byte(tokenIdentifier), 1, big.NewInt(0))
	checkAddressHasNft(t, scAddress, destinationSCAddress, nodes, []byte(tokenIdentifier), 2, big.NewInt(0))
	checkAddressHasNft(t, scAddress, scAddress, nodes, []byte(tokenIdentifier), 1, big.NewInt(1))
	checkAddressHasNft(t, scAddress, scAddress, nodes, []byte(tokenIdentifier), 2, big.NewInt(1))

	txData = []byte("transferNftViaAsyncCall" + "@" + hex.EncodeToString(destinationSCAddress) +
		"@" + hex.EncodeToString([]byte(tokenIdentifier)) + "@" + hex.EncodeToString(big.NewInt(1).Bytes()) +
		"@" + hex.EncodeToString(big.NewInt(1).Bytes()) + "@" + hex.EncodeToString([]byte("acceptAndReturnCallData")))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		scAddress,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)

	txData = []byte("transfer_nft_and_execute" + "@" + hex.EncodeToString(destinationSCAddress) +
		"@" + hex.EncodeToString([]byte(tokenIdentifier)) + "@" + hex.EncodeToString(big.NewInt(2).Bytes()) +
		"@" + hex.EncodeToString(big.NewInt(1).Bytes()) + "@" + hex.EncodeToString([]byte("acceptAndReturnCallData")))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		scAddress,
		string(txData),
		integrationTests.AdditionalGasLimit,
	)
	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, 2, nonce, round)
	time.Sleep(time.Second)

	checkAddressHasNft(t, scAddress, destinationSCAddress, nodes, []byte(tokenIdentifier), 1, big.NewInt(1))
	checkAddressHasNft(t, scAddress, destinationSCAddress, nodes, []byte(tokenIdentifier), 2, big.NewInt(1))
	checkAddressHasNft(t, scAddress, scAddress, nodes, []byte(tokenIdentifier), 1, big.NewInt(0))
	checkAddressHasNft(t, scAddress, scAddress, nodes, []byte(tokenIdentifier), 2, big.NewInt(0))
}

func TestESDTTransferNFTToSCIntraShard(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}
	nodes, leaders := esdt.CreateNodesAndPrepareBalances(1)

	defer func() {
		for _, n := range nodes {
			n.Close()
		}
	}()

	initialVal := big.NewInt(10000000000)
	integrationTests.MintAllNodes(nodes, initialVal)

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	roles := [][]byte{
		[]byte(core.ESDTRoleNFTCreate),
		[]byte(core.ESDTRoleNFTBurn),
	}
	tokenIdentifier, _ := nft.PrepareNFTWithRoles(
		t,
		nodes,
		leaders,
		nodes[0],
		&round,
		&nonce,
		core.NonFungibleESDT,
		1,
		roles,
	)

	nonceArg := hex.EncodeToString(big.NewInt(0).SetUint64(1).Bytes())
	quantityToTransfer := hex.EncodeToString(big.NewInt(1).Bytes())
	destinationSCAddress := esdt.DeployNonPayableSmartContract(t, nodes, leaders, &nonce, &round, "../../testdata/nft-receiver.wasm")
	txData := core.BuiltInFunctionESDTNFTTransfer + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + nonceArg + "@" + quantityToTransfer + "@" + hex.EncodeToString(destinationSCAddress) + "@" + hex.EncodeToString([]byte("acceptAndReturnCallData"))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		nodes[0].OwnAccount.Address,
		txData,
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, 3, nonce, round)
	time.Sleep(time.Second)

	checkAddressHasNft(t, nodes[0].OwnAccount.Address, destinationSCAddress, nodes, []byte(tokenIdentifier), 1, big.NewInt(1))
}

func TestESDTTransferNFTToSCCrossShard(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}
	nodes, leaders := esdt.CreateNodesAndPrepareBalances(2)

	defer func() {
		for _, n := range nodes {
			n.Close()
		}
	}()

	initialVal := big.NewInt(10000000000)
	integrationTests.MintAllNodes(nodes, initialVal)

	round := uint64(0)
	nonce := uint64(0)
	round = integrationTests.IncrementAndPrintRound(round)
	nonce++

	destinationSCAddress := esdt.DeployNonPayableSmartContract(t, nodes, leaders, &nonce, &round, "../../testdata/nft-receiver.wasm")

	destinationSCShardID := nodes[0].ShardCoordinator.ComputeId(destinationSCAddress)

	nodeFromOtherShard := nodes[1]
	for _, node := range nodes {
		shID := node.ShardCoordinator.ComputeId(node.OwnAccount.Address)
		if shID != destinationSCShardID {
			nodeFromOtherShard = node
			break
		}
	}

	roles := [][]byte{
		[]byte(core.ESDTRoleNFTCreate),
		[]byte(core.ESDTRoleNFTBurn),
	}
	tokenIdentifier, _ := nft.PrepareNFTWithRoles(
		t,
		nodes,
		leaders,
		nodeFromOtherShard,
		&round,
		&nonce,
		core.NonFungibleESDT,
		1,
		roles,
	)

	nonceArg := hex.EncodeToString(big.NewInt(0).SetUint64(1).Bytes())
	quantityToTransfer := hex.EncodeToString(big.NewInt(1).Bytes())

	txData := core.BuiltInFunctionESDTNFTTransfer + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + nonceArg + "@" + quantityToTransfer + "@" + hex.EncodeToString(destinationSCAddress) + "@" + hex.EncodeToString([]byte("acceptAndReturnCallData"))
	integrationTests.CreateAndSendTransaction(
		nodeFromOtherShard,
		nodes,
		big.NewInt(0),
		nodeFromOtherShard.OwnAccount.Address,
		txData,
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, 10, nonce, round)
	time.Sleep(time.Second)

	checkAddressHasNft(t, nodeFromOtherShard.OwnAccount.Address, destinationSCAddress, nodes, []byte(tokenIdentifier), 1, big.NewInt(1))

	txData = core.BuiltInFunctionESDTNFTTransfer + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + nonceArg + "@" + quantityToTransfer + "@" + hex.EncodeToString(destinationSCAddress) + "@" + hex.EncodeToString([]byte("wrongFunction"))
	integrationTests.CreateAndSendTransaction(
		nodeFromOtherShard,
		nodes,
		big.NewInt(0),
		nodeFromOtherShard.OwnAccount.Address,
		txData,
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, 10, nonce, round)
	time.Sleep(time.Second)

	checkAddressHasNft(t, nodeFromOtherShard.OwnAccount.Address, destinationSCAddress, nodes, []byte(tokenIdentifier), 1, big.NewInt(1))

	txData = core.BuiltInFunctionESDTNFTTransfer + "@" + hex.EncodeToString([]byte(tokenIdentifier)) +
		"@" + nonceArg + "@" + quantityToTransfer + "@" + hex.EncodeToString(destinationSCAddress)
	integrationTests.CreateAndSendTransaction(
		nodeFromOtherShard,
		nodes,
		big.NewInt(0),
		nodeFromOtherShard.OwnAccount.Address,
		txData,
		integrationTests.AdditionalGasLimit,
	)

	time.Sleep(time.Second)
	nonce, round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, 10, nonce, round)
	time.Sleep(time.Second)

	checkAddressHasNft(t, nodeFromOtherShard.OwnAccount.Address, destinationSCAddress, nodes, []byte(tokenIdentifier), 1, big.NewInt(1))
}

func deployAndIssueNFTSFTThroughSC(
	t *testing.T,
	nodes []*integrationTests.TestProcessorNode,
	leaders []*integrationTests.TestProcessorNode,
	nonce *uint64,
	round *uint64,
	issueFunc string,
	rolesEncoded string,
) ([]byte, string) {
	scAddress := esdt.DeployNonPayableSmartContract(t, nodes, leaders, nonce, round, "../../testdata/local-esdt-and-nft.wasm")

	issuePrice := big.NewInt(1000)
	txData := []byte(issueFunc + "@" + hex.EncodeToString([]byte("TOKEN")) +
		"@" + hex.EncodeToString([]byte("TKR")))
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		issuePrice,
		scAddress,
		string(txData),
		integrationTests.AdditionalGasLimit+core.MinMetaTxExtraGasCost,
	)

	time.Sleep(time.Second)
	nrRoundsToPropagateMultiShard := 12
	*nonce, *round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, *nonce, *round)
	time.Sleep(time.Second)

	tokenIdentifier := string(integrationTests.GetTokenIdentifier(nodes, []byte("TKR")))
	txData = []byte("setLocalRoles" + "@" + hex.EncodeToString(scAddress) +
		"@" + hex.EncodeToString([]byte(tokenIdentifier)) + rolesEncoded)
	integrationTests.CreateAndSendTransaction(
		nodes[0],
		nodes,
		big.NewInt(0),
		scAddress,
		string(txData),
		integrationTests.AdditionalGasLimit+core.MinMetaTxExtraGasCost,
	)

	time.Sleep(time.Second)
	*nonce, *round = integrationTests.WaitOperationToBeDone(t, leaders, nodes, nrRoundsToPropagateMultiShard, *nonce, *round)
	time.Sleep(time.Second)

	return scAddress, tokenIdentifier
}

func checkAddressHasNft(
	t *testing.T,
	creator []byte,
	address []byte,
	nodes []*integrationTests.TestProcessorNode,
	tickerID []byte,
	nonce uint64,
	quantity *big.Int,
) {
	esdtData := esdt.GetESDTTokenData(t, address, nodes, tickerID, nonce)

	if quantity.Cmp(big.NewInt(0)) == 0 {
		require.Nil(t, esdtData.TokenMetaData)
		return
	}

	require.NotNil(t, esdtData.TokenMetaData)
	require.Equal(t, creator, esdtData.TokenMetaData.Creator)
	require.Equal(t, quantity.Bytes(), esdtData.Value.Bytes())
}
