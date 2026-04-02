package handlers

import (
	"database/sql"
	"errors"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasvtf/goledger-challenge-besu/app/blockchain"
	"github.com/lucasvtf/goledger-challenge-besu/app/db"
)

type Handler struct {
	Chain *blockchain.Client
	DB    *sql.DB
}

type setRequest struct {
	Value *int64 `json:"value" binding:"required"`
}

func (h *Handler) Set(c *gin.Context) {
	var req setRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: value is required"})
		return
	}

	txHash, err := h.Chain.SetValue(big.NewInt(*req.Value))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tx_hash": txHash})
}

func (h *Handler) Get(c *gin.Context) {
	value, err := h.Chain.GetValue()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value": value.Int64()})
}

func (h *Handler) Sync(c *gin.Context) {
	value, err := h.Chain.GetValue()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	v := value.Int64()
	if err := db.SaveValue(h.DB, v); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value": v, "synced": true})
}

func (h *Handler) Check(c *gin.Context) {
	blockchainValue, err := h.Chain.GetValue()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dbValue, err := db.GetLatestValue(h.DB)
	if err != nil {
		if errors.Is(err, db.ErrNoSyncedValue) {
			c.JSON(http.StatusOK, gin.H{
				"blockchain_value": blockchainValue.Int64(),
				"db_value":         nil,
				"match":            false,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	match := blockchainValue.Int64() == dbValue
	c.JSON(http.StatusOK, gin.H{
		"blockchain_value": blockchainValue.Int64(),
		"db_value":         dbValue,
		"match":            match,
	})
}
