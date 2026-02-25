package gin

import (
	ginpkg "github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"simple-banking-system/module/account/storage"
)

func RegisterRoutes(r ginpkg.IRouter, db *pgxpool.Pool, transferRateLimit ginpkg.HandlerFunc) {
	store := storage.NewSQLStore(db)
	v1 := r.Group("/v1")

	accounts := v1.Group("/accounts")
	{
		accounts.POST("", CreateAccountHandler(store))
		accounts.GET("/:id", GetAccountHandler(store))
		accounts.GET("", ListAccountsHandler(store))
	}

	if transferRateLimit != nil {
		v1.POST("/transfers", transferRateLimit, TransferMoneyHandler(store))
	} else {
		v1.POST("/transfers", TransferMoneyHandler(store))
	}
}
