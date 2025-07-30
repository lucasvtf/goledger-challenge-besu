package handlers

import (
	"bytes"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"besu-api/internal/database"
	"besu-api/internal/models"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBlockchainClient struct {
	mock.Mock
}

type MockDatabaseClient struct {
	mock.Mock
}

func (m *MockBlockchainClient) GetValue() (*big.Int, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *MockBlockchainClient) SetValue(value *big.Int) (*types.Transaction, error) {
	args := m.Called(value)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockBlockchainClient) GetChainID() (*big.Int, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *MockDatabaseClient) GetValue() (*database.ContractValue, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.ContractValue), args.Error(1)
}

func (m *MockDatabaseClient) SetValue(value string) (*database.ContractValue, error) {
	args := m.Called(value)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.ContractValue), args.Error(1)
}

func (m *MockDatabaseClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestHealthHandler(t *testing.T) {
	mockBC := new(MockBlockchainClient)
	mockDB := new(MockDatabaseClient)
	handler := NewHandler(mockBC, mockDB)
	
	router := setupRouter()
	router.GET("/health", handler.HealthHandler)

	// Test
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response.Status)
	assert.Equal(t, "besu-api", response.Service)
}

func TestGetValueHandler_Success(t *testing.T) {
	mockBC := new(MockBlockchainClient)
	mockDB := new(MockDatabaseClient)
	handler := NewHandler(mockBC, mockDB)
	
	expectedValue := big.NewInt(123)
	mockBC.On("GetValue").Return(expectedValue, nil)
	
	router := setupRouter()
	router.GET("/value", handler.GetValueHandler)

	// Test
	req, _ := http.NewRequest("GET", "/value", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.GetValueResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "123", response.Value)
	assert.True(t, response.Success)
	
	mockBC.AssertExpectations(t)
}

func TestSetValueHandler_Success(t *testing.T) {
	mockBC := new(MockBlockchainClient)
	mockDB := new(MockDatabaseClient)
	handler := NewHandler(mockBC, mockDB)
	
	mockTx := types.NewTransaction(0, [20]byte{}, big.NewInt(0), 21000, big.NewInt(0), nil)
	expectedValue := big.NewInt(456)
	mockBC.On("SetValue", expectedValue).Return(mockTx, nil)
	
	router := setupRouter()
	router.POST("/value", handler.SetValueHandler)

	requestBody := models.SetValueRequest{Value: "456"}
	jsonBody, _ := json.Marshal(requestBody)

	// Test
	req, _ := http.NewRequest("POST", "/value", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.SetValueResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "456", response.Value)
	assert.True(t, response.Success)
	assert.NotEmpty(t, response.TxHash)
	
	mockBC.AssertExpectations(t)
}

func TestSetValueHandler_InvalidValue(t *testing.T) {
	mockBC := new(MockBlockchainClient)
	mockDB := new(MockDatabaseClient)
	handler := NewHandler(mockBC, mockDB)
	
	router := setupRouter()
	router.POST("/value", handler.SetValueHandler)

	requestBody := models.SetValueRequest{Value: "invalid"}
	jsonBody, _ := json.Marshal(requestBody)

	// Test
	req, _ := http.NewRequest("POST", "/value", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Message, "Invalid numeric value")
}

func TestCheckHandler_ValuesMatch(t *testing.T) {
	mockBC := new(MockBlockchainClient)
	mockDB := new(MockDatabaseClient)
	handler := NewHandler(mockBC, mockDB)
	
	blockchainValue := big.NewInt(789)
	dbValue := &database.ContractValue{
		ID:        1,
		Value:     "789",
		UpdatedAt: time.Now(),
	}
	
	mockBC.On("GetValue").Return(blockchainValue, nil)
	mockDB.On("GetValue").Return(dbValue, nil)
	
	router := setupRouter()
	router.GET("/check", handler.CheckHandler)

	// Test
	req, _ := http.NewRequest("GET", "/check", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.CheckResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "789", response.BlockchainValue)
	assert.Equal(t, "789", response.DatabaseValue)
	assert.True(t, response.Match)
	assert.True(t, response.Success)
	
	mockBC.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestCheckHandler_ValuesDifferent(t *testing.T) {
	mockBC := new(MockBlockchainClient)
	mockDB := new(MockDatabaseClient)
	handler := NewHandler(mockBC, mockDB)
	
	blockchainValue := big.NewInt(999)
	dbValue := &database.ContractValue{
		ID:        1,
		Value:     "777",
		UpdatedAt: time.Now(),
	}
	
	mockBC.On("GetValue").Return(blockchainValue, nil)
	mockDB.On("GetValue").Return(dbValue, nil)
	
	router := setupRouter()
	router.GET("/check", handler.CheckHandler)

	// Test
	req, _ := http.NewRequest("GET", "/check", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.CheckResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "999", response.BlockchainValue)
	assert.Equal(t, "777", response.DatabaseValue)
	assert.False(t, response.Match)
	assert.True(t, response.Success)
	
	mockBC.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}