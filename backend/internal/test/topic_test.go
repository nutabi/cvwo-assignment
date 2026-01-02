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

// Helper function to register a user and get auth token
func registerAndLogin(t *testing.T, router http.Handler, username, email, password string) string {
	// Register
	regBody, _ := json.Marshal(map[string]string{
		"username": username,
		"email":    email,
		"password": password,
	})
	req, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Login
	loginBody, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	req, _ = http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Extract token from response body or cookie
	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)

	// Try access_token first (gin-jwt standard)
	if token, ok := loginResponse["access_token"].(string); ok {
		return token
	}

	// Fallback to token
	if token, ok := loginResponse["token"].(string); ok {
		return token
	}

	// Try to get from cookie if not in body
	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "jwt" {
			return cookie.Value
		}
	}

	t.Fatalf("Failed to get token from login response")
	return ""
}

// TestCreateTopic tests creating a topic
func TestCreateTopic(t *testing.T) {
	router := setupTestRouter(t)
	token := registerAndLogin(t, router, "topicuser", "topic@example.com", "SecurePass123!")

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		useAuth        bool
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful topic creation",
			requestBody: map[string]interface{}{
				"title":       "My First Topic",
				"description": "This is a test topic",
			},
			useAuth:        true,
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response service.TopicInfo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "My First Topic", response.Name)
				assert.NotNil(t, response.Description)
				assert.Equal(t, "This is a test topic", *response.Description)
				assert.NotNil(t, response.Author)
				assert.Equal(t, "topicuser", response.Author.Username)
			},
		},
		{
			name: "topic without description",
			requestBody: map[string]interface{}{
				"title": "Topic Without Description",
			},
			useAuth:        true,
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response service.TopicInfo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Topic Without Description", response.Name)
			},
		},
		{
			name: "duplicate topic title",
			requestBody: map[string]interface{}{
				"title":       "My First Topic",
				"description": "This should fail due to duplicate title",
			},
			useAuth:        true,
			expectedStatus: http.StatusConflict,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], "topic title already exists")
			},
		},
		{
			name:           "missing title",
			requestBody:    map[string]interface{}{},
			useAuth:        true,
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name: "unauthorized",
			requestBody: map[string]interface{}{
				"title": "Unauthorized Topic",
			},
			useAuth:        false,
			expectedStatus: http.StatusUnauthorized,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/v1/topics", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			if tt.useAuth {
				req.Header.Set("Authorization", "Bearer "+token)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

// TestListTopics tests listing topics
func TestListTopics(t *testing.T) {
	router := setupTestRouter(t)
	token := registerAndLogin(t, router, "topiclistuser", "topiclist@example.com", "SecurePass123!")

	// Create some topics
	for i := 1; i <= 3; i++ {
		body, _ := json.Marshal(map[string]string{
			"title":       fmt.Sprintf("Topic %d", i),
			"description": fmt.Sprintf("Description for topic %d", i),
		})
		req, _ := http.NewRequest("POST", "/v1/topics", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "list all topics",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response []service.TopicInfo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(response), 3)
			},
		},
		{
			name:           "list with limit",
			queryParams:    "?limit=2",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response []service.TopicInfo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.LessOrEqual(t, len(response), 2)
			},
		},
		{
			name:           "list with offset",
			queryParams:    "?offset=1",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response []service.TopicInfo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(response), 2)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/v1/topics"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

// TestGetTopicInfo tests retrieving a single topic
func TestGetTopicInfo(t *testing.T) {
	router := setupTestRouter(t)
	token := registerAndLogin(t, router, "gettopicuser", "gettopic@example.com", "SecurePass123!")

	// Create a topic
	body, _ := json.Marshal(map[string]string{
		"title":       "Specific Topic",
		"description": "A specific topic to retrieve",
	})
	req, _ := http.NewRequest("POST", "/v1/topics", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdTopic service.TopicInfo
	json.Unmarshal(w.Body.Bytes(), &createdTopic)

	tests := []struct {
		name           string
		topicID        uint
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "get existing topic",
			topicID:        createdTopic.TopicID,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response service.TopicInfo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Specific Topic", response.Name)
				assert.NotNil(t, response.Description)
				assert.Equal(t, "A specific topic to retrieve", *response.Description)
			},
		},
		{
			name:           "get nonexistent topic",
			topicID:        99999,
			expectedStatus: http.StatusNotFound,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/topics/%d", tt.topicID), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

// TestUpdateTopic tests updating a topic
func TestUpdateTopic(t *testing.T) {
	router := setupTestRouter(t)
	token := registerAndLogin(t, router, "updatetopicuser", "updatetopic@example.com", "SecurePass123!")
	otherToken := registerAndLogin(t, router, "otheruser", "other@example.com", "SecurePass123!")

	// Create a topic
	body, _ := json.Marshal(map[string]string{
		"title":       "Original Topic",
		"description": "Original description",
	})
	req, _ := http.NewRequest("POST", "/v1/topics", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdTopic service.TopicInfo
	json.Unmarshal(w.Body.Bytes(), &createdTopic)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		useToken       string
		expectedStatus int
	}{
		{
			name: "update title only",
			requestBody: map[string]interface{}{
				"title": "Updated Title",
			},
			useToken:       token,
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "update description only",
			requestBody: map[string]interface{}{
				"description": "Updated description",
			},
			useToken:       token,
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "update both fields",
			requestBody: map[string]interface{}{
				"title":       "Fully Updated Topic",
				"description": "Fully updated description",
			},
			useToken:       token,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "no fields provided",
			requestBody:    map[string]interface{}{},
			useToken:       token,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "unauthorized user",
			requestBody: map[string]interface{}{
				"title": "Hacked Title",
			},
			useToken:       otherToken,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("PATCH", fmt.Sprintf("/v1/topics/%d", createdTopic.TopicID), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tt.useToken)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// TestDeleteTopic tests deleting a topic
func TestDeleteTopic(t *testing.T) {
	router := setupTestRouter(t)
	token := registerAndLogin(t, router, "deletetopicuser", "deletetopic@example.com", "SecurePass123!")
	otherToken := registerAndLogin(t, router, "otherdeluser", "otherdel@example.com", "SecurePass123!")

	// Create a topic to delete
	body, _ := json.Marshal(map[string]string{
		"title": "Topic to Delete",
	})
	req, _ := http.NewRequest("POST", "/v1/topics", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdTopic service.TopicInfo
	json.Unmarshal(w.Body.Bytes(), &createdTopic)

	tests := []struct {
		name           string
		topicID        uint
		useToken       string
		expectedStatus int
	}{
		{
			name:           "unauthorized user cannot delete",
			topicID:        createdTopic.TopicID,
			useToken:       otherToken,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "authorized user can delete",
			topicID:        createdTopic.TopicID,
			useToken:       token,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "delete nonexistent topic",
			topicID:        99999,
			useToken:       token,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/v1/topics/%d", tt.topicID), nil)
			req.Header.Set("Authorization", "Bearer "+tt.useToken)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
