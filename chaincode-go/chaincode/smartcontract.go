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
	Location 	string `json:"Location"`
	Time	 string `json:"Time"`
}

type Candidate struct {
	CandidateNo	 string `json:"CandidateNo"`
	VoteCount int       `json:"VoteCount"`
}                 

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
/* 	assets := []Voter{
		{RandomNo: "NULL", VotedCandidateID: "NULL", Location: "NULL", Time: "NULL"},
	}
	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.RandomNo, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	} */

	candidates := []Candidate{
		{CandidateNo: "1", VoteCount: 0}, // 오한수
		{CandidateNo: "2", VoteCount: 0}, // 최재호
		{CandidateNo: "3", VoteCount: 0}, // 기권
	}

	for _, candidateAsBytes := range candidates {
		candidateJSON, err := json.Marshal(candidateAsBytes)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(candidateAsBytes.CandidateNo, candidateJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

func (e *SmartContract) RegisterVoter(ctx contractapi.TransactionContextInterface, RandomNo string, VotedCandidateID string, Location string, Time string) error {
	fmt.Println("=============== Start Register Voter =============== ")

	randomNo := RandomNo
	location := Location
	time := Time
	votedCandidateID := VotedCandidateID

	voter := Voter{RandomNo: randomNo, Location: location, Time: time, VotedCandidateID: votedCandidateID}

	voterAsBytes, err := json.Marshal(voter)
	err = ctx.GetStub().PutState(voter.RandomNo, voterAsBytes)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}
	
	return nil
}

 func (e *SmartContract) AddVote(ctx contractapi.TransactionContextInterface, CandidateNo string) error {
	fmt.Println("=============== Start Add Vote =============== ")
	
	candidateAsBytes, err := ctx.GetStub().GetState(CandidateNo)

	if err != nil {
		return fmt.Errorf("faild: %v", err)
	}

	candidate := Candidate{}

	json.Unmarshal(candidateAsBytes, &candidate)

	candidate.CandidateNo = CandidateNo
	candidate.VoteCount++

	candidateByBytes, err := json.Marshal(candidate)
	if err != nil {
		return fmt.Errorf("faild: %v", err)
	}

	err =  ctx.GetStub().PutState(CandidateNo, candidateByBytes)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}
	return nil
}

/* // CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, RandomNo string, VotedCandidateID string, Location string, Time string) error {
	asset := Voter{

		RandomNo: RandomNo,
		VotedCandidateID: VotedCandidateID,
		Location:  Location,
		Time:      Time,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(RandomNo, assetJSON)
}
 */

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

AssetExists returns true when asset with given ID exists in world state
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
func (smartcontract *SmartContract) QueryCandidate(ctx contractapi.TransactionContextInterface, CandidateNo string) (*Candidate, error) {
	fmt.Println("=============== Start Query Candidate =============== ")

	candidateAsBytes, err := ctx.GetStub().GetState(CandidateNo)

	if err != nil {
		return nil, err
	}

	if candidateAsBytes == nil {
		return nil, err
	}

	candidate := new(Candidate)
	_ = json.Unmarshal(candidateAsBytes, candidate)

	return candidate, nil
}


 func (s *SmartContract) GetAllVoters(ctx contractapi.TransactionContextInterface) ([]*Voter, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets2 []*Voter
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset2 Voter   
		err = json.Unmarshal(queryResponse.Value, &asset2)
		if err != nil {
			return nil, err
		}
		assets2 = append(assets2, &asset2)
	}

	return assets2, nil
}


// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Candidate, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Candidate
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Candidate   
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}