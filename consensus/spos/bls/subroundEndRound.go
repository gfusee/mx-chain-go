package bls

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/display"
	"github.com/multiversx/mx-chain-go/common"
	"github.com/multiversx/mx-chain-go/consensus"
	"github.com/multiversx/mx-chain-go/consensus/spos"
	"github.com/multiversx/mx-chain-go/p2p"
)

type subroundEndRound struct {
	*spos.Subround
	processingThresholdPercentage int
	displayStatistics             func()
	hasEquivalentProof            func(headerHash []byte) bool
	getValidatedEquivalentProof   func(headerHash []byte) ([]byte, []byte)
	appStatusHandler              core.AppStatusHandler
	mutProcessingEndRound         sync.Mutex
	sentSignatureTracker          spos.SentSignaturesTracker
}

// NewSubroundEndRound creates a subroundEndRound object
func NewSubroundEndRound(
	baseSubround *spos.Subround,
	extend func(subroundId int),
	processingThresholdPercentage int,
	displayStatistics func(),
	appStatusHandler core.AppStatusHandler,
	sentSignatureTracker spos.SentSignaturesTracker,
	hasEquivalentProof func(headerHash []byte) bool,
	getValidatedEquivalentProof func(headerHash []byte) ([]byte, []byte),
) (*subroundEndRound, error) {
	err := checkNewSubroundEndRoundParams(
		baseSubround,
	)
	if err != nil {
		return nil, err
	}
	if extend == nil {
		return nil, fmt.Errorf("%w for extend function", spos.ErrNilFunctionHandler)
	}
	if hasEquivalentProof == nil {
		return nil, fmt.Errorf("%w for hasEquivalentProof function", spos.ErrNilFunctionHandler)
	}
	if getValidatedEquivalentProof == nil {
		return nil, fmt.Errorf("%w for getValidatedEquivalentProof function", spos.ErrNilFunctionHandler)
	}
	if check.IfNil(appStatusHandler) {
		return nil, spos.ErrNilAppStatusHandler
	}
	if check.IfNil(sentSignatureTracker) {
		return nil, spos.ErrNilSentSignatureTracker
	}

	srEndRound := subroundEndRound{
		Subround:                      baseSubround,
		processingThresholdPercentage: processingThresholdPercentage,
		displayStatistics:             displayStatistics,
		hasEquivalentProof:            hasEquivalentProof,
		getValidatedEquivalentProof:   getValidatedEquivalentProof,
		appStatusHandler:              appStatusHandler,
		mutProcessingEndRound:         sync.Mutex{},
		sentSignatureTracker:          sentSignatureTracker,
	}
	srEndRound.Job = srEndRound.doEndRoundJob
	srEndRound.Check = srEndRound.doEndRoundConsensusCheck
	srEndRound.Extend = extend

	return &srEndRound, nil
}

func checkNewSubroundEndRoundParams(
	baseSubround *spos.Subround,
) error {
	if baseSubround == nil {
		return spos.ErrNilSubround
	}
	if baseSubround.ConsensusState == nil {
		return spos.ErrNilConsensusState
	}

	err := spos.ValidateConsensusCore(baseSubround.ConsensusCoreHandler)

	return err
}

// receivedBlockHeaderFinalInfo method is called when a block header final info is received
func (sr *subroundEndRound) receivedBlockHeaderFinalInfo(_ context.Context, cnsDta *consensus.Message) bool {
	node := string(cnsDta.PubKey)

	if !sr.IsConsensusDataSet() {
		return false
	}

	// TODO[cleanup cns finality]: update this check
	isSenderAllowed := sr.IsNodeInConsensusGroup(node)
	if !sr.EnableEpochsHandler().IsFlagEnabled(common.ConsensusPropagationChangesFlag) {
		isSenderAllowed = sr.IsNodeLeaderInCurrentRound(node)
	}
	if !isSenderAllowed { // is NOT this node leader in current round?
		sr.PeerHonestyHandler().ChangeScore(
			node,
			spos.GetConsensusTopicID(sr.ShardCoordinator()),
			spos.LeaderPeerHonestyDecreaseFactor,
		)

		return false
	}

	// TODO[cleanup cns finality]: update this check
	isSelfSender := node == sr.SelfPubKey()
	if !sr.EnableEpochsHandler().IsFlagEnabled(common.ConsensusPropagationChangesFlag) {
		isSelfSender = sr.IsSelfLeaderInCurrentRound() || sr.IsMultiKeyLeaderInCurrentRound()
	}
	if isSelfSender {
		return false
	}

	if !sr.IsConsensusDataEqual(cnsDta.BlockHeaderHash) {
		return false
	}

	if !sr.CanProcessReceivedMessage(cnsDta, sr.RoundHandler().Index(), sr.Current()) {
		return false
	}

	if !sr.isBlockHeaderFinalInfoValid(cnsDta) {
		return false
	}

	log.Debug("step 3: block header final info has been received",
		"PubKeysBitmap", cnsDta.PubKeysBitmap,
		"AggregateSignature", cnsDta.AggregateSignature,
		"LeaderSignature", cnsDta.LeaderSignature)

	signers := computeSignersPublicKeys(sr.ConsensusGroup(), cnsDta.PubKeysBitmap)
	sr.sentSignatureTracker.ReceivedActualSigners(signers)

	sr.PeerHonestyHandler().ChangeScore(
		node,
		spos.GetConsensusTopicID(sr.ShardCoordinator()),
		spos.LeaderPeerHonestyIncreaseFactor,
	)

	return sr.doEndRoundJobByParticipant(cnsDta)
}

func (sr *subroundEndRound) isBlockHeaderFinalInfoValid(cnsDta *consensus.Message) bool {
	if check.IfNil(sr.Header) {
		return false
	}

	header := sr.Header.ShallowClone()
	err := header.SetPubKeysBitmap(cnsDta.PubKeysBitmap)
	if err != nil {
		log.Debug("isBlockHeaderFinalInfoValid.SetPubKeysBitmap", "error", err.Error())
		return false
	}

	if !sr.EnableEpochsHandler().IsFlagEnabled(common.ConsensusPropagationChangesFlag) {
		return sr.verifySignatures(header, cnsDta)
	}

	err = sr.HeaderSigVerifier().VerifySignatureForHash(header, header.GetPrevHash(), cnsDta.Signature, cnsDta.Signature)
	if err != nil {
		log.Debug("isBlockHeaderFinalInfoValid.VerifySignatureForHash", "error", err.Error())
		return false
	}

	return true
}

func (sr *subroundEndRound) verifySignatures(header data.HeaderHandler, cnsDta *consensus.Message) bool {
	err := header.SetSignature(cnsDta.AggregateSignature)
	if err != nil {
		log.Debug("verifySignatures.SetSignature", "error", err.Error())
		return false
	}

	err = header.SetLeaderSignature(cnsDta.LeaderSignature)
	if err != nil {
		log.Debug("verifySignatures.SetLeaderSignature", "error", err.Error())
		return false
	}

	err = sr.HeaderSigVerifier().VerifyLeaderSignature(header)
	if err != nil {
		log.Debug("verifySignatures.VerifyLeaderSignature", "error", err.Error())
		return false
	}
	err = sr.HeaderSigVerifier().VerifySignature(header)
	if err != nil {
		log.Debug("verifySignatures.VerifySignature", "error", err.Error())
		return false
	}

	return true
}

// receivedInvalidSignersInfo method is called when a message with invalid signers has been received
func (sr *subroundEndRound) receivedInvalidSignersInfo(_ context.Context, cnsDta *consensus.Message) bool {
	messageSender := string(cnsDta.PubKey)

	if !sr.IsConsensusDataSet() {
		return false
	}

	// TODO[cleanup cns finality]: update this check
	isSenderAllowed := sr.IsNodeInConsensusGroup(messageSender)
	if !sr.EnableEpochsHandler().IsFlagEnabled(common.ConsensusPropagationChangesFlag) {
		isSenderAllowed = sr.IsNodeLeaderInCurrentRound(messageSender)
	}
	if !isSenderAllowed { // is NOT this node leader in current round?
		sr.PeerHonestyHandler().ChangeScore(
			messageSender,
			spos.GetConsensusTopicID(sr.ShardCoordinator()),
			spos.LeaderPeerHonestyDecreaseFactor,
		)

		return false
	}

	// TODO[cleanup cns finality]: update this check
	isSelfSender := sr.IsSelfLeaderInCurrentRound() || sr.IsMultiKeyLeaderInCurrentRound()
	if !sr.EnableEpochsHandler().IsFlagEnabled(common.ConsensusPropagationChangesFlag) {
		isSelfSender = messageSender == sr.SelfPubKey()
	}
	if isSelfSender {
		return false
	}

	if !sr.IsConsensusDataEqual(cnsDta.BlockHeaderHash) {
		return false
	}

	if !sr.CanProcessReceivedMessage(cnsDta, sr.RoundHandler().Index(), sr.Current()) {
		return false
	}

	if len(cnsDta.InvalidSigners) == 0 {
		return false
	}

	err := sr.verifyInvalidSigners(cnsDta.InvalidSigners)
	if err != nil {
		log.Trace("receivedInvalidSignersInfo.verifyInvalidSigners", "error", err.Error())
		return false
	}

	log.Debug("step 3: invalid signers info has been evaluated")

	sr.PeerHonestyHandler().ChangeScore(
		messageSender,
		spos.GetConsensusTopicID(sr.ShardCoordinator()),
		spos.LeaderPeerHonestyIncreaseFactor,
	)

	return true
}

func (sr *subroundEndRound) verifyInvalidSigners(invalidSigners []byte) error {
	messages, err := sr.MessageSigningHandler().Deserialize(invalidSigners)
	if err != nil {
		return err
	}

	for _, msg := range messages {
		err = sr.verifyInvalidSigner(msg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sr *subroundEndRound) verifyInvalidSigner(msg p2p.MessageP2P) error {
	err := sr.MessageSigningHandler().Verify(msg)
	if err != nil {
		return err
	}

	cnsMsg := &consensus.Message{}
	err = sr.Marshalizer().Unmarshal(cnsMsg, msg.Data())
	if err != nil {
		return err
	}

	err = sr.SigningHandler().VerifySingleSignature(cnsMsg.PubKey, cnsMsg.BlockHeaderHash, cnsMsg.SignatureShare)
	if err != nil {
		log.Debug("verifyInvalidSigner: confirmed that node provided invalid signature",
			"pubKey", cnsMsg.PubKey,
			"blockHeaderHash", cnsMsg.BlockHeaderHash,
			"error", err.Error(),
		)
		sr.applyBlacklistOnNode(msg.Peer())
	}

	return nil
}

func (sr *subroundEndRound) applyBlacklistOnNode(peer core.PeerID) {
	sr.PeerBlacklistHandler().BlacklistPeer(peer, common.InvalidSigningBlacklistDuration)
}

func (sr *subroundEndRound) receivedHeader(headerHandler data.HeaderHandler) {
	if sr.ConsensusGroup() == nil || sr.IsSelfLeaderInCurrentRound() || sr.IsMultiKeyLeaderInCurrentRound() {
		return
	}

	sr.AddReceivedHeader(headerHandler)

	sr.doEndRoundJobByParticipant(nil)
}

// doEndRoundJob method does the job of the subround EndRound
func (sr *subroundEndRound) doEndRoundJob(_ context.Context) bool {
	// TODO[cleanup cns finality]: remove L314-L324
	if !sr.IsSelfLeaderInCurrentRound() && !sr.IsMultiKeyLeaderInCurrentRound() && !sr.EnableEpochsHandler().IsFlagEnabled(common.ConsensusPropagationChangesFlag) {
		if sr.IsNodeInConsensusGroup(sr.SelfPubKey()) || sr.IsMultiKeyInConsensusGroup() {
			err := sr.prepareBroadcastBlockDataForValidator()
			if err != nil {
				log.Warn("validator in consensus group preparing for delayed broadcast",
					"error", err.Error())
			}
		}

		return sr.doEndRoundJobByParticipant(nil)
	}

	if !sr.IsNodeInConsensusGroup(sr.SelfPubKey()) && !sr.IsMultiKeyInConsensusGroup() {
		return false
	}

	return sr.doEndRoundJobByLeader()
}

func (sr *subroundEndRound) doEndRoundJobByLeader() bool {
	if !sr.sendFinalInfo() {
		return false
	}

	leader, err := sr.getLeader()
	if err != nil {
		return false
	}

	// broadcast header
	// TODO[Sorin next PR]: decide if we send this with the delayed broadcast
	err = sr.BroadcastMessenger().BroadcastHeader(sr.Header, []byte(leader))
	if err != nil {
		log.Warn("broadcastHeader.BroadcastHeader", "error", err.Error())
	}

	startTime := time.Now()
	err = sr.BlockProcessor().CommitBlock(sr.Header, sr.Body)
	elapsedTime := time.Since(startTime)
	if elapsedTime >= common.CommitMaxTime {
		log.Warn("doEndRoundJobByLeader.CommitBlock", "elapsed time", elapsedTime)
	} else {
		log.Debug("elapsed time to commit block",
			"time [s]", elapsedTime,
		)
	}
	if err != nil {
		log.Debug("doEndRoundJobByLeader.CommitBlock", "error", err)
		return false
	}

	sr.SetStatus(sr.Current(), spos.SsFinished)

	sr.displayStatistics()

	log.Debug("step 3: Body and Header have been committed and header has been broadcast")

	err = sr.broadcastBlockDataLeader()
	if err != nil {
		log.Debug("doEndRoundJobByLeader.broadcastBlockDataLeader", "error", err.Error())
	}

	msg := fmt.Sprintf("Added proposed block with nonce  %d  in blockchain", sr.Header.GetNonce())
	log.Debug(display.Headline(msg, sr.SyncTimer().FormattedCurrentTime(), "+"))

	sr.updateMetricsForLeader()

	return true
}

func (sr *subroundEndRound) sendFinalInfo() bool {
	if !sr.shouldSendFinalData() {
		return true
	}

	bitmap := sr.GenerateBitmap(SrSignature)
	err := sr.checkSignaturesValidity(bitmap)
	if err != nil {
		log.Debug("doEndRoundJobByLeader.checkSignaturesValidity", "error", err.Error())
		return false
	}

	if check.IfNil(sr.Header) {
		log.Error("doEndRoundJobByLeader.CheckNilHeader", "error", spos.ErrNilHeader)
		return false
	}

	// Aggregate sig and add it to the block
	bitmap, sig, err := sr.aggregateSigsAndHandleInvalidSigners(bitmap)
	if err != nil {
		log.Debug("doEndRoundJobByLeader.aggregateSigsAndHandleInvalidSigners", "error", err.Error())
		return false
	}

	err = sr.Header.SetPubKeysBitmap(bitmap)
	if err != nil {
		log.Debug("doEndRoundJobByLeader.SetPubKeysBitmap", "error", err.Error())
		return false
	}

	err = sr.Header.SetSignature(sig)
	if err != nil {
		log.Debug("doEndRoundJobByLeader.SetSignature", "error", err.Error())
		return false
	}

	// TODO[cleanup cns finality]: remove this code block
	if !sr.EnableEpochsHandler().IsFlagEnabled(common.ConsensusPropagationChangesFlag) {
		// Header is complete so the leader can sign it
		leaderSignature, err := sr.signBlockHeader()
		if err != nil {
			log.Error(err.Error())
			return false
		}

		err = sr.Header.SetLeaderSignature(leaderSignature)
		if err != nil {
			log.Debug("doEndRoundJobByLeader.SetLeaderSignature", "error", err.Error())
			return false
		}
	}

	ok := sr.ScheduledProcessor().IsProcessedOKWithTimeout()
	// placeholder for subroundEndRound.doEndRoundJobByLeader script
	if !ok {
		return false
	}

	roundHandler := sr.RoundHandler()
	if roundHandler.RemainingTime(roundHandler.TimeStamp(), roundHandler.TimeDuration()) < 0 {
		log.Debug("doEndRoundJob: time is out -> cancel broadcasting final info and header",
			"round time stamp", roundHandler.TimeStamp(),
			"current time", time.Now())
		return false
	}

	// broadcast header and final info section
	aggregatedSigToBroadcast := sr.Header.GetSignature()
	bitmapToBroadcast := sr.Header.GetPubKeysBitmap()
	leaderSigToBroadcast := sr.Header.GetLeaderSignature()
	// TODO[cleanup cns finality]: remove the above lines
	if sr.EnableEpochsHandler().IsFlagEnabled(common.ConsensusPropagationChangesFlag) {
		headerHash, err := core.CalculateHash(sr.Marshalizer(), sr.Hasher(), sr.Header)
		if err != nil {
			log.Debug("doEndRoundJobByLeader: calculate header hash", "error", err.Error())
			return false
		}
		aggregatedSigToBroadcast, bitmapToBroadcast = sr.getValidatedEquivalentProof(headerHash)
		leaderSigToBroadcast = nil
	}
	sr.createAndBroadcastHeaderFinalInfo(aggregatedSigToBroadcast, bitmapToBroadcast, leaderSigToBroadcast)

	if sr.EnableEpochsHandler().IsFlagEnabled(common.ConsensusPropagationChangesFlag) {
		sr.Blockchain().SetCurrentAggregatedSignatureAndBitmap(sig, bitmap)
	}

	return true
}

func (sr *subroundEndRound) shouldSendFinalData() bool {
	// TODO[cleanup cns finality]: remove this check
	if !sr.EnableEpochsHandler().IsFlagEnabled(common.ConsensusPropagationChangesFlag) {
		return true
	}

	headerHash, err := core.CalculateHash(sr.Marshalizer(), sr.Hasher(), sr.Header)
	if err != nil {
		log.Debug("shouldSendFinalData: calculate header hash", "error", err.Error())
		return false
	}

	// TODO: check if this is the best approach. Perhaps we don't want to relay only on the first received message
	if sr.hasEquivalentProof(headerHash) {
		log.Debug("shouldSendFinalData: equivalent message already sent")
		return false
	}

	return true
}

func (sr *subroundEndRound) aggregateSigsAndHandleInvalidSigners(bitmap []byte) ([]byte, []byte, error) {
	sig, err := sr.SigningHandler().AggregateSigs(bitmap, sr.Header.GetEpoch())
	if err != nil {
		log.Debug("doEndRoundJobByLeader.AggregateSigs", "error", err.Error())

		return sr.handleInvalidSignersOnAggSigFail()
	}

	err = sr.SigningHandler().SetAggregatedSig(sig)
	if err != nil {
		log.Debug("doEndRoundJobByLeader.SetAggregatedSig", "error", err.Error())
		return nil, nil, err
	}

	err = sr.SigningHandler().Verify(sr.GetData(), bitmap, sr.Header.GetEpoch())
	if err != nil {
		log.Debug("doEndRoundJobByLeader.Verify", "error", err.Error())

		return sr.handleInvalidSignersOnAggSigFail()
	}

	return bitmap, sig, nil
}

func (sr *subroundEndRound) verifyNodesOnAggSigFail() ([]string, error) {
	invalidPubKeys := make([]string, 0)
	pubKeys := sr.ConsensusGroup()

	if check.IfNil(sr.Header) {
		return nil, spos.ErrNilHeader
	}

	for i, pk := range pubKeys {
		isJobDone, err := sr.JobDone(pk, SrSignature)
		if err != nil || !isJobDone {
			continue
		}

		sigShare, err := sr.SigningHandler().SignatureShare(uint16(i))
		if err != nil {
			return nil, err
		}

		isSuccessfull := true
		err = sr.SigningHandler().VerifySignatureShare(uint16(i), sigShare, sr.GetData(), sr.Header.GetEpoch())
		if err != nil {
			isSuccessfull = false

			err = sr.SetJobDone(pk, SrSignature, false)
			if err != nil {
				return nil, err
			}

			// use increase factor since it was added optimistically, and it proved to be wrong
			decreaseFactor := -spos.ValidatorPeerHonestyIncreaseFactor + spos.ValidatorPeerHonestyDecreaseFactor
			sr.PeerHonestyHandler().ChangeScore(
				pk,
				spos.GetConsensusTopicID(sr.ShardCoordinator()),
				decreaseFactor,
			)

			invalidPubKeys = append(invalidPubKeys, pk)
		}

		log.Trace("verifyNodesOnAggSigVerificationFail: verifying signature share", "public key", pk, "is successfull", isSuccessfull)
	}

	return invalidPubKeys, nil
}

func (sr *subroundEndRound) getFullMessagesForInvalidSigners(invalidPubKeys []string) ([]byte, error) {
	p2pMessages := make([]p2p.MessageP2P, 0)

	for _, pk := range invalidPubKeys {
		p2pMsg, ok := sr.GetMessageWithSignature(pk)
		if !ok {
			log.Trace("message not found in state for invalid signer", "pubkey", pk)
			continue
		}

		p2pMessages = append(p2pMessages, p2pMsg)
	}

	invalidSigners, err := sr.MessageSigningHandler().Serialize(p2pMessages)
	if err != nil {
		return nil, err
	}

	return invalidSigners, nil
}

func (sr *subroundEndRound) handleInvalidSignersOnAggSigFail() ([]byte, []byte, error) {
	invalidPubKeys, err := sr.verifyNodesOnAggSigFail()
	if err != nil {
		log.Debug("doEndRoundJobByLeader.verifyNodesOnAggSigFail", "error", err.Error())
		return nil, nil, err
	}

	invalidSigners, err := sr.getFullMessagesForInvalidSigners(invalidPubKeys)
	if err != nil {
		log.Debug("doEndRoundJobByLeader.getFullMessagesForInvalidSigners", "error", err.Error())
		return nil, nil, err
	}

	if len(invalidSigners) > 0 {
		sr.createAndBroadcastInvalidSigners(invalidSigners)
	}

	bitmap, sig, err := sr.computeAggSigOnValidNodes()
	if err != nil {
		log.Debug("doEndRoundJobByLeader.computeAggSigOnValidNodes", "error", err.Error())
		return nil, nil, err
	}

	return bitmap, sig, nil
}

func (sr *subroundEndRound) computeAggSigOnValidNodes() ([]byte, []byte, error) {
	threshold := sr.Threshold(sr.Current())
	numValidSigShares := sr.ComputeSize(SrSignature)

	if check.IfNil(sr.Header) {
		return nil, nil, spos.ErrNilHeader
	}

	if numValidSigShares < threshold {
		return nil, nil, fmt.Errorf("%w: number of valid sig shares lower than threshold, numSigShares: %d, threshold: %d",
			spos.ErrInvalidNumSigShares, numValidSigShares, threshold)
	}

	bitmap := sr.GenerateBitmap(SrSignature)
	err := sr.checkSignaturesValidity(bitmap)
	if err != nil {
		return nil, nil, err
	}

	sig, err := sr.SigningHandler().AggregateSigs(bitmap, sr.Header.GetEpoch())
	if err != nil {
		return nil, nil, err
	}

	err = sr.SigningHandler().SetAggregatedSig(sig)
	if err != nil {
		return nil, nil, err
	}

	return bitmap, sig, nil
}

func (sr *subroundEndRound) createAndBroadcastHeaderFinalInfo(signature []byte, bitmap []byte, leaderSignature []byte) {
	leader, err := sr.getLeader()
	if err != nil {
		return
	}

	cnsMsg := consensus.NewConsensusMessage(
		sr.GetData(),
		nil,
		nil,
		nil,
		[]byte(leader),
		nil,
		int(MtBlockHeaderFinalInfo),
		sr.RoundHandler().Index(),
		sr.ChainID(),
		bitmap,
		signature,
		leaderSignature,
		sr.GetAssociatedPid([]byte(leader)),
		nil,
	)

	index, err := sr.ConsensusGroupIndex(leader)
	if err != nil {
		log.Debug("createAndBroadcastHeaderFinalInfo.ConsensusGroupIndex", "error", err.Error())
		return
	}

	if !sr.EnableEpochsHandler().IsFlagEnabled(common.ConsensusPropagationChangesFlag) {
		err = sr.BroadcastMessenger().BroadcastConsensusMessage(cnsMsg)
		if err != nil {
			log.Debug("createAndBroadcastHeaderFinalInfo.BroadcastConsensusMessage", "error", err.Error())
			return
		}

		log.Debug("step 3: block header final info has been sent",
			"PubKeysBitmap", bitmap,
			"AggregateSignature", signature,
			"LeaderSignature", leaderSignature)
		return
	}

	sr.BroadcastMessenger().PrepareBroadcastFinalConsensusMessage(cnsMsg, index)
	log.Debug("step 3: block header final info has been sent to delayed broadcaster",
		"PubKeysBitmap", bitmap,
		"AggregateSignature", signature,
		"LeaderSignature", leaderSignature,
		"Index", index)
}

func (sr *subroundEndRound) createAndBroadcastInvalidSigners(invalidSigners []byte) {
	cnsMsg := consensus.NewConsensusMessage(
		sr.GetData(),
		nil,
		nil,
		nil,
		[]byte(sr.SelfPubKey()),
		nil,
		int(MtInvalidSigners),
		sr.RoundHandler().Index(),
		sr.ChainID(),
		nil,
		nil,
		nil,
		sr.CurrentPid(),
		invalidSigners,
	)

	// TODO[Sorin next PR]: decide if we send this with the delayed broadcast
	err := sr.BroadcastMessenger().BroadcastConsensusMessage(cnsMsg)
	if err != nil {
		log.Debug("doEndRoundJob.BroadcastConsensusMessage", "error", err.Error())
		return
	}

	log.Debug("step 3: invalid signers info has been sent")
}

func (sr *subroundEndRound) doEndRoundJobByParticipant(cnsDta *consensus.Message) bool {
	sr.mutProcessingEndRound.Lock()
	defer sr.mutProcessingEndRound.Unlock()

	if sr.RoundCanceled {
		return false
	}
	if !sr.IsConsensusDataSet() {
		return false
	}
	if !sr.IsSubroundFinished(sr.Previous()) {
		return false
	}
	if sr.IsSubroundFinished(sr.Current()) {
		return false
	}

	haveHeader, header := sr.haveConsensusHeaderWithFullInfo(cnsDta)
	if !haveHeader {
		return false
	}

	defer func() {
		sr.SetProcessingBlock(false)
	}()

	sr.SetProcessingBlock(true)

	shouldNotCommitBlock := sr.ExtendedCalled || int64(header.GetRound()) < sr.RoundHandler().Index()
	if shouldNotCommitBlock {
		log.Debug("canceled round, extended has been called or round index has been changed",
			"round", sr.RoundHandler().Index(),
			"subround", sr.Name(),
			"header round", header.GetRound(),
			"extended called", sr.ExtendedCalled,
		)
		return false
	}

	if sr.isOutOfTime() {
		return false
	}

	ok := sr.ScheduledProcessor().IsProcessedOKWithTimeout()
	if !ok {
		return false
	}

	startTime := time.Now()
	err := sr.BlockProcessor().CommitBlock(header, sr.Body)
	elapsedTime := time.Since(startTime)
	if elapsedTime >= common.CommitMaxTime {
		log.Warn("doEndRoundJobByParticipant.CommitBlock", "elapsed time", elapsedTime)
	} else {
		log.Debug("elapsed time to commit block",
			"time [s]", elapsedTime,
		)
	}
	if err != nil {
		log.Debug("doEndRoundJobByParticipant.CommitBlock", "error", err.Error())
		return false
	}

	if sr.EnableEpochsHandler().IsFlagEnabled(common.ConsensusPropagationChangesFlag) {
		signature, bitmap := sr.getValidatedEquivalentProof(cnsDta.BlockHeaderHash)
		sr.Blockchain().SetCurrentAggregatedSignatureAndBitmap(signature, bitmap)
	}

	sr.SetStatus(sr.Current(), spos.SsFinished)

	if sr.IsNodeInConsensusGroup(sr.SelfPubKey()) || sr.IsMultiKeyInConsensusGroup() {
		err = sr.setHeaderForValidator(header)
		if err != nil {
			log.Warn("doEndRoundJobByParticipant", "error", err.Error())
		}
	}

	sr.displayStatistics()

	log.Debug("step 3: Body and Header have been committed")

	headerTypeMsg := "received"
	if cnsDta != nil {
		headerTypeMsg = "assembled"
	}

	msg := fmt.Sprintf("Added %s block with nonce  %d  in blockchain", headerTypeMsg, header.GetNonce())
	log.Debug(display.Headline(msg, sr.SyncTimer().FormattedCurrentTime(), "-"))
	return true
}

func (sr *subroundEndRound) haveConsensusHeaderWithFullInfo(cnsDta *consensus.Message) (bool, data.HeaderHandler) {
	if cnsDta == nil {
		return sr.isConsensusHeaderReceived()
	}

	if check.IfNil(sr.Header) {
		return false, nil
	}

	header := sr.Header.ShallowClone()
	if !sr.EnableEpochsHandler().IsFlagEnabled(common.ConsensusPropagationChangesFlag) {
		err := header.SetPubKeysBitmap(cnsDta.PubKeysBitmap)
		if err != nil {
			return false, nil
		}

		err = header.SetSignature(cnsDta.AggregateSignature)
		if err != nil {
			return false, nil
		}

		err = header.SetLeaderSignature(cnsDta.LeaderSignature)
		if err != nil {
			return false, nil
		}

		return true, header
	}

	header.SetPreviousAggregatedSignatureAndBitmap(cnsDta.AggregateSignature, cnsDta.PubKeysBitmap)

	return true, header
}

func (sr *subroundEndRound) isConsensusHeaderReceived() (bool, data.HeaderHandler) {
	if check.IfNil(sr.Header) {
		return false, nil
	}

	consensusHeaderHash, err := core.CalculateHash(sr.Marshalizer(), sr.Hasher(), sr.Header)
	if err != nil {
		log.Debug("isConsensusHeaderReceived: calculate consensus header hash", "error", err.Error())
		return false, nil
	}

	receivedHeaders := sr.GetReceivedHeaders()

	var receivedHeaderHash []byte
	for index := range receivedHeaders {
		receivedHeader := receivedHeaders[index].ShallowClone()
		err = receivedHeader.SetLeaderSignature(nil)
		if err != nil {
			log.Debug("isConsensusHeaderReceived - SetLeaderSignature", "error", err.Error())
			return false, nil
		}

		err = receivedHeader.SetPubKeysBitmap(nil)
		if err != nil {
			log.Debug("isConsensusHeaderReceived - SetPubKeysBitmap", "error", err.Error())
			return false, nil
		}

		err = receivedHeader.SetSignature(nil)
		if err != nil {
			log.Debug("isConsensusHeaderReceived - SetSignature", "error", err.Error())
			return false, nil
		}

		receivedHeaderHash, err = core.CalculateHash(sr.Marshalizer(), sr.Hasher(), receivedHeader)
		if err != nil {
			log.Debug("isConsensusHeaderReceived: calculate received header hash", "error", err.Error())
			return false, nil
		}

		if bytes.Equal(receivedHeaderHash, consensusHeaderHash) {
			return true, receivedHeaders[index]
		}
	}

	return false, nil
}

func (sr *subroundEndRound) signBlockHeader() ([]byte, error) {
	headerClone := sr.Header.ShallowClone()
	err := headerClone.SetLeaderSignature(nil)
	if err != nil {
		return nil, err
	}

	marshalizedHdr, err := sr.Marshalizer().Marshal(headerClone)
	if err != nil {
		return nil, err
	}

	leader, errGetLeader := sr.getLeader()
	if errGetLeader != nil {
		return nil, errGetLeader
	}
	if errGetLeader != nil {
		return nil, errGetLeader
	}

	return sr.SigningHandler().CreateSignatureForPublicKey(marshalizedHdr, []byte(leader))
}

func (sr *subroundEndRound) updateMetricsForLeader() {
	sr.appStatusHandler.Increment(common.MetricCountAcceptedBlocks)
	sr.appStatusHandler.SetStringValue(common.MetricConsensusRoundState,
		fmt.Sprintf("valid block produced in %f sec", time.Since(sr.RoundHandler().TimeStamp()).Seconds()))
}

func (sr *subroundEndRound) broadcastBlockDataLeader() error {
	miniBlocks, transactions, err := sr.BlockProcessor().MarshalizedDataToBroadcast(sr.Header, sr.Body)
	if err != nil {
		return err
	}

	leader, err := sr.getLeader()
	if err != nil {
		return err
	}

	// TODO[Sorin next PR]: decide if we send this with the delayed broadcast
	return sr.BroadcastMessenger().BroadcastBlockDataLeader(sr.Header, miniBlocks, transactions, []byte(leader))
}

func (sr *subroundEndRound) setHeaderForValidator(header data.HeaderHandler) error {
	idx, pk, miniBlocks, transactions, err := sr.getIndexPkAndDataToBroadcast()
	if err != nil {
		return err
	}

	go sr.BroadcastMessenger().PrepareBroadcastHeaderValidator(header, miniBlocks, transactions, idx, pk)

	return nil
}

func (sr *subroundEndRound) prepareBroadcastBlockDataForValidator() error {
	idx, pk, miniBlocks, transactions, err := sr.getIndexPkAndDataToBroadcast()
	if err != nil {
		return err
	}

	go sr.BroadcastMessenger().PrepareBroadcastBlockDataValidator(sr.Header, miniBlocks, transactions, idx, pk)

	return nil
}

// doEndRoundConsensusCheck method checks if the consensus is achieved
func (sr *subroundEndRound) doEndRoundConsensusCheck() bool {
	if sr.RoundCanceled {
		return false
	}

	return sr.IsSubroundFinished(sr.Current())
}

// computeSignersPublicKeys will extract from the provided consensus group slice only the strings that matched with the bitmap
func computeSignersPublicKeys(consensusGroup []string, bitmap []byte) []string {
	nbBitsBitmap := len(bitmap) * 8
	consensusGroupSize := len(consensusGroup)
	size := consensusGroupSize
	if consensusGroupSize > nbBitsBitmap {
		size = nbBitsBitmap
	}

	result := make([]string, 0, len(consensusGroup))

	for i := 0; i < size; i++ {
		indexRequired := (bitmap[i/8] & (1 << uint16(i%8))) > 0
		if !indexRequired {
			continue
		}

		pubKey := consensusGroup[i]
		result = append(result, pubKey)
	}

	return result
}

func (sr *subroundEndRound) checkSignaturesValidity(bitmap []byte) error {
	if !sr.hasProposerSignature(bitmap) {
		return spos.ErrMissingProposerSignature
	}

	consensusGroup := sr.ConsensusGroup()
	signers := computeSignersPublicKeys(consensusGroup, bitmap)
	for _, pubKey := range signers {
		isSigJobDone, err := sr.JobDone(pubKey, SrSignature)
		if err != nil {
			return err
		}

		if !isSigJobDone {
			return spos.ErrNilSignature
		}
	}

	return nil
}

func (sr *subroundEndRound) hasProposerSignature(bitmap []byte) bool {
	// TODO[cleanup cns finality]: remove this check
	if !sr.EnableEpochsHandler().IsFlagEnabled(common.ConsensusPropagationChangesFlag) {
		return true
	}

	proposerIndex := 0
	return bitmap[proposerIndex/8]&(1<<uint8(proposerIndex%8)) > 0
}

func (sr *subroundEndRound) isOutOfTime() bool {
	startTime := sr.RoundTimeStamp
	maxTime := sr.RoundHandler().TimeDuration() * time.Duration(sr.processingThresholdPercentage) / 100
	if sr.RoundHandler().RemainingTime(startTime, maxTime) < 0 {
		log.Debug("canceled round, time is out",
			"round", sr.SyncTimer().FormattedCurrentTime(), sr.RoundHandler().Index(),
			"subround", sr.Name())

		sr.RoundCanceled = true
		return true
	}

	return false
}

func (sr *subroundEndRound) getIndexPkAndDataToBroadcast() (int, []byte, map[uint32][]byte, map[string][][]byte, error) {
	minIdx := sr.getMinConsensusGroupIndexOfManagedKeys()

	idx, err := sr.SelfConsensusGroupIndex()
	if err == nil {
		if idx < minIdx {
			minIdx = idx
		}
	}

	if minIdx == sr.ConsensusGroupSize() {
		return -1, nil, nil, nil, err
	}

	miniBlocks, transactions, err := sr.BlockProcessor().MarshalizedDataToBroadcast(sr.Header, sr.Body)
	if err != nil {
		return -1, nil, nil, nil, err
	}

	consensusGroup := sr.ConsensusGroup()
	pk := []byte(consensusGroup[minIdx])

	return minIdx, pk, miniBlocks, transactions, nil
}

func (sr *subroundEndRound) getMinConsensusGroupIndexOfManagedKeys() int {
	minIdx := sr.ConsensusGroupSize()

	for idx, validator := range sr.ConsensusGroup() {
		if !sr.IsKeyManagedByCurrentNode([]byte(validator)) {
			continue
		}

		if idx < minIdx {
			minIdx = idx
		}
	}

	return minIdx
}

func (sr *subroundEndRound) getLeader() (string, error) {
	leader := sr.SelfPubKey()
	// TODO[cleanup cns finality]: only use sr.SelfPubKey
	if !sr.EnableEpochsHandler().IsFlagEnabled(common.ConsensusPropagationChangesFlag) {
		var errGetLeader error
		leader, errGetLeader = sr.GetLeader()
		if errGetLeader != nil {
			log.Debug("GetLeader", "error", errGetLeader)
			return "", errGetLeader
		}
	}

	return leader, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (sr *subroundEndRound) IsInterfaceNil() bool {
	return sr == nil
}
