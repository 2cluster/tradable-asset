package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	ledgerapi "github.com/2cluster/tradable-asset/ledger-api"
)

type State uint

const (
	ISSUED State = iota + 1
	NEGOTIATION
	SIGNING
	ACTIVE
	REDEEMED
)

func (state State) String() string {
	names := []string{"ISSUED", "NEGOTIATION", "SIGNING", "ACTIVE", "REDEEMED"}

	if state < ISSUED || state > REDEEMED {
		return "UNKNOWN"
	}

	return names[state-1]
}


type Identity struct {
	ID				string `json:"id"`
	Org				string `json:"org"`
	Name			string `json:"name"`
	EthAdr			string `json:"ethAdr"`
}


func CreateKey(name string, id string) string {
	return ledgerapi.MakeKey(name, id)
}


type assetAlias Asset
type jsonAsset struct {
	*assetAlias
	State State  `json:"currentState"`
	Class  string `json:"class"`
	Key   string `json:"key"`
}


type Asset struct {
	ID             string `json:"ID"`
	Color          string `json:"color"`
	Size           int    `json:"size"`
	Owner		   Identity `json:"owner"`
	AppraisedValue int    `json:"appraisedValue"`

	Lender	   	   Identity `json:"lender"`
	Lsignature	   bool `json:"lSignature"`
	Borrower	   Identity `json:"borrower"`
	Bsignature	   bool `json:"BSignature"`

	state          State  `metadata:"currentState"`
	class          string `metadata:"class"`
	key            string `metadata:"key"`
}

type buyRequest struct {
	Owner		   Identity `json:"owner"`
	Buyer		   Identity `json:"buyer"`
	Message		   string `json:"message"`
}

type HistoryQueryResult struct {
	Record    *Asset    `json:"record"`
	TxId     string     `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}


func (cp *Asset) UnmarshalJSON(data []byte) error {
	jcp := jsonAsset{assetAlias: (*assetAlias)(cp)}

	err := json.Unmarshal(data, &jcp)

	if err != nil {
		return err
	}

	cp.state = jcp.State

	return nil
}

// MarshalJSON special handler for managing JSON marshalling
func (cp Asset) MarshalJSON() ([]byte, error) {
	jcp := jsonAsset{assetAlias: (*assetAlias)(&cp), State: cp.state, Class:"Asset", Key: ledgerapi.MakeKey(cp.Owner.Name, cp.ID)}

	return json.Marshal(&jcp)
}



// GetState returns the state
func (cp *Asset) GetState() State {
	return cp.state
}

// SetIssued returns the state to issued
func (cp *Asset) SetIssued() {
	cp.state = ISSUED
}

// SetIssued returns the state to issued
func (cp *Asset) SetNegotiation() {
	cp.state = NEGOTIATION
}

// SetIssued returns the state to issued
func (cp *Asset) SetSigning() {
	cp.state = SIGNING
}

// SetActive sets the state to trading
func (cp *Asset) SetActive() {
	cp.state = ACTIVE
}

// SetRedeemed sets the state to redeemed
func (cp *Asset) SetRedeemed() {
	cp.state = REDEEMED
}

// IsIssued returns true if state is issued
func (cp *Asset) IsIssued() bool {
	return cp.state == ISSUED
}

// IsActive returns true if state is issued
func (cp *Asset) IsActive() bool {
	return cp.state == ISSUED
}

// IsActive returns true if state is issued
func (cp *Asset) IsSigning() bool {
	return cp.state == SIGNING
}


// IsNegotiation returns true if state is trading
func (cp *Asset) IsNegotiation() bool {
	return cp.state == NEGOTIATION
}

// IsRedeemed returns true if state is redeemed
func (cp *Asset) IsRedeemed() bool {
	return cp.state == REDEEMED
}

// GetSplitKey returns values which should be used to form key
func (cp *Asset) GetSplitKey() []string {
	return []string{cp.Owner.Name, cp.ID}
}

// Serialize formats the commercial paper as JSON bytes
func (cp *Asset) Serialize() ([]byte, error) {
	return json.Marshal(cp)
}

// Deserialize formats the commercial paper from JSON bytes
func Deserialize(bytes []byte, cp *Asset) error {
	err := json.Unmarshal(bytes, cp)

	if err != nil {
		return fmt.Errorf("Error deserializing commercial paper. %s", err.Error())
	}

	return nil
}