package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRateLimitOnAuthEndpoints(t *testing.T) {
	r := setupTestRouter(t)

	// Test rate limiting on login endpoint
	t.Run("rate_limit_on_login", func(t *testing.T) {
		loginPayload := `{"username":"testuser","password":"TestPassword123"}`

		// Make 5 requests (should succeed)
		for i := 0; i < 5; i++ {
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(loginPayload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// Requests should not be rate limited (might be 401 unauthorized, but not 429)
			assert.NotEqual(t, http.StatusTooManyRequests, w.Code, "Request %d should not be rate limited", i+1)
		}

		// 6th request should be rate limited
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(loginPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code, "6th request should be rate limited")
		assert.Contains(t, w.Body.String(), "Too many authentication attempts")
	})

	t.Run("rate_limit_on_register", func(t *testing.T) {
		// Create a new router for this test to reset rate limits
		r := setupTestRouter(t)

		// Make 5 registration attempts
		for i := 0; i < 5; i++ {
			payload := `{"username":"user` + string(rune(i+48)) + `","email":"user` + string(rune(i+48)) + `@test.com","password":"TestPass123"}`
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", strings.NewReader(payload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.NotEqual(t, http.StatusTooManyRequests, w.Code, "Request %d should not be rate limited", i+1)
		}

		// 6th request should be rate limited
		payload := `{"username":"user6","email":"user6@test.com","password":"TestPass123"}`
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code, "6th request should be rate limited")
		assert.Contains(t, w.Body.String(), "Too many authentication attempts")
	})
}
