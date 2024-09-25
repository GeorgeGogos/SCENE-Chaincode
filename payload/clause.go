package payload

import (
	"fmt"
)

type Clause struct {
	ClauseId   string `json:"clause_id"`
	LicenseId  string `json:"license_id"`
	Hash 	   string `json:"hash,omitempty"`
}

func (i Clause) String() string {
	return fmt.Sprintf("Clauses (ClauseId=%s, LicenseId=%s, Hash=%s)",
		i.ClauseId, i.LicenseId, i.Hash)
}

func (i Clause) Validate(p ContractPayload) error {
	if i.ClauseId == "" {
		return fmt.Errorf("Error validating Clauses: ClauseId cannot be an empty string.")
	}
	if i.LicenseId == "" || i.LicenseId != p.LicenseSaleId  {
		return fmt.Errorf("Error validating Clauses: LicenseId does not match the License Sales IDs.")
	}
	return nil
}

func (i Clause) ValidateHash(p ContractPayload) error {
	if i.Hash == "" {
		return fmt.Errorf("Error validating Clauses: Hash cannot be an empty string.")
	}
	return nil
}
