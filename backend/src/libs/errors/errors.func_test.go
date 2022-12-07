package errors

import (
	"errors"
	"net/http"
	"testing"

	"bank-service/src/libs/i18n"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	t.Run("ErrFieldValidation", func(t *testing.T) {
		// Fixtures
		field := "field"
		validation := "validation"
		valid := "valid"
		template := map[string]interface{}{
			"Field":      field,
			"Validation": validation,
			"Valid":      valid,
		}

		// Action
		errAddingTemplate := errFieldValidation.WithTemplate(template)
		errWithTemplate := ErrFieldValidation(field, validation, valid)

		// Assert data
		assert.EqualError(t, errWithTemplate, errAddingTemplate.Error())
	})
}

func TestGetStatusCode(t *testing.T) {
	t.Run("Using my error", func(t *testing.T) {
		err := NewMyError(999, i18n.Message{MessageID: "myError instance"})
		status := GetStatusCode(err)
		assert.Equal(t, 999, status)
	})
	t.Run("Using an external error", func(t *testing.T) {
		err := errors.New("This error is not a myError instance")
		status := GetStatusCode(err)
		assert.Equal(t, http.StatusInternalServerError, status)
	})
}

func TestGetErrorMessage(t *testing.T) {
	t.Run("Using my error", func(t *testing.T) {
		err := ErrInternalServer
		message := GetErrorMessage(err)
		assert.Equal(t, ErrInternalServer.Error(), message)
	})
	t.Run("Using an external error", func(t *testing.T) {
		err := errors.New("This error is not a myError instance")
		message := GetErrorMessage(err)
		assert.Equal(t, ErrInternalServer.Error(), message)
	})
}

func TestGetAction(t *testing.T) {
	t.Run("Using my error", func(t *testing.T) {
		err := NewMyError(
			http.StatusForbidden,
			i18n.Message{MessageID: "ERRORS.MISSING_POSTBOARDING"},
		).SetAction("Some action")
		action := GetAction(err)
		assert.Equal(t, "Some action", *action)
	})
	t.Run("Using an external error", func(t *testing.T) {
		err := errors.New("This error is not a myError instance")
		action := GetAction(err)
		assert.Nil(t, action)
	})
}

func TestGetKeyBody(t *testing.T) {
	t.Run("Using my error", func(t *testing.T) {
		keyBody := "custom_key"
		err := MyError{
			keyBody: keyBody,
		}
		got := GetKeyBody(err)
		assert.Equal(t, keyBody, got)
	})
	t.Run("Using an external error", func(t *testing.T) {
		err := errors.New("This error is not a myError instance")
		got := GetKeyBody(err)
		assert.Equal(t, defaultKeyBody, got)
	})
}
