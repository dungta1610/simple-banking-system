package gin

import (
	"errors"
	"net/http"

	ginpkg "github.com/gin-gonic/gin"

	"simple-banking-system/module/account/biz"
	"simple-banking-system/module/account/model"
)

func CreateAccountHandler(store biz.CreateAccountStore) ginpkg.HandlerFunc {
	return func(c *ginpkg.Context) {
		var req model.CreateAccountRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{
				"error": "invalid request body",
			})
			return
		}

		createAccountBiz := biz.NewCreateAccountBiz(store)
		account, err := createAccountBiz.CreateAccount(c.Request.Context(), &req)

		if err != nil {
			switch {
			case errors.Is(err, model.ErrInvalidRequest):
				c.JSON(http.StatusBadRequest, ginpkg.H{
					"error": err.Error(),
				})
				return
			default:
				c.JSON(http.StatusInternalServerError, ginpkg.H{
					"error": err.Error(),
				})
				return
			}
		}

		c.JSON(http.StatusCreated, ginpkg.H{
			"data": account,
		})
	}
}
