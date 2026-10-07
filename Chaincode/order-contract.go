package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type OrderContract struct {
	contractapi.Contract
}

type Order struct {
	OrderID    string `json:"orderID"`
	Make       string `json:"make"`
	Model      string `json:"model"`
	Color      string `json:"color"`
	DealerName string `json:"dealerName"`
}

func (o *OrderContract) CreateOrder(
	ctx contractapi.TransactionContextInterface,
	orderID string,
) (string, error) {

	transientData, err := ctx.GetStub().GetTransient()
	if err != nil {
		return "", fmt.Errorf("failed to get transient data: %v", err)
	}

	makeBytes, ok := transientData["make"]
	if !ok {
		return "", fmt.Errorf("make was not provided")
	}

	modelBytes, ok := transientData["model"]
	if !ok {
		return "", fmt.Errorf("model was not provided")
	}

	colorBytes, ok := transientData["color"]
	if !ok {
		return "", fmt.Errorf("color was not provided")
	}

	dealerBytes, ok := transientData["dealerName"]
	if !ok {
		return "", fmt.Errorf("dealerName was not provided")
	}

	order := Order{
		OrderID:    orderID,
		Make:       string(makeBytes),
		Model:      string(modelBytes),
		Color:      string(colorBytes),
		DealerName: string(dealerBytes),
	}

	orderJSON, err := json.Marshal(order)
	if err != nil {
		return "", fmt.Errorf("failed to marshal order: %v", err)
	}

	err = ctx.GetStub().PutPrivateData("orderCollection", orderID, orderJSON)
	if err != nil {
		return "", fmt.Errorf("failed to put order into private collection: %v", err)
	}

	return fmt.Sprintf("successfully created order %s", orderID), nil
}

func (o *OrderContract) ReadOrder(
	ctx contractapi.TransactionContextInterface,
	orderID string,
) (*Order, error) {

	orderJSON, err := ctx.GetStub().GetPrivateData("orderCollection", orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to read order: %v", err)
	}

	if orderJSON == nil {
		return nil, fmt.Errorf("order %s does not exist", orderID)
	}

	var order Order

	err = json.Unmarshal(orderJSON, &order)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal order: %v", err)
	}

	return &order, nil
}
