package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	RPCURL          string
	ContractAddress string
	PrivateKey      string // dev-only default: Alice's pre-funded key for local Besu devnet
	DatabaseURL     string
	Port            string
}

func Load() *Config {
	// Load .env file if present; ignore error if missing
	_ = godotenv.Load()

	return &Config{
		RPCURL:          getEnv("RPC_URL", "http://localhost:8545"),
		ContractAddress: getEnv("CONTRACT_ADDRESS", ""),
		PrivateKey:      strings.TrimPrefix(getEnv("PRIVATE_KEY", "8f2a55949038a9610f50fb23b5883af3b4ecb3c3bb792cbcefbd1542c692be63"), "0x"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/besu?sslmode=disable"),
		Port:            getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
