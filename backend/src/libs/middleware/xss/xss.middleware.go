package xss

import (
	"bytes"
	goerrors "errors"
	"io/ioutil"
	"net/http"
	"reflect"

	"bank-service/src/libs/errors"
	utilsHttp "bank-service/src/libs/http"
	"bank-service/src/libs/logger"

	"github.com/jinzhu/copier"
	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/square/go-jose.v2/json"
)

var (
	getBodyRequest = utilsHttp.GetBodyRequest // to inject spy
	jsonMarshal    = json.Marshal
)

/*
IXSSMiddleware interface for xss middleware
*/
type IXSSMiddleware interface {
	Handler(next http.Handler) http.Handler
}

type xssMiddleware struct{}

//NewXSSMiddleware is a constructor for middleware struct
func NewXSSMiddleware() IXSSMiddleware {
	return new(xssMiddleware)
}

/*
sanitizeBody takes a pointer to an http.Request and
removes any malicious xss code found in any string inside the body
*/
func (xss *xssMiddleware) sanitizeBody(r *http.Request) error {
	var body interface{}
	err := getBodyRequest(r, &body)
	if err != nil && !goerrors.Is(err, errors.ErrInvalidRequestBody) {
		logger.GetInstance().Error(err)
		return err
	}
	if body != nil {
		sanitizedBody := xss.sanitize(body)
		bodyjson, err := jsonMarshal(sanitizedBody)
		if err != nil {
			logger.GetInstance().Error(err)
			return err
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyjson))
	}
	return nil
}

/*
sanitizeParams takes a pointer to an http.Request and
removes any malicious xss code found in the query params
*/
func (xss *xssMiddleware) sanitizeParams(r *http.Request) {
	params := r.URL.Query()
	if len(params) > 0 {
		for k, v := range params {
			params[k] = xss.sanitize(v).([]string)
		}
		r.URL.RawQuery = params.Encode()
	}
}

/*
sanitizePathParams takes a pointer to an http.Request and
removes any malicious xss code found in the path params
*/
func (xss *xssMiddleware) sanitizePathParams(r *http.Request) {
	pathParams := r.URL.Path
	if pathParams != "" {
		sanizedPathParams := xss.sanitize(pathParams)
		r.URL.Path = sanizedPathParams.(string)
	}
}

/*
Handler can be applied to a route to sanitize the request.
It doesn't throw an error in case of xss attack
Just removes any malicious xss code found in body, query or params.
*/
func (xss *xssMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := xss.sanitizeBody(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//xss.sanitizeParams(r)
		xss.sanitizePathParams(r)
		next.ServeHTTP(w, r)
	})
}

/*
sanitize takes a pointer to any type of input and sanitizes it
It's a recursive function, so it will sanitize any nested structs
*/
func (xss *xssMiddleware) sanitize(input interface{}) interface{} {
	if input == nil {
		return input
	}
	switch reflect.TypeOf(input).Kind() {
	// if it's a map[string]something we copy it to a map[string]interface{} and iterate over it, calling this method recursively
	case reflect.Map:
		auxMap := make(map[string]interface{})
		copier.Copy(&auxMap, input)
		for k, v := range auxMap {
			auxMap[k] = xss.sanitize(v)
		}
		copier.Copy(&input, auxMap)
	case reflect.Slice, reflect.Array:
		switch reflect.TypeOf(input).Elem().Kind() {
		case reflect.Int, reflect.Bool, reflect.Float64:
			// if it's a slice of primitives, we don't sanitize it
			return input
		case reflect.String:
			// if it's a slice of string we iterate over it sanitizing the strings
			for i, v := range input.([]string) {
				input.([]string)[i] = bluemonday.StrictPolicy().Sanitize(v)
			}
		default:
			// if it's a slice of interface{} we copy it to a []interface{} and iterate over it, calling this method recursively
			var auxArr []interface{}
			copier.Copy(&auxArr, input)
			for i, v := range auxArr {
				auxArr[i] = xss.sanitize(v)
			}
			input = auxArr
		}
	case reflect.String:
		// if it's a string, we finally sanitize it
		input = bluemonday.StrictPolicy().Sanitize(input.(string))
	}
	return input
}
