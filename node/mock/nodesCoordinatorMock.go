package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/sharding"
)

type NodesCoordinatorMock struct {
	ComputeValidatorsGroupCalled        func(randomness []byte, round uint64, shardId uint32, epoch uint32) ([]sharding.Validator, error)
	GetValidatorsPublicKeysCalled       func(randomness []byte, round uint64, shardId uint32, epoch uint32) ([]string, error)
	GetValidatorsRewardsAddressesCalled func(randomness []byte, round uint64, shardId uint32, epoch uint32) ([]string, error)
}

func (ncm *NodesCoordinatorMock) GetAllValidatorsPublicKeys(_ uint32) (map[uint32][][]byte, error) {
	return nil, nil
}

func (ncm *NodesCoordinatorMock) ComputeConsensusGroup(
	randomness []byte,
	round uint64,
	shardId uint32,
	epoch uint32,
) (validatorsGroup []sharding.Validator, err error) {

	if ncm.ComputeValidatorsGroupCalled != nil {
		return ncm.ComputeValidatorsGroupCalled(randomness, round, shardId, epoch)
	}

	list := []sharding.Validator{
		NewValidatorMock(big.NewInt(0), 0, []byte("A"), []byte("AA")),
		NewValidatorMock(big.NewInt(0), 0, []byte("B"), []byte("BB")),
		NewValidatorMock(big.NewInt(0), 0, []byte("C"), []byte("CC")),
		NewValidatorMock(big.NewInt(0), 0, []byte("D"), []byte("DD")),
		NewValidatorMock(big.NewInt(0), 0, []byte("E"), []byte("EE")),
		NewValidatorMock(big.NewInt(0), 0, []byte("F"), []byte("FF")),
		NewValidatorMock(big.NewInt(0), 0, []byte("G"), []byte("GG")),
		NewValidatorMock(big.NewInt(0), 0, []byte("H"), []byte("HH")),
		NewValidatorMock(big.NewInt(0), 0, []byte("I"), []byte("II")),
	}

	return list, nil
}

func (ncm *NodesCoordinatorMock) GetConsensusValidatorsPublicKeys(
	randomness []byte,
	round uint64,
	shardId uint32,
	epoch uint32,
) ([]string, error) {
	if ncm.GetValidatorsPublicKeysCalled != nil {
		return ncm.GetValidatorsPublicKeysCalled(randomness, round, shardId, epoch)
	}

	validators, err := ncm.ComputeConsensusGroup(randomness, round, shardId, epoch)
	if err != nil {
		return nil, err
	}

	pubKeys := make([]string, 0)

	for _, v := range validators {
		pubKeys = append(pubKeys, string(v.PubKey()))
	}

	return pubKeys, nil
}

func (ncm *NodesCoordinatorMock) GetConsensusValidatorsRewardsAddresses(
	randomness []byte,
	round uint64,
	shardId uint32,
	epoch uint32,
) ([]string, error) {
	if ncm.GetValidatorsPublicKeysCalled != nil {
		return ncm.GetValidatorsRewardsAddressesCalled(randomness, round, shardId, epoch)
	}

	validators, err := ncm.ComputeConsensusGroup(randomness, round, shardId, epoch)
	if err != nil {
		return nil, err
	}

	addresses := make([]string, 0)
	for _, v := range validators {
		addresses = append(addresses, string(v.Address()))
	}

	return addresses, nil
}

func (ncm *NodesCoordinatorMock) ConsensusGroupSize(_ uint32) int {
	panic("implement me")
}

func (ncm *NodesCoordinatorMock) SetNodesPerShards(
	_ map[uint32][]sharding.Validator,
	_ map[uint32][]sharding.Validator,
	_ uint32,
) error {
	return nil
}

func (ncm *NodesCoordinatorMock) GetSelectedPublicKeys(_ []byte, _ uint32, _ uint32) (publicKeys []string, err error) {
	panic("implement me")
}

func (ncm *NodesCoordinatorMock) GetValidatorWithPublicKey(_ []byte, _ uint32) (sharding.Validator, uint32, error) {
	panic("implement me")
}

func (ncm *NodesCoordinatorMock) GetValidatorsIndexes(_ []string, _ uint32) ([]uint64, error) {
	panic("implement me")
}

func (ncm *NodesCoordinatorMock) GetOwnPublicKey() []byte {
	panic("implement me")
}

// IsInterfaceNil returns true if there is no value under the interface
func (ncm *NodesCoordinatorMock) IsInterfaceNil() bool {
	if ncm == nil {
		return true
	}
	return false
}
