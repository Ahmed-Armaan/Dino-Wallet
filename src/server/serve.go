package server

import (
	"log"

	"github.com/Ahmed-Armaan/Dino-Wallet.git/database"
	"github.com/Ahmed-Armaan/Dino-Wallet.git/handlers"
	"github.com/gin-gonic/gin"
)

func Serve(db database.DataBaseStore) {
	r := gin.Default()

	r.POST("/top_up", handlers.Topup(db))
	r.POST("/bonus", handlers.Bonus(db))
	r.POST("/purchase", handlers.Purchase(db))
	r.GET("/balance", handlers.GetBalance(db))
	r.GET("/ledger", handlers.GetLedger(db))

	if err := r.Run(); err != nil {
		log.Fatal("Failed to start server")
	}
}
