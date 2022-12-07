package errors

import (
	"bank-service/src/libs/i18n"
)

const (
	defaultKeyBody = "error"
)

//MyError for custom errors
type MyError struct {
	statusCode int
	err        i18n.Message
	action     *string //Slug for needs-action header
	keyBody    string  //To use a different key on JSON response
	data       map[string]interface{}
}

//Error is the method that is necessary to be implemented for using MyError as an error
func (m MyError) Error() string {
	if m.data != nil {
		m.err.TemplateData = m.data
	}
	return i18n.T(m.err)
}

/*
Is overrided method so we can compare a target error with a specify error
that was created from MyError struct
*/
func (m MyError) Is(target error) bool {
	t, ok := target.(MyError)
	if !ok {
		return false
	}
	return m.statusCode == t.statusCode && m.err.MessageID == t.err.MessageID
}

/*
WithTemplate returns a new error (MyError) with template from an existing MyError
*/
func (m MyError) WithTemplate(template map[string]interface{}) error {
	m.data = template
	return m
}

//NewMyError constructor for myError struct
func NewMyError(status int, err i18n.Message) MyError {
	return MyError{
		statusCode: status,
		err:        err,
		keyBody:    defaultKeyBody,
	}
}

/*
SetAction sets action and returns MyError
*/
func (m MyError) SetAction(action string) MyError {
	m.action = &action
	return m
}

/*
SetKeyBody sets key body and returns MyError
*/
func (m MyError) SetKeyBody(keyBody string) MyError {
	m.keyBody = keyBody
	return m
}

/*
GetData return data field
*/
func (m MyError) GetData() map[string]interface{} {
	return m.data
}
