/*
 * SPDX-License-Identifier: Apache-2.0
 */

 package chaincode

 import (
	 "github.com/hyperledger/fabric-contract-api-go/contractapi"
 )
 
 // TransactionContextInterface an interface to
 // describe the minimum required functions for
 // a transaction context in the commercial
 // paper
 type TransactionContextInterface interface {
	 contractapi.TransactionContextInterface
	 GetAssetList() ListInterface
 }
 
 type TransactionContext struct {
	 contractapi.TransactionContext
	 assetList *list
 }
 
 func (tc *TransactionContext) GetAssetList() ListInterface {
	 if tc.assetList == nil {
		 tc.assetList = newList(tc)
	 }
 
	 return tc.assetList
 }
 