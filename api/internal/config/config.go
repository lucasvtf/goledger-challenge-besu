package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BesuRPCURL      string
	ContractAddress string
	PrivateKey      string
	Port            string
	DatabaseURL     string
	DatabaseHost    string
	DatabasePort    string
	DatabaseUser    string
	DatabasePass    string
	DatabaseName    string
}

const ContractABI = `[
	{
		"inputs": [],
		"name": "get",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "_value",
				"type": "uint256"
			}
		],
		"name": "set",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]`

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}
	
	config := &Config{
		BesuRPCURL:      getEnv("BESU_RPC_URL"),
		ContractAddress: getEnv("CONTRACT_ADDRESS"),
		PrivateKey:      getEnv("PRIVATE_KEY"),
		Port:            getEnv("PORT"),
		DatabaseURL:     getEnv("DATABASE_URL"),
		DatabaseHost:    getEnv("DB_HOST"),
		DatabasePort:    getEnv("DB_PORT"),
		DatabaseUser:    getEnv("DB_USER"),
		DatabasePass:    getEnv("DB_PASS"),
		DatabaseName:    getEnv("DB_NAME"),
	}
	
	return config
}

func getEnv(key string) string {
	return os.Getenv(key)
}