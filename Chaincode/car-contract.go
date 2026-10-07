package main

import (
	"encoding/json"
	"fmt"
 "github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// CarContract provides functions for managing cars.
type CarContract struct {
	contractapi.Contract
}

// Car represents a car stored in the ledger.
type Car struct {
	AssetType         string `json:"AssetType"`
	CarId             string `json:"CarId"`
	Color             string `json:"Color"`
	DateOfManufacture string `json:"DateOfManufacture"`
	OwnedBy           string `json:"OwnedBy"`
	Make              string `json:"Make"`
	Model             string `json:"Model"`
	Status            string `json:"Status"`
}

// CarExists returns true when a car with the given ID exists.
func (c *CarContract) CarExists(ctx contractapi.TransactionContextInterface, carID string) (bool, error) {
	data, err := ctx.GetStub().GetState(carID)

	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return data != nil, nil
}

// CreateCar creates a new car.
func (c *CarContract) CreateCar(
	ctx contractapi.TransactionContextInterface,
	carID string,
	make string,
	model string,
	color string,
	manufacturerName string,
	dateOfManufacture string,
) (string, error) {

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()

	if err != nil {
		return "", err
	}

	if clientOrgID != "Org1MSP" {
		return "", fmt.Errorf(
			"User under following MSPID: %v can't perform this action",
			clientOrgID,
		)
	}

	exists, err := c.CarExists(ctx, carID)

	if err != nil {
		return "", fmt.Errorf(
			"Could not read from world state. %s",
			err,
		)
	}

	if exists {
		return "", fmt.Errorf(
			"The asset %s already exists",
			carID,
		)
	}

	car := Car{
		AssetType:         "car",
		CarId:              carID,
		Color:              color,
		DateOfManufacture: dateOfManufacture,
		Make:               make,
		Model:              model,
		OwnedBy:            manufacturerName,
		Status:             "In Factory",
	}

	bytes, err := json.Marshal(car)

	if err != nil {
		return "", fmt.Errorf(
			"failed to marshal car: %v",
			err,
		)
	}

	err = ctx.GetStub().PutState(carID, bytes)

	if err != nil {
		return "", err
	}
type EventData struct {
    Type  string `json:"type"`
    Model string `json:"model"`
}

addCarEventData := EventData{
    Type:  "Car creation",
    Model: model,
}

eventDataByte, err := json.Marshal(addCarEventData)
if err != nil {
    return "", fmt.Errorf("failed to marshal event data: %v", err)
}

err = ctx.GetStub().SetEvent("CreateCar", eventDataByte)
if err != nil {
    return "", fmt.Errorf("failed to set CreateCar event: %v", err)
}
	return fmt.Sprintf(
		"successfully added car %v",
		carID,
	), nil
}

// ReadCar retrieves a car from the world state.
func (c *CarContract) ReadCar(
	ctx contractapi.TransactionContextInterface,
	carID string,
) (*Car, error) {

	exists, err := c.CarExists(ctx, carID)

	if err != nil {
		return nil, fmt.Errorf(
			"Could not read from world state. %s",
			err,
		)
	}

	if !exists {
		return nil, fmt.Errorf(
			"The asset %s does not exist",
			carID,
		)
	}

	bytes, err := ctx.GetStub().GetState(carID)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to read car: %v",
			err,
		)
	}

	car := new(Car)

	err = json.Unmarshal(bytes, car)

	if err != nil {
		return nil, fmt.Errorf(
			"Could not unmarshal world state data to type Car",
		)
	}

	return car, nil
}

// DeleteCar removes a car from the world state.
func (c *CarContract) DeleteCar(
	ctx contractapi.TransactionContextInterface,
	carID string,
) (string, error) {

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()

	if err != nil {
		return "", err
	}

	if clientOrgID != "Org1MSP" {
		return "", fmt.Errorf(
			"User under following MSPID: %v cannot perform this action",
			clientOrgID,
		)
	}

	exists, err := c.CarExists(ctx, carID)

	if err != nil {
		return "", fmt.Errorf(
			"Could not read from world state. %s",
			err,
		)
	}

	if !exists {
		return "", fmt.Errorf(
			"The asset %s does not exist",
			carID,
		)
	}

	err = ctx.GetStub().DelState(carID)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Car with id %v is deleted from the world state.",
		carID,
	), nil
}
func (c *CarContract) GetAllCars(ctx contractapi.TransactionContextInterface) ([]*Car, error) {
    queryString := `{"selector":{"AssetType":"car"},"sort":[{"CarId":"desc"}]}`

    resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
    if err != nil {
        return nil, err
    }

    defer resultsIterator.Close()

    return carResultIteratorFunction(resultsIterator)
}

func carResultIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*Car, error) {
    var cars []*Car

    for resultsIterator.HasNext() {
        queryResult, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }

        var car Car

        err = json.Unmarshal(queryResult.Value, &car)
        if err != nil {
            return nil, err
        }

        cars = append(cars, &car)
    }

    return cars, nil
}
func (c *CarContract) GetCarsByRange(
    ctx contractapi.TransactionContextInterface,
    startKey string,
    endKey string,
) ([]*Car, error) {

    resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
    if err != nil {
        return nil, err
    }

    defer resultsIterator.Close()

    return carResultIteratorFunction(resultsIterator)
}


 type CarPaginationResult struct {
    Cars     []*Car `json:"cars"`
    Bookmark string `json:"bookmark"`
}

func (c *CarContract) GetCarsWithPagination(
    ctx contractapi.TransactionContextInterface,
    pageSize int32,
    bookmark string,
) (*CarPaginationResult, error) {

    queryString := `{"selector":{"AssetType":"car"}}`

    resultsIterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(
        queryString,
        pageSize,
        bookmark,
    )

    if err != nil {
        return nil, err
    }

    defer resultsIterator.Close()

    cars, err := carResultIteratorFunction(resultsIterator)
    if err != nil {
        return nil, err
    }

    return &CarPaginationResult{
        Cars:     cars,
        Bookmark: metadata.Bookmark,
    }, nil
}
