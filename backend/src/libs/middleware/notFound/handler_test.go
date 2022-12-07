package notfound

import (
	"net/http"
	"net/http/httptest"
	"testing"

	myErrors "bank-service/src/libs/errors"
	utilsHttp "bank-service/src/libs/http"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCustomNotFoundHanlder(t *testing.T) {
	t.Run("Should response a custom body on url not found", func(t *testing.T) {
		router := new(mux.Router)
		CustomNotFoundHandler(router)
		ts := httptest.NewServer(router)
		defer ts.Close()

		//Action
		resp, err := http.DefaultClient.Get(ts.URL + "/not_found_url")

		//Data Assertion
		bodyResponse, _ := utilsHttp.GetBodyBankResponse(resp, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.Equal(t, myErrors.ErrURLNotFound.Error(), bodyResponse.Message)
		assert.Equal(t, myErrors.ErrURLNotFound.Error(), bodyResponse.Errors[0]["error"])
	})
}
