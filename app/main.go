package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/lucasvtf/goledger-challenge-besu/app/blockchain"
	"github.com/lucasvtf/goledger-challenge-besu/app/config"
	"github.com/lucasvtf/goledger-challenge-besu/app/db"
	"github.com/lucasvtf/goledger-challenge-besu/app/handlers"
)

func main() {
	cfg := config.Load()

	if cfg.ContractAddress == "" {
		log.Fatal("CONTRACT_ADDRESS environment variable is required")
	}

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	if err := db.Migrate(database); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	chain, err := blockchain.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to blockchain node: %v", err)
	}
	defer chain.Close()

	log.Printf("Connected to database and blockchain node")
	log.Printf("Contract address: %s", cfg.ContractAddress)
	log.Printf("RPC URL: %s", cfg.RPCURL)

	h := &handlers.Handler{
		Chain: chain,
		DB:    database,
	}

	r := gin.Default()
	r.POST("/set", h.Set)
	r.GET("/get", h.Get)
	r.POST("/sync", h.Sync)
	r.GET("/check", h.Check)

	log.Printf("Starting server on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
