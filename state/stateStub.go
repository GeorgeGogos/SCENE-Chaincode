package state

import (
	"fmt"

	"github.com/hyperledger-labs/cckit/router"

	"github.com/GeorgeGogos/SCENE-Chaincode/payload"
)

type StateStub struct {
	context router.Context
}

func NewStateStub(c router.Context) *StateStub {
	return &StateStub{
		context: c,
	}
}

func (s StateStub) NewContract(payload payload.ContractPayload, owner string) error {

	contract_State := &ContractState{
		Licensee:			owner,
		LicenseSaleId:     	payload.LicenseSaleId,
		ProductId:  		payload.ProductId,
		LicenseStatus: 	"Pending",
		Orgs:           	payload.Orgs,
		Parties:			payload.Parties,
		Clauses:         	payload.Clause,
	}
	fmt.Printf("%s\n", contract_State.String())
	if err := s.context.State().Insert(contract_State); err != nil {
		retErr := fmt.Errorf("Error: Insert() returned error: %s", err.Error())
		return retErr
	}
	return nil
}
