package validator

import (
	"strings"

	"bank-service/src/libs/errors"
	"bank-service/src/utils/helpers"

	"github.com/go-playground/validator"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

//init is only called once (implicit)
func init() {
	validate = validator.New()
}

//ValidateStruct receives an obj from struct that uses validator tags
func ValidateStruct(obj interface{}) error {
	if err := validate.Struct(obj); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			err := validationErrors[0]
			return createCustomError(err.Field(), err.Tag(), err.Param())
		}
		return err
	}

	return nil
}

//ValidateVar receives an obj from struct that uses validator tags
func ValidateVar(value interface{}, name, tag string) error {
	if err := validate.Var(value, tag); err != nil {
		err := err.(validator.ValidationErrors)[0]
		return createCustomError(name, err.Tag(), err.Param())
	}
	return nil
}

//createCustomError creates a new ErrFieldValidation error using received values
func createCustomError(field, validation, valid string) error {
	return errors.ErrFieldValidation(field, validation, valid)
}

/*
ValidateFieldIsOneOf receives a field name, the string value and a list to find
*/
func ValidateFieldIsOneOf(fieldName, value string, list []string) error {
	if !helpers.StringInSlice(value, list) {
		return createCustomError(fieldName, "One of", strings.Join(list, ", "))
	}

	return nil
}

/*
ValidateSliceValuesAreOneOf validates that each element in values is in list
*/
func ValidateSliceValuesAreOneOf(fieldName string, values, list []string) error {
	for _, val := range values {
		if err := ValidateFieldIsOneOf(fieldName, val, list); err != nil {
			return err
		}
	}
	return nil
}
