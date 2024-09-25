package chaincode

import (
	"encoding/json"

	"github.com/GeorgeGogos/SCENE-Chaincode/payload"

	"github.com/GeorgeGogos/SCENE-Chaincode/output"
	"github.com/GeorgeGogos/SCENE-Chaincode/state"

	"fmt"

	logging "github.com/CERTH-ITI-DLT-Lab/hlf-cc-logging"

	"github.com/hyperledger-labs/cckit/router"
)

func ProposeContract(c router.Context) (interface{}, error) {
	contractPayload := c.Param("contractPayload").(payload.ContractPayload) // Assert the chaincode parameter

	logging.CCLoggerInstance.Printf("Received input: %s. Attempting to validate contract request...\n", contractPayload.String())
	if err := contractPayload.Validate(); err != nil {
		retErr := fmt.Errorf("Error: Validate() returned error: %s", err.Error())
		logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
		return nil, retErr
	}

	for i := 0; i < len(contractPayload.Clause); i++ {
		if err := contractPayload.Clause[i].Validate(contractPayload); err != nil {
			retErr := fmt.Errorf("Error: Validate() returned error: %s", err.Error())
			logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
			return nil, retErr
		}
	}

	logging.CCLoggerInstance.Printf("Checking ACL rules\n")
	owner, err := OnlyContractOrgs(c)
	if err != nil {
		retErr := fmt.Errorf("The user invoking the Contract does not belong in the ACL: %s", err.Error())
		return nil, retErr
	} else if owner != string(contractPayload.Parties[0]) && owner != string(contractPayload.Parties[1]) {
		retErr := fmt.Errorf("The Org invoking the chaincode does not match the Orgs in payload")
		logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
		return nil, retErr
	}

	stateStub := state.NewStateStub(c)
	if err := stateStub.NewContract(contractPayload, owner); err != nil {
		retErr := fmt.Errorf("Error: CreateContract returned error: %s", err.Error())
		logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
		return nil, retErr
	}
	logging.CCLoggerInstance.Printf("Successfully created a Contract between Orgs: %s, %s", contractPayload.Parties[0], contractPayload.Parties[1])
	return nil, nil
}

func AcceptContract(c router.Context) (interface{}, error) {
	licenseSaleID := c.ParamString("license_sale_id")
	logging.CCLoggerInstance.Printf("Received input: %s. Attempting to validate contract request...\n", licenseSaleID)

	if stateContract, err := c.State().Get(state.ContractState{LicenseSaleId: licenseSaleID}, &state.ContractState{}); err != nil {

		retErr := fmt.Errorf("The requested License does not exists: %s", err.Error())
		logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
		return nil, retErr
	} else {
		fmt.Printf("Data of stateContract: %s\n", stateContract)
		acceptedContract := stateContract.(state.ContractState)

		logging.CCLoggerInstance.Printf("Checking ACL rules\n")
		if owner, err := OnlyContractOrgs(c); err != nil {
			retErr := fmt.Errorf("The user invoking the Contract does not belong in the ACL: %s", err.Error())
			return nil, retErr
		} else if owner != string(acceptedContract.Parties[0]) && owner != string(acceptedContract.Parties[1]) {
			retErr := fmt.Errorf("The party invoking the chaincode does not match the Orgs in payload")
			logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
			return nil, retErr
		} else if owner == acceptedContract.Licensee {
			retErr := fmt.Errorf("The party invoking the chaincode cannot be the one accepting it")
			logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
			return nil, retErr
		}
		if acceptedContract.LicenseStatus == "Rejected" || acceptedContract.LicenseStatus == "Accepted" {
			retErr := fmt.Errorf("Error in Contract payload.Contract Status is not 'Pending'.")
			logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
			return nil, retErr
		} else {
			acceptedContract.LicenseStatus = "Accepted"
			fmt.Printf("Data of acceptedContract: %s\n", acceptedContract)
			if err := c.State().Put(acceptedContract); err != nil {
				retErr := fmt.Errorf("Error: Put() returned error: %s", err.Error())
				logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
				return nil, retErr
			}

		}

	}

	return nil, nil
}

func RejectContract(c router.Context) (interface{}, error) {
	licenseSaleID := c.ParamString("license_sale_id")
	logging.CCLoggerInstance.Printf("Received input: %s. Attempting to validate contract request...\n", licenseSaleID)

	if stateContract, err := c.State().Get(state.ContractState{LicenseSaleId: licenseSaleID}, &state.ContractState{}); err != nil {

		retErr := fmt.Errorf("The requested License does not exists: %s", err.Error())
		logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
		return nil, retErr
	} else {
		fmt.Printf("Data of stateContract: %s\n", stateContract)
		rejectedContract := stateContract.(state.ContractState)

		logging.CCLoggerInstance.Printf("Checking ACL rules\n")
		if owner, err := OnlyContractOrgs(c); err != nil {
			retErr := fmt.Errorf("The user invoking the Contract does not belong in the ACL: %s", err.Error())
			return nil, retErr
		} else if owner != string(rejectedContract.Parties[0]) && owner != string(rejectedContract.Parties[1]) {
			retErr := fmt.Errorf("The Org invoking the chaincode does not match the Orgs in payload")
			logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
			return nil, retErr
		} else if owner == rejectedContract.Licensee {
			retErr := fmt.Errorf("The Org invoking the chaincode cannot be the one rejecting it")
			logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
			return nil, retErr
		}

		if rejectedContract.LicenseStatus == "Rejected" || rejectedContract.LicenseStatus == "Accepted" {
			retErr := fmt.Errorf("Error in Contract payload. Contract Status is not 'Pending'.")
			logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
			return nil, retErr
		} else {
			rejectedContract.LicenseStatus = "Rejected"
			fmt.Printf("Data of rejectedContract: %s\n", rejectedContract)
			if err := c.State().Put(rejectedContract); err != nil {
				retErr := fmt.Errorf("Error: Put() returned error: %s", err.Error())
				logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
				return nil, retErr
			}

		}

	}

	return nil, nil
}

func DissolveContract(c router.Context) (interface{}, error) {
	licenseSaleID := c.ParamString("license_sale_id")

	if exists, err := c.State().Exists(state.ContractState{LicenseSaleId: licenseSaleID}); err != nil {
		retErr := fmt.Errorf("Error: Exists() returned error: %s", err.Error())
		logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
		return nil, retErr
	} else if !exists {
		retErr := fmt.Errorf("Error: Invalid delete operation, contract with ID: %s does not exist in contract state", licenseSaleID)
		logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
		return nil, retErr
	} else {
		stateContract, _ := c.State().Get(state.ContractState{LicenseSaleId: licenseSaleID}, &state.ContractState{})
		deletedContract := stateContract.(state.ContractState)

		logging.CCLoggerInstance.Printf("Checking ACL rules\n")
		if owner, err := OnlyContractOrgs(c); err != nil {
			retErr := fmt.Errorf("The user invoking the Contract does not belong in the ACL: %s", err.Error())
			return nil, retErr
		} else if owner != string(deletedContract.Parties[0]) && owner != string(deletedContract.Parties[1]) {
			retErr := fmt.Errorf("The Org invoking the chaincode does not match the Orgs in payload")
			logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
			return nil, retErr
		}
	}

	if err := c.State().Delete(&state.ContractState{LicenseSaleId: licenseSaleID}); err != nil {
		retErr := fmt.Errorf("Error: Delete() returned error: %s", err.Error())
		logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
		return nil, retErr

	}

	return nil, nil
}

func GetContractByID(c router.Context) (interface{}, error) {
	licenseSaleID := c.ParamString("license_sale_id")
	if owner, err := OnlyContractOrgs(c); err != nil {
		retErr := fmt.Errorf("The user invoking the Contract does not belong in the ACL: %s", err.Error())
		return nil, retErr
	} else {
		if stateContract, err := c.State().Get(state.ContractState{LicenseSaleId: licenseSaleID}, &state.ContractState{}); err != nil {
			retErr := fmt.Errorf("The requested Contract does not exist: %s", err.Error())
			logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
			return nil, retErr
		} else {
			stateContractStruct := stateContract.(state.ContractState)
			outputContract := output.OutputContract{
				LicenseSaleId:     	stateContractStruct.LicenseSaleId,
				ProductId:   		stateContractStruct.ProductId,
				LicenseStatus: 		stateContractStruct.LicenseSaleId,
				Orgs:           	stateContractStruct.Orgs,
				Parties:			stateContractStruct.Parties,
				Clauses:         	stateContractStruct.Clauses,
			}
			if owner != string(stateContractStruct.Parties[0]) && owner != string(stateContractStruct.Parties[1]) {
				retErr := fmt.Errorf("The party invoking the chaincode does not match the parties in payload")
				logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
				return nil, retErr
			} else {
				if stateContractStruct.LicenseStatus != "Accepted" {
					retErr := fmt.Errorf("There is no contract with this ID")
					logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
					return nil, retErr
				} else {
					logging.CCLoggerInstance.Printf("Attemting to marshal output...\n")
					marshaledOutput, err := json.Marshal(outputContract)
					if err != nil {
						retErr := fmt.Errorf("error: json.Marshal() of output keys returned error: %s", err.Error())
						logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
						return nil, retErr
					}
					logging.CCLoggerInstance.Printf("Query successfully completed! Returning output: %s\n", string(marshaledOutput))
					return marshaledOutput, nil
				}

			}

		}
	}
}

func GetContracts(c router.Context) (interface{}, error) {
	if querylist, err := c.State().List([]string{state.ContractStateEntity}, &state.ContractState{}); err != nil {
		retErr := fmt.Errorf("Error: List() returned error in list function: %s", err.Error())
		logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
		return nil, retErr
	} else {
		logging.CCLoggerInstance.Printf("Checking ACL rules\n")
		if owner, err := OnlyContractOrgs(c); err != nil {
			retErr := fmt.Errorf("The user invoking the Contract does not belong in the ACL: %s", err.Error())
			return nil, retErr
		} else {

			queriedInterfaceArray := querylist.([]interface{})
			var outputList []output.OutputContract
			for _, curQueriedObj := range queriedInterfaceArray {
				stateContractStruct := curQueriedObj.(state.ContractState)
				outputContract := output.OutputContract{
					LicenseSaleId:     	stateContractStruct.LicenseSaleId,
					ProductId:   		stateContractStruct.ProductId,
					LicenseStatus: 		stateContractStruct.LicenseSaleId,
					Orgs:           	stateContractStruct.Orgs,
					Parties:			stateContractStruct.Parties,
					Clauses:         	stateContractStruct.Clauses,
				}
				if (owner == string(stateContractStruct.Parties[0]) || owner == string(stateContractStruct.Parties[1])) && stateContractStruct.LicenseStatus == "Accepted" {
					outputList = append(outputList, outputContract)
				}
			}
			if len(outputList) == 0 {
				emptyResultArray := make([]state.ContractState, 0)
				logging.CCLoggerInstance.Printf("Query successfully completed! Returning output: %v\n", emptyResultArray)
				return emptyResultArray, nil
			}
			logging.CCLoggerInstance.Printf("Attemting to marshal output...\n")
			marshaledOutput, err := json.Marshal(outputList)
			if err != nil {
				retErr := fmt.Errorf("error: json.Marshal() of output keys returned error: %s", err.Error())
				logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
				return nil, retErr
			}
			logging.CCLoggerInstance.Printf("Query successfully completed! Returning output: %s\n", string(marshaledOutput))
			return marshaledOutput, nil
		}

	}
}

func GetContractIDs(c router.Context) (interface{}, error) {
	if querylist, err := c.State().List([]string{state.ContractStateEntity}, &state.ContractState{}); err != nil {
		retErr := fmt.Errorf("Error: List() returned error in list function: %s", err.Error())
		logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
		return nil, retErr
	} else {
		logging.CCLoggerInstance.Printf("Checking ACL rules\n")
		if owner, err := OnlyContractOrgs(c); err != nil {
			retErr := fmt.Errorf("The user invoking the Contract does not belong in the ACL: %s", err.Error())
			return nil, retErr
		} else {

			queriedInterfaceArray := querylist.([]interface{})
			var outputList []string
			for _, curQueriedObj := range queriedInterfaceArray {
				stateContractStruct := curQueriedObj.(state.ContractState)
				if (owner == string(stateContractStruct.Parties[0]) || owner == string(stateContractStruct.Parties[1])) && stateContractStruct.LicenseStatus == "Accepted" {
					outputList = append(outputList, stateContractStruct.LicenseSaleId)
				}
			}
			if len(outputList) == 0 {
				emptyResultArray := make([]string, 0)
				logging.CCLoggerInstance.Printf("Query successfully completed! Returning output: %v\n", emptyResultArray)
				return emptyResultArray, nil
			}
			logging.CCLoggerInstance.Printf("Attemting to marshal output...\n")
			marshaledOutput, err := json.Marshal(outputList)
			if err != nil {
				retErr := fmt.Errorf("error: json.Marshal() of output keys returned error: %s", err.Error())
				logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
				return nil, retErr
			}
			logging.CCLoggerInstance.Printf("Query successfully completed! Returning output: %s\n", string(marshaledOutput))
			return marshaledOutput, nil
		}

	}
}
/*
func UpdateContractItem(c router.Context) (interface{}, error) {
	contractID := c.ParamString("license_sale_id")
	itemPayload := c.Param("itemPayload").(payload.Item) // Assert the chaincode parameter

	logging.CCLoggerInstance.Printf("Checking Contract Exists()\n")
	if exists, err := c.State().Exists(state.ContractState{ContractId: contractID}); err != nil {
		retErr := fmt.Errorf("Error: Exists() returned error: %s", err.Error())
		logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
		return nil, retErr
	} else if !exists {
		retErr := fmt.Errorf("Error: Invalid delete operation, contract with ID: %s does not exist in contract state", contractID)
		logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
		return nil, retErr
	}

	if stateContract, err := c.State().Get(state.ContractState{ContractId: contractID}, &state.ContractState{}); err != nil {
		retErr := fmt.Errorf("The requested Contract does not exist: %s", err.Error())
		logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
		return nil, retErr
	} else {
		stateContractStruct := stateContract.(state.ContractState)
		stateToPayloadStruct := payload.ContractPayload{
			ContractId:   stateContractStruct.ContractId,
			ContractType: stateContractStruct.ContractType,
			Orgs:         stateContractStruct.Orgs,
			Items:        stateContractStruct.Items,
		}
		logging.CCLoggerInstance.Printf("Checking ACL rules\n")
		if owner, err := OnlyContractOrgs(c); err != nil {
			retErr := fmt.Errorf("The user invoking the Contract does not belong in the ACL: %s", err.Error())
			return nil, retErr
		} else if owner != string(itemPayload.OrgId) || stateContractStruct.ContractStatus != "Accepted" {
			retErr := fmt.Errorf("The Org invoking the chaincode does not match the Item's Org")
			logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
			return nil, retErr
		}
		logging.CCLoggerInstance.Printf("Received input: %s. Attempting to validate contract request...\n", itemPayload.String())
		if err := itemPayload.Validate(stateToPayloadStruct); err != nil {
			retErr := fmt.Errorf("Error: Validate() returned error: %s", err.Error())
			logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
			return nil, retErr
		}
		logging.CCLoggerInstance.Printf("Checking Items array length\n")
		if len(stateContractStruct.Items) == 0 {
			retErr := fmt.Errorf("The requested Contract does not include any Items")
			logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
			return nil, retErr
		} else {
			itemFound := false
			for i, curQueriedObj := range stateContractStruct.Items {
				if curQueriedObj.ObjectId == itemPayload.ObjectId && curQueriedObj.OrgId == itemPayload.OrgId && curQueriedObj.ObjectType == itemPayload.ObjectType {
					itemFound = true
					curQueriedObj.Enabled = itemPayload.Enabled
					curQueriedObj.Write = itemPayload.Write
					stateContractStruct.Items[i] = curQueriedObj
				}
			}
			if !itemFound {
				retErr := fmt.Errorf("The requested Contract does not include the requested Item")
				logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
				return nil, retErr
			} else {
				if err := c.State().Put(stateContractStruct); err != nil {
					retErr := fmt.Errorf("Error: Put() returned error: %s", err.Error())
					return nil, retErr
				}
				fmt.Printf("Updated Contract: %s\n", &stateContractStruct)
				return nil, nil
			}
		}
	}
}

func DeleteContractItem(c router.Context) (interface{}, error) {
	contractID := c.ParamString("license_sale_id")
	itemID := c.ParamString("item_ID")

	logging.CCLoggerInstance.Printf("Checking Contract Exists()\n")
	if exists, err := c.State().Exists(state.ContractState{ContractId: contractID}); err != nil {
		retErr := fmt.Errorf("Error: Exists() returned error: %s", err.Error())
		logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
		return nil, retErr
	} else if !exists {
		retErr := fmt.Errorf("Error: Invalid delete operation, contract with ID: %s does not exist in contract state", contractID)
		logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
		return nil, retErr
	}

	if stateContract, err := c.State().Get(state.ContractState{ContractId: contractID}, &state.ContractState{}); err != nil {
		retErr := fmt.Errorf("The requested Contract does not exist: %s", err.Error())
		logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
		return nil, retErr
	} else {
		stateContractStruct := stateContract.(state.ContractState)

		logging.CCLoggerInstance.Printf("Checking Items array length\n")
		if len(stateContractStruct.Items) == 0 {
			retErr := fmt.Errorf("The requested Contract does not include any Items")
			logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
			return nil, retErr
		} else {
			itemFound := false
			for i, curQueriedObj := range stateContractStruct.Items {
				if curQueriedObj.ObjectId == itemID {
					itemFound = true
					logging.CCLoggerInstance.Printf("Checking ACL rules\n")
					if owner, err := OnlyContractOrgs(c); err != nil {
						retErr := fmt.Errorf("The user invoking the Contract does not belong in the ACL: %s", err.Error())
						return nil, retErr
					} else if owner != string(curQueriedObj.OrgId) || stateContractStruct.ContractStatus != "Accepted" {
						retErr := fmt.Errorf("The Org invoking the chaincode does not match the Item's Org")
						logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
						return nil, retErr
					}
					stateContractStruct.Items = append(stateContractStruct.Items[:i], stateContractStruct.Items[i+1:]...)
					break
				}
			}
			if !itemFound {
				retErr := fmt.Errorf("The requested Contract does not include the requested Item")
				logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
				return nil, retErr
			} else {
				if err := c.State().Put(stateContractStruct); err != nil {
					retErr := fmt.Errorf("Error: Put() returned error: %s", err.Error())
					return nil, retErr
				}
				fmt.Printf("Updated Contract: %s\n", &stateContractStruct)
				return nil, nil
			}
		}
	}
}

func AddContractItem(c router.Context) (interface{}, error) {
	contractID := c.ParamString("license_sale_id")
	itemPayload := c.Param("itemPayload").(payload.Item) // Assert the chaincode parameter

	logging.CCLoggerInstance.Printf("Checking Contract Exists()\n")
	if exists, err := c.State().Exists(state.ContractState{ContractId: contractID}); err != nil {
		retErr := fmt.Errorf("Error: Exists() returned error: %s", err.Error())
		logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
		return nil, retErr
	} else if !exists {
		retErr := fmt.Errorf("Error: Invalid delete operation, contract with ID: %s does not exist in contract state", contractID)
		logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
		return nil, retErr
	}

	if stateContract, err := c.State().Get(state.ContractState{ContractId: contractID}, &state.ContractState{}); err != nil {
		retErr := fmt.Errorf("The requested Contract does not exist: %s", err.Error())
		logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
		return nil, retErr
	} else {
		stateContractStruct := stateContract.(state.ContractState)
		stateItemsArray := stateContractStruct.Items
		stateToPayloadStruct := payload.ContractPayload{
			ContractId:   stateContractStruct.ContractId,
			ContractType: stateContractStruct.ContractType,
			Orgs:         stateContractStruct.Orgs,
			Items:        stateContractStruct.Items,
		}
		logging.CCLoggerInstance.Printf("Checking ACL rules\n")
		if owner, err := OnlyContractOrgs(c); err != nil {
			retErr := fmt.Errorf("The user invoking the Contract does not belong in the ACL: %s", err.Error())
			return nil, retErr
		} else if owner != string(itemPayload.OrgId) || stateContractStruct.ContractStatus != "Accepted" {
			retErr := fmt.Errorf("The Org invoking the chaincode does not match the Item's Org")
			logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
			return nil, retErr
		}
		logging.CCLoggerInstance.Printf("Received input: %s. Attempting to validate contract request...\n", itemPayload.String())
		if err := itemPayload.Validate(stateToPayloadStruct); err != nil {
			retErr := fmt.Errorf("Error: Validate() returned error: %s", err.Error())
			logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
			return nil, retErr
		}
		logging.CCLoggerInstance.Printf("Checking Items array length\n")

		itemFound := false
		for _, curQueriedObj := range stateItemsArray {
			if curQueriedObj.ObjectId == itemPayload.ObjectId {
				itemFound = true
			}
		}
		if itemFound {
			retErr := fmt.Errorf("The requested Item is already included into the Contract")
			logging.CCLoggerInstance.Printf("%s\n", retErr.Error())
			return nil, retErr
		} else {
			stateContractStruct.Items = append(stateContractStruct.Items, itemPayload)
			if err := c.State().Put(stateContractStruct); err != nil {
				retErr := fmt.Errorf("Error: Put() returned error: %s", err.Error())
				return nil, retErr
			}
			fmt.Printf("Updated Contract: %s\n", &stateContractStruct)
			return nil, nil
		}

	}
}
*/