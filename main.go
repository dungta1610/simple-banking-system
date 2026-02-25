package main

import (
	"context"
	"log"
	"net/http"
	"os"

	ginpkg "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"simple-banking-system/component/postgres"
	accountgin "simple-banking-system/module/account/transport/gin"
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

	pool, err := postgres.NewPool(ctx, postgres.Config{
		DSN: dsn,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	r := ginpkg.Default()

	r.GET("/health", func(c *ginpkg.Context) {
		if err := pool.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, ginpkg.H{
				"status": "down",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, ginpkg.H{
			"status": "ok",
		})
	})

	accountgin.RegisterRoutes(r, pool)

	log.Printf("server is running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
