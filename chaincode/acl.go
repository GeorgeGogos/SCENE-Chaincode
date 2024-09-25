package chaincode

import (
	"fmt"
	"crypto/sha256"
	"encoding/base64"
	logging "github.com/CERTH-ITI-DLT-Lab/hlf-cc-logging"
	"github.com/hyperledger-labs/cckit/router"
)

func OnlyContractOrgs(c router.Context) (string, error) {
	invokerID, err := GetInvokerIDFromContext(c)
	if err != nil {
		return "", fmt.Errorf("ACL::onlyContractOrgs() Error: GetOwnerIDFromContext() returned error: %s", err.Error())
	}
	logging.CCLoggerInstance.Println("ACL::onlyContractOrgs(): Access granted!")
	return invokerID, nil				
}

func GetInvokerIDFromContext(c router.Context) (string, error) {
	cid, err := c.Client()
	if err != nil {
		return "", fmt.Errorf("Error: Client() returned error: %s", err.Error())
	}
	mspID, err := cid.GetMSPID()
	if err != nil {
		return "", fmt.Errorf("Error: GetMSPID() returned error: %s", err.Error())
	}
	clientID, err := cid.GetID()
	if err != nil {
		return "", fmt.Errorf("Error: GetID() returned error: %s", err.Error())
	}
	x509HashIDURLEncoded := base64.RawURLEncoding.
		EncodeToString(sha256.New().Sum([]byte(fmt.Sprintf("x509:%s:%s", mspID, clientID))))
	return fmt.Sprintf("userid:%s", x509HashIDURLEncoded), nil
}