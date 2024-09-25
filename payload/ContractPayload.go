package payload

import (
	"encoding/json"
	"fmt"
)

type ContractPayload struct {
	LicenseSaleId   string   	`json:"license_sale_id"`
	ProductId 		string   	`json:"product_id"`
	Orgs         	[]string 	`json:"orgs"`
	Parties			[]string	`json:"parties"`
	Clause        	[]Clause	`json:"licenses,omitempty"`
}
type ContractPayloadAllias ContractPayload

func (p ContractPayload) String() string {
	marshaledClauses, _ := json.Marshal(p.Clause)
	return fmt.Sprintf("ContractPayload (LicenseSaleId=%s, ProductId=%s, Orgs=%s, Parties=%s, Clauses=%s)",
		p.LicenseSaleId, p.ProductId, p.Orgs, p.Parties, string(marshaledClauses))
}

func (p ContractPayload) Validate() error {
	if p.LicenseSaleId == "" {
		return fmt.Errorf("Error validating Contract payload: ContractID cannot be an empty string.")
	}
	if p.ProductId == "" {
		return fmt.Errorf("Error validating Contract payload: ProductId cannot be an empty string.")
	}
	if len(p.Orgs) != 2 {
		return fmt.Errorf("Error validating Contract payload: contracted Orgs must be two (2).")
	}
	for i := 0; i < len(p.Orgs); i++ {
		if p.Orgs[i] == "" {
			return fmt.Errorf("Error validating Contract payload: contracted Orgs cannot be an empty string.")
		}
	}
	if len(p.Parties) != 2 {
		return fmt.Errorf("Error validating Contract payload: contracted Parties must be two (2).")
	}
	for i := 0; i < len(p.Parties); i++ {
		if p.Parties[i] == "" {
			return fmt.Errorf("Error validating Contract payload: contracted Parties cannot be an empty string.")
		}
	}
	return nil
}
