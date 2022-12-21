package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
type Voter struct {
	RandomNo       string `json:"RandomNo"`
	VotedCandidateID string `json:"VotedCandidateID"`
	location 	string `json:"Location"`
	Time	 string `json:"Time"`

}

// Candidate defined as struct
type Candidate struct {
	CandidateID string `json:"CandidateID"`
	Name        string `json:"Name"`
	TotalVote   int    `json:"TotalVote"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	candidates := []Candidate{
		Candidate{CandidateID: "1", Name: "오한수", TotalVote: 0},
		Candidate{CandidateID: "2", Name: "최재호", TotalVote: 0},
		Candidate{CandidateID: "3", Name: "기권", TotalVote: 0},
	}

	for _, candidate:= range candidates {
		candidateJSON, err := json.Marshal(candidate)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(candidate.CandidateID, candidateJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

func (s *SmartContract) addVote(ctx contractapi.TransactionContextInterface, RandomNo string, VotedCandidateID string, Location string, Time string) error {
	fmt.Println("=============== Start Add Vote =============== ")

	candidateID := VotedCandidateID

	voterAsBytes, err := ctx.GetStub().GetState("VOTER" + RandomNo)
	candidateAsBytes, err := ctx.GetStub().GetState("CANDIDATE" + candidateID)

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	voter := Voter{}
	candidate := Candidate{}

	json.Unmarshal(voterAsBytes, &voter)
	json.Unmarshal(candidateAsBytes, &candidate)

	if voter.VotedCandidateID != "" {
		return fmt.Errorf("Voter already voted a candidate")
	}

	voter.VotedCandidateID = candidateID
	candidate.TotalVote++

	voterByBytes, _ := json.Marshal(voter)
	candidateByBytes, _ := json.Marshal(candidate)

	err = ctx.GetStub().PutState("VOTER"+voter.RandomNo, voterByBytes)
	err = ctx.GetStub().PutState("CANDIDATE"+candidate.CandidateID, candidateByBytes)

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	fmt.Println("=============== End Add Vote =============== ")
	return nil
}

/* // CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, v_num int, c_name string, location string, time string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := Asset{
		id:             id,
		v_num:          v_num,
		c_name:         c_name,
		location:       location,
		time:	 time,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
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
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, v_num int, c_name string, location string, time string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	asset := Asset{
		id:             id,
		v_num:          v_num,
		c_name:         c_name,
		location:       location,
		time:	 	time,
	}	
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, location string) error {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return err
	}

	asset.location = location
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}
 */
// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Candidate, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var candidates []*Candidate
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var candidate Candidate
		err = json.Unmarshal(queryResponse.Value, &candidate)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, &candidate) 
	}

	return candidates, nil
}