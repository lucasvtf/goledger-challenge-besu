#!/bin/bash

# This script removes all docker containers, docker networks and temporary files related 
# to the Besu network created by startBesu.sh

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/colors.sh"

echo -e "${YELLOW}Cleaning old files...${NC}"
echo
if [ "$OS" = "Darwin" ]; then
    rm -rf tmpFiles/
    rm -rf networkFiles/
    rm -rf genesis/
    rm -rf nodes/
    rm -rf config/qbftConfigFile.json
    rm -f .env.network
    rm -f minimal/bootnodes.*
else
    sudo rm -rf tmpFiles/
    sudo rm -rf networkFiles/
    sudo rm -rf genesis/
    sudo rm -rf nodes/
    sudo rm -rf config/qbftConfigFile.json
    sudo rm -f .env.network
    sudo rm -f minimal/bootnodes.*
fi

echo -e "${YELLOW}Removing all previous besu node containers...${NC}"
docker rm -f $(docker ps -f name=besu. -aq) 2>/dev/null || true
echo

echo -e "${YELLOW}Removing docker besu_test_network...${NC}"
docker network rm besu_test_network 2>/dev/null || true
echo