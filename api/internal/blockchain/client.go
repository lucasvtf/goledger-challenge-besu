package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"besu-api/internal/config"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	client     *ethclient.Client
	contract   common.Address
	abi        abi.ABI
	privateKey *ecdsa.PrivateKey
	auth       *bind.TransactOpts
}

func NewClient(cfg *config.Config) (*Client, error) {
	client, err := ethclient.Dial(cfg.BesuRPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Besu: %v", err)
	}

	contractABI, err := abi.JSON(strings.NewReader(config.ContractABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(cfg.PrivateKey, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %v", err)
	}

	
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %v", err)
	}

	auth.GasLimit = uint64(300000)
	auth.GasPrice = big.NewInt(0)  

	return &Client{
		client:     client,
		contract:   common.HexToAddress(cfg.ContractAddress),
		abi:        contractABI,
		privateKey: privateKey,
		auth:       auth,
	}, nil
}

func (c *Client) GetChainID() (*big.Int, error) {
	return c.client.ChainID(context.Background())
}

func (c *Client) SetValue(value *big.Int) (*types.Transaction, error) {
	data, err := c.abi.Pack("set", value)
	if err != nil {
		return nil, fmt.Errorf("failed to pack set method: %v", err)
	}

	fromAddress := crypto.PubkeyToAddress(c.privateKey.PublicKey)
	nonce, err := c.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	gasLimit := uint64(300000)
	gasPrice := big.NewInt(0)
	
	tx := types.NewTransaction(nonce, c.contract, big.NewInt(0), gasLimit, gasPrice, data)

	chainID, err := c.client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), c.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}

	err = c.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %v", err)
	}

	return signedTx, nil
}

func (c *Client) GetValue() (*big.Int, error) {
	data, err := c.abi.Pack("get")
	if err != nil {
		return nil, fmt.Errorf("failed to pack get method: %v", err)
	}

	call := ethereum.CallMsg{
		To:   &c.contract,
		Data: data,
	}

	result, err := c.client.CallContract(context.Background(), call, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %v", err)
	}

	var value *big.Int
	err = c.abi.UnpackIntoInterface(&value, "get", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack result: %v", err)
	}

	return value, nil
}