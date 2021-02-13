package chaincode

import (
	"encoding/json"
	"fmt"

	// "github.com/hyperledger/fabric-chaincode-go/pkg/statebased"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	Type           string `json:"ObjectType"`
	ID             string `json:"ID"`
	Color          string `json:"color"`
	Size           int    `json:"size"`
	Owner          string `json:"owner"`
	AppraisedValue int    `json:"appraisedValue"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) string {

	return "Successfully initialized!"
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, appraisedValue int) (*Asset, error) {
	
	clientOrgID, err := getClientOrgID(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get verified OrgID: %v", err)
	}
	
	exists, err := assetExists(ctx, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("the asset %s already exists", id)
	}

	asset := Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          clientOrgID,
		AppraisedValue: appraisedValue,
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return nil, err
	}


	err = ctx.GetStub().PutState(id, assetJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to put asset in public data: %v", err)
	}

	return &asset, nil
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, appraisedValue int) (*Asset, error) {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return nil, err
	}

	clientOrgID, err := getClientOrgID(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get verified OrgID: %v", err)
	}

	if clientOrgID != asset.Owner {
		return nil, fmt.Errorf("%s is not the owner of this asset",clientOrgID)
	}

	// overwriting original asset with new asset
	updatedAsset := Asset{
		ID:             asset.ID,
		Color:          color,
		Size:           size,
		Owner:          asset.Owner,
		AppraisedValue: appraisedValue,
	}


	updatedJSON, err := json.Marshal(updatedAsset)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(id, updatedJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to put asset in public data: %v", err)
	}

	return &updatedAsset, nil
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return nil, err
	}

	clientOrgID, err := getClientOrgID(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get verified OrgID: %v", err)
	}

	if clientOrgID != asset.Owner {
		return nil, fmt.Errorf("%s is not the owner of this asset",clientOrgID)
	}

	err = ctx.GetStub().DelState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete asset: %v", err)
	}
	return new(Asset), nil
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) (*Asset, error) {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return nil, err
	}

	clientOrgID, err := getClientOrgID(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get verified OrgID: %v", err)
	}

	if clientOrgID != asset.Owner {
		return nil, fmt.Errorf("%s is not the owner of this asset",clientOrgID)
	}

	asset.Owner = newOwner
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(id, assetJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to put asset in public data: %v", err)
	}

	return asset, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func assetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

func getClientOrgID(ctx contractapi.TransactionContextInterface, verifyOrg bool) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("failed getting client's orgID: %v", err)
	}

	if verifyOrg {
		err = verifyClientOrgMatchesPeerOrg(clientOrgID)
		if err != nil {
			return "", err
		}
	}

	return clientOrgID, nil
}

func verifyClientOrgMatchesPeerOrg(clientOrgID string) error {
	peerOrgID, err := shim.GetMSPID()
	if err != nil {
		return fmt.Errorf("failed getting peer's orgID: %v", err)
	}

	if clientOrgID != peerOrgID {
		return fmt.Errorf("client from org %s is not authorized to read or write private data from an org %s peer",
			clientOrgID,
			peerOrgID,
		)
	}

	return nil
}
