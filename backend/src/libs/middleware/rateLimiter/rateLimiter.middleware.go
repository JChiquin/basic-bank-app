package ratelimiter

import (
	"net/http"
	"time"

	"bank-service/src/libs/logger"

	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sethvargo/go-limiter/memorystore"
)

const (
	defaultLimit    = 20
	defaultInterval = time.Second * 20
)

/*
NewCustomRateLimiterMiddleware is a partial application of NewRateLimiterMiddleware,
fixing the third parameter to the authorizationKey function
*/
func NewCustomRateLimiterMiddleware(tokens uint64, interval time.Duration) func(next http.Handler) http.Handler {
	return NewRateLimiterMiddleware(tokens, interval, authorizationKey())
}

/*
NewDefaultLimiterMiddleware is a partial application of NewRateLimiterMiddleware,
fixing the first two parameters to default values and the third to the authorizationKey function
*/
func NewDefaultRateLimiterMiddleware() func(next http.Handler) http.Handler {
	return NewRateLimiterMiddleware(defaultLimit, defaultInterval, authorizationKey())
}

/*
NewRateLimiterMiddleware takes a number of tokens, an interval and a function as parameters,
and uses them to configure a rate limit middleware.
It returns a Handler function.
*/
func NewRateLimiterMiddleware(tokens uint64, interval time.Duration, fn httplimit.KeyFunc) func(next http.Handler) http.Handler {
	store, _ := configureStore(tokens, interval) // doesn't return an error
	middleware, err := httplimit.NewMiddleware(store, fn)
	if err != nil {
		logger.GetInstance().Panic("rate limit middleware: " + err.Error())
	}
	return middleware.Handle
}

/*
configureStore takes a number of tokens and an interval to create a new
memory store
*/
func configureStore(tokens uint64, interval time.Duration) (limiter.Store, error) {
	return memorystore.New(&memorystore.Config{
		Tokens:   tokens,   // Number of tokens allowed per interval.
		Interval: interval, // Interval until tokens reset.
	})
}

/*
authorizationKey is a function that gets and returns the Authorization header
from a request
*/
func authorizationKey(headers ...string) func(r *http.Request) (string, error) {
	return func(r *http.Request) (string, error) {
		return r.Header.Get("Authorization"), nil
	}
}
