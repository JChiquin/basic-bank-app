package env

import (
	"bank-service/src/utils/constant"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	t.Run("updateEnvVars", func(t *testing.T) {
		// Fixture
		AppEnv = "develop"
		temp := constant.IntervalBetweenEnvVarUpdate
		constant.IntervalBetweenEnvVarUpdate = time.Second * 1
		oldPort := BankServiceRestPort
		newPort := "9090"

		// Action
		setupUpdateEnvVars()

		// Check it didn't change the first time
		assert.Equal(t, oldPort, BankServiceRestPort)

		os.Setenv("BANK_SERVICE_REST_PORT", newPort)
		time.Sleep(constant.IntervalBetweenEnvVarUpdate * 2) // Wait for function to run again

		// Data assertion
		assert.Equal(t, newPort, BankServiceRestPort)

		t.Cleanup(func() {
			constant.IntervalBetweenEnvVarUpdate = temp
			BankServiceRestPort = oldPort
			AppEnv = "testing"
		})
	})
}
