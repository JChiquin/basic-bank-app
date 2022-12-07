package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bank-service/src/libs/env"

	"github.com/stretchr/testify/assert"
)

func TestCors(t *testing.T) {
	t.Run("CORS should be enabled", func(t *testing.T) {
		tempWhiteList := env.WhiteList
		env.WhiteList = "http://localhost:9000"
		handler := func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
			})
		}()
		handler = SetCors(handler)
		ts := httptest.NewServer(handler)
		defer ts.Close()
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/ping", nil)
		req.Header.Set("Origin", "http://localhost:9000")
		res, err := ts.Client().Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
		assert.NotEmpty(t, res.Header.Get("Access-Control-Allow-Origin"))
		assert.NotEmpty(t, res.Header.Get("Access-Control-Allow-Credentials"))
		assert.NotEmpty(t, res.Header.Get("Access-Control-Expose-Headers"))
		assert.Empty(t, res.Header.Get("Access-Control-Allow-Methods"))
		assert.Empty(t, res.Header.Get("Access-Control-Allow-Headers"))

		t.Cleanup(func() {
			env.WhiteList = tempWhiteList
		})
	})
}
