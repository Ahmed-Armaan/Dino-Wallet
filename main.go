package main

import (
	"fmt"
	"log"

	"github.com/Ahmed-Armaan/Dino-Wallet.git/database"
	"github.com/Ahmed-Armaan/Dino-Wallet.git/server"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Could not load ENV, check error if not in production")
	}

	db, err := database.DbInit()
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	server.Serve(db)
}
