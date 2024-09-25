package state

import (
	"encoding/json"
	"fmt"

	"github.com/GeorgeGogos/SCENE-Chaincode/payload"
)

const ContractStateEntity = `ContractState`

type ContractState struct {
	Licensee		    string         		`json:"licensee"`
	LicenseSaleId   	string         		`json:"license_sale_id"`
	ProductId   		string         		`json:"product_id"`
	LicenseStatus 		string         		`json:"license_status"`
	Orgs           		[]string       		`json:"orgs"`
	Parties           	[]string       		`json:"parties"`
	Clauses          	[]payload.Clause 	`json:"clauses,omitempty"`
}

func (s ContractState) Key() ([]string, error) {
	return []string{ContractStateEntity, s.LicenseSaleId}, nil
}

func (s ContractState) String() string {
	marshaledClauses, _ := json.Marshal(s.Clauses)
	return fmt.Sprintf("ContractState (Licensee=%s, LicenseSaleId=%s, ProductId=%s, LicenseStatus=%s, Orgs=%s, Parties=%s, Clauses=%s)", 
		s.Licensee, s.LicenseSaleId, s.ProductId, s.LicenseStatus, s.Orgs, s.Parties, string(marshaledClauses))
}
