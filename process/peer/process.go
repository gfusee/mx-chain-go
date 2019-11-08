package peer

import (
	"errors"
	"sync"

	"github.com/ElrondNetwork/elrond-go/data"
	"github.com/ElrondNetwork/elrond-go/data/block"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/dataRetriever"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/ElrondNetwork/elrond-go/storage"
)

// DataPool indicates the main funcitonality needed in order to fetch the required blocks from the pool
type DataPool interface {
	MetaBlocks() storage.Cacher
	IsInterfaceNil() bool
}

// ArgValidatorStatisticsProcessor holds all dependencies for the validatorStatistics
type ArgValidatorStatisticsProcessor struct {
	InitialNodes     []*sharding.InitialNode
	Marshalizer      marshal.Marshalizer
	NodesCoordinator sharding.NodesCoordinator
	ShardCoordinator sharding.Coordinator
	DataPool         DataPool
	StorageService   dataRetriever.StorageService
	AdrConv          state.AddressConverter
	PeerAdapter      state.AccountsAdapter
}

type validatorStatistics struct {
	marshalizer      marshal.Marshalizer
	dataPool         DataPool
	storageService   dataRetriever.StorageService
	nodesCoordinator sharding.NodesCoordinator
	shardCoordinator sharding.Coordinator
	adrConv          state.AddressConverter
	peerAdapter      state.AccountsAdapter
	prevShardInfo    map[uint32]block.ShardData
	mutPrevShardInfo sync.RWMutex
}

// NewValidatorStatisticsProcessor instantiates a new validatorStatistics structure responsible of keeping account of
//  each validator actions in the consensus process
func NewValidatorStatisticsProcessor(arguments ArgValidatorStatisticsProcessor) (*validatorStatistics, error) {
	if arguments.PeerAdapter == nil || arguments.PeerAdapter.IsInterfaceNil() {
		return nil, process.ErrNilPeerAccountsAdapter
	}
	if arguments.AdrConv == nil || arguments.AdrConv.IsInterfaceNil() {
		return nil, process.ErrNilAddressConverter
	}
	if arguments.DataPool == nil || arguments.DataPool.IsInterfaceNil() {
		return nil, process.ErrNilDataPoolHolder
	}
	if arguments.StorageService == nil || arguments.StorageService.IsInterfaceNil() {
		return nil, process.ErrNilStorage
	}
	if arguments.NodesCoordinator == nil || arguments.NodesCoordinator.IsInterfaceNil() {
		return nil, process.ErrNilNodesCoordinator
	}
	if arguments.ShardCoordinator == nil || arguments.ShardCoordinator.IsInterfaceNil() {
		return nil, process.ErrNilShardCoordinator
	}
	if arguments.Marshalizer == nil || arguments.Marshalizer.IsInterfaceNil() {
		return nil, process.ErrNilMarshalizer
	}

	vs := &validatorStatistics{
		peerAdapter:      arguments.PeerAdapter,
		adrConv:          arguments.AdrConv,
		nodesCoordinator: arguments.NodesCoordinator,
		shardCoordinator: arguments.ShardCoordinator,
		dataPool:         arguments.DataPool,
		storageService:   arguments.StorageService,
		marshalizer:      arguments.Marshalizer,
		prevShardInfo:    make(map[uint32]block.ShardData, 0),
	}

	err := vs.SaveInitialState(arguments.InitialNodes)
	if err != nil {
		return nil, err
	}

	return vs, nil
}

// SaveInitialState takes an initial peer list, validates it and sets up the initial state for each of the peers
func (p *validatorStatistics) SaveInitialState(in []*sharding.InitialNode) error {
	for _, node := range in {
		err := p.initializeNode(node)
		if err != nil {
			return err
		}
	}

	_, err := p.peerAdapter.Commit()
	if err != nil {
		return err
	}

	return nil
}

// IsNodeValid calculates if a node that's present in the initial validator list
//  contains all the required information in order to be able to participate in consensus
func (p *validatorStatistics) IsNodeValid(node *sharding.InitialNode) bool {
	if len(node.PubKey) == 0 {
		return false
	}
	if len(node.Address) == 0 {
		return false
	}

	return true
}

// UpdatePeerState takes the in a header, updates the peer state for all of the
//  consensus members and returns the new root hash
func (p *validatorStatistics) UpdatePeerState(header data.HeaderHandler) ([]byte, error) {
	if header.GetNonce() == 0 {
		return p.peerAdapter.RootHash()
	}

	consensusGroup, err := p.nodesCoordinator.ComputeValidatorsGroup(header.GetPrevRandSeed(), header.GetRound(), header.GetShardID())
	if err != nil {
		return nil, err
	}

	err = p.updateValidatorInfo(consensusGroup, header.GetShardID())
	if err != nil {
		return nil, err
	}

	if header.GetNonce() == 1 {
		return p.peerAdapter.RootHash()
	}

	previousHeader, err := process.GetMetaHeader(header.GetPrevHash(), p.dataPool.MetaBlocks(), p.marshalizer, p.storageService)
	if err != nil {
		return nil, err
	}

	err = p.checkForMissedBlocks(
		header.GetRound(),
		previousHeader.GetRound(),
		previousHeader.GetPrevRandSeed(),
		previousHeader.GetShardID(),
	)
	if err != nil {
		return nil, err
	}

	err = p.updateShardDataPeerState(header, previousHeader)
	if err != nil {
		return nil, err
	}

	return p.peerAdapter.RootHash()
}

// Commit commits the validator statistics trie and returns the root hash
func (p *validatorStatistics) Commit() ([]byte, error) {
	return p.peerAdapter.Commit()
}

// RootHash returns the root hash of the validator statistics trie
func (p *validatorStatistics) RootHash() ([]byte, error) {
	return p.peerAdapter.RootHash()
}

func (p *validatorStatistics) checkForMissedBlocks(currentHeaderRound, previousHeaderRound uint64,
	prevRandSeed []byte, shardId uint32) error {
	if currentHeaderRound-previousHeaderRound <= 1 {
		return nil
	}

	for i := previousHeaderRound + 1; i < currentHeaderRound; i++ {
		consensusGroup, err := p.nodesCoordinator.ComputeValidatorsGroup(prevRandSeed, i, shardId)
		if err != nil {
			return err
		}

		leaderPeerAcc, err := p.getPeerAccount(consensusGroup[0].Address())
		if err != nil {
			return err
		}

		err = leaderPeerAcc.DecreaseLeaderSuccessRateWithJournal()
		if err != nil {
			return err
		}
	}

	return nil
}

// RevertPeerState takes the current and previous headers and undos the peer state
//  for all of the consensus members
func (p *validatorStatistics) RevertPeerState(header data.HeaderHandler) error {
	_ = p.peerAdapter.RecreateTrie(header.GetValidatorStatsRootHash())
	return nil
}

// RevertPeerStateToSnapshot reverts the applied changes to the peerAdapter
func (p *validatorStatistics) RevertPeerStateToSnapshot(snapshot int) error {
	return p.peerAdapter.RevertToSnapshot(snapshot)
}

func (p *validatorStatistics) updateShardDataPeerState(header, previousHeader data.HeaderHandler) error {
	metaHeader, ok := header.(*block.MetaBlock)
	if !ok {
		return process.ErrInvalidMetaHeader
	}
	prevMetaHeader, ok := header.(*block.MetaBlock)
	if !ok {
		return process.ErrInvalidMetaHeader
	}

	err := p.loadPreviousShardHeaders(metaHeader, prevMetaHeader)
	if err !=  nil {
		return err
	}

	for _, h := range metaHeader.ShardInfo {

		shardConsensus, err := p.nodesCoordinator.ComputeValidatorsGroup(h.PrevRandSeed, h.Round, h.ShardId)
		if err != nil {
			return err
		}

		err = p.updateValidatorInfo(shardConsensus, h.ShardId)
		if err != nil {
			return err
		}

		if h.Nonce == 1 {
			continue
		}

		p.mutPrevShardInfo.RLock()
		prevShardData, ok := p.prevShardInfo[h.ShardId]
		p.mutPrevShardInfo.RUnlock()
		if !ok {
			return process.ErrMissingPrevShardData
		}

		err = p.checkForMissedBlocks(
			h.Round,
			prevShardData.Round,
			prevShardData.PrevRandSeed,
			h.ShardId,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *validatorStatistics) initializeNode(node *sharding.InitialNode) error {
	if !p.IsNodeValid(node) {
		return process.ErrInvalidInitialNodesState
	}

	peerAccount, err := p.generatePeerAccount(node)
	if err != nil {
		return err
	}

	err = p.savePeerAccountData(peerAccount, node)
	if err != nil {
		return err
	}

	return nil
}

func (p *validatorStatistics) generatePeerAccount(node *sharding.InitialNode) (*state.PeerAccount, error) {
	address, err := p.adrConv.CreateAddressFromHex(node.Address)
	if err != nil {
		return nil, err
	}

	acc, err := p.peerAdapter.GetAccountWithJournal(address)
	if err != nil {
		return nil, err
	}

	peerAccount, ok := acc.(*state.PeerAccount)
	if !ok {
		return nil, process.ErrInvalidPeerAccount
	}

	return peerAccount, nil
}

func (p *validatorStatistics) savePeerAccountData(peerAccount *state.PeerAccount, data *sharding.InitialNode) error {
	err := peerAccount.SetAddressWithJournal([]byte(data.Address))
	if err != nil {
		return err
	}

	err = peerAccount.SetSchnorrPublicKeyWithJournal([]byte(data.Address))
	if err != nil {
		return err
	}

	err = peerAccount.SetBLSPublicKeyWithJournal([]byte(data.PubKey))
	if err != nil {
		return err
	}

	return nil
}

func (p *validatorStatistics) updateValidatorInfo(validatorList []sharding.Validator, shardId uint32) error {
	lenValidators := len(validatorList)
	for i := 0; i < lenValidators; i++ {
		peerAcc, err := p.getPeerAccount(validatorList[i].Address())
		if err != nil {
			return err
		}

		isLeader := i == 0
		if isLeader {
			err = peerAcc.IncreaseLeaderSuccessRateWithJournal()
		} else {
			err = peerAcc.IncreaseValidatorSuccessRateWithJournal()
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *validatorStatistics) getPeerAccount(address []byte) (state.PeerAccountHandler, error) {
	addressContainer, err := p.adrConv.CreateAddressFromPublicKeyBytes(address)
	if err != nil {
		return nil, err
	}

	account, err := p.peerAdapter.GetExistingAccount(addressContainer)
	if err != nil {
		return nil, err
	}

	peerAccount, ok := account.(state.PeerAccountHandler)
	if !ok {
		return nil, process.ErrInvalidPeerAccount
	}

	return peerAccount, nil
}

// loadPreviousShardHeaders loads the previous shard headers for a given metablock. For the metachain it's easy
//  since it has all the shard headers in it's storage, but for the shard it's a bit trickier and we need
//  to iterate through past metachain headers until we find all the ShardData's we are interested in
func (p *validatorStatistics) loadPreviousShardHeaders(currentHeader, previousHeader *block.MetaBlock) error {

	if p.shardCoordinator.SelfId() > p.shardCoordinator.NumberOfShards() {
		return p.loadPreviousShardHeadersMeta(currentHeader)
	}

	p.mutPrevShardInfo.Lock()
	defer p.mutPrevShardInfo.Unlock()

	p.prevShardInfo = make(map[uint32]block.ShardData)
	missingShardIds := make([]uint32, 0)

	for _, currentShardData := range currentHeader.ShardInfo {
		prevShardData := p.getMatchingShardData(currentShardData.ShardId, previousHeader.ShardInfo)
		if prevShardData != nil {
			p.prevShardInfo[currentShardData.ShardId] = *prevShardData
		} else {
			missingShardIds = append(missingShardIds, currentShardData.ShardId)
		}
	}

	searchHeader := &block.MetaBlock{}
	*searchHeader = *previousHeader
	for len(missingShardIds) > 0 {
		recursiveHeader, err := process.GetMetaHeader(searchHeader.GetPrevHash(), p.dataPool.MetaBlocks(), p.marshalizer, p.storageService)
		if err != nil {
			return err
		}
		for i, shardId := range missingShardIds {
			prevShardData := p.getMatchingShardData(shardId, recursiveHeader.ShardInfo)
			if prevShardData != nil {
				p.prevShardInfo[shardId] = *prevShardData
				missingShardIds = append(missingShardIds[:i], missingShardIds[i+1:]...)
			}
		}
		*searchHeader = *recursiveHeader
	}
	return nil
}

func (p *validatorStatistics) loadPreviousShardHeadersMeta(header *block.MetaBlock) error {
	p.mutPrevShardInfo.Lock()
	defer p.mutPrevShardInfo.Unlock()

	metaDataPool, ok := p.dataPool.(dataRetriever.MetaPoolsHolder)
	if !ok {
		return errors.New("woooot")
	}

	for _, shardData := range header.ShardInfo {
		if shardData.Nonce == 1 {
			continue
		}

		previousHeader, err := process.GetShardHeader(shardData.PrevHash, metaDataPool.ShardHeaders(), p.marshalizer,
			p.storageService)
		if err != nil {
			return err
		}

		p.prevShardInfo[shardData.ShardId] = block.ShardData{
			ShardId: previousHeader.ShardId,
			Nonce: previousHeader.Nonce,
			Round: previousHeader.Round,
			PrevRandSeed: previousHeader.PrevRandSeed,
			PrevHash: previousHeader.PrevHash,
		}
	}
	return nil
}

func (p *validatorStatistics) getMatchingShardData(shardId uint32, shardInfo []block.ShardData) *block.ShardData {
	for _, prevShardData := range shardInfo {
		if shardId == prevShardData.ShardId {
			return &prevShardData
		}
	}

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (p *validatorStatistics) IsInterfaceNil() bool {
	if p == nil {
		return true
	}
	return false
}
