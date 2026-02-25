package gin

import (
	"errors"
	"net/http"

	ginpkg "github.com/gin-gonic/gin"

	"simple-banking-system/module/account/biz"
	"simple-banking-system/module/account/model"
)

func TransferMoneyHandler(store biz.TransferMoneyStore) ginpkg.HandlerFunc {
	return func(c *ginpkg.Context) {
		var req model.CreateTransferRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{
				"error": "invalid request body",
			})
			return
		}

		transferBiz := biz.NewTransferMoneyBiz(store)
		result, err := transferBiz.TransferMoney(c.Request.Context(), &req)

		if err != nil {
			switch {
			case errors.Is(err, model.ErrInvalidRequest):
				c.JSON(http.StatusBadRequest, ginpkg.H{
					"error": err.Error(),
				})
				return

			case errors.Is(err, model.ErrSameAccountTransfer):
				c.JSON(http.StatusBadRequest, ginpkg.H{
					"error": err.Error(),
				})
				return

			case errors.Is(err, model.ErrAccountNotFound):
				c.JSON(http.StatusNotFound, ginpkg.H{
					"error": err.Error(),
				})
				return

			case errors.Is(err, model.ErrInsufficientFunds):
				c.JSON(http.StatusConflict, ginpkg.H{
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
			"data": result,
		})
	}
}
