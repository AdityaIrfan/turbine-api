package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

// Rate limiters for GET and POST requests per token
var (
	getTokenLimiters  = make(map[string]*rate.Limiter)
	postTokenLimiters = make(map[string]*rate.Limiter)
	mu                = sync.Mutex{}
)

// Create or get rate limiter for a specific token and method
func getRateLimiter(token string, method string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	var limiter *rate.Limiter
	var limitersMap map[string]*rate.Limiter

	switch method {
	case http.MethodGet:
		limitersMap = getTokenLimiters
		_, exists := limitersMap[token]
		if !exists {
			limiter = rate.NewLimiter(rate.Every(5*time.Minute), 3) // 10 GET requests per minute
			limitersMap[token] = limiter
		}
	case http.MethodPost:
		limitersMap = postTokenLimiters
		_, exists := limitersMap[token]
		if !exists {
			limiter = rate.NewLimiter(rate.Every(1*time.Minute), 5) // 5 POST requests per minute
			limitersMap[token] = limiter
		}
	default:
		return nil
	}

	return limiter
}

// Rate limiter middleware based on bearer token
func RateLimiterMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if len(token) > 0 && len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:] // Extract the token after "Bearer "
		} else {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		}

		method := c.Request().Method
		limiter := getRateLimiter(token, method)
		if limiter == nil {
			return next(c)
		}

		if !limiter.Allow() {
			return c.JSON(http.StatusTooManyRequests, map[string]string{"message": "Too Many Requests"})
		}

		return next(c)
	}
}
