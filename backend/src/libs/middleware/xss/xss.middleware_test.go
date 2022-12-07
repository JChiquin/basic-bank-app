package xss

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	utilsHttp "bank-service/src/libs/http"

	"github.com/stretchr/testify/assert"
)

type bodyDTO struct {
	Text    string   `json:"text"`
	Number  int      `json:"number"`
	Boolean bool     `json:"boolean"`
	Xss     string   `json:"xss"`
	Numbers []int    `json:"numbers"`
	Words   []string `json:"words"`
}

var (
	xss                  = "<img src onerror='alert(\"Hacked\")'></>"
	errGetBodyRequest    = errors.New("error getting request body")
	getBodyRequestBackup = getBodyRequest
	getBodyRequestSpy    = func(req *http.Request, data interface{}) error {
		return errGetBodyRequest
	}
	errJsonMarshal    = errors.New("error on json marshal")
	jsonMarshalBackup = jsonMarshal
	jsonMarshalSpy    = func(v interface{}) ([]byte, error) {
		return nil, errJsonMarshal
	}
)

func TestSanitize(t *testing.T) {
	safeString := "hello world"
	testCases := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "safe string",
			input:    safeString,
			expected: safeString,
		},
		{
			name:     "unsafe string",
			input:    xss,
			expected: "",
		},
		{
			name:     "safe array of string",
			input:    []string{"hello", "world"},
			expected: []string{"hello", "world"},
		},
		{
			name:     "safe array of int",
			input:    []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "safe array of booleans",
			input:    []bool{true, false, true},
			expected: []bool{true, false, true},
		},
		{
			name:     "safe array of floats",
			input:    []float64{1.1, 2.2, 3.3},
			expected: []float64{1.1, 2.2, 3.3},
		},
		{
			name:     "safe array of interface",
			input:    []interface{}{"hello", "world", true, 2, 4, 5.6, map[string]interface{}{"key": "value"}},
			expected: []interface{}{"hello", "world", true, 2, 4, 5.6, map[string]interface{}{"key": "value"}},
		},
		{
			name:     "unsafe array of interface",
			input:    []interface{}{"hello", xss, "world", true, 2, 4, 5.6, map[string]interface{}{"key": "value", "xss": xss, "key3": "value3"}},
			expected: []interface{}{"hello", "", "world", true, 2, 4, 5.6, map[string]interface{}{"key": "value", "xss": "", "key3": "value3"}},
		},
		{
			name:     "unsafe string in array",
			input:    []string{"please stop", xss, safeString},
			expected: []string{"please stop", "", safeString},
		},
		{
			name:     "unsafe string in array of interface",
			input:    []interface{}{safeString, xss, safeString},
			expected: []interface{}{safeString, "", safeString},
		},
		{
			name:     "unsafe string in map",
			input:    map[string]string{"key1": safeString, "key2": xss, "key3": safeString},
			expected: map[string]string{"key1": safeString, "key2": "", "key3": safeString},
		},
		{
			name:     "unsafe string in nested map",
			input:    map[string]map[string]string{"key1": {"key1": safeString, "key2": xss, "key3": safeString}},
			expected: map[string]map[string]string{"key1": {"key1": safeString, "key2": "", "key3": safeString}},
		},
		{
			name:     "unsafe string in nested array",
			input:    [][]string{{"hello", "world"}, {xss, safeString}},
			expected: []interface{}{[]string{"hello", "world"}, []string{"", safeString}},
		},
		{
			name: "unsafe string in map with booleans",
			input: map[string]interface{}{
				"key":  xss,
				"key1": true,
				"key2": false,
			},
			expected: map[string]interface{}{
				"key":  "",
				"key1": true,
				"key2": false,
			},
		},
		{
			name: "unsafe string in map with numbers",
			input: map[string]interface{}{
				"key":  xss,
				"key1": 10.5,
				"key2": -100,
				"key3": 1e10,
			},
			expected: map[string]interface{}{
				"key":  "",
				"key1": 10.5,
				"key2": -100,
				"key3": 1e10,
			},
		},
		{
			name: "unsafe in map and string in array inside a map",
			input: map[string]interface{}{
				"key":  []string{xss, safeString},
				"key2": "Hello world",
				"key3": xss,
			},
			expected: map[string]interface{}{
				"key":  []string{"", safeString},
				"key2": "Hello world",
				"key3": "",
			},
		},
		{
			name: "Nil and empty things",
			input: map[string]interface{}{
				"key1": nil,
				"key2": 0,
				"key3": "",
				"key4": false,
				"key5": []interface{}{},
				"key6": map[interface{}]interface{}{},
				"key7": []interface{}{nil},
				"key8": []interface{}{nil, []interface{}{}, map[interface{}]interface{}{}, "", 0, false},
			},
			expected: map[string]interface{}{
				"key1": nil,
				"key2": 0,
				"key3": "",
				"key4": false,
				"key5": []interface{}{},
				"key6": map[interface{}]interface{}{},
				"key7": []interface{}{nil},
				"key8": []interface{}{nil, []interface{}{}, map[interface{}]interface{}{}, "", 0, false},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			xssMiddleware := new(xssMiddleware)

			got := xssMiddleware.sanitize(tC.input)
			assert.Equal(t, tC.expected, got)
		})
	}
	t.Run("Test json with null values", func(t *testing.T) {
		bodyJSON := []byte(`{
			"key1": null,
			"key2": 0,
			"key3": "",
			"key4": false,
			"key5": [],
			"key6": {},
			"key7": [null],
			"key8": [null, [], {}, "", 0, false]
		   }`)
		var body interface{}
		json.Unmarshal(bodyJSON, &body)
		xssMiddleware := new(xssMiddleware)
		got := xssMiddleware.sanitize(body)
		assert.Equal(t, body, got)
	})
}

func TestSanitizeBody(t *testing.T) {
	xssMiddleware := new(xssMiddleware)
	body := bodyDTO{
		Text:    "valid text",
		Number:  1,
		Boolean: true,
		Xss:     xss,
		Numbers: []int{1, 2, 3},
		Words:   []string{"one", xss, "three"},
	}
	t.Run("Should success on", func(t *testing.T) {
		t.Run("sanitizing a valid body", func(t *testing.T) {
			// fixture
			bodyJson, _ := json.Marshal(body)
			expectedBody := bodyDTO{
				Text:    "valid text",
				Number:  1,
				Boolean: true,
				Xss:     "",
				Numbers: []int{1, 2, 3},
				Words:   []string{"one", "", "three"},
			}

			// action
			req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(bodyJson))
			err := xssMiddleware.sanitizeBody(req)

			var got bodyDTO
			utilsHttp.GetBodyRequest(req, &got)

			// assertion
			assert.NoError(t, err)
			assert.Equal(t, expectedBody, got)
		})
		t.Run("sanitizing a empty body", func(t *testing.T) {
			// action
			req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(nil))
			err := xssMiddleware.sanitizeBody(req)

			var got bodyDTO
			utilsHttp.GetBodyRequest(req, &got)

			// assertion
			assert.NoError(t, err)
			assert.Empty(t, got)
		})
	})
	t.Run("Should Fail on", func(t *testing.T) {
		t.Run("error getting body from request", func(t *testing.T) {
			// inject spy
			getBodyRequest = getBodyRequestSpy

			bodyJson, _ := json.Marshal(body)
			req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(bodyJson))
			err := xssMiddleware.sanitizeBody(req)

			//Data Assertion
			assert.ErrorIs(t, err, errGetBodyRequest)
			t.Cleanup(func() {
				getBodyRequest = getBodyRequestBackup
			})
		})
		t.Run("error making the new body", func(t *testing.T) {
			// inject spy
			jsonMarshal = jsonMarshalSpy

			bodyJson, _ := json.Marshal(body)
			req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(bodyJson))
			err := xssMiddleware.sanitizeBody(req)

			//Data Assertion
			assert.ErrorIs(t, err, errJsonMarshal)
			t.Cleanup(func() {
				jsonMarshal = jsonMarshalBackup
			})
		})
	})
}

func TestSanitizeParams(t *testing.T) {
	xssMiddleware := new(xssMiddleware)
	t.Run("Should success on", func(t *testing.T) {
		t.Run("sanitizing valid query params body", func(t *testing.T) {
			params := map[string]string{
				"text": "valid text",
				"xss":  xss,
			}
			expectedParams := url.Values{
				"text": []string{"valid text"},
				"xss":  []string{""},
			}
			req, _ := http.NewRequest("GET", "", bytes.NewBuffer(nil))
			// append query to req
			q := req.URL.Query()
			for k, v := range params {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()

			// action
			xssMiddleware.sanitizeParams(req)

			// assertion
			assert.Equal(t, expectedParams, req.URL.Query())
		})
		t.Run("sanitizing empty query params", func(t *testing.T) {
			params := map[string]string{}
			expectedParams := url.Values{}
			req, _ := http.NewRequest("GET", "", bytes.NewBuffer(nil))
			// append query to req
			q := req.URL.Query()
			for k, v := range params {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()

			// action
			xssMiddleware.sanitizeParams(req)

			// assertion
			assert.Equal(t, expectedParams, req.URL.Query())
		})
	})
}

func TestSanitizePathParams(t *testing.T) {
	xssMiddleware := new(xssMiddleware)
	t.Run("Should success on", func(t *testing.T) {
		t.Run("sanitizing valid path params", func(t *testing.T) {
			params := "/valid/text/with/xss/" + xss
			expectedParams := "/valid/text/with/xss/"

			req, _ := http.NewRequest("GET", params, bytes.NewBuffer(nil))

			// action
			xssMiddleware.sanitizePathParams(req)

			// assertion
			assert.Equal(t, expectedParams, req.URL.Path)
		})
		t.Run("sanitizing empty path params", func(t *testing.T) {
			expectedParams := "/"

			req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(nil))

			// action
			xssMiddleware.sanitizePathParams(req)

			// assertion
			assert.Equal(t, expectedParams, req.URL.Path)
		})
	})
}

func TestHandler(t *testing.T) {
	xssMiddleware := NewXSSMiddleware()
	body := bodyDTO{
		Text:    "valid text",
		Number:  1,
		Boolean: true,
		Xss:     xss,
		Numbers: []int{1, 2, 3},
		Words:   []string{"one", xss, "three"},
	}
	t.Run("Should success on", func(t *testing.T) {
		t.Run("Sanitizing a request", func(t *testing.T) {
			// fixture
			bodyJson, _ := json.Marshal(body)
			expectedBody := bodyDTO{
				Text:    "valid text",
				Number:  1,
				Boolean: true,
				Xss:     "",
				Numbers: []int{1, 2, 3},
				Words:   []string{"one", "", "three"},
			}

			// action
			ts := httptest.NewServer(
				xssMiddleware.Handler(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						var body bodyDTO
						utilsHttp.GetBodyRequest(r, &body)
						assert.Equal(t, expectedBody, body)
					})))
			defer ts.Close()
			req, _ := http.NewRequest("POST", ts.URL, bytes.NewBuffer(bodyJson))
			resp, _ := ts.Client().Do(req)

			//Data Assertion
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	})
	t.Run("Should fail on", func(t *testing.T) {
		t.Run("sanitizeBody fails", func(t *testing.T) {
			// inject spy
			getBodyRequest = getBodyRequestSpy
			ts := httptest.NewServer(
				xssMiddleware.Handler(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						var body bodyDTO
						utilsHttp.GetBodyRequest(r, &body)
						assert.Nil(t, body)
					})))
			defer ts.Close()
			bodyJson, _ := json.Marshal(body)
			req, _ := http.NewRequest("GET", ts.URL, bytes.NewBuffer(bodyJson))
			resp, _ := ts.Client().Do(req)

			//Data Assertion
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			t.Cleanup(func() {
				getBodyRequest = getBodyRequestBackup
			})
		})
	})
}
