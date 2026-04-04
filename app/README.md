# GoLedger Challenge Besu — Go Application

A Go REST API that interacts with a Hyperledger Besu blockchain network and a PostgreSQL database. The application reads and writes values to the `SimpleStorage` smart contract and synchronizes on-chain state to an SQL database.

## Prerequisites

- [Go](https://golang.org/dl/) 1.21+
- [Docker](https://www.docker.com/) (for Besu network and PostgreSQL)
- Running Besu network with deployed `SimpleStorage` contract (see [root README](../README.md))

## Quick Start

### 1. Start the Besu network and deploy the contract

From the project root:

```bash
make devnet-deploy
```

Note the contract address printed to stdout.

### 2. Start PostgreSQL

From the project root:

```bash
docker-compose up -d
```

This starts a PostgreSQL 16 container with the correct database and credentials.

### 3. Run the application

```bash
cd app
export CONTRACT_ADDRESS=<contract address from step 1>
go run .
```

The server starts on port 8080 by default.

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `CONTRACT_ADDRESS` | *(required)* | Deployed SimpleStorage contract address |
| `RPC_URL` | `http://localhost:8545` | Besu JSON-RPC endpoint |
| `PRIVATE_KEY` | Alice's pre-funded key | Private key for signing transactions (no `0x` prefix) |
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/besu?sslmode=disable` | PostgreSQL connection string |
| `PORT` | `8080` | HTTP server port |

## API Endpoints

### POST /set

Set a new value on the smart contract.

```bash
curl -X POST http://localhost:8080/set \
  -H "Content-Type: application/json" \
  -d '{"value": 42}'
```

Response:
```json
{"tx_hash": "0x..."}
```

### GET /get

Read the current value from the blockchain.

```bash
curl http://localhost:8080/get
```

Response:
```json
{"value": 42}
```

### POST /sync

Read the value from the blockchain and save it to the database.

```bash
curl -X POST http://localhost:8080/sync
```

Response:
```json
{"value": 42, "synced": true}
```

### GET /check

Compare the database value with the blockchain value.

```bash
curl http://localhost:8080/check
```

Response:
```json
{"blockchain_value": 42, "db_value": 42, "match": true}
```

If no sync has been performed yet, `db_value` will be `null` and `match` will be `false`.

## Architecture

```
app/
├── main.go              # Entry point: config, DB, router, server
├── config/config.go     # Environment-driven configuration
├── blockchain/client.go # Besu interaction via go-ethereum
├── db/db.go             # PostgreSQL connection and queries
└── handlers/handlers.go # Gin HTTP handlers
```

### Data Flow

1. **SET**: Client sends value via REST API -> Go app sends `set(uint256)` transaction to Besu -> waits for mining confirmation -> returns tx hash
2. **GET**: Client requests value -> Go app calls `get()` on the contract (read-only) -> returns value
3. **SYNC**: Go app calls `get()` on the contract -> inserts the value into `contract_state` table in PostgreSQL
4. **CHECK**: Go app reads both blockchain (`get()`) and database (latest row in `contract_state`) -> compares and returns match status

### Blockchain Interaction

The application uses the [go-ethereum](https://github.com/ethereum/go-ethereum) library to interact with the Besu network:

- **Write operations** (`SetValue`): Uses `bind.NewBoundContract` + `bind.NewKeyedTransactorWithChainID` to sign and send transactions, then `bind.WaitMined` to wait for confirmation.
- **Read operations** (`GetValue`): Uses `bind.NewBoundContract` + `boundContract.Call` with `bind.CallOpts` for gas-free view calls.

The contract ABI is embedded as a constant string in `blockchain/client.go`.

### Database

PostgreSQL stores the synchronized contract state in a `contract_state` table:

```sql
CREATE TABLE contract_state (
    id SERIAL PRIMARY KEY,
    value BIGINT NOT NULL,
    synced_at TIMESTAMP DEFAULT NOW()
);
```

The table is auto-created on application startup via the `Migrate()` function.

## Stopping the Environment

```bash
# Stop PostgreSQL (from project root)
docker-compose down

# Stop PostgreSQL and remove data
docker-compose down -v

# Stop Besu network
make stop-devnet
```
