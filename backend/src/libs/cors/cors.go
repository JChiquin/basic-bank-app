package cors

import (
	"net/http"
	"strings"

	"bank-service/src/libs/env"

	"github.com/gorilla/handlers"
)

/*
SetCors recieves a handler and returns a handler with CORS support
*/
func SetCors(handler http.Handler) http.Handler {
	credentialsOk := handlers.AllowCredentials()
	headersOk := handlers.AllowedHeaders([]string{
		"Accept",
		"Accept-Language",
		"Content-Type",
		"Content-Language",
		"Origin",
		"Authorization",
		"X-Requested-With",
	})
	originsOk := handlers.AllowedOrigins(strings.Split(env.WhiteList, ","))
	methodsOk := handlers.AllowedMethods([]string{"GET", "PUT", "PATCH", "POST", "DELETE", "OPTIONS", "HEAD"})
	exposeHeadersOk := handlers.ExposedHeaders([]string{
		"X-notifications-unreaded",
		"X-pagination-total-count",
		"X-pagination-page-count",
		"X-pagination-current-page",
		"X-pagination-page-size",
		"needs-action",
	})
	return handlers.CORS(originsOk, headersOk, methodsOk, exposeHeadersOk, credentialsOk)(handler)
}
