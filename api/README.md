# 🚀 GOLEDGER CHALLENGE REST API

A robust REST API built with Go that interacts with a **smart contract** deployed on the **Hyperledger Besu blockchain**. The API allows interaction with the smart contract, stores contract values in a PostgreSQL database, and provides endpoints to retrieve, set, and synchronize values between the blockchain and database.

## 🌟 Features

- **🔗 Blockchain Integration**: Direct interaction with Hyperledger Besu nodes
- **📄 Smart Contract Support**: Read and write operations on simple storage smart contract
- **🗄️ Database Synchronization**: PostgreSQL integration with automatic data sync
- **🚀 REST API**: Clean and documented HTTP endpoints with simple-storage routes
- **🔧 Environment Configuration**: Flexible configuration using `.env` files
- **🧪 Comprehensive Testing**: Unit and integration tests included
- **🐳 Docker Support**: PostgreSQL containerization for easy development
- **📊 Value Comparison**: Built-in endpoints to compare blockchain vs database values

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   REST API      │    │   Smart         │    │   PostgreSQL    │
│   (Gin)         │◄──►│   Contract      │    │   Database      │
│                 │    │   (Besu)        │    │                 │
│ • Set Value     │    │ • Store Value   │    │ • Data          │
│ • Get Value     │    │ • Retrieve      │    │   Persistence   │
│ • Sync Data     │    │   Value         │    │ • Sync          │
│ • Check Match   │    │ • Transactions  │    │   Storage       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🗄️ Database

The application uses a **PostgreSQL** database to store the smart contract variable value. The database schema includes a table to store contract values with timestamps for tracking changes.

## 🛠️ Technologies

- **Go**: Programming language for high-performance backend development
- **PostgreSQL**: Relational database for data persistence
- **Gin**: Web framework for building REST APIs
- **Hyperledger Besu**: Enterprise-grade Ethereum blockchain client
- **Docker**: Containerization for development environment
- **Testify**: Testing framework for Go applications

## 📋 Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Make (optional, for using Makefile commands)
- Git

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/JoaoVFerreira/goledger-challenge-besu.git
cd goledger-challenge-besu
```

### 2. Environment Setup

```bash
# Copy environment template
cp .env.example .env

# Edit the configuration file with your values
nano .env
```

Use the `.env.example` file as a base, which contains the following configuration:

```env
CONTRACT_ADDRESS=your_contract_address
BLOCKCHAIN_NODE=http://your_besu_node_url
PRIVATE_KEY=your_private_key
DATABASE_URL=postgres://username:password@localhost:5432/your_db_name
PORT=8080
```

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Start Infrastructure (Optional)

```bash
# Start PostgreSQL database (if using Docker)
docker-compose up -d

# Start Besu blockchain network (if using local setup)
cd besu
./startDev.sh
```

### 5. Run the Application

```bash
# Using Go directly
go run ./cmd/main.go

# Or using Makefile
make run
```

The API will be available at `http://localhost:8080`

## 📡 API Endpoints

### 1. Set Contract Value
```http
POST /simple-storage/set/value
Content-Type: application/json

{
  "value": "2500"
}
```
**Description:** Sets a new value for the smart contract variable. This value is sent to the smart contract on the blockchain.

**Success Response:**
```json
{
  "message": "value set successfully"
}
```

**Error Response:**
```json
{
  "error": "invalid request"
}
```

### 2. Get Contract Value
```http
GET /simple-storage/get/value
```
**Description:** Returns the current value of the smart contract variable directly from the blockchain.

**Success Response:**
```json
{
  "value": 2500
}
```

**Error Response:**
```json
{
  "error": "Failed to retrieve contract value"
}
```

### 3. Sync Value to Database
```http
GET /simple-storage/sync/value
```
**Description:** Synchronizes the smart contract variable value to the SQL database. This operation stores the current smart contract value in the database.

**Success Response:**
```json
{
  "message": "Value synchronized successfully",
  "value": 2500
}
```

**Error Response:**
```json
{
  "error": "Failed to synchronize value with database"
}
```

### 4. Check Value Consistency
```http
GET /simple-storage/check/value
```
**Description:** Compares the value stored in the database with the current smart contract variable value. Returns `true` if they are equal, otherwise returns `false`.

**Success Response:**
```json
{
  "isEqual": true
}
```

**Error Response:**
```json
{
  "error": "Failed to compare values"
}
```

## 🧪 Testing

### Run All Tests
```bash
# Using Makefile
make test

# Using Go directly
go test ./... -v
```

### Test Coverage
```bash
make test-coverage
```

### Integration Tests
The application includes comprehensive integration tests that verify:
- API endpoint functionality
- Blockchain interaction
- Database operations
- Complete workflow scenarios

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `CONTRACT_ADDRESS` | Smart contract address | `0x1234...` |
| `BLOCKCHAIN_NODE` | Besu RPC endpoint | `http://localhost:8545` |
| `PRIVATE_KEY` | Transaction signing key | `0xabc123...` |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://user:pass@localhost:5432/db` |
| `PORT` | API server port | `8080` |

## 🏃‍♂️ Development

### Available Make Commands
```bash
make build       # Build the application
make run         # Run the application
make test        # Run all tests
make clean       # Clean build artifacts
make docker-up   # Start Docker services
make docker-down # Stop Docker services
```

## 🐳 Docker Development

### Start Services
```bash
# Start PostgreSQL
docker-compose up -d postgres

# Start all services
docker-compose up -d
```