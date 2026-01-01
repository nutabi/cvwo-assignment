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

// TestCreateComment tests creating a comment on a post
func TestCreateComment(t *testing.T) {
	router := setupTestRouter(t)
	token := registerAndLogin(t, router, "commentuser", "comment@example.com", "SecurePass123!")

	// Create a topic first
	topicBody, _ := json.Marshal(map[string]string{
		"title": "Topic for Comments",
	})
	req, _ := http.NewRequest("POST", "/v1/topics", bytes.NewBuffer(topicBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdTopic service.TopicInfo
	json.Unmarshal(w.Body.Bytes(), &createdTopic)

	// Create a post
	postBody, _ := json.Marshal(map[string]string{
		"title":   "Post for Comments",
		"content": "This post will have comments",
	})
	req, _ = http.NewRequest("POST", fmt.Sprintf("/v1/topics/%d/posts", createdTopic.TopicID), bytes.NewBuffer(postBody))
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
		requestBody    map[string]interface{}
		useAuth        bool
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "successful comment creation",
			postID: createdPost.PostID,
			requestBody: map[string]interface{}{
				"content": "This is a great post!",
			},
			useAuth:        true,
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response service.CommentInfo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "This is a great post!", response.Content)
				assert.Equal(t, createdPost.PostID, response.PostID)
				assert.NotNil(t, response.Author)
				assert.Equal(t, "commentuser", response.Author.Username)
			},
		},
		{
			name:   "missing content",
			postID: createdPost.PostID,
			requestBody: map[string]interface{}{
				"title": "No content field",
			},
			useAuth:        true,
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name:   "empty content",
			postID: createdPost.PostID,
			requestBody: map[string]interface{}{
				"content": "",
			},
			useAuth:        true,
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name:   "unauthorized",
			postID: createdPost.PostID,
			requestBody: map[string]interface{}{
				"content": "Unauthorized comment",
			},
			useAuth:        false,
			expectedStatus: http.StatusUnauthorized,
			checkResponse:  nil,
		},
		{
			name:   "nonexistent post",
			postID: 99999,
			requestBody: map[string]interface{}{
				"content": "Comment on nonexistent post",
			},
			useAuth:        true,
			expectedStatus: http.StatusCreated,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", fmt.Sprintf("/v1/posts/%d/comments", tt.postID), bytes.NewBuffer(body))
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

// TestListComments tests listing comments with various filters
func TestListComments(t *testing.T) {
	router := setupTestRouter(t)
	token1 := registerAndLogin(t, router, "commentlistuser1", "commentlist1@example.com", "SecurePass123!")
	token2 := registerAndLogin(t, router, "commentlistuser2", "commentlist2@example.com", "SecurePass123!")

	// Create topic and posts
	topicBody, _ := json.Marshal(map[string]string{"title": "Topic for Comment Listing"})
	req, _ := http.NewRequest("POST", "/v1/topics", bytes.NewBuffer(topicBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token1)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var topic service.TopicInfo
	json.Unmarshal(w.Body.Bytes(), &topic)

	// Create post 1
	post1Body, _ := json.Marshal(map[string]string{
		"title":   "Post 1",
		"content": "First post",
	})
	req, _ = http.NewRequest("POST", fmt.Sprintf("/v1/topics/%d/posts", topic.TopicID), bytes.NewBuffer(post1Body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token1)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var post1 service.PostInfo
	json.Unmarshal(w.Body.Bytes(), &post1)

	// Create post 2
	post2Body, _ := json.Marshal(map[string]string{
		"title":   "Post 2",
		"content": "Second post",
	})
	req, _ = http.NewRequest("POST", fmt.Sprintf("/v1/topics/%d/posts", topic.TopicID), bytes.NewBuffer(post2Body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token1)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var post2 service.PostInfo
	json.Unmarshal(w.Body.Bytes(), &post2)

	// Create comments on post 1 by user 1
	for i := 1; i <= 3; i++ {
		commentBody, _ := json.Marshal(map[string]string{
			"content": fmt.Sprintf("Comment %d on Post 1 by User 1", i),
		})
		req, _ = http.NewRequest("POST", fmt.Sprintf("/v1/posts/%d/comments", post1.PostID), bytes.NewBuffer(commentBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token1)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	// Create comments on post 1 by user 2
	for i := 1; i <= 2; i++ {
		commentBody, _ := json.Marshal(map[string]string{
			"content": fmt.Sprintf("Comment %d on Post 1 by User 2", i),
		})
		req, _ = http.NewRequest("POST", fmt.Sprintf("/v1/posts/%d/comments", post1.PostID), bytes.NewBuffer(commentBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token2)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	// Create comments on post 2
	for i := 1; i <= 2; i++ {
		commentBody, _ := json.Marshal(map[string]string{
			"content": fmt.Sprintf("Comment %d on Post 2", i),
		})
		req, _ = http.NewRequest("POST", fmt.Sprintf("/v1/posts/%d/comments", post2.PostID), bytes.NewBuffer(commentBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token1)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	// Get user IDs
	req, _ = http.NewRequest("GET", "/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token1)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var user1 service.UserProfile
	json.Unmarshal(w.Body.Bytes(), &user1)

	tests := []struct {
		name          string
		queryParams   string
		expectedCount int
		checkResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:          "list all comments",
			queryParams:   "",
			expectedCount: 7, // 3 + 2 + 2
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var comments []service.CommentInfo
				err := json.Unmarshal(w.Body.Bytes(), &comments)
				assert.NoError(t, err)
				assert.Equal(t, 7, len(comments))
			},
		},
		{
			name:          "filter by post 1",
			queryParams:   fmt.Sprintf("?post_id=%d", post1.PostID),
			expectedCount: 5, // 3 + 2
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var comments []service.CommentInfo
				err := json.Unmarshal(w.Body.Bytes(), &comments)
				assert.NoError(t, err)
				assert.Equal(t, 5, len(comments))
				for _, comment := range comments {
					assert.Equal(t, post1.PostID, comment.PostID)
				}
			},
		},
		{
			name:          "filter by post 2",
			queryParams:   fmt.Sprintf("?post_id=%d", post2.PostID),
			expectedCount: 2,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var comments []service.CommentInfo
				err := json.Unmarshal(w.Body.Bytes(), &comments)
				assert.NoError(t, err)
				assert.Equal(t, 2, len(comments))
				for _, comment := range comments {
					assert.Equal(t, post2.PostID, comment.PostID)
				}
			},
		},
		{
			name:          "filter by user 1",
			queryParams:   fmt.Sprintf("?user_id=%d", user1.UserID),
			expectedCount: 5, // 3 on post1 + 2 on post2
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var comments []service.CommentInfo
				err := json.Unmarshal(w.Body.Bytes(), &comments)
				assert.NoError(t, err)
				assert.Equal(t, 5, len(comments))
				for _, comment := range comments {
					assert.Equal(t, user1.UserID, comment.Author.UserID)
				}
			},
		},
		{
			name:          "pagination - limit 3",
			queryParams:   "?limit=3",
			expectedCount: 3,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var comments []service.CommentInfo
				err := json.Unmarshal(w.Body.Bytes(), &comments)
				assert.NoError(t, err)
				assert.Equal(t, 3, len(comments))
			},
		},
		{
			name:          "pagination - offset 3",
			queryParams:   "?limit=10&offset=3",
			expectedCount: 4,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var comments []service.CommentInfo
				err := json.Unmarshal(w.Body.Bytes(), &comments)
				assert.NoError(t, err)
				assert.Equal(t, 4, len(comments))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/v1/comments"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

// TestGetComment tests retrieving a specific comment
func TestGetComment(t *testing.T) {
	router := setupTestRouter(t)
	token := registerAndLogin(t, router, "getcommentuser", "getcomment@example.com", "SecurePass123!")

	// Create a topic
	topicBody, _ := json.Marshal(map[string]string{"title": "Topic for Get Comment"})
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
	var post service.PostInfo
	json.Unmarshal(w.Body.Bytes(), &post)

	// Create a comment
	commentBody, _ := json.Marshal(map[string]string{
		"content": "Specific comment to retrieve",
	})
	req, _ = http.NewRequest("POST", fmt.Sprintf("/v1/posts/%d/comments", post.PostID), bytes.NewBuffer(commentBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdComment service.CommentInfo
	json.Unmarshal(w.Body.Bytes(), &createdComment)

	tests := []struct {
		name           string
		commentID      uint
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "get existing comment",
			commentID:      createdComment.CommentID,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response service.CommentInfo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Specific comment to retrieve", response.Content)
				assert.Equal(t, post.PostID, response.PostID)
			},
		},
		{
			name:           "get nonexistent comment",
			commentID:      99999,
			expectedStatus: http.StatusNotFound,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/comments/%d", tt.commentID), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

// TestUpdateComment tests updating a comment
func TestUpdateComment(t *testing.T) {
	router := setupTestRouter(t)
	token := registerAndLogin(t, router, "updatecommentuser", "updatecomment@example.com", "SecurePass123!")
	otherToken := registerAndLogin(t, router, "otherusercomment", "othercomment@example.com", "SecurePass123!")

	// Create a topic
	topicBody, _ := json.Marshal(map[string]string{"title": "Topic for Update Comment"})
	req, _ := http.NewRequest("POST", "/v1/topics", bytes.NewBuffer(topicBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var topic service.TopicInfo
	json.Unmarshal(w.Body.Bytes(), &topic)

	// Create a post
	postBody, _ := json.Marshal(map[string]string{
		"title":   "Post for Comment Update",
		"content": "Post content",
	})
	req, _ = http.NewRequest("POST", fmt.Sprintf("/v1/topics/%d/posts", topic.TopicID), bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var post service.PostInfo
	json.Unmarshal(w.Body.Bytes(), &post)

	// Create a comment
	commentBody, _ := json.Marshal(map[string]string{
		"content": "Original comment",
	})
	req, _ = http.NewRequest("POST", fmt.Sprintf("/v1/posts/%d/comments", post.PostID), bytes.NewBuffer(commentBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdComment service.CommentInfo
	json.Unmarshal(w.Body.Bytes(), &createdComment)

	tests := []struct {
		name           string
		commentID      uint
		requestBody    map[string]interface{}
		useToken       string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:      "successful update",
			commentID: createdComment.CommentID,
			requestBody: map[string]interface{}{
				"content": "Updated comment content",
			},
			useToken:       token,
			expectedStatus: http.StatusNoContent,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				// Verify the update by fetching the comment
				req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/comments/%d", createdComment.CommentID), nil)
				w2 := httptest.NewRecorder()
				router.ServeHTTP(w2, req)
				var updated service.CommentInfo
				json.Unmarshal(w2.Body.Bytes(), &updated)
				assert.Equal(t, "Updated comment content", updated.Content)
			},
		},
		{
			name:      "missing content field",
			commentID: createdComment.CommentID,
			requestBody: map[string]interface{}{
				"text": "Wrong field",
			},
			useToken:       token,
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse:  nil,
		},
		{
			name:           "unauthorized",
			commentID:      createdComment.CommentID,
			requestBody:    map[string]interface{}{"content": "Unauthorized update"},
			useToken:       "",
			expectedStatus: http.StatusUnauthorized,
			checkResponse:  nil,
		},
		{
			name:      "forbidden - different user",
			commentID: createdComment.CommentID,
			requestBody: map[string]interface{}{
				"content": "Update by wrong user",
			},
			useToken:       otherToken,
			expectedStatus: http.StatusForbidden,
			checkResponse:  nil,
		},
		{
			name:      "nonexistent comment",
			commentID: 99999,
			requestBody: map[string]interface{}{
				"content": "Update nonexistent",
			},
			useToken:       token,
			expectedStatus: http.StatusNotFound,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("PATCH", fmt.Sprintf("/v1/comments/%d", tt.commentID), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			if tt.useToken != "" {
				req.Header.Set("Authorization", "Bearer "+tt.useToken)
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

// TestDeleteComment tests deleting a comment
func TestDeleteComment(t *testing.T) {
	router := setupTestRouter(t)
	token := registerAndLogin(t, router, "deletecommentuser", "deletecomment@example.com", "SecurePass123!")
	otherToken := registerAndLogin(t, router, "otherdelcomment", "otherdelcomment@example.com", "SecurePass123!")

	// Create a topic
	topicBody, _ := json.Marshal(map[string]string{"title": "Topic for Delete Comment"})
	req, _ := http.NewRequest("POST", "/v1/topics", bytes.NewBuffer(topicBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var topic service.TopicInfo
	json.Unmarshal(w.Body.Bytes(), &topic)

	// Create a post
	postBody, _ := json.Marshal(map[string]string{
		"title":   "Post for Comment Deletion",
		"content": "Post content",
	})
	req, _ = http.NewRequest("POST", fmt.Sprintf("/v1/topics/%d/posts", topic.TopicID), bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var post service.PostInfo
	json.Unmarshal(w.Body.Bytes(), &post)

	// Create a comment for successful deletion
	commentBody1, _ := json.Marshal(map[string]string{
		"content": "Comment to delete",
	})
	req, _ = http.NewRequest("POST", fmt.Sprintf("/v1/posts/%d/comments", post.PostID), bytes.NewBuffer(commentBody1))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var comment1 service.CommentInfo
	json.Unmarshal(w.Body.Bytes(), &comment1)

	// Create a comment for forbidden test
	commentBody2, _ := json.Marshal(map[string]string{
		"content": "Comment by first user",
	})
	req, _ = http.NewRequest("POST", fmt.Sprintf("/v1/posts/%d/comments", post.PostID), bytes.NewBuffer(commentBody2))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var comment2 service.CommentInfo
	json.Unmarshal(w.Body.Bytes(), &comment2)

	tests := []struct {
		name           string
		commentID      uint
		useToken       string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "successful deletion",
			commentID:      comment1.CommentID,
			useToken:       token,
			expectedStatus: http.StatusNoContent,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				// Verify deletion by trying to get the comment
				req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/comments/%d", comment1.CommentID), nil)
				w2 := httptest.NewRecorder()
				router.ServeHTTP(w2, req)
				assert.Equal(t, http.StatusNotFound, w2.Code)
			},
		},
		{
			name:           "unauthorized",
			commentID:      comment2.CommentID,
			useToken:       "",
			expectedStatus: http.StatusUnauthorized,
			checkResponse:  nil,
		},
		{
			name:           "forbidden - different user",
			commentID:      comment2.CommentID,
			useToken:       otherToken,
			expectedStatus: http.StatusForbidden,
			checkResponse:  nil,
		},
		{
			name:           "nonexistent comment",
			commentID:      99999,
			useToken:       token,
			expectedStatus: http.StatusNotFound,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/v1/comments/%d", tt.commentID), nil)
			if tt.useToken != "" {
				req.Header.Set("Authorization", "Bearer "+tt.useToken)
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
