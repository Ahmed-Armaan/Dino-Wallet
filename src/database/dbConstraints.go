package database

import (
	"gorm.io/gorm"
)

func SetConstraints(db *gorm.DB) error {
	if err := db.Exec(`
		ALTER TABLE accounts
		DROP CONSTRAINT IF EXISTS balance_non_negative;

		ALTER TABLE accounts
		ADD CONSTRAINT balance_non_negative
		CHECK (allow_negative OR balance >= 0);
	`).Error; err != nil {
		return err
	}

	return nil
}
