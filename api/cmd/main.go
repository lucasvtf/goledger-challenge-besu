package main

import (
	"log"

	"besu-api/internal/blockchain"
	"besu-api/internal/config"
	"besu-api/internal/database"
	"besu-api/internal/handlers"
	"besu-api/internal/interfaces"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	var bc interfaces.BlockchainInterface
	bcClient, err := blockchain.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize blockchain client: %v", err)
	}
	bc = bcClient

	log.Println("Testing blockchain connection...")
	chainID, err := bc.GetChainID()
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}
	log.Printf("Connected to blockchain with Chain ID: %v", chainID)

	var db interfaces.DatabaseInterface
	dbClient, err := database.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database client: %v", err)
	}
	db = dbClient
	defer db.Close()

	log.Println("Connected to database successfully")

	handler := handlers.NewHandler(bc, db)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	api := r.Group("/api/v1")
	{
		api.GET("/health", handler.HealthHandler)
		api.GET("/value", handler.GetValueHandler)
		api.POST("/value", handler.SetValueHandler)
		api.POST("/sync", handler.SyncHandler)
		api.GET("/check", handler.CheckHandler)
	}

	log.Printf("🚀 Server starting on port %s", cfg.Port)
	log.Printf("📍 Besu RPC: %s", cfg.BesuRPCURL)
	log.Printf("📄 Contract: %s", cfg.ContractAddress)
	log.Printf("💾 Database: %s:%s/%s", cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseName)
	log.Printf("🔗 Endpoints:")
	log.Printf("   GET  /api/v1/health - Health check")
	log.Printf("   GET  /api/v1/value  - Get contract value")
	log.Printf("   POST /api/v1/value  - Set contract value")
	log.Printf("   POST /api/v1/sync   - Sync blockchain to database")
	log.Printf("   GET  /api/v1/check  - Check if values match")

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}