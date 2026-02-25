package gin

import (
	ginpkg "github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"simple-banking-system/module/account/storage"
)

func RegisterRoutes(r ginpkg.IRouter, db *pgxpool.Pool) {
	store := storage.NewSQLStore(db)

	v1 := r.Group("/v1")

	accounts := v1.Group("/accounts")
	{
		accounts.POST("", CreateAccountHandler(store))
		accounts.GET("/:id", GetAccountHandler(store))
		accounts.GET("", ListAccountsHandler(store))
	}

	v1.POST("/transfers", TransferMoneyHandler(store))
}
