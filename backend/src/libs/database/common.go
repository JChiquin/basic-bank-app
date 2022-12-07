package database

import (
	"time"

	"bank-service/src/libs/logger"

	"gorm.io/gorm"
)

/*
setConnectionMaxLifetime receives a gorm connection, gets its generic database and
sets the connection max lifetime
*/
func setConnectionMaxLifetime(db *gorm.DB, duration time.Duration) {
	// Get generic database object sql.DB to use its functions
	sqlDB, err := db.DB()
	if err != nil {
		logger.GetInstance().Errorf("Error setting connection max lifetime: %s", err)
		return
	}

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(duration)
}
