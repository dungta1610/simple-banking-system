package gin

import (
	"errors"
	"net/http"

	ginpkg "github.com/gin-gonic/gin"

	"simple-banking-system/module/account/biz"
	"simple-banking-system/module/account/model"
)

func ListAccountsHandler(store biz.ListAccountsStore) ginpkg.HandlerFunc {
	return func(c *ginpkg.Context) {
		var query model.ListAccountsQuery

		if err := c.ShouldBindQuery(&query); err != nil {
			c.JSON(http.StatusBadRequest, ginpkg.H{
				"error": "invalid query params",
			})
			return
		}

		listAccountsBiz := biz.NewListAccountsBiz(store)
		accounts, err := listAccountsBiz.ListAccounts(c.Request.Context(), &query)

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

		c.JSON(http.StatusOK, ginpkg.H{
			"data": accounts,
		})
	}
}
