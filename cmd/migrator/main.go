package main

import (
	_ "embed"
	"fmt"
	"log"
	"time"

	"github.com/Ahmed-Armaan/Dino-Wallet.git/database"
	"github.com/joho/godotenv"
)

//go:embed seed.sql
var seedSql string

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env loaded")
	}

	var db database.DataBaseStore
	var err error

	// Try to connect to database for 2 minutes
	for attempt := range 60 {
		if db, err = database.DbInit(); err != nil {
			fmt.Printf("Migration attemp %d Failed\n", attempt+1)
		} else if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	fmt.Println("Migration Successful")

	if err = db.Seed(seedSql); err != nil {
		log.Fatalf("database seed failed: %s\n", err)
	}
	fmt.Println("Seed Successful")
}
