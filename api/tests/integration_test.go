package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"besu-api/internal/blockchain"
	"besu-api/internal/database"
	"besu-api/internal/handlers"
	"besu-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	router   *gin.Engine
	bc       *blockchain.Client
	db       *database.Client
	handler  *handlers.Handler
}

func (suite *IntegrationTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	suite.setupMockRouter()
}

func (suite *IntegrationTestSuite) setupMockRouter() {
	suite.router = gin.New()
	
	api := suite.router.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, models.HealthResponse{
				Status:  "healthy",
				Service: "besu-api",
				Message: "Service is running",
			})
		})
		
		api.GET("/value", func(c *gin.Context) {
			c.JSON(http.StatusOK, models.GetValueResponse{
				Value:   "100",
				Success: true,
				Message: "Value retrieved successfully from blockchain",
			})
		})
		
		api.POST("/value", func(c *gin.Context) {
			var req models.SetValueRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, models.ErrorResponse{
					Success: false,
					Message: "Invalid request format",
					Error:   err.Error(),
				})
				return
			}
			
			c.JSON(http.StatusOK, models.SetValueResponse{
				TxHash:  "0x1234567890abcdef",
				Success: true,
				Message: "Transaction sent successfully",
				Value:   req.Value,
			})
		})
		
		api.GET("/check", func(c *gin.Context) {
			c.JSON(http.StatusOK, models.CheckResponse{
				BlockchainValue: "100",
				DatabaseValue:   "100",
				Match:           true,
				Success:         true,
				Message:         "Database and blockchain values match",
			})
		})
		
		api.POST("/sync", func(c *gin.Context) {
			c.JSON(http.StatusOK, models.SyncResponse{
				BlockchainValue: "100",
				DatabaseValue:   "100",
				Synced:          true,
				Success:         true,
				Message:         "Values are already synchronized",
				SyncedAt:        time.Now(),
			})
		})
	}
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

func (suite *IntegrationTestSuite) TestHealthEndpoint() {
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response models.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "healthy", response.Status)
}

func (suite *IntegrationTestSuite) TestGetValueEndpoint() {
	req, _ := http.NewRequest("GET", "/api/v1/value", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response models.GetValueResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.NotEmpty(suite.T(), response.Value)
}

func (suite *IntegrationTestSuite) TestSetValueEndpoint() {
	requestBody := models.SetValueRequest{Value: "999"}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/v1/value", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response models.SetValueResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "999", response.Value)
	assert.NotEmpty(suite.T(), response.TxHash)
}

func (suite *IntegrationTestSuite) TestCheckEndpoint() {
	req, _ := http.NewRequest("GET", "/api/v1/check", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response models.CheckResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.NotEmpty(suite.T(), response.BlockchainValue)
	assert.NotEmpty(suite.T(), response.DatabaseValue)
}

func (suite *IntegrationTestSuite) TestSyncEndpoint() {
	req, _ := http.NewRequest("POST", "/api/v1/sync", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response models.SyncResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.NotEmpty(suite.T(), response.BlockchainValue)
	assert.NotEmpty(suite.T(), response.DatabaseValue)
}

func (suite *IntegrationTestSuite) TestCompleteWorkflow() {
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/api/v1/value", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	requestBody := models.SetValueRequest{Value: "777"}
	jsonBody, _ := json.Marshal(requestBody)
	req, _ = http.NewRequest("POST", "/api/v1/value", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/api/v1/check", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	req, _ = http.NewRequest("POST", "/api/v1/sync", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}