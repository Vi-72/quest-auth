//go:build integration

package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"quest-auth/cmd"
	"quest-auth/tests/integration/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthHTTP_Register_Success(t *testing.T) {
	// Arrange
	testDB := utils.NewTestDatabase(t)
	defer testDB.Cleanup(t)

	router := setupTestRouter(t, testDB)

	registerRequest := map[string]interface{}{
		"email":    "test@example.com",
		"phone":    "+1234567890",
		"name":     "John Doe",
		"password": "securepassword123",
	}

	requestBody, _ := json.Marshal(registerRequest)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusCreated, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.NotEmpty(t, response["UserID"])
	assert.Equal(t, "test@example.com", response["Email"])
	assert.Equal(t, "+1234567890", response["Phone"])
	assert.Equal(t, "John Doe", response["Name"])
}

func TestAuthHTTP_Register_InvalidEmail_ReturnsBadRequest(t *testing.T) {
	// Arrange
	testDB := utils.NewTestDatabase(t)
	defer testDB.Cleanup(t)

	router := setupTestRouter(t, testDB)

	registerRequest := map[string]interface{}{
		"email":    "invalid-email",
		"phone":    "+1234567890",
		"name":     "John Doe",
		"password": "securepassword123",
	}

	requestBody, _ := json.Marshal(registerRequest)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "bad-request", response["type"])
	assert.Contains(t, response["detail"], "email")
}

func TestAuthHTTP_Register_DuplicateEmail_ReturnsBadRequest(t *testing.T) {
	// Arrange
	testDB := utils.NewTestDatabase(t)
	defer testDB.Cleanup(t)

	router := setupTestRouter(t, testDB)

	// Сначала регистрируем пользователя
	registerRequest := map[string]interface{}{
		"email":    "test@example.com",
		"phone":    "+1234567890",
		"name":     "John Doe",
		"password": "securepassword123",
	}

	requestBody, _ := json.Marshal(registerRequest)
	req1 := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(requestBody))
	req1.Header.Set("Content-Type", "application/json")

	recorder1 := httptest.NewRecorder()
	router.ServeHTTP(recorder1, req1)
	require.Equal(t, http.StatusCreated, recorder1.Code)

	// Пытаемся зарегистрировать с тем же email
	registerRequest2 := map[string]interface{}{
		"email":    "test@example.com", // Тот же email
		"phone":    "+9876543210",      // Другой телефон
		"name":     "Jane Smith",
		"password": "anotherpassword123",
	}

	requestBody2, _ := json.Marshal(registerRequest2)
	req2 := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(requestBody2))
	req2.Header.Set("Content-Type", "application/json")

	// Act
	recorder2 := httptest.NewRecorder()
	router.ServeHTTP(recorder2, req2)

	// Assert
	assert.Equal(t, http.StatusBadRequest, recorder2.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder2.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "bad-request", response["type"])
	assert.Contains(t, response["detail"], "email already exists")
}

func TestAuthHTTP_Login_Success(t *testing.T) {
	// Arrange
	testDB := utils.NewTestDatabase(t)
	defer testDB.Cleanup(t)

	router := setupTestRouter(t, testDB)

	// Сначала регистрируем пользователя
	registerUser(t, router, "test@example.com", "+1234567890", "John Doe", "securepassword123")

	loginRequest := map[string]interface{}{
		"email":    "test@example.com",
		"password": "securepassword123",
	}

	requestBody, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.NotEmpty(t, response["UserID"])
	assert.Equal(t, "test@example.com", response["Email"])
	assert.Equal(t, "+1234567890", response["Phone"])
	assert.Equal(t, "John Doe", response["Name"])
}

func TestAuthHTTP_Login_InvalidCredentials_ReturnsUnauthorized(t *testing.T) {
	// Arrange
	testDB := utils.NewTestDatabase(t)
	defer testDB.Cleanup(t)

	router := setupTestRouter(t, testDB)

	// Регистрируем пользователя
	registerUser(t, router, "test@example.com", "+1234567890", "John Doe", "securepassword123")

	loginRequest := map[string]interface{}{
		"email":    "test@example.com",
		"password": "wrongpassword",
	}

	requestBody, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "unauthorized", response["type"])
	assert.Contains(t, response["detail"], "invalid credentials")
}

func TestAuthHTTP_HealthCheck_Success(t *testing.T) {
	// Arrange
	testDB := utils.NewTestDatabase(t)
	defer testDB.Cleanup(t)

	router := setupTestRouter(t, testDB)

	req := httptest.NewRequest("GET", "/health", nil)

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "ok", response["status"])
}

// Helper functions

func setupTestRouter(t *testing.T, testDB *utils.TestDatabase) http.Handler {
	configs := cmd.Config{
		HttpPort:            "8080",
		DbHost:              "localhost",
		DbPort:              "5432",
		DbUser:              "username",
		DbPassword:          "secret",
		DbName:              testDB.Name,
		DbSslMode:           "disable",
		EventGoroutineLimit: 1,
	}

	compositionRoot := cmd.NewCompositionRoot(configs, testDB.DB)
	router := cmd.NewRouter(compositionRoot)
	return router
}

func registerUser(t *testing.T, router http.Handler, email, phone, name, password string) string {
	registerRequest := map[string]interface{}{
		"email":    email,
		"phone":    phone,
		"name":     name,
		"password": password,
	}

	requestBody, _ := json.Marshal(registerRequest)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusCreated, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	return response["UserID"].(string)
}
