package interfaces

import (
	"math/big"

	"besu-api/internal/database"

	"github.com/ethereum/go-ethereum/core/types"
)

type BlockchainInterface interface {
	GetValue() (*big.Int, error)
	SetValue(value *big.Int) (*types.Transaction, error)
	GetChainID() (*big.Int, error)
}

type DatabaseInterface interface {
	GetValue() (*database.ContractValue, error)
	SetValue(value string) (*database.ContractValue, error)
	Close() error
}