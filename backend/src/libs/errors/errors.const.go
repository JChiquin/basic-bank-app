package errors

import (
	"net/http"

	"bank-service/src/libs/i18n"
)

var (
	//ErrInternalServer indicates an internal server error, this is used a default error when it's a general error
	ErrInternalServer = NewMyError(http.StatusInternalServerError, i18n.Message{MessageID: "ERRORS.INTERNAL_SERVER"})

	//ErrPageTooHigh indicates that page param is higher than the valid number
	ErrPageTooHigh = NewMyError(http.StatusBadRequest, i18n.Message{MessageID: "ERRORS.PAGE_TOO_LARGE"})

	//ErrPageSizeTooHigh indicates that pageSize param is higher than the valid number
	ErrPageSizeTooHigh = NewMyError(http.StatusBadRequest, i18n.Message{MessageID: "ERRORS.PAGE_SIZE_TOO_LARGE"})

	//ErrInvalidRequestBody indicates that the request body has an invalid syntax
	ErrInvalidRequestBody = NewMyError(http.StatusBadRequest, i18n.Message{MessageID: "ERRORS.INVALID_REQUEST_BODY"})

	//ErrNotFound indicates an entity not found error
	ErrNotFound = NewMyError(http.StatusNotFound, i18n.Message{MessageID: "ERRORS.NOT_FOUND"})

	//ErrUserExists indicates the user al ready exists
	ErrUserExists = NewMyError(http.StatusConflict, i18n.Message{MessageID: "ERRORS.USER_DUPLICATED"})

	//ErrUnauthorized indicates that client doesn't pass the JWT middleware
	ErrUnauthorized = NewMyError(http.StatusUnauthorized, i18n.Message{MessageID: "ERRORS.UNAUTHORIZED"})

	//ErrURLNotFound indicates a requested URL doesn't exist
	ErrURLNotFound = NewMyError(http.StatusNotFound, i18n.Message{MessageID: "ERRORS.URL_NOT_FOUND"})

	//ErrOrderNotApprovable indicates the current order cannot be approved
	ErrOrderNotApprovable = NewMyError(http.StatusConflict, i18n.Message{MessageID: "ERRORS.ORDER_NOT_APPROVABLE"})
)

//Private errors
var (
	errUnsupportedFieldValue = NewMyError(http.StatusBadRequest, i18n.Message{MessageID: "ERRORS.UNSUPPORTED_FIELD_VALUE"})
	errFieldValidation       = NewMyError(http.StatusBadRequest, i18n.Message{MessageID: "ERRORS.FIELD_VALIDATION"})
)
