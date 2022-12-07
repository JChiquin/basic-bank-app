package src

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	utilsHttp "bank-service/src/libs/http"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestSetupHandler(t *testing.T) {
	t.Run("Should not panics", func(t *testing.T) {
		assert.NotPanics(t, func() { SetupHandler() })
	})
}

func TestPingEndpoint(t *testing.T) {
	t.Run("Should response http 200 status", func(t *testing.T) {
		muxRouter := mux.NewRouter()
		pingEndpoint(muxRouter)
		ts := httptest.NewServer(muxRouter)
		defer ts.Close()
		req, _ := http.NewRequest("GET", ts.URL+"/ping", nil)
		res, _ := ts.Client().Do(req)

		//Data Assertion
		type data struct {
			Service  string    `json:"service"`
			Datetime time.Time `json:"datetime"`
		}

		bodyResult, _ := utilsHttp.GetBodyBankResponse(res, &data{})
		bodyData := bodyResult.Data.(*data)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "bank-service is online", bodyData.Service)
		assert.Empty(t, bodyResult.Errors)
		assert.Equal(t, http.StatusText(http.StatusOK), bodyResult.Message)
		assert.WithinDuration(t, time.Now(), bodyData.Datetime, 1*time.Second)
	})
}
