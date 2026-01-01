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

// TestCreatePost tests creating a post in a topic
func TestCreatePost(t *testing.T) {
	router := setupTestRouter(t)
	token := registerAndLogin(t, router, "postuser", "post@example.com", "SecurePass123!")

	// Create a topic first
	topicBody, _ := json.Marshal(map[string]string{
		"title": "Topic for Posts",
	})
	req, _ := http.NewRequest("POST", "/v1/topics", bytes.NewBuffer(topicBody))
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
		requestBody    map[string]interface{}
		useAuth        bool
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:    "successful post creation with content",
			topicID: createdTopic.TopicID,
			requestBody: map[string]interface{}{
				"title":   "My First Post",
				"content": "This is the post content",
			},
			useAuth:        true,
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response service.PostInfo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "My First Post", response.Title)
				assert.NotNil(t, response.Content)
				assert.Equal(t, "This is the post content", *response.Content)
				assert.Equal(t, createdTopic.TopicID, response.TopicID)
				assert.NotNil(t, response.Author)
				assert.Equal(t, "postuser", response.Author.Username)
			},
		},
		{
			name:    "post without content",
			topicID: createdTopic.TopicID,
			requestBody: map[string]interface{}{
				"title": "Post Without Content",
			},
			useAuth:        true,
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response service.PostInfo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Post Without Content", response.Title)
			},
		},
		{
			name:    "missing title",
			topicID: createdTopic.TopicID,
			requestBody: map[string]interface{}{
				"content": "Content without title",
			},
			useAuth:        true,
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name:    "unauthorized",
			topicID: createdTopic.TopicID,
			requestBody: map[string]interface{}{
				"title": "Unauthorized Post",
			},
			useAuth:        false,
			expectedStatus: http.StatusUnauthorized,
			checkResponse:  nil,
		},
		{
			name:    "nonexistent topic",
			topicID: 99999,
			requestBody: map[string]interface{}{
				"title": "Post in Nonexistent Topic",
			},
			useAuth:        true,
			expectedStatus: http.StatusCreated,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", fmt.Sprintf("/v1/topics/%d/posts", tt.topicID), bytes.NewBuffer(body))
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

// TestListPosts tests listing posts with various filters
func TestListPosts(t *testing.T) {
	router := setupTestRouter(t)
	token := registerAndLogin(t, router, "postlistuser", "postlist@example.com", "SecurePass123!")

	// Create topics
	topic1Body, _ := json.Marshal(map[string]string{"title": "Topic 1"})
	req, _ := http.NewRequest("POST", "/v1/topics", bytes.NewBuffer(topic1Body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var topic1 service.TopicInfo
	json.Unmarshal(w.Body.Bytes(), &topic1)

	topic2Body, _ := json.Marshal(map[string]string{"title": "Topic 2"})
	req, _ = http.NewRequest("POST", "/v1/topics", bytes.NewBuffer(topic2Body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var topic2 service.TopicInfo
	json.Unmarshal(w.Body.Bytes(), &topic2)

	// Create posts in different topics
	for i := 1; i <= 2; i++ {
		body, _ := json.Marshal(map[string]string{
			"title":   fmt.Sprintf("Post %d in Topic 1", i),
			"content": fmt.Sprintf("Content for post %d", i),
		})
		req, _ := http.NewRequest("POST", fmt.Sprintf("/v1/topics/%d/posts", topic1.TopicID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	body, _ := json.Marshal(map[string]string{
		"title": "Post in Topic 2",
	})
	req, _ = http.NewRequest("POST", fmt.Sprintf("/v1/topics/%d/posts", topic2.TopicID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "list all posts",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response []service.PostInfo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(response), 3)
			},
		},
		{
			name:           "filter by topic_id",
			queryParams:    fmt.Sprintf("?topic_id=%d", topic1.TopicID),
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response []service.PostInfo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(response), 2)
				for _, post := range response {
					assert.Equal(t, topic1.TopicID, post.TopicID)
				}
			},
		},
		{
			name:           "list with limit",
			queryParams:    "?limit=2",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response []service.PostInfo
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
				var response []service.PostInfo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(response), 2)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/v1/posts"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

// TestGetPostInfo tests retrieving a single post
func TestGetPostInfo(t *testing.T) {
	router := setupTestRouter(t)
	token := registerAndLogin(t, router, "getpostuser", "getpost@example.com", "SecurePass123!")

	// Create a topic
	topicBody, _ := json.Marshal(map[string]string{"title": "Topic for Get Post"})
	req, _ := http.NewRequest("POST", "/v1/topics", bytes.NewBuffer(topicBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var topic service.TopicInfo
	json.Unmarshal(w.Body.Bytes(), &topic)

	// Create a post
	postBody, _ := json.Marshal(map[string]string{
		"title":   "Specific Post",
		"content": "Content for specific post",
	})
	req, _ = http.NewRequest("POST", fmt.Sprintf("/v1/topics/%d/posts", topic.TopicID), bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdPost service.PostInfo
	json.Unmarshal(w.Body.Bytes(), &createdPost)

	tests := []struct {
		name           string
		postID         uint
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "get existing post",
			postID:         createdPost.PostID,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response service.PostInfo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Specific Post", response.Title)
				assert.NotNil(t, response.Content)
				assert.Equal(t, "Content for specific post", *response.Content)
			},
		},
		{
			name:           "get nonexistent post",
			postID:         99999,
			expectedStatus: http.StatusNotFound,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/posts/%d", tt.postID), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

// TestUpdatePost tests updating a post
func TestUpdatePost(t *testing.T) {
	router := setupTestRouter(t)
	token := registerAndLogin(t, router, "updatepostuser", "updatepost@example.com", "SecurePass123!")
	otherToken := registerAndLogin(t, router, "otheruserpost", "otherpost@example.com", "SecurePass123!")

	// Create a topic
	topicBody, _ := json.Marshal(map[string]string{"title": "Topic for Update Post"})
	req, _ := http.NewRequest("POST", "/v1/topics", bytes.NewBuffer(topicBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var topic service.TopicInfo
	json.Unmarshal(w.Body.Bytes(), &topic)

	// Create a post
	postBody, _ := json.Marshal(map[string]string{
		"title":   "Original Post",
		"content": "Original content",
	})
	req, _ = http.NewRequest("POST", fmt.Sprintf("/v1/topics/%d/posts", topic.TopicID), bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdPost service.PostInfo
	json.Unmarshal(w.Body.Bytes(), &createdPost)

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
			name: "update content only",
			requestBody: map[string]interface{}{
				"content": "Updated content",
			},
			useToken:       token,
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "update both fields",
			requestBody: map[string]interface{}{
				"title":   "Fully Updated Post",
				"content": "Fully updated content",
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
			req, _ := http.NewRequest("PATCH", fmt.Sprintf("/v1/posts/%d", createdPost.PostID), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tt.useToken)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// TestDeletePost tests deleting a post
func TestDeletePost(t *testing.T) {
	router := setupTestRouter(t)
	token := registerAndLogin(t, router, "deletepostuser", "deletepost@example.com", "SecurePass123!")
	otherToken := registerAndLogin(t, router, "otherdelpost", "otherdelpost@example.com", "SecurePass123!")

	// Create a topic
	topicBody, _ := json.Marshal(map[string]string{"title": "Topic for Delete Post"})
	req, _ := http.NewRequest("POST", "/v1/topics", bytes.NewBuffer(topicBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var topic service.TopicInfo
	json.Unmarshal(w.Body.Bytes(), &topic)

	// Create a post to delete
	postBody, _ := json.Marshal(map[string]string{
		"title": "Post to Delete",
	})
	req, _ = http.NewRequest("POST", fmt.Sprintf("/v1/topics/%d/posts", topic.TopicID), bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdPost service.PostInfo
	json.Unmarshal(w.Body.Bytes(), &createdPost)

	tests := []struct {
		name           string
		postID         uint
		useToken       string
		expectedStatus int
	}{
		{
			name:           "unauthorized user cannot delete",
			postID:         createdPost.PostID,
			useToken:       otherToken,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "authorized user can delete",
			postID:         createdPost.PostID,
			useToken:       token,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "delete nonexistent post",
			postID:         99999,
			useToken:       token,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/v1/posts/%d", tt.postID), nil)
			req.Header.Set("Authorization", "Bearer "+tt.useToken)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
