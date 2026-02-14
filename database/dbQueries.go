package database

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssetAmount struct {
	Asset   AssetType `json:"asset"`
	Balance int64     `json:"balance"`
}

type UserBalance struct {
	Balances []AssetAmount `json:"balances"`
}

type LedgerItemRow struct {
	LedgerEntry `gorm:"embedded"`
	Account     AccountType `json:"account"`
	User        *string     `json:"user"`
}

func (db *DataBaseHolder) DbTopUp(username string, assetType AssetType, amount int, idempotencyKey uuid.UUID) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}
	const maxRetries = 5

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := db.db.Transaction(func(tx *gorm.DB) error {
			var user User
			if err := tx.
				Where("user_name = ?", username).
				First(&user).Error; err != nil {
				return err
			}

			if err := tx.Exec(`
				SELECT pg_advisory_xact_lock(hashtext(?))
			`, user.ID.String()).Error; err != nil {
				return err
			}

			var asset Asset
			if err := tx.
				Where("code = ?", assetType).
				First(&asset).Error; err != nil {
				return err
			}

			var userAccount Account
			if err := tx.
				Where("user_id = ? AND asset_id = ? AND type = ?",
					user.ID, asset.ID, AccountNormal).
				First(&userAccount).Error; err != nil {
				return err
			}

			var treasuryAccount Account
			if err := tx.
				Where("user_id IS NULL AND asset_id = ? AND type = ?",
					asset.ID, AccountTreasury).
				First(&treasuryAccount).Error; err != nil {
				return err
			}

			ledgerTx := LedgerTransaction{
				IdempotencyKey: idempotencyKey,
			}
			if err := tx.Create(&ledgerTx).Error; err != nil {
				return err
			}

			entries := []LedgerEntry{
				{
					TransactionID: ledgerTx.ID,
					AccountID:     userAccount.ID,
					Amount:        int64(amount),
				},
				{
					TransactionID: ledgerTx.ID,
					AccountID:     treasuryAccount.ID,
					Amount:        -int64(amount),
				},
			}

			if err := tx.Create(&entries).Error; err != nil {
				return err
			}

			if err := tx.Model(&Account{}).
				Where("id = ?", userAccount.ID).
				Update("balance", gorm.Expr("balance + ?", amount)).
				Error; err != nil {
				return err
			}

			if err := tx.Model(&Account{}).
				Where("id = ?", treasuryAccount.ID).
				Update("balance", gorm.Expr("balance - ?", amount)).
				Error; err != nil {
				return err
			}

			return nil
		})

		if err == nil {
			return nil
		}
		time.Sleep(time.Duration(attempt) * 50 * time.Millisecond)
	}

	return errors.New("max retries exceeded")
}

func (db *DataBaseHolder) GiveRandomBonus(username string, idempotencyKey uuid.UUID) error {
	const maxRetries = 5

	assets := []AssetType{
		AssetGold,
		AssetGem,
		AssetCoins,
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomAsset := assets[rng.Intn(len(assets))]
	randomAmount := rng.Intn(51)

	if randomAmount == 0 {
		return nil
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := db.db.Transaction(func(tx *gorm.DB) error {
			var user User
			if err := tx.
				Where("user_name = ?", username).
				First(&user).Error; err != nil {
				fmt.Println(err)
				return err
			}

			if err := tx.Exec(`
				SELECT pg_advisory_xact_lock(hashtext(?))
			`, user.ID.String()).Error; err != nil {
				fmt.Println(err)
				return err
			}

			var asset Asset
			if err := tx.
				Where("code = ?", randomAsset).
				First(&asset).Error; err != nil {
				fmt.Println(err)
				return err
			}

			var userAccount Account
			if err := tx.
				Where("user_id = ? AND asset_id = ? AND type = ?",
					user.ID, asset.ID, AccountNormal).
				First(&userAccount).Error; err != nil {
				fmt.Println(err)
				return err
			}

			var treasuryAccount Account
			if err := tx.
				Where("user_id IS NULL AND asset_id = ? AND type = ?",
					asset.ID, AccountTreasury).
				First(&treasuryAccount).Error; err != nil {
				fmt.Println(err)
				return err
			}

			ledgerTx := LedgerTransaction{
				IdempotencyKey: idempotencyKey,
			}
			if err := tx.Create(&ledgerTx).Error; err != nil {
				fmt.Println(err)
				return err
			}

			entries := []LedgerEntry{
				{
					TransactionID: ledgerTx.ID,
					AccountID:     userAccount.ID,
					Amount:        int64(randomAmount),
				},
				{
					TransactionID: ledgerTx.ID,
					AccountID:     treasuryAccount.ID,
					Amount:        -int64(randomAmount),
				},
			}

			if err := tx.Create(&entries).Error; err != nil {
				fmt.Println(err)
				return err
			}

			if err := tx.Model(&Account{}).
				Where("id = ?", userAccount.ID).
				Update("balance", gorm.Expr("balance + ?", randomAmount)).
				Error; err != nil {
				fmt.Println(err)
				return err
			}

			if err := tx.Model(&Account{}).
				Where("id = ?", treasuryAccount.ID).
				Update("balance", gorm.Expr("balance - ?", randomAmount)).
				Error; err != nil {
				fmt.Println(err)
				return err
			}

			return nil
		})

		if err == nil {
			return nil
		}

		time.Sleep(time.Duration(attempt) * 50 * time.Millisecond)
	}

	return errors.New("max retries exceeded")
}

func (db *DataBaseHolder) Purchase(username string, assetType AssetType, amount int, idempotencyKey uuid.UUID) error {
	const maxRetries = 5
	if amount <= 0 {
		return errors.New("invalid amount")
	}
	ErrInsufficientBal := errors.New("insufficient balance")

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := db.db.Transaction(func(tx *gorm.DB) error {
			var user User
			if err := tx.
				Where("user_name = ?", username).
				First(&user).Error; err != nil {
				return err
			}

			if err := tx.Exec(`
				SELECT pg_advisory_xact_lock(hashtext(?))
			`, user.ID.String()).Error; err != nil {
				return err
			}

			var asset Asset
			if err := tx.
				Where("code = ?", assetType).
				First(&asset).Error; err != nil {
				return err
			}

			var userAccount Account
			if err := tx.
				Where("user_id = ? AND asset_id = ? AND type = ?",
					user.ID, asset.ID, AccountNormal).
				First(&userAccount).Error; err != nil {
				return err
			}

			var treasuryAccount Account
			if err := tx.
				Where("user_id IS NULL AND asset_id = ? AND type = ?",
					asset.ID, AccountTreasury).
				First(&treasuryAccount).Error; err != nil {
				return err
			}

			if !userAccount.AllowNegative && userAccount.Balance < int64(amount) {
				return ErrInsufficientBal
			}

			ledgerTx := LedgerTransaction{
				IdempotencyKey: idempotencyKey,
			}
			if err := tx.Create(&ledgerTx).Error; err != nil {
				return err
			}

			entries := []LedgerEntry{
				{
					TransactionID: ledgerTx.ID,
					AccountID:     userAccount.ID,
					Amount:        -int64(amount),
				},
				{
					TransactionID: ledgerTx.ID,
					AccountID:     treasuryAccount.ID,
					Amount:        int64(amount),
				},
			}

			if err := tx.Create(&entries).Error; err != nil {
				return err
			}

			if err := tx.Model(&Account{}).
				Where("id = ?", userAccount.ID).
				Update("balance", gorm.Expr("balance - ?", amount)).
				Error; err != nil {
				return err
			}

			if err := tx.Model(&Account{}).
				Where("id = ?", treasuryAccount.ID).
				Update("balance", gorm.Expr("balance + ?", amount)).
				Error; err != nil {
				return err
			}

			return nil
		})

		if err == nil || errors.Is(err, ErrInsufficientBal) {
			return err
		}
		time.Sleep(time.Duration(attempt) * 50 * time.Millisecond)
	}

	return errors.New("max retries exceeded")
}

func (db *DataBaseHolder) Ledger(pageToken int, pageSize int) ([]LedgerItemRow, error) {
	ledger := []LedgerItemRow{}

	if err := db.db.Model(&LedgerEntry{}).
		Select("ledger_entries.*, accounts.type AS account, users.user_name AS user").
		Joins("JOIN accounts ON accounts.id = ledger_entries.account_id").
		Joins("LEFT JOIN users ON accounts.user_id = users.id").
		Order("ledger_entries.created_at DESC").
		Offset(pageToken * pageSize).
		Limit(pageSize).
		Scan(&ledger).Error; err != nil {
		return nil, err
	}

	return ledger, nil
}

func (db *DataBaseHolder) Balance(username string) (*UserBalance, error) {
	var user User
	if err := db.db.Model(&User{}).
		Select("id").
		Where("user_name = ?", username).
		First(&user).Error; err != nil {
		return nil, err
	}

	var assetAmounts []AssetAmount
	if err := db.db.Model(&Account{}).
		Select("assets.code AS asset, accounts.balance as balance").
		Joins("JOIN assets ON assets.id = accounts.asset_id").
		Where("accounts.user_id = ?", user.ID).
		Find(&assetAmounts).Error; err != nil {
		return nil, err
	}

	balance := UserBalance{
		Balances: assetAmounts,
	}

	return &balance, nil
}

func (db *DataBaseHolder) Seed(seedSql string) error {
	if err := db.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(seedSql).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
