package scene

import (
	"github.com/GeorgeGogos/SCENE-Chaincode/chaincode"
	"github.com/GeorgeGogos/SCENE-Chaincode/payload"

	logging "github.com/CERTH-ITI-DLT-Lab/hlf-cc-logging"

	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger-labs/cckit/router/param"
)

func NewCC() *router.Chaincode {
	logging.InitCCLogger()
	r := router.New(`scene_chaincode`).Use(logging.SetContextMiddlewareFunc())

	r.Init(func(context router.Context) (i interface{}, e error) {
		// No implementation required with this example
		// It could be where data migration is performed, if necessary
		return nil, nil
	})

	r.
		// Read methods
		Query(`GetContractByID`, chaincode.GetContractByID, param.String("license_sale_id")).
		Query(`GetContracts`, chaincode.GetContracts).
		Query(`GetContractIDs`, chaincode.GetContractIDs).

		// Transaction methods
		Invoke(`ProposeContract`, chaincode.ProposeContract, param.Struct("contractPayload", &payload.ContractPayload{})).
		Invoke(`AcceptContract`, chaincode.AcceptContract, param.String("license_sale_id")).
		Invoke(`RejectContract`, chaincode.RejectContract, param.String("license_sale_id")).
		Invoke(`DissolveContract`, chaincode.DissolveContract, param.String("license_sale_id"))
		/*Invoke(`UpdateContractClause`, chaincode.UpdateContractClause, param.String("license_sale_id"), param.Struct("clausePayload", &payload.Clause{})).
		Invoke(`DeleteContractClause`, chaincode.DeleteContractClause, param.String("license_sale_id"), param.String("clause_ID")).
		Invoke(`AddContractClause`, chaincode.AddContractClause, param.String("license_sale_id"), param.Struct("clausePayload", &payload.Clause{}))*/

	return router.NewChaincode(r)

}
