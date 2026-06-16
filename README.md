# Hyperledger Fabric Go Automobile Gateway

A decentralized application (DApp) implementation showcasing a manufacturer dashboard connected to a Hyperledger Fabric network.

---

## 1. Project Overview

This repository contains a **Sample Client Application** (`sample_app`) that provides a web-based Manufacturer Dashboard. Through this dashboard, manufacturers can:
*   **Register Automobile Assets**: Commit new vehicles to the blockchain ledger with key parameters (Asset ID, Make, Model, Color, Manufacture Date, and Manufacturer Name).
*   **Query Ledger State**: Retrieve cryptographically verified vehicle details from the distributed ledger in real-time.

The application contains two main components:
1.  **Frontend Dashboard**: A premium, responsive dark-themed user interface built with HTML, CSS, and Vanilla JavaScript featuring dynamic DOM rendering and animated toast notifications.
2.  **Go Backend Service**: A REST API built on the Gin framework, responsible for handling web requests and interacting with the Fabric Gateway peer node.

---

## 2. Directory Structure

The expected project directory layout under the workspace `KBA-CHF` is as follows:

```text
KBA-CHF/
├── fabric-samples/          # Official Hyperledger Fabric samples & test network
└── KBA-Automobile/          # Core Automobile blockchain project
    ├── Chaincode/           # Smart contract code in Go/Java
    ├── Client/              # Core transaction & connection scripts
    ├── Events/              # Event listeners and handlers
```

---

## 3. Project Setup & Prerequisites

### Step 1: Create the Workspace Directory
Create a workspace directory named `KBA-CHF` to house both the Hyperledger Fabric test network assets and the automobile application:
```bash
mkdir KBA-CHF
cd KBA-CHF
```

### Step 2: Download Required Project Folders
Download the following folders and extract them directly inside the `KBA-CHF` directory:

*   **`fabric-samples`**: Includes binary files, Docker configuration files, and the test-network scripts.
    *   *Download Link*: [Google Drive Folder](https://drive.google.com/drive/folders/18EcPdneH4po4XURpfOHNTgivRrfPUHxK?usp=sharing)
*   **`KBA-Automobile`**: Contains the chaincode, transaction endpoints, and event code.
    *   *Download Link*: [Google Drive Folder](https://drive.google.com/drive/folders/1w8eH57tbF8S2Lwn0og5yF14Aduzs_vSU?usp=sharing)

Verify that both directories are successfully located in the same parent directory before proceeding.

---

## 4. Hyperledger Fabric Network Lifecycle

### Step 1: Bring Up the Test Network
Navigate to the test-network directory and spin up the Fabric network with CA support and CouchDB:
```bash
cd fabric-samples/test-network

./network.sh up createChannel \
  -c autochannel \
  -ca \
  -s couchdb
```

#### Startup Parameter Guide
| Parameter | Description |
| :--- | :--- |
| `up` | Starts the Hyperledger Fabric nodes (orderers and peers). |
| `createChannel` | Spins up the network, creates a channel, and joins the peers to it. |
| `-c autochannel` | Creates a channel named `autochannel`. |
| `-ca` | Configures and runs Certificate Authorities (CA) for identity generation. |
| `-s couchdb` | Deploys CouchDB instances as state databases instead of LevelDB. |

#### Verify Container Status
Confirm all core network containers are running:
```bash
docker ps
```
You should see active containers for the **Orderer**, **Peer0 Org1**, **Peer0 Org2**, **CouchDB instances**, and **Certificate Authorities**.

---

### Step 2: Add Organization 3 (Org3)
Add a third organization (Org3) to the running channel:
```bash
cd addOrg3

./addOrg3.sh up \
  -c autochannel \
  -ca \
  -s couchdb

cd ..
```

#### Org3 Parameter Guide
| Parameter | Description |
| :--- | :--- |
| `up` | Dynamically appends Organization 3 components to the running network. |
| `-c autochannel` | Targets the active channel `autochannel`. |
| `-ca` | Configures and launches Certificate Authorities for Org3. |
| `-s couchdb` | Integrates CouchDB instances for Org3 state database. |

#### Verify Extended Container Status
Verify that Org3 peers have successfully joined:
```bash
docker ps -a
```
You should see additional peer and CouchDB containers associated with Organization 3.

---

### Step 3: Deploy the Smart Contract (Chaincode)
Deploy the `KBA-Automobile` smart contract to the channel using private data configurations:
```bash
./network.sh deployCC \
  -ccn KBA-Automobile \
  -ccp ../../KBA-Automobile/Chaincode/ \
  -ccl go \
  -c autochannel \
  -cccg ../../KBA-Automobile/Chaincode/collections.json
```

#### Chaincode Deployment Parameter Guide
| Parameter | Description |
| :--- | :--- |
| `-ccn KBA-Automobile` | Specifies the deployed chaincode identifier (`KBA-Automobile`). |
| `-ccp <path>` | Sets the path to the smart contract source directory. |
| `-ccl go` | Sets the smart contract programming language to Go. |
| `-c autochannel` | Installs and commits the contract on channel `autochannel`. |
| `-cccg collections.json` | Passes the Private Data Collection definition file. |

#### Verify Chaincode Deployment
Verify that the chaincode container is running:
```bash
docker ps -a
```
Look for container names matching:
`dev-peer0.org1.example.com-KBA-Automobile_1.0-<hash>`

You can also check query commit records using:
```bash
peer lifecycle chaincode querycommitted \
  -C autochannel \
  -n KBA-Automobile
```

---

## 5. Setting Up the Web Application

To execute the manufacturer client backend and UI dashboard, follow these steps.

### Step 1: Copy Core Client Files
Create a new directory named `sample_app` (representing your `SAMPLE-APP`) inside the `KBA-Automobile/` project:
```text
KBA-Automobile/
├── Chaincode/
├── Client/
├── Events/
└── sample_app/
```
Copy the following files from your `Client/` module into the `sample_app/` folder:
*   `client.go`
*   `profile.go`
*   `connect.go`

Also add your backend routing `main.go` and static web assets directory `public/` into the `sample_app/` folder.

### Step 2: Initialize & Configure Go Project
Navigate into the `sample_app/` directory and initialize the Go modules:
```bash
cd sample_app
go mod init sampleapp
```

Fetch and install the required dependencies (such as the Gin web framework and Fabric Gateway client libraries):
```bash
go mod tidy
```

---

## 6. Running the Application

### Option A: Local Sandbox Mode (No Fabric Network Required)
For frontend development, interface testing, or API verification, you can run the application in **Mock Mode**. This operates a local in-memory ledger database simulating blockchain interactions instantly without needing Docker or a live Fabric peer connection:

```bash
# Enable Mock Mode
export MOCK_FABRIC=true

# Start the application
go run main.go connect.go client.go profile.go
```

Open your browser and navigate to **`http://localhost:8080`** to access the dashboard.

### Option B: Production Gateway Mode (Live Fabric Connection)
To connect the manufacturer dashboard to your live, running Hyperledger Fabric test network, configure your cryptographic MSP details via environment variables and run:

```bash
# Disable Mock Mode
export MOCK_FABRIC=false

# Connection parameters
export FABRIC_MSP_ID="Org1MSP"
export FABRIC_PEER_ENDPOINT="localhost:7051"
export FABRIC_PEER_HOST_OVERRIDE="peer0.org1.example.com"

# Paths to organization cryptographic material
export FABRIC_CERT_PATH="/path/to/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem"
export FABRIC_KEY_PATH="/path/to/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore"
export FABRIC_TLS_CERT_PATH="/path/to/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"

# Start the application
go run main.go connect.go client.go profile.go
```

---

## 7. Tearing Down the Network

Once development, debugging, or testing activities are complete, stop and remove all Hyperledger Fabric containers, volumes, and temporary network assets:

```bash
cd fabric-samples/test-network
./network.sh down
```

#### Teardown Operation Summary
*   Stops and deletes all peer, orderer, CA, and CouchDB containers.
*   Deletes generated certificates, cryptographic keys, and channel artifacts.
*   Removes local Docker chaincode containers and associated volumes.

#### Verify Cleanup
Confirm that no Fabric containers remain:
```bash
docker ps -a
```
The console output should be clear of Fabric network components.