package handlers

import (
	"net/http"

	"github.com/Ahmed-Armaan/Dino-Wallet.git/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BonusReq struct {
	User           string    `json:"user"`
	IdempotencyKey uuid.UUID `json:"idempotency_key"`
}

func Bonus(db database.DataBaseStore) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req BonusReq

		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}

		if req.User == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "user is required",
			})
			return
		}

		if req.IdempotencyKey == uuid.Nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "idempotency_key is required",
			})
			return
		}

		if err := db.GiveRandomBonus(req.User, req.IdempotencyKey); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "bonus processed",
		})
	}
}
