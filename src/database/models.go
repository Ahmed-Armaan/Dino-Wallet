package database

import (
	"time"

	"github.com/google/uuid"
)

type AccountType string
type AssetType string

const (
	AccountNormal   AccountType = "normal"
	AccountTreasury AccountType = "treasury"
	AccountRevenue  AccountType = "revenue"
)

const (
	AssetGold  AssetType = "gold"
	AssetGem   AssetType = "gem"
	AssetCoins AssetType = "coins"
)

type Asset struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Code      AssetType `gorm:"unique;not null"`
	CreatedAt time.Time
}

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserName  string    `gorm:"unique;not null"`
	CreatedAt time.Time
}

type Account struct {
	ID            uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID        *uuid.UUID  `gorm:"type:uuid"`
	AssetID       uuid.UUID   `gorm:"type:uuid;not null"`
	Type          AccountType `gorm:"type:text;not null"`
	Balance       int64       `gorm:"not null;check:allow_negative OR balance >= 0"`
	AllowNegative bool        `gorm:"not null"`
	CreatedAt     time.Time
}

type LedgerTransaction struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	IdempotencyKey uuid.UUID `gorm:"type:uuid;unique;not null"`
	CreatedAt      time.Time
}

type LedgerEntry struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TransactionID uuid.UUID `gorm:"type:uuid;not null"`
	AccountID     uuid.UUID `gorm:"type:uuid;not null"`
	Amount        int64     `gorm:"not null"`
	CreatedAt     time.Time
}
