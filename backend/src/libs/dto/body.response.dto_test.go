package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBodyResponse(t *testing.T) {
	t.Run("Should Succeed", func(t *testing.T) {
		t.Run("New Body", func(t *testing.T) {
			// Fixture
			message := "mock message"
			errors := []map[string]string{{"error": "some error"}}
			data := "some data"

			// Run Foo
			got := NewBodyResponse(message, errors, data)

			// Data assertion
			assert.Equal(t, message, got.Message)
			assert.Equal(t, errors, got.Errors)
			assert.Equal(t, data, got.Data)
		})
	})
}
