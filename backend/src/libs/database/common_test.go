package database

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//Setup and teardown
func TestMain(m *testing.M) {
	//setup
	SetupBankGormDB()
	code := m.Run() //run tests
	os.Exit(code)
}
func TestSetConnectionMaxLifetime(t *testing.T) {
	t.Run("Should success on", func(t *testing.T) {
		t.Run("Without duration, so it's never closed", func(t *testing.T) {
			gormDB := openBankBD()

			//Action
			setConnectionMaxLifetime(gormDB, 0)

			//Data Assertion
			sqlDB, _ := gormDB.DB()
			assert.Equal(t, int64(0), sqlDB.Stats().MaxLifetimeClosed)
		})
		t.Run("A very small time, so it's closed after that time", func(t *testing.T) {
			gormDB := openBankBD()

			//Action
			setConnectionMaxLifetime(gormDB, time.Millisecond)

			//Data Assertion
			sqlDB, _ := gormDB.DB()
			assert.Equal(t, int64(0), sqlDB.Stats().MaxLifetimeClosed, "None should be dead, because time has not passed")

			//Channel to wait for the death of the connection
			c := make(chan int64, 1)
			go func() {
				for {
					time.Sleep(time.Second / 10) //Try every 1/10 seconds
					if count := sqlDB.Stats().MaxLifetimeClosed; count > 0 {
						c <- count
						break
					}
				}
			}()
			select {
			case count := <-c:
				assert.Equal(t, int64(1), count, "There should be only one dead connection")
			case <-time.After(3 * time.Second):
				t.Error("Test timeout")
			}
		})
	})
	t.Run("Should fail on", func(t *testing.T) {
		fakeConnection := &gorm.DB{
			Config: &gorm.Config{
				ConnPool: &gorm.PreparedStmtDB{},
			},
		}
		loggerBuffer := &bytes.Buffer{}
		logrus.SetOutput(loggerBuffer) //To catch logs

		//Action
		setConnectionMaxLifetime(fakeConnection, 0)

		//Data Assertion
		assert.Contains(t, loggerBuffer.String(), "Error setting connection max lifetime")

		t.Cleanup(func() {
			logrus.SetOutput(logrus.StandardLogger().Out)
		})
	})
}

func openBankBD() *gorm.DB {
	db, _ := gorm.Open(postgres.Open(CreateBankConnectionString()), nil)
	return db
}
