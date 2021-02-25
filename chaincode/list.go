/*
 * SPDX-License-Identifier: Apache-2.0
 */

 package chaincode

 import ledgerapi "github.com/2cluster/tradable-asset/ledger-api"
 
 // ListInterface defines functionality needed
 // to interact with the world state on behalf
 // of a commercial paper
 type ListInterface interface {
	 AddAsset(*Asset) error
	 GetAsset(string, string) (*Asset, error)
	 UpdateAsset(*Asset) error
	 DeleteAsset(string, string) error
 }
 
 type list struct {
	 stateList ledgerapi.StateListInterface
 }
 
 func (cpl *list) AddAsset(asset *Asset) error {
	 return cpl.stateList.AddState(asset)
 }

 func (cpl *list) DeleteAsset(owner string, assetID string) error {
	return cpl.stateList.DelState(CreateKey(owner, assetID))
}
 
 func (cpl *list) GetAsset(owner string, assetID string) (*Asset, error) {
	 cp := new(Asset)
 
	 err := cpl.stateList.GetState(CreateKey(owner, assetID), cp)
 
	 if err != nil {
		 return nil, err
	 }

	 if cp.ID == "" {
		 return nil, nil
	 }
 
	 return cp, nil
 }
 
 func (cpl *list) UpdateAsset(paper *Asset) error {
	 return cpl.stateList.UpdateState(paper)
 }
 
 // NewList create a new list from context
 func newList(ctx TransactionContextInterface) *list {
	 stateList := new(ledgerapi.StateList)
	 stateList.Ctx = ctx
	 stateList.Class = "Asset"
	 stateList.Deserialize = func(bytes []byte, state ledgerapi.StateInterface) error {
		 return Deserialize(bytes, state.(*Asset))
	 }
 
	 list := new(list)
	 list.stateList = stateList
 
	 return list
 }
 