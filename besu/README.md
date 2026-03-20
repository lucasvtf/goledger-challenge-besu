# Hyperledger Besu Test Network

The `startBesu.sh` script sets up and starts a local 4-node Hyperledger Besu network using the QBFT consensus mechanism. It stops any previously running network, downloads the Besu binary if needed, and delegates to `minimal/minimalNetwork.sh` to generate keys, the genesis file, and start all containers.

> **⚠️ Important Note**: If the Besu binary is not found at `./bin/besu`, the script automatically downloads it from the official Hyperledger Besu repository. An internet connection is required.

## Usage

```bash
./startBesu.sh
```

The script will:
1. Stop any existing Besu network (via `stopBesu.sh`)
2. Download the Besu binary if not already present
3. Generate blockchain config, keys, and genesis file
4. Start the Docker network `besu_test_network`
5. Start 4 Besu nodes, waiting for each to be responsive before starting the next

Node RPC ports:

| Node | RPC Port | P2P Port |
|------|----------|----------|
| besu.node-1 | 8545 | 30303 |
| besu.node-2 | 8547 | 30304 |
| besu.node-3 | 8549 | 30305 |
| besu.node-4 | 8551 | 30306 |

---

The `stopBesu.sh` script stops and cleans up the entire Besu test network.

## Usage

```bash
./stopBesu.sh
```

This script will:
- Remove all temporary files (`genesis/`, `nodes/`, `minimal/bootnodes.*`)
- Stop and remove all Besu Docker containers
- Remove the `besu_test_network` Docker network

---

## External Dependencies

* **Java JDK** — required to run the Besu binary (`./bin/besu`):
  * **macOS**: Java 21+. Install via Homebrew: `brew install openjdk@21`, or download manually from [Oracle](https://www.oracle.com/java/technologies/downloads/).
  * **Linux/Unix**: Java 17+. Download from [Oracle](https://www.oracle.com/java/technologies/downloads/) or install via your package manager.
  * See the [official Besu installation guide](https://besu.hyperledger.org/private-networks/get-started/install/binary-distribution) for full details.
* [Docker](https://docs.docker.com/engine/)
* [jq](https://jqlang.org/download/)
* [wget](https://www.gnu.org/software/wget/) / [curl](https://curl.se/download.html) (used to download the Besu binary; `curl` is the fallback on macOS)