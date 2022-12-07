package errors

import (
	"errors"
	"net/http"
)

/*
GetStatusCode checks if it's a MyError and returns its status code
if not, returns StatusInternalServerError
*/
func GetStatusCode(err error) int {
	if errors.As(err, &MyError{}) {
		return err.(MyError).statusCode
	}
	return http.StatusInternalServerError
}

/*
GetErrorMessage checks if it's a MyError and returns its error method
if not, returns a generic internal server error
*/
func GetErrorMessage(err error) string {
	if errors.As(err, &MyError{}) {
		return err.(MyError).Error()
	}
	return ErrInternalServer.Error()
}

/*
GetAction checks if it's a MyError and returns its action
if not, returns nil
*/
func GetAction(err error) *string {
	if errors.As(err, &MyError{}) {
		return err.(MyError).action
	}
	return nil
}

/*
GetKeyBody checks if it's a MyError and returns its key body
if not, returns a default key body
*/
func GetKeyBody(err error) string {
	if errors.As(err, &MyError{}) {
		return err.(MyError).keyBody
	}
	return defaultKeyBody
}

//ErrFieldValidation indicates a field validation error
func ErrFieldValidation(field, validation, valid string) error {
	return errFieldValidation.WithTemplate(
		map[string]interface{}{
			"Field":      field,
			"Validation": validation,
			"Valid":      valid,
		})
}

//ErrFieldValidation returns a errFieldValidation with custom template
func ErrUnsupportedFieldValue(field string) error {
	return errUnsupportedFieldValue.WithTemplate(
		map[string]interface{}{
			"Field": field,
		})
}
