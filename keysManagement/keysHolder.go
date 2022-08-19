package keysManagement

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	crypto "github.com/ElrondNetwork/elrond-go-crypto"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/common"
	"github.com/ElrondNetwork/elrond-go/config"
)

const minRoundsWithoutReceivedMessages = 1

var log = logger.GetOrCreate("keysManagement")

type virtualPeersHolder struct {
	mut                              sync.RWMutex
	data                             map[string]*peerInfo
	pids                             map[core.PeerID]struct{}
	keyGenerator                     crypto.KeyGenerator
	p2pIdentityGenerator             P2PIdentityGenerator
	isMainMachine                    bool
	maxRoundsWithoutReceivedMessages int
	namedIdentitiesMap               map[string]namedIdentity
}

// ArgsVirtualPeersHolder represents the argument for the virtual peers holder
type ArgsVirtualPeersHolder struct {
	KeyGenerator                     crypto.KeyGenerator
	P2PIdentityGenerator             P2PIdentityGenerator
	IsMainMachine                    bool
	MaxRoundsWithoutReceivedMessages int
	NamedIdentities                  []config.NamedIdentity
}

// NewVirtualPeersHolder creates a new instance of a virtual peers holder
func NewVirtualPeersHolder(args ArgsVirtualPeersHolder) (*virtualPeersHolder, error) {
	err := checkVirtualPeersHolderArgs(args)
	if err != nil {
		return nil, err
	}

	holder := &virtualPeersHolder{
		data:                             make(map[string]*peerInfo),
		pids:                             make(map[core.PeerID]struct{}),
		keyGenerator:                     args.KeyGenerator,
		p2pIdentityGenerator:             args.P2PIdentityGenerator,
		isMainMachine:                    args.IsMainMachine,
		maxRoundsWithoutReceivedMessages: args.MaxRoundsWithoutReceivedMessages,
		namedIdentitiesMap:               createNamedIdentitiesMap(args.NamedIdentities),
	}

	return holder, nil
}

func checkVirtualPeersHolderArgs(args ArgsVirtualPeersHolder) error {
	if check.IfNil(args.KeyGenerator) {
		return errNilKeyGenerator
	}
	if check.IfNil(args.P2PIdentityGenerator) {
		return errNilP2PIdentityGenerator
	}
	if args.MaxRoundsWithoutReceivedMessages < minRoundsWithoutReceivedMessages {
		return fmt.Errorf("%w for MaxRoundsWithoutReceivedMessages, minimum %d, got %d",
			errInvalidValue, minRoundsWithoutReceivedMessages, args.MaxRoundsWithoutReceivedMessages)
	}
	if len(args.NamedIdentities) == 0 {
		return errMissingNamedIdentity
	}

	return nil
}

func createNamedIdentitiesMap(namedIdentities []config.NamedIdentity) map[string]namedIdentity {
	namedIdentitiesMap := make(map[string]namedIdentity)

	for _, identity := range namedIdentities {
		for idx, blsKey := range identity.BLSKeys {
			bls, err := hex.DecodeString(blsKey)
			if err != nil {
				continue
			}

			blsStr := string(bls)
			namedIdentitiesMap[blsStr] = namedIdentity{
				name:     fmt.Sprintf("%s-%d", identity.NodeName, idx),
				identity: identity.Identity,
			}

			log.Trace(fmt.Sprintf("Found named identity: %s, %s, %s",
				namedIdentitiesMap[blsStr].name,
				namedIdentitiesMap[blsStr].identity,
				core.GetTrimmedPk(blsKey)),
			)
		}
	}

	return namedIdentitiesMap
}

// AddVirtualPeer will try to add a new virtual peer providing the private key bytes.
// It errors if the generated public key is already contained by the struct
// It will auto-generate some fields like the machineID and pid
func (holder *virtualPeersHolder) AddVirtualPeer(privateKeyBytes []byte) error {
	privateKey, err := holder.keyGenerator.PrivateKeyFromByteArray(privateKeyBytes)
	if err != nil {
		return fmt.Errorf("%w for provided bytes %s", err, hex.EncodeToString(privateKeyBytes))
	}

	publicKey := privateKey.GeneratePublic()
	publicKeyBytes, err := publicKey.ToByteArray()
	if err != nil {
		return fmt.Errorf("%w for provided bytes %s", err, hex.EncodeToString(privateKeyBytes))
	}

	p2pPrivateKeyBytes, pid, err := holder.p2pIdentityGenerator.CreateRandomP2PIdentity()
	if err != nil {
		return err
	}

	holder.mut.Lock()
	defer holder.mut.Unlock()

	pInfo := &peerInfo{
		pid:                pid,
		p2pPrivateKeyBytes: p2pPrivateKeyBytes,
		privateKey:         privateKey,
		machineID:          generateRandomMachineID(),
		namedIdentity:      holder.namedIdentitiesMap[string(publicKeyBytes)],
	}

	_, found := holder.data[string(publicKeyBytes)]
	if found {
		return fmt.Errorf("%w for provided bytes %s and generated public key %s",
			errDuplicatedKey, hex.EncodeToString(privateKeyBytes), hex.EncodeToString(publicKeyBytes))
	}

	holder.data[string(publicKeyBytes)] = pInfo
	holder.pids[pid] = struct{}{}

	log.Debug("added new key definition",
		"hex public key", hex.EncodeToString(publicKeyBytes),
		"pid", pid.Pretty(),
		"machine ID", pInfo.machineID)

	return nil
}

func (holder *virtualPeersHolder) getPeerInfo(pkBytes []byte) *peerInfo {
	holder.mut.RLock()
	defer holder.mut.RUnlock()

	return holder.data[string(pkBytes)]
}

func generateRandomMachineID() string {
	buff := make([]byte, common.MaxMachineIDLen/2)
	_, _ = rand.Read(buff)

	return hex.EncodeToString(buff)
}

// GetPrivateKey returns the associated private key with the provided public key bytes. Errors if the key is not found
func (holder *virtualPeersHolder) GetPrivateKey(pkBytes []byte) (crypto.PrivateKey, error) {
	pInfo := holder.getPeerInfo(pkBytes)
	if pInfo == nil {
		return nil, fmt.Errorf("%w in GetPrivateKey for public key %s",
			errMissingPublicKeyDefinition, hex.EncodeToString(pkBytes))
	}

	return pInfo.privateKey, nil
}

// GetP2PIdentity returns the associated p2p identity with the provided public key bytes: the private key and the peer ID
func (holder *virtualPeersHolder) GetP2PIdentity(pkBytes []byte) ([]byte, core.PeerID, error) {
	pInfo := holder.getPeerInfo(pkBytes)
	if pInfo == nil {
		return nil, "", fmt.Errorf("%w in GetP2PIdentity for public key %s",
			errMissingPublicKeyDefinition, hex.EncodeToString(pkBytes))
	}

	return pInfo.p2pPrivateKeyBytes, pInfo.pid, nil
}

// GetMachineID returns the associated machine ID with the provided public key bytes
func (holder *virtualPeersHolder) GetMachineID(pkBytes []byte) (string, error) {
	pInfo := holder.getPeerInfo(pkBytes)
	if pInfo == nil {
		return "", fmt.Errorf("%w in GetMachineID for public key %s",
			errMissingPublicKeyDefinition, hex.EncodeToString(pkBytes))
	}

	return pInfo.machineID, nil
}

// IncrementRoundsWithoutReceivedMessages increments the number of rounds without received messages on a provided public key
func (holder *virtualPeersHolder) IncrementRoundsWithoutReceivedMessages(pkBytes []byte) error {
	if holder.isMainMachine {
		return nil
	}

	pInfo := holder.getPeerInfo(pkBytes)
	if pInfo == nil {
		return fmt.Errorf("%w in IncrementRoundsWithoutReceivedMessages for public key %s",
			errMissingPublicKeyDefinition, hex.EncodeToString(pkBytes))
	}

	pInfo.incrementRoundsWithoutReceivedMessages()

	return nil
}

// ResetRoundsWithoutReceivedMessages resets the number of rounds without received messages on a provided public key
func (holder *virtualPeersHolder) ResetRoundsWithoutReceivedMessages(pkBytes []byte) error {
	if holder.isMainMachine {
		return nil
	}

	pInfo := holder.getPeerInfo(pkBytes)
	if pInfo == nil {
		return fmt.Errorf("%w in ResetRoundsWithoutReceivedMessages for public key %s",
			errMissingPublicKeyDefinition, hex.EncodeToString(pkBytes))
	}

	pInfo.resetRoundsWithoutReceivedMessages()

	return nil
}

// GetManagedKeysByCurrentNode returns all keys that will be managed by this node
func (holder *virtualPeersHolder) GetManagedKeysByCurrentNode() map[string]crypto.PrivateKey {
	holder.mut.RLock()
	defer holder.mut.RUnlock()

	allManagedKeys := make(map[string]crypto.PrivateKey)
	for pk, pInfo := range holder.data {
		isSlaveAndMainFailed := !holder.isMainMachine && !pInfo.isNodeActiveOnMainMachine(holder.maxRoundsWithoutReceivedMessages)
		shouldAddToMap := holder.isMainMachine || isSlaveAndMainFailed
		if !shouldAddToMap {
			continue
		}

		allManagedKeys[pk] = pInfo.privateKey
	}

	return allManagedKeys
}

// IsKeyManagedByCurrentNode returns true if the key is managed by the current node
func (holder *virtualPeersHolder) IsKeyManagedByCurrentNode(pkBytes []byte) bool {
	pInfo := holder.getPeerInfo(pkBytes)
	if pInfo == nil {
		return false
	}

	if holder.isMainMachine {
		return true
	}

	return !pInfo.isNodeActiveOnMainMachine(holder.maxRoundsWithoutReceivedMessages)
}

// IsKeyRegistered returns true if the key is registered (not necessarily managed by the current node)
func (holder *virtualPeersHolder) IsKeyRegistered(pkBytes []byte) bool {
	pInfo := holder.getPeerInfo(pkBytes)
	return pInfo != nil
}

// IsPidManagedByCurrentNode returns true if the peer id is managed by the current node
func (holder *virtualPeersHolder) IsPidManagedByCurrentNode(pid core.PeerID) bool {
	holder.mut.RLock()
	defer holder.mut.RUnlock()

	_, found := holder.pids[pid]

	return found
}

// IsInterfaceNil returns true if there is no value under the interface
func (holder *virtualPeersHolder) IsInterfaceNil() bool {
	return holder == nil
}
