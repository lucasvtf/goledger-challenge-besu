package models

import "time"

type GetValueResponse struct {
	Value   string `json:"value"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type SetValueResponse struct {
	TxHash  string `json:"tx_hash"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Value   string `json:"value"`
}

type SetValueRequest struct {
	Value string `json:"value" binding:"required"`
}

type SyncResponse struct {
	BlockchainValue string    `json:"blockchain_value"`
	DatabaseValue   string    `json:"database_value"`
	Synced          bool      `json:"synced"`
	Success         bool      `json:"success"`
	Message         string    `json:"message"`
	SyncedAt        time.Time `json:"synced_at"`
}

type CheckResponse struct {
	BlockchainValue string `json:"blockchain_value"`
	DatabaseValue   string `json:"database_value"`
	Match           bool   `json:"match"`
	Success         bool   `json:"success"`
	Message         string `json:"message"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Message string `json:"message"`
}