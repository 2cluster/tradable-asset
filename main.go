// SPDX-License-Identifier: Undefined


package main

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/2cluster/tradable-asset/chaincode"
	// "github.com/spf13/viper"
)

func main() {

	contract := new(chaincode.SmartContract)
	contract.TransactionContextHandler = new(chaincode.TransactionContext)
	// contract.Name = "Dealblock"
	// contract.Info.Version = "0.0.1"

	chaincode, err := contractapi.NewChaincode(contract)

	if err != nil {
		panic(fmt.Sprintf("Error creating chaincode. %s", err.Error()))
	}

	// chaincode.Info.Title = "AssetChaincode"
	// chaincode.Info.Version = "0.0.1"

	err = chaincode.Start()

	if err != nil {
		panic(fmt.Sprintf("Error starting chaincode. %s", err.Error()))
	}
}