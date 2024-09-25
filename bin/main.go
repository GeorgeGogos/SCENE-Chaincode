package main

import (
	"fmt"

	scene "github.com/GeorgeGogos/SCENE-Chaincode"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

func main() {
	cc := scene.NewCC()
	if err := shim.Start(cc); err != nil {
		fmt.Printf("Error while attempting to start chaincode: %s\n", err.Error())
	}
}
