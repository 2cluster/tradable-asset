package chaincode

import (
	"encoding/json"
	"fmt"

	eth "github.com/2cluster/ethclient/client"
	"github.com/ethereum/go-ethereum/common"
	// "log"
	// ledgerapi "github.com/2cluster/tradable-asset/ledger-api"
)

func (s *SmartContract) RequestToBuy(ctx TransactionContextInterface, transientID string, assetid string, assetorg string) (error) {

	exists, err := assetExists(ctx, assetorg, assetid)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s:%s does not exist", assetorg, assetid)
	}
	
	asset, err := ctx.GetAssetList().GetAsset(assetorg, assetid)
	if err != nil {
		return err
	}

	if asset.GetState().String() != "ISSUED" {
		return fmt.Errorf("this asset is not available for buy requests")
	}

	buyer, err := getIdentity(ctx)
	if err != nil {
		return err
	}

	// Get new asset from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("error getting transient: %v", err)
	}

	// Asset properties are private, therefore they get passed in transient field, instead of func args
	transientData, ok := transientMap[transientID]
	if !ok {
		//log error to stdout
		return fmt.Errorf("asset not found in the transient map input")
	}

	var req buyRequest

	err = json.Unmarshal(transientData, &req)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	
	req.Buyer = *buyer
	req.Owner = (Identity)(asset.Owner)

	if len(req.Message) == 0 {
		return fmt.Errorf("Message field must be a non-empty string")
	}

	orgCollection, err := getCollectionName(asset.Owner.Org)
	if err != nil {
		return fmt.Errorf("failed to infer private collection name for the org: %v", err)
	}

	requestBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal asset into JSON: %v", err)
	}

	err = ctx.GetStub().PutPrivateData(orgCollection, asset.ID, requestBytes)
	if err != nil {
		return fmt.Errorf("failed to put asset into private data collecton: %v", err)
	}

	asset.SetNegotiation()

	err = ctx.GetAssetList().UpdateAsset(asset)
	if err != nil {
		return err
	}
	return nil
}

func (s *SmartContract) GetRequests(ctx TransactionContextInterface) ([]*buyRequest, error) {

	identity, err := getIdentity(ctx)
	if err != nil {
		return nil, err
	}

	ownerCollection, err := getCollectionName(identity.Org)
	if err != nil {
		return nil, fmt.Errorf("failed to infer private collection name for the org: %v", err)
	}

	queryString := fmt.Sprintf(`{"selector":{"owner": {"id": "%s"}}}`, identity.ID)

	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult(ownerCollection, queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to read asset from owner's Collection: %v", err)
	}

	defer resultsIterator.Close()

	return constructQueryResponseFromIteratorRequests(resultsIterator)

}

func (s *SmartContract) AgreeToSell(ctx TransactionContextInterface, assetid string) (error) {
	seller, err := getIdentity(ctx)
	if err != nil {
		return err
	}

	ownerCollection, err := getCollectionName(seller.Org)
	if err != nil {
		return fmt.Errorf("failed to infer private collection name for the org: %v", err)
	}

	bytes, err := ctx.GetStub().GetPrivateData(ownerCollection, assetid)
	if err != nil {
		return fmt.Errorf("failed to read asset from owner's Collection: %v", err)
	}
	var request buyRequest
	err = json.Unmarshal(bytes, &request)
	if err != nil {
		return fmt.Errorf("Error deserializing commercial paper. %s", err.Error())
	}
	
	asset, err := ctx.GetAssetList().GetAsset(seller.Name, assetid)
	if err != nil {
		return err
	}

	asset.Lender = request.Buyer
	asset.Borrower = request.Owner

	asset.SetSigning()

	err = ctx.GetAssetList().UpdateAsset(asset)
	if err != nil {
		return err
	}

	return nil
}


// func (s *SmartContract) Sign(ctx TransactionContextInterface, ownerName string, assetid string) (error) {
// 	caller, err := getIdentity(ctx)
// 	if err != nil {
// 		return err
// 	}
	
// 	asset, err := ctx.GetAssetList().GetAsset(ownerName, assetid)
// 	if err != nil {
// 		return err
// 	}

// 	if caller.ID == asset.Lender.ID {
// 		return fmt.Errorf("failed to add org to endorsement policy: %v", err)
// 	}
// 	if caller.ID == asset.Borrower.ID {
// 		return fmt.Errorf("failed to add org to endorsement policy: %v", err)
// 	}

// 	asset.Lender = request.Buyer
// 	asset.Borrower = request.Owner

// 	asset.SetSigning()

// 	err = ctx.GetAssetList().UpdateAsset(asset)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }



func (s *SmartContract) Query(ctx TransactionContextInterface) string {
	caller, err := getIdentity(ctx)
	if err != nil {
		fmt.Println(err)
	}

	client, err:= eth.NewClient(caller.Name, "f238a37e42b7062bdbc062a1833a6361f9a6d0e324a95ca2f7c4c3034e67ee5c")
	if err != nil {
		fmt.Println(err)
	}

	client.BindContract(common.HexToAddress("0xC382b4aF66EDb6Aa717B6A07330d41364b787B02"))
	balance := client.QueryBalance(client.Account.Address)
	return balance
	
}