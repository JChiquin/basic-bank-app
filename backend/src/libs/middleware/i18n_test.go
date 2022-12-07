package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	utilsMock "bank-service/src/utils/test/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLanguageMiddleware(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockHTTPHandler := new(utilsMock.MockHTTPHandler)

		mockHTTPHandler.On("ServeHTTP", mock.Anything, mock.Anything).Return()
		ts := httptest.NewServer(LanguageMiddleware(mockHTTPHandler))
		defer ts.Close()
		req, _ := http.NewRequest("GET", ts.URL, nil)
		res, _ := ts.Client().Do(req)

		//Mock Assertion: Behavioral
		mockHTTPHandler.AssertNumberOfCalls(t, "ServeHTTP", 1)

		//Data Assertion
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}
