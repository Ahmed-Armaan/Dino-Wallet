package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Ahmed-Armaan/Dino-Wallet.git/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TopUpReq struct {
	User           string             `json:"user"`
	Asset          database.AssetType `json:"asset"`
	Amount         int                `json:"amount"`
	IdempotencyKey uuid.UUID          `json:"idempotency_key"`
}

func Topup(db database.DataBaseStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := c.Request.Body
		defer c.Request.Body.Close()

		body, err := io.ReadAll(data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "failed to read body",
			})
			return
		}

		topUpReq := TopUpReq{}
		if err := json.Unmarshal(body, &topUpReq); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "failed to read body",
			})
			return
		}

		if err := db.DbTopUp(topUpReq.User, topUpReq.Asset, topUpReq.Amount, topUpReq.IdempotencyKey); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Transaction Failed",
			})
			return
		}

		c.Status(200)
	}
}
