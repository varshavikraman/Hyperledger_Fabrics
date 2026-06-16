package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// Mock database for when MOCK_FABRIC=true is active
var (
	mockLedger = make(map[string]string)
	mockMu     sync.RWMutex
)

func init() {
	// Initialize the mock database with sample data
	mockLedger["CAR0"] = `{"carId":"CAR0","make":"Toyota","model":"Prius","color":"Blue","dateOfManufacture":"2023-01-01","manufacturerName":"Toyota Inc"}`
	mockLedger["CAR1"] = `{"carId":"CAR1","make":"Ford","model":"Mustang","color":"Red","dateOfManufacture":"2024-05-12","manufacturerName":"Ford Motor Company"}`
	mockLedger["CAR2"] = `{"carId":"CAR2","make":"Tesla","model":"Model Y","color":"White","dateOfManufacture":"2025-10-20","manufacturerName":"Tesla Inc"}`
}

// submitTxnFn interacts with the Hyperledger Fabric network or fallback mock ledger
func submitTxnFn(
	org string,
	channel string,
	chaincodeID string,
	contractName string,
	action string,
	transient map[string][]byte,
	functionName string,
	args ...string,
) string {
	profile := loadProfile()

	if profile.UseMock {
		return handleMockTxn(functionName, args...)
	}

	// Real Fabric Gateway Connection
	gw, err := getGateway()
	if err != nil {
		log.Printf("Gateway connection error: %v", err)
		return fmt.Sprintf(`{"error":"Gateway connection error: %v"}`, err)
	}

	network := gw.GetNetwork(channel)
	var contract *client.Contract
	if contractName != "" {
		contract = network.GetContractWithName(chaincodeID, contractName)
	} else {
		contract = network.GetContract(chaincodeID)
	}

	var txnResult []byte
	if action == "invoke" {
		// Invoke transaction (state-modifying)
		proposal, err := contract.NewProposal(
			functionName,
			client.WithArguments(args...),
			client.WithTransient(transient),
		)
		if err != nil {
			log.Printf("Failed to create proposal: %v", err)
			return fmt.Sprintf(`{"error":"Failed to create proposal: %v"}`, err)
		}

		transaction, err := proposal.Endorse()
		if err != nil {
			log.Printf("Failed to endorse transaction: %v", err)
			return fmt.Sprintf(`{"error":"Failed to endorse transaction: %v"}`, err)
		}

		txnResult = transaction.Result()
		_, err = transaction.Submit()
		if err != nil {
			log.Printf("Failed to submit transaction: %v", err)
			return fmt.Sprintf(`{"error":"Failed to submit transaction: %v"}`, err)
		}
	} else {
		// Query action (evaluate only)
		txnResult, err = contract.EvaluateTransaction(
			functionName,
			args...,
		)
		if err != nil {
			log.Printf("Failed to evaluate transaction: %v", err)
			return fmt.Sprintf(`{"error":"Failed to evaluate transaction: %v"}`, err)
		}
	}

	return string(txnResult)
}

func handleMockTxn(functionName string, args ...string) string {
	mockMu.Lock()
	defer mockMu.Unlock()

	log.Printf("[MOCK LEDGER] Executing %s with args: %v", functionName, args)

	switch functionName {
	case "CreateCar":
		if len(args) < 6 {
			return `{"error":"Invalid arguments. Expected 6 arguments for CreateCar."}`
		}
		carID := args[0]
		car := Car{
			CarId:        carID,
			Make:         args[1],
			Model:        args[2],
			Color:        args[3],
			Manufacturer: args[4],
			Date:         args[5],
		}
		carJSON, err := json.Marshal(car)
		if err != nil {
			return fmt.Sprintf(`{"error":"Failed to marshal car: %s"}`, err.Error())
		}
		mockLedger[carID] = string(carJSON)
		return string(carJSON)

	case "ReadCar":
		if len(args) < 1 {
			return `{"error":"Invalid arguments. Expected 1 argument for ReadCar."}`
		}
		carID := args[0]
		carData, exists := mockLedger[carID]
		if !exists {
			return fmt.Sprintf(`{"error":"Car with ID '%s' not found on the mock ledger."}`, carID)
		}
		return carData

	default:
		return fmt.Sprintf(`{"error":"Unsupported function: %s"}`, functionName)
	}
}
