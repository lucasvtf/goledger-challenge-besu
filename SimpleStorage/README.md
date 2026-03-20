# SimpleStorage — Foundry Project

This Foundry project compiles, tests, and deploys the `SimpleStorage` smart contract to the local Hyperledger Besu network.

## Directory structure

```
SimpleStorage/
├── src/
│   └── SimpleStorage.sol   # The smart contract
├── script/
│   └── SimpleStorage.s.sol # Deployment script
├── test/
│   └── SimpleStorage.t.sol # Test suite
├── lib/
│   └── forge-std/          # Foundry standard library (submodule)
├── .env.example            # Environment variable template
└── foundry.toml            # Foundry configuration
```

## Prerequisites

- [Foundry](https://book.getfoundry.sh/getting-started/installation) (`forge`, `cast`)

## Build

```shell
forge build
```

## Test

```shell
forge test -v
```

## Deploy to the local Besu network

> The local Besu network must be running before deploying. See [`besu/README.md`](../besu/README.md).

Copy the environment template and fill in your deployer private key:

```shell
cp .env.example .env
```

The Besu dev network pre-funds the account below, whose key can be used directly:

```
PRIVATE_KEY=0x8f2a55949038a9610f50fb23b5883af3b4ecb3c3bb792cbcefbd1542c692be63
```

Then deploy:

```shell
source .env && forge script script/SimpleStorage.s.sol:SimpleStorageScript \
    --rpc-url besu \
    --broadcast
```

The deployed contract address is printed to stdout on success.

