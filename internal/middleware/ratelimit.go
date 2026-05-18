package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/ilhamazhar/golang-gpt/pkg/response"
)

func RateLimit(limiter *redis_rate.Limiter, rate redis_rate.Limit, keyFn func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := keyFn(c)
		res, err := limiter.Allow(c.Request.Context(), key, rate)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "rate limiter unavailable", nil)
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", res.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", int(res.ResetAfter.Seconds())))

		if res.Allowed == 0 {
			c.Header("Retry-After", fmt.Sprintf("%d", int(res.RetryAfter.Seconds())))
			response.Fail(c, http.StatusTooManyRequests, "too many requests", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

// IPKey rate-limits by the caller's IP — use on public/unauthenticated routes.
func IPKey(prefix string) func(*gin.Context) string {
	return func(c *gin.Context) string {
		return fmt.Sprintf("%s:ip:%s", prefix, c.ClientIP())
	}
}

// UserKey rate-limits by authenticated user ID — falls back to IP when no claims.
func UserKey(prefix string) func(*gin.Context) string {
	return func(c *gin.Context) string {
		if claims := ClaimsFromContext(c); claims != nil {
			return fmt.Sprintf("%s:user:%s", prefix, claims.UserID.String())
		}
		return fmt.Sprintf("%s:ip:%s", prefix, c.ClientIP())
	}
}
