package handlers

import (
	"net/http"

	"github.com/Ahmed-Armaan/Dino-Wallet.git/database"
	"github.com/gin-gonic/gin"
)

func GetBalance(db database.DataBaseStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		userName := c.Query("user")

		balances, err := db.Balance(userName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "database error",
			})
			return
		}

		c.JSON(200, gin.H{
			"balance": balances.Balances,
		})
	}
}
