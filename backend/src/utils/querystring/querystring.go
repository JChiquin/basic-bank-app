package querystring

import (
	myErrors "bank-service/src/libs/errors"
	"bank-service/src/utils/constant"
	"errors"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/gorilla/schema"
)

//newQueryStringDecoder returns a decoder with some custom options
func newQueryStringDecoder() *schema.Decoder {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	decoder.RegisterConverter(time.Time{}, timeConverter)
	return decoder
}

//Decode decodes url.Values into dst
func Decode(dst interface{}, values url.Values) error {
	decoder := newQueryStringDecoder()
	err := decoder.Decode(dst, values)
	if err == nil {
		return nil
	}
	multiError := schema.MultiError{}
	if errors.As(err, &multiError) {
		for _, err2 := range multiError {
			switch t := err2.(type) {
			case schema.ConversionError:
				return myErrors.ErrUnsupportedFieldValue(t.Key)
			}
		}
	}
	return myErrors.ErrInvalidRequestBody
}

//timeConverter is a custom time converter, to remove unnecessary quotation marks
func timeConverter(value string) reflect.Value {
	sanitizedValue := strings.ReplaceAll(value, `"`, ``)
	if v, err := time.Parse(time.RFC3339, sanitizedValue); err == nil {
		return reflect.ValueOf(v)
	} else if v, err := time.Parse(constant.FormatYYYYzMMzDD, sanitizedValue); err == nil { // Try different format to parse if previous one fails (date only)
		v = v.In(time.Now().Location()) // Add location to parsed date
		return reflect.ValueOf(v)
	}
	return reflect.Value{} // this is the same as the private const invalidType
}
