package middleware

import (
	"net/http"

	"bank-service/src/libs/i18n"
)

/*
LanguageMiddleware takes the Accept-Language header and set that language
*/
func LanguageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i18n.SetLanguage(r.Header.Get("Accept-Language"))
		next.ServeHTTP(w, r)
	})
}
