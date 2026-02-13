package handlers

import (
	"net/http"
	"strconv"

	"github.com/Ahmed-Armaan/Dino-Wallet.git/database"
	"github.com/gin-gonic/gin"
)

func GetLedger(db database.DataBaseStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageSize := c.Query("pageSize")
		pageToken := c.Query("pageToken")

		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "pageSize must be integer",
			})
			return
		}

		pageTokenInt, err := strconv.Atoi(pageToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "pageToken must be integer",
			})
			return
		}

		ledger, err := db.Ledger(pageTokenInt, pageSizeInt)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "database error",
			})
			return
		}

		c.JSON(200, ledger)
	}
}
