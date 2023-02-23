package bls

import (
	"fmt"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-go/consensus"
	"github.com/multiversx/mx-chain-go/consensus/spos"
	"github.com/multiversx/mx-chain-go/errors"
	"github.com/multiversx/mx-chain-go/outport"
)

// factory defines the data needed by this factory to create all the subrounds and give them their specific
// functionality
type factory struct {
	consensusCore  spos.ConsensusCoreHandler
	consensusState *spos.ConsensusState
	worker         spos.WorkerHandler

	appStatusHandler core.AppStatusHandler
	outportHandler   outport.OutportHandler
	chainID          []byte
	currentPid       core.PeerID
	consensusModel   consensus.ConsensusModel
}

// NewSubroundsFactory creates a new factory object
func NewSubroundsFactory(
	consensusDataContainer spos.ConsensusCoreHandler,
	consensusState *spos.ConsensusState,
	worker spos.WorkerHandler,
	chainID []byte,
	currentPid core.PeerID,
	appStatusHandler core.AppStatusHandler,
	consensusModel consensus.ConsensusModel,
) (*factory, error) {
	err := checkNewFactoryParams(
		consensusDataContainer,
		consensusState,
		worker,
		chainID,
		appStatusHandler,
	)
	if err != nil {
		return nil, err
	}

	fct := factory{
		consensusCore:    consensusDataContainer,
		consensusState:   consensusState,
		worker:           worker,
		appStatusHandler: appStatusHandler,
		chainID:          chainID,
		currentPid:       currentPid,
		consensusModel:   consensusModel,
	}

	return &fct, nil
}

func checkNewFactoryParams(
	container spos.ConsensusCoreHandler,
	state *spos.ConsensusState,
	worker spos.WorkerHandler,
	chainID []byte,
	appStatusHandler core.AppStatusHandler,
) error {
	err := spos.ValidateConsensusCore(container)
	if err != nil {
		return err
	}
	if state == nil {
		return spos.ErrNilConsensusState
	}
	if check.IfNil(worker) {
		return spos.ErrNilWorker
	}
	if check.IfNil(appStatusHandler) {
		return spos.ErrNilAppStatusHandler
	}
	if len(chainID) == 0 {
		return spos.ErrInvalidChainID
	}

	return nil
}

// SetOutportHandler method will update the value of the factory's outport
func (fct *factory) SetOutportHandler(driver outport.OutportHandler) {
	fct.outportHandler = driver
}

// GenerateSubrounds will generate the subrounds used in BLS consensus
func (fct *factory) GenerateSubrounds() error {
	fct.initConsensusThreshold()
	fct.consensusCore.Chronology().RemoveAllSubrounds()
	fct.worker.RemoveAllReceivedMessagesCalls()

	err := fct.generateStartRoundSubround()
	if err != nil {
		return err
	}

	err = fct.generateBlockSubround()
	if err != nil {
		return err
	}

	err = fct.generateSignatureSubround()
	if err != nil {
		return err
	}

	err = fct.generateEndRoundSubround()
	if err != nil {
		return err
	}

	return nil
}

func (fct *factory) getTimeDuration() time.Duration {
	return fct.consensusCore.RoundHandler().TimeDuration()
}

func (fct *factory) generateStartRoundSubround() error {
	subround, err := spos.NewSubround(
		-1,
		SrStartRound,
		SrBlock,
		int64(float64(fct.getTimeDuration())*srStartStartTime),
		int64(float64(fct.getTimeDuration())*srStartEndTime),
		getSubroundName(SrStartRound),
		fct.consensusState,
		fct.worker.GetConsensusStateChangedChannel(),
		fct.worker.ExecuteStoredMessages,
		fct.consensusCore,
		fct.chainID,
		fct.currentPid,
		fct.appStatusHandler,
	)
	if err != nil {
		return err
	}

	subroundStartRoundInstance, err := NewSubroundStartRound(
		subround,
		fct.worker.Extend,
		processingThresholdPercent,
		fct.worker.ExecuteStoredMessages,
		fct.worker.ResetConsensusMessages,
	)
	if err != nil {
		return err
	}

	err = subroundStartRoundInstance.SetOutportHandler(fct.outportHandler)
	if err != nil {
		return err
	}

	fct.consensusCore.Chronology().AddSubround(subroundStartRoundInstance)

	return nil
}

func (fct *factory) generateBlockSubround() error {
	subroundBlockInstance, err := fct.generateBlockSubroundV1()
	if err != nil {
		return err
	}

	switch fct.consensusModel {
	case consensus.ConsensusModelV1:
		fct.worker.AddReceivedMessageCall(MtBlockBodyAndHeader, subroundBlockInstance.receivedBlockBodyAndHeader)
		fct.worker.AddReceivedMessageCall(MtBlockBody, subroundBlockInstance.receivedBlockBody)
		fct.worker.AddReceivedMessageCall(MtBlockHeader, subroundBlockInstance.receivedBlockHeader)
		fct.consensusCore.Chronology().AddSubround(subroundBlockInstance)

		return nil
	case consensus.ConsensusModelV2:
		subroundBlockV2Instance, errV2 := NewSubroundBlockV2(subroundBlockInstance)
		if errV2 != nil {
			return errV2
		}

		fct.worker.AddReceivedMessageCall(MtBlockBodyAndHeader, subroundBlockV2Instance.receivedBlockBodyAndHeader)
		fct.worker.AddReceivedMessageCall(MtBlockBody, subroundBlockV2Instance.receivedBlockBody)
		fct.worker.AddReceivedMessageCall(MtBlockHeader, subroundBlockV2Instance.receivedBlockHeader)
		fct.consensusCore.Chronology().AddSubround(subroundBlockV2Instance)

		return nil
	default:
		return fmt.Errorf("%w model %v", errors.ErrUnimplementedConsensusModel, fct.consensusModel)
	}
}

func (fct *factory) generateBlockSubroundV1() (*subroundBlock, error) {
	subround, err := spos.NewSubround(
		SrStartRound,
		SrBlock,
		SrSignature,
		int64(float64(fct.getTimeDuration())*srBlockStartTime),
		int64(float64(fct.getTimeDuration())*srBlockEndTime),
		getSubroundName(SrBlock),
		fct.consensusState,
		fct.worker.GetConsensusStateChangedChannel(),
		fct.worker.ExecuteStoredMessages,
		fct.consensusCore,
		fct.chainID,
		fct.currentPid,
		fct.appStatusHandler,
	)
	if err != nil {
		return nil, err
	}

	subroundBlockInstance, err := NewSubroundBlock(
		subround,
		fct.worker.Extend,
		processingThresholdPercent,
	)
	if err != nil {
		return nil, err
	}

	return subroundBlockInstance, nil
}

func (fct *factory) generateSignatureSubround() error {
	subroundSignatureInstance, err := fct.generateSignatureSubroundV1()
	if err != nil {
		return err
	}

	switch fct.consensusModel {
	case consensus.ConsensusModelV1:
		fct.worker.AddReceivedMessageCall(MtSignature, subroundSignatureInstance.receivedSignature)
		fct.consensusCore.Chronology().AddSubround(subroundSignatureInstance)

		return nil
	case consensus.ConsensusModelV2:
		subroundSignatureV2Instance, errV2 := NewSubroundSignatureV2(subroundSignatureInstance)
		if errV2 != nil {
			return errV2
		}

		fct.worker.AddReceivedMessageCall(MtSignature, subroundSignatureV2Instance.receivedSignature)
		fct.consensusCore.Chronology().AddSubround(subroundSignatureV2Instance)

		return nil
	default:
		return fmt.Errorf("%w model %v", errors.ErrUnimplementedConsensusModel, fct.consensusModel)
	}
}

func (fct *factory) generateSignatureSubroundV1() (*subroundSignature, error) {
	subround, err := spos.NewSubround(
		SrBlock,
		SrSignature,
		SrEndRound,
		int64(float64(fct.getTimeDuration())*srSignatureStartTime),
		int64(float64(fct.getTimeDuration())*srSignatureEndTime),
		getSubroundName(SrSignature),
		fct.consensusState,
		fct.worker.GetConsensusStateChangedChannel(),
		fct.worker.ExecuteStoredMessages,
		fct.consensusCore,
		fct.chainID,
		fct.currentPid,
		fct.appStatusHandler,
	)
	if err != nil {
		return nil, err
	}

	subroundSignatureInstance, err := NewSubroundSignature(
		subround,
		fct.worker.Extend,
		fct.appStatusHandler,
	)
	if err != nil {
		return nil, err
	}

	return subroundSignatureInstance, nil
}

func (fct *factory) generateEndRoundSubround() error {
	subroundEndRoundInstance, err := fct.generateEndRoundSubroundV1()
	if err != nil {
		return err
	}

	switch fct.consensusModel {
	case consensus.ConsensusModelV1:
		fct.worker.AddReceivedMessageCall(MtBlockHeaderFinalInfo, subroundEndRoundInstance.receivedBlockHeaderFinalInfo)
		fct.worker.AddReceivedMessageCall(MtInvalidSigners, subroundEndRoundInstance.receivedInvalidSignersInfo)
		fct.worker.AddReceivedHeaderHandler(subroundEndRoundInstance.receivedHeader)
		fct.consensusCore.Chronology().AddSubround(subroundEndRoundInstance)

		return nil
	case consensus.ConsensusModelV2:
		subroundSignatureV2Instance, errV2 := NewSubroundEndRoundV2(subroundEndRoundInstance)
		if errV2 != nil {
			return errV2
		}

		fct.worker.AddReceivedMessageCall(MtBlockHeaderFinalInfo, subroundSignatureV2Instance.receivedBlockHeaderFinalInfo)
		fct.worker.AddReceivedMessageCall(MtInvalidSigners, subroundSignatureV2Instance.receivedInvalidSignersInfo)
		fct.worker.AddReceivedHeaderHandler(subroundSignatureV2Instance.receivedHeader)
		fct.consensusCore.Chronology().AddSubround(subroundSignatureV2Instance)

		return nil
	default:
		return fmt.Errorf("%w model %v", errors.ErrUnimplementedConsensusModel, fct.consensusModel)
	}
}

func (fct *factory) generateEndRoundSubroundV1() (*subroundEndRound, error) {
	subround, err := spos.NewSubround(
		SrSignature,
		SrEndRound,
		-1,
		int64(float64(fct.getTimeDuration())*srEndStartTime),
		int64(float64(fct.getTimeDuration())*srEndEndTime),
		getSubroundName(SrEndRound),
		fct.consensusState,
		fct.worker.GetConsensusStateChangedChannel(),
		fct.worker.ExecuteStoredMessages,
		fct.consensusCore,
		fct.chainID,
		fct.currentPid,
		fct.appStatusHandler,
	)
	if err != nil {
		return nil, err
	}

	subroundEndRoundInstance, err := NewSubroundEndRound(
		subround,
		fct.worker.Extend,
		spos.MaxThresholdPercent,
		fct.worker.DisplayStatistics,
		fct.appStatusHandler,
	)
	if err != nil {
		return nil, err
	}

	return subroundEndRoundInstance, nil
}

func (fct *factory) initConsensusThreshold() {
	pBFTThreshold := core.GetPBFTThreshold(fct.consensusState.ConsensusGroupSize())
	pBFTFallbackThreshold := core.GetPBFTFallbackThreshold(fct.consensusState.ConsensusGroupSize())
	fct.consensusState.SetThreshold(SrBlock, 1)
	fct.consensusState.SetThreshold(SrSignature, pBFTThreshold)
	fct.consensusState.SetFallbackThreshold(SrBlock, 1)
	fct.consensusState.SetFallbackThreshold(SrSignature, pBFTFallbackThreshold)
}

// IsInterfaceNil returns true if there is no value under the interface
func (fct *factory) IsInterfaceNil() bool {
	return fct == nil
}
