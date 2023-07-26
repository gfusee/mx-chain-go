package smartContract

import (
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-go/process"
)

type sovereignSCProcessFactory struct {
	scrProcessorCreator SCRProcessorCreator
}

// NewSovereignSCProcessFactory creates a new smart contract process factory
func NewSovereignSCProcessFactory(creator SCRProcessorCreator) (*sovereignSCProcessFactory, error) {
	if check.IfNil(creator) {
		return nil, process.ErrNilSCProcessorCreator
	}
	return &sovereignSCProcessFactory{}, nil
}

// CreateSCProcessor creates a new smart contract processor
func (scpf *sovereignSCProcessFactory) CreateSCProcessor(args ArgsNewSmartContractProcessor) (SCRProcessorHandler, error) {
	sp, err := scpf.scrProcessorCreator.CreateSCRProcessor(args)
	if err != nil {
		return nil, err
	}

	scProc, ok := sp.(*scProcessor)
	if !ok {
		return nil, process.ErrWrongTypeAssertion
	}

	return NewSovereignSCRProcessor(scProc)
}

// IsInterfaceNil returns true if there is no value under the interface
func (scpf *sovereignSCProcessFactory) IsInterfaceNil() bool {
	return scpf == nil
}
