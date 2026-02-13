package database

import (
	"os"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DataBaseHolder struct {
	db *gorm.DB
}

type DataBaseStore interface {
	DbTopUp(username string, assetType AssetType, amount int, idempotencyKey uuid.UUID) error
	GiveRandomBonus(username string, idempotencyKey uuid.UUID) error
	Purchase(username string, assetType AssetType, amount int, idempotencyKey uuid.UUID) error
	Balance(username string) (*UserBalance, error)
	Ledger(pageToken int, pageSize int) ([]LedgerItemRow, error)
}

func DbInit() (DataBaseStore, error) {
	dsn := os.Getenv("DATABASE_URL")
	db := &DataBaseHolder{}

	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	DB.AutoMigrate(&Asset{}, &User{}, &Account{}, &LedgerTransaction{}, &LedgerEntry{})

	if err := SetConstraints(DB); err != nil {
		return db, err
	}

	db.db = DB
	return db, nil
}
