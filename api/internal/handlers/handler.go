package handlers

import (
	"log"
	"math/big"
	"net/http"

	"besu-api/internal/interfaces"
	"besu-api/internal/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	blockchain interfaces.BlockchainInterface
	database   interfaces.DatabaseInterface
}

func NewHandler(bc interfaces.BlockchainInterface, db interfaces.DatabaseInterface) *Handler {
	return &Handler{
		blockchain: bc,
		database:   db,
	}
}

func (h *Handler) HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, models.HealthResponse{
		Status:  "healthy",
		Service: "besu-api",
		Message: "Service is running",
	})
}

func (h *Handler) GetValueHandler(c *gin.Context) {
	value, err := h.blockchain.GetValue()
	if err != nil {
		log.Printf("Error getting value from blockchain: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve value from blockchain",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.GetValueResponse{
		Value:   value.String(),
		Success: true,
		Message: "Value retrieved successfully from blockchain",
	})
}

func (h *Handler) SetValueHandler(c *gin.Context) {
	var req models.SetValueRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	value, ok := new(big.Int).SetString(req.Value, 10)
	if !ok {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid numeric value",
			Error:   "Failed to parse value as number",
		})
		return
	}

	tx, err := h.blockchain.SetValue(value)
	if err != nil {
		log.Printf("Error setting value on blockchain: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to set value on blockchain",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SetValueResponse{
		TxHash:  tx.Hash().Hex(),
		Success: true,
		Message: "Transaction sent successfully",
		Value:   req.Value,
	})
}

func (h *Handler) CheckHandler(c *gin.Context) {
	blockchainValue, err := h.blockchain.GetValue()
	if err != nil {
		log.Printf("Error getting value from blockchain: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve value from blockchain",
			Error:   err.Error(),
		})
		return
	}

	dbValue, err := h.database.GetValue()
	if err != nil {
		log.Printf("Error getting value from database: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve value from database",
			Error:   err.Error(),
		})
		return
	}

	blockchainValueStr := blockchainValue.String()
	
	isMatch := dbValue.Value == blockchainValueStr
	
	var message string
	if isMatch {
		message = "Database and blockchain values match"
		log.Printf("CHECK: Values match - blockchain=%s, database=%s", 
			blockchainValueStr, dbValue.Value)
	} else {
		message = "Database and blockchain values do not match"
		log.Printf("CHECK: Values differ - blockchain=%s, database=%s", 
			blockchainValueStr, dbValue.Value)
	}

	c.JSON(http.StatusOK, models.CheckResponse{
		BlockchainValue: blockchainValueStr,
		DatabaseValue:   dbValue.Value,
		Match:           isMatch,
		Success:         true,
		Message:         message,
	})
}

func (h *Handler) SyncHandler(c *gin.Context) {
	blockchainValue, err := h.blockchain.GetValue()
	if err != nil {
		log.Printf("Error getting value from blockchain: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve value from blockchain",
			Error:   err.Error(),
		})
		return
	}

	dbValue, err := h.database.GetValue()
	if err != nil {
		log.Printf("Error getting value from database: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve value from database",
			Error:   err.Error(),
		})
		return
	}

	blockchainValueStr := blockchainValue.String()
	
	if dbValue.Value == blockchainValueStr {
		log.Printf("Values already synchronized: blockchain=%s, database=%s", 
			blockchainValueStr, dbValue.Value)
		
		c.JSON(http.StatusOK, models.SyncResponse{
			BlockchainValue: blockchainValueStr,
			DatabaseValue:   dbValue.Value,
			Synced:          true,
			Success:         true,
			Message:         "Values are already synchronized",
			SyncedAt:        dbValue.UpdatedAt,
		})
		return
	}

	updatedValue, err := h.database.SetValue(blockchainValueStr)
	if err != nil {
		log.Printf("Error updating value in database: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to update value in database",
			Error:   err.Error(),
		})
		return
	}

	log.Printf("Synchronized value: blockchain=%s, database=%s -> %s", 
		blockchainValueStr, dbValue.Value, updatedValue.Value)

	c.JSON(http.StatusOK, models.SyncResponse{
		BlockchainValue: blockchainValueStr,
		DatabaseValue:   updatedValue.Value,
		Synced:          true,
		Success:         true,
		Message:         "Value synchronized successfully from blockchain to database",
		SyncedAt:        updatedValue.UpdatedAt,
	})
}