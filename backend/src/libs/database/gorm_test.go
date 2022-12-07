package database

import (
	"bank-service/src/libs/env"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func resetOnceBank() {
	once = sync.Once{}
}

func TestSetupBankGormDB(t *testing.T) {
	t.Run("Should success on", func(t *testing.T) {
		t.Run("Should success on", func(t *testing.T) {
			resetOnceBank()
			db := SetupBankGormDB()
			sqlDB, _ := db.DB()
			errPing := sqlDB.Ping()
			//Data Assertion
			assert.NotNil(t, db)
			assert.NoError(t, errPing)
		})
		t.Run("Wait for postgres", func(t *testing.T) {
			// Smaller time & wrong DB name
			oldDelta := env.BankServiceSecondsBetweenAttempts
			env.BankServiceSecondsBetweenAttempts = time.Second / 2
			oldValue := env.BankServicePostgresqlNameTest
			env.BankServicePostgresqlNameTest = "Bank_SERVICE_POSTGRESQL_NAME_not_found"
			var db *gorm.DB
			var errPing error
			wait := make(chan bool)
			go func() {
				resetOnceBank()
				db = SetupBankGormDB()
				sqlDB, _ := db.DB()
				errPing = sqlDB.Ping()
				wait <- true
			}()
			time.Sleep(env.BankServiceSecondsBetweenAttempts)
			env.BankServicePostgresqlNameTest = oldValue
			<-wait

			//Data Assertion
			assert.NotNil(t, db)
			assert.NoError(t, errPing)
			t.Cleanup(func() {
				env.BankServicePostgresqlNameTest = oldValue
				env.BankServiceSecondsBetweenAttempts = oldDelta
			})
		})
	})
}

func TestGetBankGormConnection(t *testing.T) {
	t.Run("Should success when the connection is already open", func(t *testing.T) {
		resetOnceBank()
		db := SetupBankGormDB()
		dbSingleton := GetBankGormConnection()
		sqlDB, _ := dbSingleton.DB()
		errPing := sqlDB.Ping()
		//Data Assertion
		assert.Equal(t, db, dbSingleton)
		assert.NoError(t, errPing)
	})
}
