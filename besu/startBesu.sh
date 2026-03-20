#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/colors.sh"

# Stop any existing Besu network
echo -e "${YELLOW}Stopping any existing Besu network...${NC}"
if [ -f "$SCRIPT_DIR/stopBesu.sh" ]; then
    bash "$SCRIPT_DIR/stopBesu.sh"
fi

# Detect operating system
echo -e "${YELLOW}Detecting operating system...${NC}"
export OS="$(uname)"

# Check if besu binary is installed and download it if not
echo -e "${YELLOW}Checking if besu binary is installed...${NC}"
if ! [ -x "$(command -v ./bin/besu)" ]; then
    if [ "$OS" = "Darwin" ]; then
        wget -P . https://github.com/hyperledger/besu/releases/download/25.4.1/besu-25.4.1.tar.gz || curl -L -o besu-25.4.1.tar.gz https://github.com/hyperledger/besu/releases/download/25.4.1/besu-25.4.1.tar.gz
    else
        wget -P . https://github.com/hyperledger/besu/releases/download/25.4.1/besu-25.4.1.tar.gz
    fi
    tar --strip-components=1 -xzf besu-25.4.1.tar.gz
    rm besu-25.4.1.tar.gz
fi

# Start the minimal network
bash "$SCRIPT_DIR/minimal/minimalNetwork.sh"
