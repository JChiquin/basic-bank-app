package validator

import (
	"reflect"
	"strings"
	"testing"

	"bank-service/src/libs/errors"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"
)

func TestValidateStruct(t *testing.T) {
	//Fixtures
	type foo struct {
		Name      string      `json:"name" validate:"required"`
		SomeValue interface{} `json:"some_value" validate:"gt=1"`
	}
	validName := "test"
	validSomeValue := 10

	testCases := []struct {
		TestName string
		Input    interface{}
		Expected error
	}{
		{
			TestName: "Pass all validations",
			Input: foo{
				Name:      validName,
				SomeValue: validSomeValue,
			},
			Expected: nil,
		},
		{
			TestName: "Fail required validation",
			Input: foo{
				Name:      "",
				SomeValue: validSomeValue,
			},
			Expected: errors.ErrFieldValidation("Name", "required", ""),
		},
		{
			TestName: "Fail greater validation",
			Input: foo{
				Name:      validName,
				SomeValue: 0,
			},
			Expected: errors.ErrFieldValidation("SomeValue", "gt", "1"),
		},
		{
			TestName: "Using something that's not a struct",
			Input:    validName, //string value
			Expected: &validator.InvalidValidationError{Type: reflect.TypeOf(validName)},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.TestName, func(t *testing.T) {
			err := ValidateStruct(tC.Input)
			assert.Equal(t, tC.Expected, err)
		})
	}

	t.Run("Should fail on InvalidValidationError", func(t *testing.T) {
		//Fixtures
		err := ValidateStruct(validName)

		assert.Error(t, err)
	})
}

func TestValidateFieldIsOneOf(t *testing.T) {
	//Fixtures
	fieldName := "name"
	list := []string{"Jorge", "Chiquin", "Valderrama"}
	t.Run("Should success on value is in list", func(t *testing.T) {
		err := ValidateFieldIsOneOf(fieldName, list[0], list)
		assert.NoError(t, err)
	})
	t.Run("Should fail on value is not in list", func(t *testing.T) {
		err := ValidateFieldIsOneOf(fieldName, "not found value", list)

		expectedErr := errors.ErrFieldValidation(fieldName, "One of", strings.Join(list, ", "))
		assert.Equal(t, expectedErr, err)
	})
}

func TestValidateSliceValuesAreOneOf(t *testing.T) {
	fieldName := "name"
	list := []string{"name", "boke", "dib"}
	t.Run("Should success on all values in list", func(t *testing.T) {
		err := ValidateSliceValuesAreOneOf(fieldName, list, list)
		assert.NoError(t, err)
	})
	t.Run("Should fail on not all values in list", func(t *testing.T) {
		err := ValidateSliceValuesAreOneOf(fieldName, []string{"name", "f", "boke"}, list)
		assert.Error(t, err)
	})
}

func TestValidateVar(t *testing.T) {
	//Fixtures
	name := "some_field_name"
	t.Run("Should success on", func(t *testing.T) {
		testCases := []struct {
			TestName string
			Value    interface{}
			Tag      string
		}{
			{
				TestName: "Valid number",
				Value:    10,
				Tag:      "required,gt=1,lt=20",
			},
			{
				TestName: "Valid string",
				Value:    "some value",
				Tag:      "required,len=10",
			},
			{
				TestName: "Valid UUID4",
				Value:    "9b65d3ae-9e5b-439b-8a1b-9dd577b3e1f9",
				Tag:      "uuid4",
			},
		}
		for _, tC := range testCases {
			t.Run(tC.TestName, func(t *testing.T) {

				//Action
				err := ValidateVar(tC.Value, name, tC.Tag)

				//Data Assertion
				assert.NoError(t, err)
			})
		}
	})
	t.Run("Should fail on", func(t *testing.T) {
		testCases := []struct {
			TestName    string
			Value       interface{}
			Tag         string
			ExpectedErr error
		}{
			{
				TestName:    "Invalid number",
				Value:       20,
				Tag:         "required,gt=1,lt=20",
				ExpectedErr: errors.ErrFieldValidation(name, "lt", "20"),
			},
			{
				TestName:    "Invalid string",
				Value:       "hello",
				Tag:         "required,len=10",
				ExpectedErr: errors.ErrFieldValidation(name, "len", "10"),
			},
			{
				TestName:    "Invalid UUID4",
				Value:       "319af744-3ef4-11eb-b378-0242ac130002", //This is a UUID1
				Tag:         "required,uuid4",
				ExpectedErr: errors.ErrFieldValidation(name, "uuid4", ""),
			},
		}
		for _, tC := range testCases {
			t.Run(tC.TestName, func(t *testing.T) {

				//Action
				err := ValidateVar(tC.Value, name, tC.Tag)

				//Data Assertion
				assert.Equal(t, tC.ExpectedErr, err)
			})
		}
	})
}
