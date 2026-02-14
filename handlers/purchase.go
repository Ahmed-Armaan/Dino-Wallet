package handlers

import (
	"net/http"

	"github.com/Ahmed-Armaan/Dino-Wallet.git/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PurchaseReq struct {
	User           string             `json:"user"`
	Asset          database.AssetType `json:"asset"`
	Amount         int                `json:"amount"`
	IdempotencyKey uuid.UUID          `json:"idempotency_key"`
}

func Purchase(db database.DataBaseStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req PurchaseReq

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

		if req.Amount <= 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid amount",
			})
			return
		}

		if req.IdempotencyKey == uuid.Nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "idempotency_key is required",
			})
			return
		}

		if err := db.Purchase(req.User, req.Asset, req.Amount, req.IdempotencyKey); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "purchase successful",
		})
	}
}
