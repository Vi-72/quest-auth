package casesteps

import (
	"context"

	"gorm.io/gorm"
)

// CleanupAuthTables truncates auth-related tables before tests
func CleanupAuthTables(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec("TRUNCATE TABLE events CASCADE").Error; err != nil {
		return err
	}
	if err := db.WithContext(ctx).Exec("TRUNCATE TABLE users CASCADE").Error; err != nil {
		return err
	}
	return nil
}
