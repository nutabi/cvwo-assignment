package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nutabi/cvwo-assignment/backend/internal/service"
	"github.com/stretchr/testify/assert"
)

// TestUserRegistration tests the user registration endpoint
func TestUserRegistration(t *testing.T) {
	router := setupTestRouter(t)

	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful registration",
			requestBody: map[string]string{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "SecurePass123!",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response service.UserProfile
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "testuser", response.Username)
				assert.NotNil(t, response.Email)
				assert.Equal(t, "test@example.com", *response.Email)
			},
		},
		{
			name: "invalid username",
			requestBody: map[string]string{
				"username": "ab",
				"email":    "test@example.com",
				"password": "SecurePass123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse:  nil,
		},
		{
			name: "invalid email",
			requestBody: map[string]string{
				"username": "testuser2",
				"email":    "invalid-email",
				"password": "SecurePass123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse:  nil,
		},
		{
			name: "invalid password",
			requestBody: map[string]string{
				"username": "testuser3",
				"email":    "test3@example.com",
				"password": "weak",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse:  nil,
		},
		{
			name: "duplicate username",
			requestBody: map[string]string{
				"username": "testuser",
				"email":    "another@example.com",
				"password": "SecurePass123!",
			},
			expectedStatus: http.StatusConflict,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

// TestUserLogin tests the login endpoint
func TestUserLogin(t *testing.T) {
	router := setupTestRouter(t)

	// First, register a user
	regBody, _ := json.Marshal(map[string]string{
		"username": "loginuser",
		"email":    "login@example.com",
		"password": "SecurePass123!",
	})
	req, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful login",
			requestBody: map[string]string{
				"username": "loginuser",
				"password": "SecurePass123!",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "access_token")
			},
		},
		{
			name: "incorrect password",
			requestBody: map[string]string{
				"username": "loginuser",
				"password": "WrongPassword123!",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse:  nil,
		},
		{
			name: "nonexistent user",
			requestBody: map[string]string{
				"username": "nonexistent",
				"password": "SecurePass123!",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

// TestPublicUserProfile tests retrieving a public user profile
func TestPublicUserProfile(t *testing.T) {
	router := setupTestRouter(t)

	// Register a user first
	regBody, _ := json.Marshal(map[string]string{
		"username": "publicuser",
		"email":    "public@example.com",
		"password": "SecurePass123!",
	})
	req, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var userProfile service.UserProfile
	json.Unmarshal(w.Body.Bytes(), &userProfile)

	tests := []struct {
		name           string
		userID         uint
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "existing user",
			userID:         userProfile.UserID,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response service.UserProfile
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "publicuser", response.Username)
				// Email should not be visible in public profile
				assert.Nil(t, response.Email)
			},
		},
		{
			name:           "nonexistent user",
			userID:         99999,
			expectedStatus: http.StatusNotFound,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/users/%d", tt.userID), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

// TestCurrentUserProfile tests retrieving the current authenticated user's profile
func TestCurrentUserProfile(t *testing.T) {
	router := setupTestRouter(t)

	// Register and login a user
	regBody, _ := json.Marshal(map[string]string{
		"username": "currentuser",
		"email":    "current@example.com",
		"password": "SecurePass123!",
	})
	req, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Login to get token
	loginBody, _ := json.Marshal(map[string]string{
		"username": "currentuser",
		"password": "SecurePass123!",
	})
	req, _ = http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)

	token := ""
	if tkn, ok := loginResponse["access_token"].(string); ok {
		token = tkn
	} else if tkn, ok := loginResponse["token"].(string); ok {
		token = tkn
	} else {
		// Try to get from cookie
		cookies := w.Result().Cookies()
		for _, cookie := range cookies {
			if cookie.Name == "jwt" {
				token = cookie.Value
				break
			}
		}
	}

	t.Run("get current user profile", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/users/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response service.UserProfile
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "currentuser", response.Username)
		// Email should be visible for current user
		assert.NotNil(t, response.Email)
		assert.Equal(t, "current@example.com", *response.Email)
	})

	t.Run("unauthorized access", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/users/me", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

// TestUpdateUserProfile tests updating the current user's profile
func TestUpdateUserProfile(t *testing.T) {
	router := setupTestRouter(t)

	// Register and login a user
	regBody, _ := json.Marshal(map[string]string{
		"username": "updateuser",
		"email":    "update@example.com",
		"password": "SecurePass123!",
	})
	req, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Login to get token
	loginBody, _ := json.Marshal(map[string]string{
		"username": "updateuser",
		"password": "SecurePass123!",
	})
	req, _ = http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)

	token := ""
	if tkn, ok := loginResponse["access_token"].(string); ok {
		token = tkn
	} else if tkn, ok := loginResponse["token"].(string); ok {
		token = tkn
	} else {
		// Try to get from cookie
		cookies := w.Result().Cookies()
		for _, cookie := range cookies {
			if cookie.Name == "jwt" {
				token = cookie.Value
				break
			}
		}
	}

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "update bio only",
			requestBody: map[string]interface{}{
				"bio": "This is my new bio",
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "update avatar_url only",
			requestBody: map[string]interface{}{
				"avatar_url": "https://example.com/avatar.jpg",
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "update both fields",
			requestBody: map[string]interface{}{
				"bio":        "Updated bio again",
				"avatar_url": "https://example.com/avatar2.jpg",
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "no fields provided",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("PATCH", "/v1/users/me", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
