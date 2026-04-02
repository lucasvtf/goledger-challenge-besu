package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/lucasvtf/goledger-challenge-besu/app/config"
)

const simpleStorageABI = `[{"inputs":[{"internalType":"uint256","name":"x","type":"uint256"}],"name":"set","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"get","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"storedData","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`

var parsedABI abi.ABI

func init() {
	var err error
	parsedABI, err = abi.JSON(strings.NewReader(simpleStorageABI))
	if err != nil {
		panic(fmt.Sprintf("failed to parse ABI: %v", err))
	}
}

// Client holds a persistent connection to the Besu node and serializes
// write transactions via a mutex to prevent nonce conflicts.
type Client struct {
	eth             *ethclient.Client
	cfg             *config.Config
	contractAddress common.Address
	txMu            sync.Mutex
}

func NewClient(cfg *config.Config) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	eth, err := ethclient.DialContext(ctx, cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to node: %w", err)
	}

	return &Client{
		eth:             eth,
		cfg:             cfg,
		contractAddress: common.HexToAddress(cfg.ContractAddress),
	}, nil
}

func (c *Client) Close() {
	c.eth.Close()
}

func (c *Client) SetValue(value *big.Int) (string, error) {
	c.txMu.Lock()
	defer c.txMu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	chainID, err := c.eth.ChainID(ctx)
	if err != nil {
		return "", fmt.Errorf("error querying chain ID: %w", err)
	}

	boundContract := bind.NewBoundContract(c.contractAddress, parsedABI, c.eth, c.eth, c.eth)

	priv, err := crypto.HexToECDSA(c.cfg.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("error loading private key: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(priv, chainID)
	if err != nil {
		return "", fmt.Errorf("error creating transactor: %w", err)
	}

	tx, err := boundContract.Transact(auth, "set", value)
	if err != nil {
		return "", fmt.Errorf("error sending transaction: %w", err)
	}

	waitCtx, waitCancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer waitCancel()

	receipt, err := bind.WaitMined(waitCtx, c.eth, tx)
	if err != nil {
		return "", fmt.Errorf("error waiting for transaction to be mined: %w", err)
	}

	if receipt.Status == 0 {
		return "", fmt.Errorf("transaction reverted: %s", tx.Hash().Hex())
	}

	return tx.Hash().Hex(), nil
}

func (c *Client) GetValue() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	boundContract := bind.NewBoundContract(c.contractAddress, parsedABI, c.eth, c.eth, c.eth)

	caller := bind.CallOpts{
		Pending: false,
		Context: ctx,
	}

	var output []interface{}
	err := boundContract.Call(&caller, &output, "get")
	if err != nil {
		return nil, fmt.Errorf("error calling contract: %w", err)
	}

	if len(output) == 0 {
		return nil, fmt.Errorf("no output from contract call")
	}

	value, ok := output[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected output type: %T", output[0])
	}

	return value, nil
}
