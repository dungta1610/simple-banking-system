package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"simple-banking-system/component/postgres"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN is required")
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

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		if err := pool.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "down"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
