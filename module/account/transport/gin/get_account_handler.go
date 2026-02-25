package gin

import (
	"errors"
	"net/http"
	"strconv"

	ginpkg "github.com/gin-gonic/gin"

	"simple-banking-system/module/account/biz"
	"simple-banking-system/module/account/model"
)

func GetAccountHandler(store biz.GetAccountStore) ginpkg.HandlerFunc {
	return func(c *ginpkg.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{
				"error": "invalid account id",
			})
			return
		}

		getAccountBiz := biz.NewGetAccountBiz(store)
		account, err := getAccountBiz.GetAccount(c.Request.Context(), id)

		if err != nil {
			switch {
			case errors.Is(err, model.ErrInvalidRequest):
				c.JSON(http.StatusBadRequest, ginpkg.H{
					"error": err.Error(),
				})
				return
			case errors.Is(err, model.ErrAccountNotFound):
				c.JSON(http.StatusNotFound, ginpkg.H{
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

		c.JSON(http.StatusOK, ginpkg.H{
			"data": account,
		})
	}
}
