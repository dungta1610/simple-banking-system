package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	ginpkg "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"simple-banking-system/component/postgres"
	"simple-banking-system/component/ratelimit"
	rediscmp "simple-banking-system/component/redis"
	accountgin "simple-banking-system/module/account/transport/gin"
	"simple-banking-system/module/account/transport/gin/middleware"
)

func main() {
	_ = godotenv.Load()
	dsn := os.Getenv("DB_DSN")

	if dsn == "" {
		log.Fatal("DB_DSN is required")
	}

	redisAddr := os.Getenv("REDIS_ADDR")

	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	ctx := context.Background()

	pool, err := postgres.NewPool(ctx, postgres.Config{DSN: dsn})

	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	redisClient, err := rediscmp.NewClient(ctx, rediscmp.Config{
		Addr: redisAddr,
		DB:   0,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer redisClient.Close()

	transferLimiter := ratelimit.NewRedisLimiter(redisClient, "transfer_limit:")
	transferRateLimitMW := middleware.RateLimit(transferLimiter, 10, time.Minute)
	r := ginpkg.Default()

	r.GET("/health", func(c *ginpkg.Context) {
		if err := pool.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, ginpkg.H{
				"status": "down",
				"error":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, ginpkg.H{"status": "ok"})
	})

	accountgin.RegisterRoutes(r, pool, transferRateLimitMW)

	log.Printf("server is running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
