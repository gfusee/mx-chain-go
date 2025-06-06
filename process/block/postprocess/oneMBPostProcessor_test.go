package postprocess

import (
	"bytes"
	"sort"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/stretchr/testify/assert"

	"github.com/multiversx/mx-chain-go/dataRetriever"
	"github.com/multiversx/mx-chain-go/process"
	"github.com/multiversx/mx-chain-go/process/mock"
	"github.com/multiversx/mx-chain-go/testscommon/economicsmocks"
	"github.com/multiversx/mx-chain-go/testscommon/hashingMocks"
	"github.com/multiversx/mx-chain-go/testscommon/storage"
)

func TestNewOneMBPostProcessor_NilHasher(t *testing.T) {
	t.Parallel()

	irp, err := NewOneMiniBlockPostProcessor(
		nil,
		&mock.MarshalizerMock{},
		mock.NewMultiShardsCoordinatorMock(5),
		&storage.ChainStorerStub{},
		block.TxBlock,
		dataRetriever.TransactionUnit,
		&economicsmocks.EconomicsHandlerMock{},
	)

	assert.Nil(t, irp)
	assert.Equal(t, process.ErrNilHasher, err)
}

func TestNewOneMBPostProcessor_NilMarshalizer(t *testing.T) {
	t.Parallel()

	irp, err := NewOneMiniBlockPostProcessor(
		&hashingMocks.HasherMock{},
		nil,
		mock.NewMultiShardsCoordinatorMock(5),
		&storage.ChainStorerStub{},
		block.TxBlock,
		dataRetriever.TransactionUnit,
		&economicsmocks.EconomicsHandlerMock{},
	)

	assert.Nil(t, irp)
	assert.Equal(t, process.ErrNilMarshalizer, err)
}

func TestNewOneMBPostProcessor_NilShardCoord(t *testing.T) {
	t.Parallel()

	irp, err := NewOneMiniBlockPostProcessor(
		&hashingMocks.HasherMock{},
		&mock.MarshalizerMock{},
		nil,
		&storage.ChainStorerStub{},
		block.TxBlock,
		dataRetriever.TransactionUnit,
		&economicsmocks.EconomicsHandlerMock{},
	)

	assert.Nil(t, irp)
	assert.Equal(t, process.ErrNilShardCoordinator, err)
}

func TestNewOneMBPostProcessor_NilStorer(t *testing.T) {
	t.Parallel()

	irp, err := NewOneMiniBlockPostProcessor(
		&hashingMocks.HasherMock{},
		&mock.MarshalizerMock{},
		mock.NewMultiShardsCoordinatorMock(5),
		nil,
		block.TxBlock,
		dataRetriever.TransactionUnit,
		&economicsmocks.EconomicsHandlerMock{},
	)

	assert.Nil(t, irp)
	assert.Equal(t, process.ErrNilStorage, err)
}

func TestNewOneMBPostProcessor_NilEconomicsFeeHandler(t *testing.T) {
	t.Parallel()

	irp, err := NewOneMiniBlockPostProcessor(
		&hashingMocks.HasherMock{},
		&mock.MarshalizerMock{},
		mock.NewMultiShardsCoordinatorMock(5),
		&storage.ChainStorerStub{},
		block.TxBlock,
		dataRetriever.TransactionUnit,
		nil,
	)

	assert.Nil(t, irp)
	assert.Equal(t, process.ErrNilEconomicsFeeHandler, err)
}

func TestNewOneMBPostProcessor_OK(t *testing.T) {
	t.Parallel()

	irp, err := NewOneMiniBlockPostProcessor(
		&hashingMocks.HasherMock{},
		&mock.MarshalizerMock{},
		mock.NewMultiShardsCoordinatorMock(5),
		&storage.ChainStorerStub{},
		block.TxBlock,
		dataRetriever.TransactionUnit,
		&economicsmocks.EconomicsHandlerMock{},
	)

	assert.Nil(t, err)
	assert.NotNil(t, irp)
}

func TestOneMBPostProcessor_CreateAllInterMiniBlocks(t *testing.T) {
	t.Parallel()

	irp, _ := NewOneMiniBlockPostProcessor(
		&hashingMocks.HasherMock{},
		&mock.MarshalizerMock{},
		mock.NewMultiShardsCoordinatorMock(5),
		&storage.ChainStorerStub{},
		block.TxBlock,
		dataRetriever.TransactionUnit,
		&economicsmocks.EconomicsHandlerMock{},
	)

	mbs := irp.CreateAllInterMiniBlocks()
	assert.Equal(t, 0, len(mbs))
}

func TestOneMBPostProcessor_CreateAllInterMiniBlocksOneMinBlock(t *testing.T) {
	t.Parallel()

	irp, _ := NewOneMiniBlockPostProcessor(
		&hashingMocks.HasherMock{},
		&mock.MarshalizerMock{},
		mock.NewMultiShardsCoordinatorMock(5),
		&storage.ChainStorerStub{},
		block.TxBlock,
		dataRetriever.TransactionUnit,
		&economicsmocks.EconomicsHandlerMock{},
	)

	txs := make([]data.TransactionHandler, 0)
	txs = append(txs, &transaction.Transaction{})
	txs = append(txs, &transaction.Transaction{})

	// with no InitProcessedResults, means that the transactions are added as scheduled transactions, not as
	// processing results from the execution of other transactions or miniblocks
	err := irp.AddIntermediateTransactions(txs, nil)
	assert.Nil(t, err)

	mbs := irp.CreateAllInterMiniBlocks()
	assert.Equal(t, 1, len(mbs))
}

func TestOneMBPostProcessor_VerifyNilBody(t *testing.T) {
	t.Parallel()

	irp, _ := NewOneMiniBlockPostProcessor(
		&hashingMocks.HasherMock{},
		&mock.MarshalizerMock{},
		mock.NewMultiShardsCoordinatorMock(5),
		&storage.ChainStorerStub{},
		block.TxBlock,
		dataRetriever.TransactionUnit,
		&economicsmocks.EconomicsHandlerMock{},
	)

	err := irp.VerifyInterMiniBlocks(&block.Body{})
	assert.Nil(t, err)
}

func TestOneMBPostProcessor_VerifyTooManyBlock(t *testing.T) {
	t.Parallel()

	irp, _ := NewOneMiniBlockPostProcessor(
		&hashingMocks.HasherMock{},
		&mock.MarshalizerMock{},
		mock.NewMultiShardsCoordinatorMock(5),
		&storage.ChainStorerStub{},
		block.TxBlock,
		dataRetriever.TransactionUnit,
		&economicsmocks.EconomicsHandlerMock{},
	)

	txs := make([]data.TransactionHandler, 0)
	txs = append(txs, &transaction.Transaction{SndAddr: []byte("snd"), RcvAddr: []byte("recvaddr1")})
	txs = append(txs, &transaction.Transaction{SndAddr: []byte("snd"), RcvAddr: []byte("recvaddr2")})
	txs = append(txs, &transaction.Transaction{SndAddr: []byte("snd"), RcvAddr: []byte("recvaddr3")})
	txs = append(txs, &transaction.Transaction{SndAddr: []byte("snd"), RcvAddr: []byte("recvaddr4")})
	txs = append(txs, &transaction.Transaction{SndAddr: []byte("snd"), RcvAddr: []byte("recvaddr5")})

	err := irp.AddIntermediateTransactions(txs, nil)
	assert.Nil(t, err)

	miniBlock := &block.MiniBlock{
		SenderShardID:   0,
		ReceiverShardID: 0,
		Type:            block.TxBlock}

	for i := 0; i < len(txs); i++ {
		txHash, _ := core.CalculateHash(&mock.MarshalizerMock{}, &hashingMocks.HasherMock{}, txs[i])
		miniBlock.TxHashes = append(miniBlock.TxHashes, txHash)
	}

	sort.Slice(miniBlock.TxHashes, func(a, b int) bool {
		return bytes.Compare(miniBlock.TxHashes[a], miniBlock.TxHashes[b]) < 0
	})

	body := &block.Body{}
	body.MiniBlocks = append(body.MiniBlocks, miniBlock)
	body.MiniBlocks = append(body.MiniBlocks, miniBlock)

	err = irp.VerifyInterMiniBlocks(body)
	assert.Equal(t, process.ErrTooManyReceiptsMiniBlocks, err)
}

func TestOneMBPostProcessor_VerifyNilMiniBlocks(t *testing.T) {
	t.Parallel()

	irp, _ := NewOneMiniBlockPostProcessor(
		&hashingMocks.HasherMock{},
		&mock.MarshalizerMock{},
		mock.NewMultiShardsCoordinatorMock(5),
		&storage.ChainStorerStub{},
		block.TxBlock,
		dataRetriever.TransactionUnit,
		&economicsmocks.EconomicsHandlerMock{},
	)

	miniBlock := &block.MiniBlock{
		SenderShardID:   0,
		ReceiverShardID: 0,
		Type:            block.TxBlock}
	body := &block.Body{}
	body.MiniBlocks = append(body.MiniBlocks, miniBlock)

	err := irp.VerifyInterMiniBlocks(body)
	assert.Equal(t, process.ErrNilMiniBlocks, err)
}

func TestOneMBPostProcessor_VerifyOk(t *testing.T) {
	t.Parallel()

	irp, _ := NewOneMiniBlockPostProcessor(
		&hashingMocks.HasherMock{},
		&mock.MarshalizerMock{},
		mock.NewMultiShardsCoordinatorMock(5),
		&storage.ChainStorerStub{},
		block.TxBlock,
		dataRetriever.TransactionUnit,
		&economicsmocks.EconomicsHandlerMock{},
	)

	txs := make([]data.TransactionHandler, 0)
	txs = append(txs, &transaction.Transaction{SndAddr: []byte("snd"), RcvAddr: []byte("recvaddr1")})
	txs = append(txs, &transaction.Transaction{SndAddr: []byte("snd"), RcvAddr: []byte("recvaddr2")})
	txs = append(txs, &transaction.Transaction{SndAddr: []byte("snd"), RcvAddr: []byte("recvaddr3")})
	txs = append(txs, &transaction.Transaction{SndAddr: []byte("snd"), RcvAddr: []byte("recvaddr4")})
	txs = append(txs, &transaction.Transaction{SndAddr: []byte("snd"), RcvAddr: []byte("recvaddr5")})

	err := irp.AddIntermediateTransactions(txs, nil)
	assert.Nil(t, err)

	miniBlock := &block.MiniBlock{
		SenderShardID:   0,
		ReceiverShardID: 0,
		Type:            block.TxBlock}

	for i := 0; i < len(txs); i++ {
		txHash, _ := core.CalculateHash(&mock.MarshalizerMock{}, &hashingMocks.HasherMock{}, txs[i])
		miniBlock.TxHashes = append(miniBlock.TxHashes, txHash)
	}

	sort.Slice(miniBlock.TxHashes, func(a, b int) bool {
		return bytes.Compare(miniBlock.TxHashes[a], miniBlock.TxHashes[b]) < 0
	})

	body := &block.Body{}
	body.MiniBlocks = append(body.MiniBlocks, miniBlock)

	err = irp.VerifyInterMiniBlocks(body)
	assert.Nil(t, err)
}
