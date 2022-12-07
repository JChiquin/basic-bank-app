package errors

import (
	"errors"
	"testing"

	"bank-service/src/libs/i18n"

	"github.com/stretchr/testify/assert"
)

func TestIs(t *testing.T) {
	testCases := []struct {
		TestName string
		Target   error
		Source   error
		Expected bool
	}{
		{
			TestName: "Same error",
			Source:   ErrInternalServer,
			Target:   ErrInternalServer,
			Expected: true,
		},
		{
			TestName: "Target is MyError instance, but not the same source",
			Source:   ErrInternalServer,
			Target:   errFieldValidation,
			Expected: false,
		},
		{
			TestName: "Other general error",
			Source:   ErrInternalServer,
			Target:   errors.New("other error"),
			Expected: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.TestName, func(t *testing.T) {
			assert.Equal(t, tC.Expected, errors.Is(tC.Source, tC.Target))
		})
	}
}

func TestWithTemplate(t *testing.T) {
	testCases := []struct {
		TestName string
		Source   MyError
		Template map[string]interface{}
	}{
		{
			TestName: "Without",
			Source:   NewMyError(100, i18n.Message{MessageID: "FOO.BAR"}),
		},
		{
			TestName: "With one template",
			Template: map[string]interface{}{
				"Field": "Value to interpolate",
			},
			Source: NewMyError(100, i18n.Message{MessageID: "FOO.BAR"}),
		},
		{
			TestName: "With two template",
			Template: map[string]interface{}{
				"Field":  "Value to interpolate",
				"Field2": "Value2 to interpolate",
			},
			Source: NewMyError(100, i18n.Message{MessageID: "FOO.BAR"}),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.TestName, func(t *testing.T) {
			resultError := tC.Source.WithTemplate(tC.Template)
			assert.Equal(t, tC.Template, resultError.(MyError).data)
			if tC.Template != nil {
				assert.NotEqual(t, tC.Template, tC.Source.data)
				assert.NotEqual(t, tC.Source, resultError)
			}
		})
	}
}

func TestGetData(t *testing.T) {
	template := map[string]interface{}{
		"Field":  "Value to interpolate",
		"Field2": "Value2 to interpolate",
	}
	testCases := []struct {
		TestName string
		Source   MyError
		Expected map[string]interface{}
	}{
		{
			TestName: "Without",
			Source:   ErrInternalServer,
			Expected: nil,
		},
		{
			TestName: "With data",
			Source:   ErrInternalServer.WithTemplate(template).(MyError),
			Expected: template,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.TestName, func(t *testing.T) {
			got := tC.Source.GetData()
			assert.Equal(t, tC.Expected, got)
		})
	}
}

func TestSetAction(t *testing.T) {
	testCases := []struct {
		TestName string
		Source   MyError
		Action   string
	}{
		{
			TestName: "Without needed action",
			Source:   NewMyError(100, i18n.Message{MessageID: "FOO.BAR"}),
		},
		{
			TestName: "With needed action",
			Source:   NewMyError(100, i18n.Message{MessageID: "FOO.BAR"}),
			Action:   "some-action",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.TestName, func(t *testing.T) {
			resultError := tC.Source.SetAction(tC.Action)
			assert.Equal(t, &tC.Action, resultError.action)
			assert.NotEqual(t, &tC.Action, tC.Source.action)
			assert.NotEqual(t, tC.Source, resultError)
		})
	}
}

func TestSetKeyBody(t *testing.T) {
	testCases := []struct {
		TestName string
		Source   MyError
		KeyBody  string
	}{
		{
			TestName: "Without key body",
			Source:   NewMyError(100, i18n.Message{MessageID: "FOO.BAR"}),
		},
		{
			TestName: "With key body",
			Source:   NewMyError(100, i18n.Message{MessageID: "FOO.BAR"}),
			KeyBody:  "custom_key",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.TestName, func(t *testing.T) {
			resultError := tC.Source.SetKeyBody(tC.KeyBody)
			assert.Equal(t, tC.KeyBody, resultError.keyBody)
			assert.NotEqual(t, tC.KeyBody, tC.Source.keyBody)
		})
	}
}

func TestError(t *testing.T) {
	testCases := []struct {
		TestName string
		Input    error
		Expected string
	}{
		{
			TestName: "Error without template",
			Input:    ErrInternalServer,
			Expected: i18n.T(i18n.Message{MessageID: "ERRORS.INTERNAL_SERVER"}),
		},
		{
			TestName: "Error with template",
			Input:    ErrFieldValidation("foo", "bar", "10"),
			Expected: i18n.T(i18n.Message{
				MessageID: "ERRORS.FIELD_VALIDATION",
				TemplateData: map[string]interface{}{
					"Field":      "foo",
					"Validation": "bar",
					"Valid":      "10",
				},
			}),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.TestName, func(t *testing.T) {
			assert.Equal(t, tC.Expected, tC.Input.Error())
		})
	}
}
