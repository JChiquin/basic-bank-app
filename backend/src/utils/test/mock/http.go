package mock

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

/*
MockHTTPHandler mocks an http handler
*/
type MockHTTPHandler struct {
	mock.Mock
}

func (mock *MockHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mock.Called(w, r)
}

/*
MHTTPHandle generates a mock-server to serve func to test
*/
func MHTTPHandle(
	method string,
	path string,
	f func(http.ResponseWriter, *http.Request),
	params string,
	query url.Values,
	body interface{},
) *http.Response {
	if query == nil {
		query = url.Values{}
	}
	var reqBody []byte
	if body != nil {
		switch value := body.(type) {
		case []byte: //For use body as JSON
			reqBody = value
		case string: //For use body as JSON
			reqBody = []byte(value)
		default:
			reqBody, _ = json.Marshal(body)
		}
	}
	r := mux.NewRouter()
	r.HandleFunc(path, f).Methods(method)
	ts := httptest.NewServer(r)
	defer ts.Close()
	req, _ := http.NewRequest(method, ts.URL+params, bytes.NewBuffer(reqBody))
	req.URL.RawQuery = query.Encode()
	res, _ := ts.Client().Do(req)
	return res
}
