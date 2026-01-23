package middleware

import (
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
)

// NewAuthRateLimiter creates a rate limiter for authentication endpoints
// to prevent brute force attacks.
// Limits: 5 requests per minute per IP address
func NewAuthRateLimiter() gin.HandlerFunc {
	// Define the key function to identify clients by IP address
	keyFunc := func(c *gin.Context) string {
		return c.ClientIP()
	}

	// Define the error handler for when rate limit is exceeded
	errorHandler := func(c *gin.Context, info ratelimit.Info) {
		c.JSON(429, gin.H{
			"error":       "Too many authentication attempts. Please try again later.",
			"retry_after": info.ResetTime.Sub(time.Now()).Seconds(),
		})
	}

	// Create store using in-memory storage
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Minute,
		Limit: 5,
	})

	// Create and return the rate limiter middleware
	return ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})
}

// NewGeneralRateLimiter creates a rate limiter for general API endpoints
// Limits: 100 requests per minute per IP address
func NewGeneralRateLimiter() gin.HandlerFunc {
	keyFunc := func(c *gin.Context) string {
		return c.ClientIP()
	}

	errorHandler := func(c *gin.Context, info ratelimit.Info) {
		c.JSON(429, gin.H{
			"error":       "Rate limit exceeded. Please try again later.",
			"retry_after": info.ResetTime.Sub(time.Now()).Seconds(),
		})
	}

	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Minute,
		Limit: 100,
	})

	return ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})
}
