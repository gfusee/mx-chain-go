package track

import (
	"github.com/multiversx/mx-chain-core-go/hashing"
	"github.com/multiversx/mx-chain-core-go/marshal"

	"github.com/multiversx/mx-chain-go/common"
	"github.com/multiversx/mx-chain-go/dataRetriever"
	"github.com/multiversx/mx-chain-go/process"
	"github.com/multiversx/mx-chain-go/sharding"
)

// ArgBlockProcessor holds all dependencies required to process tracked blocks in order to create new instances of
// block processor
type ArgBlockProcessor struct {
	HeaderValidator                       process.HeaderConstructionValidator
	RequestHandler                        process.RequestHandler
	ShardCoordinator                      sharding.Coordinator
	BlockTracker                          blockTrackerHandler
	CrossNotarizer                        blockNotarizerHandler
	SelfNotarizer                         blockNotarizerHandler
	CrossNotarizedHeadersNotifier         blockNotifierHandler
	SelfNotarizedFromCrossHeadersNotifier blockNotifierHandler
	SelfNotarizedHeadersNotifier          blockNotifierHandler
	FinalMetachainHeadersNotifier         blockNotifierHandler
	RoundHandler                          process.RoundHandler
	EnableEpochsHandler                   common.EnableEpochsHandler
	ProofsPool                            process.ProofsPool
	Marshaller                            marshal.Marshalizer
	Hasher                                hashing.Hasher
	HeadersPool                           dataRetriever.HeadersPool
	IsImportDBMode                        bool
}
