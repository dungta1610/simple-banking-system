package middleware

import (
	"net/http"
	"time"

	"simple-banking-system/component/ratelimit"

	"github.com/gin-gonic/gin"
)

func RateLimit(limiter ratelimit.Limiter, limit int64, window time.Duration) gin.HandlerFunc {
	if limit <= 0 {
		limit = 10
	}

	if window <= 0 {
		window = time.Minute
	}

	return func(c *gin.Context) {
		if limiter == nil {
			c.Next()
			return
		}

		ip := c.ClientIP()
		path := c.FullPath()

		if path == "" {
			path = c.Request.URL.Path
		}

		if !limiter.IsAllowed(c.Request.Context(), ip, path, limit, window) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests, retry later",
			})

			return
		}

		c.Next()
	}
}
