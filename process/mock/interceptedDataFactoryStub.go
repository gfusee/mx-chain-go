package mock

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-go/process"
)

// InterceptedDataFactoryStub -
type InterceptedDataFactoryStub struct {
	CreateCalled func(buff []byte) (process.InterceptedData, error)
}

// Create -
func (idfs *InterceptedDataFactoryStub) Create(buff []byte, _ core.PeerID) (process.InterceptedData, error) {
	return idfs.CreateCalled(buff)
}

// IsInterfaceNil -
func (idfs *InterceptedDataFactoryStub) IsInterfaceNil() bool {
	return idfs == nil
}
