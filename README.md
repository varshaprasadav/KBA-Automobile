# KBA Automobile Blockchain Application

A Hyperledger Fabric based Automobile application for managing cars, orders, private data, car history, and car registration.

The project uses Hyperledger Fabric, Go chaincode, Fabric Gateway, CouchDB, and a Gin-based web application.

## Project Features

The application supports:

* Create a car
* Read a car
* Check whether a car exists
* Get all cars
* Query cars by range
* Query cars with pagination
* Create private orders
* Read private orders
* Find matching orders
* Match an order with a car
* Register a car
* Get car history
* CouchDB queries and indexes
* Chaincode events
* Manufacturer web dashboard
* Go client application
* Gin REST API

## Technology Stack

* Hyperledger Fabric 2.5.16
* Go
* Fabric Gateway
* Gin Web Framework
* CouchDB
* Docker
* Docker Compose
* HTML
* CSS
* JavaScript
* Ubuntu on WSL 2
* Docker Desktop

## Network Configuration

This project uses the Hyperledger Fabric test network.

### Organizations

* Org1MSP
* Org2MSP

### Channel

```text
mychannel
```

### Peers

```text
Org1 Peer: localhost:7051
Org2 Peer: localhost:9051
```

### Orderer

```text
localhost:7050
```

### Chaincode

```text
KBA-Automobile
```

### Chaincode Contract

```text
CarContract
OrderContract
```

## Project Structure

The project is divided into the Fabric chaincode and the client/web application.

```text
KBA-Automobile/
│
├── Chaincode/
│   ├── car-contract.go
│   ├── order-contract.go
│   ├── collections.json
│   ├── go.mod
│   └── go.sum
│
└── Client/
    ├── client.go
    ├── connect.go
    ├── profile.go
    ├── event.go
    ├── main.go
    ├── go.mod
    ├── go.sum
    │
    ├── templates/
    │   └── index.html
    │
    └── public/
        ├── scripts/
        │   └── index.js
        │
        └── styles/
            └── style.css
```

The Fabric test network is maintained separately in the `fabric-samples` directory.

```text
~/fabric-samples/
├── bin/
├── config/
└── test-network/
```

## Prerequisites

Install the following:

* Ubuntu on WSL 2
* Docker Desktop
* Docker Desktop WSL integration
* Go
* Git
* Hyperledger Fabric binaries and samples

Verify Docker:

```bash
docker --version
```

Verify Go:

```bash
go version
```

Verify Fabric peer:

```bash
peer version
```

The project was developed using Hyperledger Fabric 2.5.16.

## 1. Start the Fabric Test Network

Open Ubuntu and go to the Fabric test network:

```bash
cd ~/fabric-samples/test-network
```

Set the Fabric configuration path and add Fabric binaries to PATH:

```bash
export FABRIC_CFG_PATH=$HOME/fabric-samples/config
export PATH=$PATH:$HOME/fabric-samples/bin
```

Start the network with `mychannel`, Certificate Authorities, and CouchDB:

```bash
./network.sh up createChannel -c mychannel -ca -s couchdb
```

The command starts:

* Orderer
* Org1 peer
* Org2 peer
* CouchDB
* Certificate Authorities

Check the containers:

```bash
docker ps
```

## 2. Verify the Channel

The application uses:

```text
mychannel
```

Do not replace `mychannel` with `autochannel` unless a separate channel named `autochannel` has been created and the chaincode has been deployed to it.

## 3. Prepare the Chaincode

Go to the chaincode directory:

```bash
cd ~/fabric-samples/KBA-Automobile/Chaincode
```

Install the Go dependencies:

```bash
go mod tidy
```

Build the chaincode:

```bash
go build
```

## 4. Deploy KBA-Automobile Chaincode

Return to the test network:

```bash
cd ~/fabric-samples/test-network
```

Deploy the chaincode:

```bash
./network.sh deployCC \
-ccn KBA-Automobile \
-ccp ../KBA-Automobile/Chaincode \
-ccl go \
-c mychannel \
-ccv 8.0 \
-ccs 8 \
-cccg ../KBA-Automobile/Chaincode/collections.json
```

The deployed chaincode uses:

```text
Chaincode name: KBA-Automobile
Channel: mychannel
Version: 8.0
Sequence: 8
```

The exact version and sequence should be increased for future chaincode upgrades.

## 5. Chaincode Contracts

### CarContract

The `CarContract` manages automobile information.

A car contains information such as:

```text
CarId
Make
Model
Color
DateOfManufacture
OwnedBy
Status
```

Example:

```json
{
    "AssetType": "car",
    "CarId": "Car-10",
    "Color": "Red",
    "DateOfManufacture": "27/10/2023",
    "Make": "Tata",
    "Model": "Punch",
    "OwnedBy": "fac01",
    "Status": "In Factory"
}
```

### OrderContract

The `OrderContract` manages private automobile orders.

Order information is stored using Fabric private data.

The private collection is defined in:

```text
Chaincode/collections.json
```

## 6. Private Data Collection

The project uses an `orderCollection` private data collection.

The collection allows Org1 and Org2 to access the private order data.

The collection definition is stored in:

```text
collections.json
```

## 7. Chaincode Events

The `CreateCar` transaction generates a chaincode event.

The event name is:

```text
CreateCar
```

The event contains information such as:

```json
{
    "type": "Car creation",
    "model": "Punch"
}
```

The client application contains an event listener in:

```text
event.go
```

The event listener can be started using:

```go
chaincodeEventListener(
    "org1",
    "mychannel",
    "KBA-Automobile",
)
```

When a new car is created, the listener receives an event similar to:

```text
Received Chaincode Event: CreateCar
Data: {"type":"Car creation","model":"Punch"}
```

## 8. Client Application

The client application is located in:

```text
~/KBA-Automobile-Client
```

The main client files are:

```text
client.go
connect.go
profile.go
event.go
main.go
```

### Client configuration

The client connects to Org1 using:

```text
MSP ID: Org1MSP
Peer: localhost:7051
```

Org2 is also configured:

```text
MSP ID: Org2MSP
Peer: localhost:9051
```

The client uses Fabric Gateway to communicate with the blockchain network.

## 9. Install Client Dependencies

Go to the client directory:

```bash
cd ~/KBA-Automobile-Client
```

Run:

```bash
go mod tidy
```

Build the client:

```bash
go build
```

## 10. Create a Car Using the Client

The client uses the following transaction:

```text
CreateCar
```

Example:

```go
result := submitTxnFn(
    "org1",
    "mychannel",
    "KBA-Automobile",
    "CarContract",
    "invoke",
    make(map[string][]byte),
    "CreateCar",
    "Car-10",
    "Tata",
    "Punch",
    "Red",
    "fac01",
    "27/10/2023",
)
```

Run:

```bash
go run .
```

A successful transaction returns:

```text
*** Transaction submitted successfully: successfully added car Car-10
```

## 11. Query a Car Using the Client

The client can call:

```text
ReadCar
```

Example:

```go
result := submitTxnFn(
    "org1",
    "mychannel",
    "KBA-Automobile",
    "CarContract",
    "query",
    make(map[string][]byte),
    "ReadCar",
    "Car-10",
)
```

## 12. Manufacturer Web Dashboard

The project includes a web application using the Gin framework.

The web application contains:

```text
templates/index.html
public/scripts/index.js
public/styles/style.css
```

The dashboard provides:

### Create Car

The manufacturer can enter:

* Car ID
* Make
* Model
* Color
* Date of Manufacture
* Manufacturer Name

### Query Car

The user can enter a Car ID and retrieve the corresponding car information.

## 13. Start the Web Application

Go to:

```bash
cd ~/KBA-Automobile-Client
```

Make sure the Go dependencies are installed:

```bash
go mod tidy
```

Build the application:

```bash
go build
```

Start the Gin server:

```bash
go run .
```

The server runs on:

```text
http://localhost:8080
```

Open the following address in a browser:

```text
http://localhost:8080
```

The Manufacturer Dashboard should be displayed.

## 14. Web Application API

### Create Car

The browser sends a POST request to:

```text
/api/car
```

Example request:

```json
{
    "carId": "Car-11",
    "make": "Tata",
    "model": "Nexon",
    "color": "White",
    "dateOfManufacture": "2023-10-28",
    "manufacturerName": "Factory-1"
}
```

The Gin server passes this information to:

```text
CreateCar
```

on the Fabric network.

### Query Car

The browser sends a GET request:

```text
/api/car/Car-10
```

The server calls:

```text
ReadCar
```

and returns the result to the browser.

## 15. Environment Variables for Fabric CLI

From:

```bash
cd ~/fabric-samples/test-network
```

set:

```bash
export FABRIC_CFG_PATH=$HOME/fabric-samples/config
export PATH=$PATH:$HOME/fabric-samples/bin
```

Set the Orderer TLS certificate:

```bash
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

Set Org1 peer TLS certificate:

```bash
export ORG1_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
```

Set Org2 peer TLS certificate:

```bash
export ORG2_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
```

Enable TLS:

```bash
export CORE_PEER_TLS_ENABLED=true
```

## 16. Org1 Environment

```bash
export CORE_PEER_LOCALMSPID=Org1MSP
```

```bash
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
```

```bash
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
```

```bash
export CORE_PEER_ADDRESS=localhost:7051
```

## 17. CreateCar Using Peer CLI

Example:

```bash
peer chaincode invoke \
-o localhost:7050 \
--ordererTLSHostnameOverride orderer.example.com \
--tls \
--cafile "$ORDERER_CA" \
-C mychannel \
-n KBA-Automobile \
--peerAddresses localhost:7051 \
--tlsRootCertFiles "$ORG1_PEER_TLSROOTCERT" \
--peerAddresses localhost:9051 \
--tlsRootCertFiles "$ORG2_PEER_TLSROOTCERT" \
-c '{"function":"CreateCar","Args":["Car-11","Tata","Nexon","White","Factory-1","28/10/2023"]}'
```

## 18. Query ReadCar Using Peer CLI

```bash
peer chaincode query \
-C mychannel \
-n KBA-Automobile \
-c '{"function":"ReadCar","Args":["Car-11"]}'
```

## 19. Org2 Environment

```bash
export CORE_PEER_LOCALMSPID=Org2MSP
```

```bash
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
```

```bash
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
```

```bash
export CORE_PEER_ADDRESS=localhost:9051
```

## 20. Create Private Order

Set the order information:

```bash
export MAKE=$(echo -n "Tata" | base64 | tr -d '\n')
```

```bash
export MODEL=$(echo -n "Tiago" | base64 | tr -d '\n')
```

```bash
export COLOR=$(echo -n "Blue" | base64 | tr -d '\n')
```

```bash
export DEALER_NAME=$(echo -n "Popular" | base64 | tr -d '\n')
```

Create the order:

```bash
peer chaincode invoke \
-o localhost:7050 \
--ordererTLSHostnameOverride orderer.example.com \
--tls \
--cafile "$ORDERER_CA" \
-C mychannel \
-n KBA-Automobile \
--peerAddresses localhost:7051 \
--tlsRootCertFiles "$ORG1_PEER_TLSROOTCERT" \
--peerAddresses localhost:9051 \
--tlsRootCertFiles "$ORG2_PEER_TLSROOTCERT" \
-c '{"Args":["OrderContract:CreateOrder","ORD-01"]}' \
--transient "{\"make\":\"$MAKE\",\"model\":\"$MODEL\",\"color\":\"$COLOR\",\"dealerName\":\"$DEALER_NAME\"}"
```

## 21. Read Private Order

```bash
peer chaincode query \
-C mychannel \
-n KBA-Automobile \
-c '{"Args":["OrderContract:ReadOrder","ORD-01"]}'
```

## 22. Matching Orders

To find orders matching a car:

```bash
peer chaincode query \
-C mychannel \
-n KBA-Automobile \
-c '{"Args":["GetMatchingOrders","Car-02"]}'
```

## 23. Match Order

```bash
peer chaincode invoke \
-o localhost:7050 \
--ordererTLSHostnameOverride orderer.example.com \
--tls \
--cafile "$ORDERER_CA" \
-C mychannel \
-n KBA-Automobile \
--peerAddresses localhost:7051 \
--tlsRootCertFiles "$ORG1_PEER_TLSROOTCERT" \
--peerAddresses localhost:9051 \
--tlsRootCertFiles "$ORG2_PEER_TLSROOTCERT" \
-c '{"function":"MatchOrder","Args":["Car-02","ORD-01"]}'
```

## 24. Car Registration

Car registration can be performed using:

```text
RegisterCar
```

Example:

```bash
peer chaincode invoke \
-o localhost:7050 \
--ordererTLSHostnameOverride orderer.example.com \
--tls \
--cafile "$ORDERER_CA" \
-C mychannel \
-n KBA-Automobile \
--peerAddresses localhost:7051 \
--tlsRootCertFiles "$ORG1_PEER_TLSROOTCERT" \
--peerAddresses localhost:9051 \
--tlsRootCertFiles "$ORG2_PEER_TLSROOTCERT" \
-c '{"function":"RegisterCar","Args":["Car-02","Bob","KL-01-7777"]}'
```

## 25. Car History

To retrieve the history of a car:

```bash
peer chaincode query \
-C mychannel \
-n KBA-Automobile \
-c '{"Args":["GetCarHistory","Car-02"]}'
```

## 26. Query All Cars

The chaincode provides a function for retrieving all cars:

```text
GetAllCars
```

Example:

```bash
peer chaincode query \
-C mychannel \
-n KBA-Automobile \
-c '{"Args":["GetAllCars"]}'
```

## 27. Query Cars by Range

The chaincode provides:

```text
GetCarsByRange
```

This can be used to retrieve cars within a specified key range.

## 28. Pagination

The chaincode provides:

```text
GetCarsWithPagination
```

Pagination uses CouchDB query results and a bookmark to retrieve cars in multiple pages.

The result contains:

```json
{
    "cars": [],
    "bookmark": ""
}
```

## 29. CouchDB Index

The project uses CouchDB for rich queries.

CouchDB indexes are included with the chaincode and are used to improve query operations such as retrieving cars by asset type.

## 30. Important Channel Difference

The original KBA Automobile tutorial uses:

```text
autochannel
```

This implementation uses:

```text
mychannel
```

Therefore, commands in this project use:

```text
-C mychannel
```

and the Fabric Gateway client uses:

```go
"mychannel"
```

Do not change these values to `autochannel` unless the network has been created using `autochannel` and the chaincode has been deployed to that channel.

## 31. Important Security Notes

Do not commit the following to GitHub:

* Private keys
* User certificates
* Fabric crypto material
* Wallet files
* `.env` files containing secrets
* API keys
* Passwords
* Generated Docker data
* `node_modules`
* Large generated Fabric directories

The `fabric-samples` directory should normally be installed separately instead of being copied into the application repository.

## 32. Recommended Git Files

A `.gitignore` file should contain entries similar to:

```gitignore
# Go binaries
client
*.exe

# Environment files
.env

# Logs
*.log

# OS files
.DS_Store
Thumbs.db

# Go test files
*.test

# Generated Fabric data
organizations/
system-genesis-block/
channel-artifacts/
log.txt
```

Do not add generated Fabric crypto material to the Git repository.

## 33. Running the Complete Application

### Terminal 1 — Fabric Network

```bash
cd ~/fabric-samples/test-network
```

```bash
export FABRIC_CFG_PATH=$HOME/fabric-samples/config
export PATH=$PATH:$HOME/fabric-samples/bin
```

Check the network:

```bash
docker ps
```

### Terminal 2 — Client/Web Application

```bash
cd ~/KBA-Automobile-Client
```

Install dependencies:

```bash
go mod tidy
```

Build:

```bash
go build
```

Start:

```bash
go run .
```

Open:

```text
http://localhost:8080
```

## 34. Stopping the Network

When you have finished testing, go to:

```bash
cd ~/fabric-samples/test-network
```

To stop the Fabric network:

```bash
./network.sh down
```

The `down` command removes the running test-network containers and related generated network state.

## 35. Development Workflow

The recommended development flow is:

```text
1. Start Fabric network
        ↓
2. Create mychannel
        ↓
3. Deploy KBA-Automobile chaincode
        ↓
4. Start CouchDB
        ↓
5. Run Go client
        ↓
6. Start Gin web application
        ↓
7. Open Manufacturer Dashboard
        ↓
8. Create or query cars
        ↓
9. Fabric processes the transaction
        ↓
10. Ledger is updated
```

For chaincode events:

```text
CreateCar
    ↓
Chaincode SetEvent()
    ↓
Fabric commits transaction
    ↓
Client event listener
    ↓
CreateCar event received
```

## 36. Current Implementation

The application currently contains:

```text
Channel:
mychannel

Organizations:
Org1MSP
Org2MSP

Chaincode:
KBA-Automobile

Contracts:
CarContract
OrderContract

Web Framework:
Gin

Blockchain Client:
Fabric Gateway

Database:
CouchDB

Event:
CreateCar

Web Server:
localhost:8080
```

## 37. Example Successful Chaincode Event

A successful `CreateCar` transaction can produce:

```text
*** Transaction submitted successfully: successfully added car Car-10
```

The event listener receives:

```text
Received Chaincode Event: CreateCar
Data: {"type":"Car creation","model":"Punch"}
```

This confirms that the chaincode transaction and event listener are working together successfully.
