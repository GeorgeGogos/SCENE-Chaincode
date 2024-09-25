package output

import (
	"encoding/json"
	"fmt"

	"github.com/GeorgeGogos/SCENE-Chaincode/payload"
)

type OutputContract struct {
	LicenseSaleId   	string         		`json:"license_sale_id"`
	ProductId   		string         		`json:"product_id"`
	LicenseStatus 		string         		`json:"license_status"`
	Orgs           		[]string       		`json:"orgs"`
	Parties           	[]string       		`json:"parties"`
	Clauses          	[]payload.Clause 	`json:"clauses,omitempty"`
}

func (s OutputContract) String() string {
	marshaledClauses, _ := json.Marshal(s.Clauses)
	return fmt.Sprintf("OutputContract (LicenseSaleId=%s, ProductId=%s, LicenseStatus=%s, Orgs=%s, Clauses=%s)",
		s.LicenseSaleId, s.ProductId, s.LicenseStatus, s.Orgs, string(marshaledClauses))
}
