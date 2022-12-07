package ratelimiter

//Unit tests

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	mockHandler = func(response http.ResponseWriter, request *http.Request) {}
	firstToken  = "first jwt"
	secondToken = "second jwt"
)

func TestNewDefaultRateLimiterMiddleware(t *testing.T) {
	t.Run("Should success on", func(t *testing.T) {
		t.Run("Requests within limit", func(t *testing.T) {
			testCases := []struct {
				TestName string
				Token    string
				Calls    int
				Passed   int
				Rejected int
			}{
				{
					TestName: "Amount of requests not close to the limit",
					Token:    firstToken,
					Calls:    defaultLimit - 15,
					Passed:   defaultLimit - 15,
					Rejected: 0,
				},
				{
					TestName: "Requests are half of limit",
					Token:    secondToken,
					Calls:    defaultLimit / 2,
					Passed:   defaultLimit / 2,
					Rejected: 0,
				},
				{
					TestName: "Requests are close to the limit",
					Token:    firstToken,
					Calls:    defaultLimit - 4,
					Passed:   defaultLimit - 4,
					Rejected: 0,
				},
				{
					TestName: "Requests are exactly on the limit",
					Token:    firstToken,
					Calls:    defaultLimit,
					Passed:   defaultLimit,
					Rejected: 0,
				},
				{
					TestName: "Requests pass the limit by 1",
					Token:    secondToken,
					Calls:    defaultLimit + 1,
					Passed:   defaultLimit,
					Rejected: 1,
				},
				{
					TestName: "Requests pass the limit by a lot",
					Token:    firstToken,
					Calls:    defaultLimit * 2,
					Passed:   defaultLimit,
					Rejected: defaultLimit,
				},
			}

			for _, tC := range testCases {
				t.Run(tC.TestName, func(t *testing.T) {
					passed := 0
					rejected := 0

					ts := httptest.NewServer(NewDefaultRateLimiterMiddleware()(http.HandlerFunc(mockHandler)))
					defer ts.Close()
					for i := 0; i < tC.Calls; i++ {
						req, _ := http.NewRequest("GET", ts.URL, nil)
						req.Header.Add("Authorization", tC.Token)
						res, _ := ts.Client().Do(req)
						if res.StatusCode == http.StatusOK {
							passed++
						} else {
							rejected++
						}
					}
					//Assertion
					assert.Equal(t, tC.Passed, passed)
					assert.Equal(t, tC.Rejected, rejected)
				})
			}
		})
	})
}
func TestNewCustomRateLimiterMiddleware(t *testing.T) {
	t.Run("Should success on", func(t *testing.T) {
		t.Run("Requests within limit", func(t *testing.T) {
			testCases := []struct {
				TestName string
				Limit    uint64
				Calls    int
				Passed   int
				Rejected int
			}{
				{
					TestName: "Amount of requests not close to the limit",
					Limit:    15,
					Calls:    2,
					Passed:   2,
					Rejected: 0,
				},
				{
					TestName: "Requests are half of limit",
					Limit:    16,
					Calls:    8,
					Passed:   8,
					Rejected: 0,
				},
				{
					TestName: "Requests are close to the limit",
					Limit:    10,
					Calls:    7,
					Passed:   7,
					Rejected: 0,
				},
				{
					TestName: "Requests are exactly on the limit",
					Limit:    12,
					Calls:    12,
					Passed:   12,
					Rejected: 0,
				},
				{
					TestName: "Requests pass the limit by 1",
					Limit:    5,
					Calls:    6,
					Passed:   5,
					Rejected: 1,
				},
				{
					TestName: "Requests pass the limit by a lot",
					Limit:    6,
					Calls:    15,
					Passed:   6,
					Rejected: 9,
				},
			}

			for _, tC := range testCases {
				t.Run(tC.TestName, func(t *testing.T) {
					passed := 0
					rejected := 0

					ts := httptest.NewServer(NewCustomRateLimiterMiddleware(tC.Limit, time.Second*30)(http.HandlerFunc(mockHandler)))
					defer ts.Close()

					//Action
					for i := 0; i < tC.Calls; i++ {
						res, _ := http.DefaultClient.Get(ts.URL)
						if res.StatusCode == http.StatusOK {
							passed++
						} else {
							rejected++
						}
					}

					//Assertion
					assert.Equal(t, tC.Passed, passed)
					assert.Equal(t, tC.Rejected, rejected)
				})
			}
		})
		t.Run("Requests reach limit but are reset before next request", func(t *testing.T) {
			passed := 0
			rejected := 0
			sleep := time.Second * 5

			ts := httptest.NewServer(NewCustomRateLimiterMiddleware(20, sleep)(http.HandlerFunc(mockHandler)))
			defer ts.Close()
			for i := 0; i < defaultLimit+10; i++ {
				if i == defaultLimit+1 {
					time.Sleep(sleep) // Wait for requests to reset
				}
				req, _ := http.NewRequest("GET", ts.URL, nil)
				res, _ := ts.Client().Do(req)
				if res.StatusCode == http.StatusOK {
					passed++
				} else {
					rejected++
				}
			}
			//Assertion
			assert.Equal(t, defaultLimit+9, passed)
			assert.Equal(t, 1, rejected)
		})
		t.Run("One user reaches the requests limit and the other one doesn't", func(t *testing.T) {
			firstUserPassed := 0
			firstUserRejected := 0
			secondUserPassed := 0
			secondUserRejected := 0
			sleep := time.Second * 10

			ts := httptest.NewServer(NewCustomRateLimiterMiddleware(20, sleep)(http.HandlerFunc(mockHandler)))
			defer ts.Close()
			// 0 - 29 = 30
			for i := 0; i < defaultLimit+10; i++ {
				req, _ := http.NewRequest("GET", ts.URL, nil)
				if i < defaultLimit {
					req.Header.Set("Authorization", firstToken)
					res, _ := ts.Client().Do(req)
					if res.StatusCode == http.StatusOK {
						firstUserPassed++
					} else {
						firstUserRejected++
					}
				}
				req.Header.Set("Authorization", secondToken)
				res, _ := ts.Client().Do(req)
				if res.StatusCode == http.StatusOK {
					secondUserPassed++
				} else {
					secondUserRejected++
				}
			}

			//Assertion
			assert.Equal(t, defaultLimit, firstUserPassed)
			assert.Equal(t, 0, firstUserRejected)
			assert.Equal(t, defaultLimit, secondUserPassed)
			assert.Equal(t, 10, secondUserRejected)
		})
	})
}

func TestNewRateLimiterMiddleware(t *testing.T) {
	keyFunc := func(headers ...string) func(r *http.Request) (string, error) {
		return func(r *http.Request) (string, error) {
			return r.Header.Get("User"), nil
		}
	}
	t.Run("Should success on", func(t *testing.T) {
		t.Run("Requests within limit", func(t *testing.T) {
			firstUserKey := "first user"
			secondUserKey := "boke user"
			type User struct {
				Calls    int
				Passed   int
				Rejected int
			}
			testCases := []struct {
				TestName   string
				FirstUser  User
				SecondUser User
			}{
				{
					TestName: "Amount of requests not close to the limit for either user",
					FirstUser: User{
						Calls:    defaultLimit - 18,
						Passed:   defaultLimit - 18,
						Rejected: 0,
					},
					SecondUser: User{
						Calls:    defaultLimit - 14,
						Passed:   defaultLimit - 14,
						Rejected: 0,
					},
				},
				{
					TestName: "Requests are half of limit for each user",
					FirstUser: User{
						Calls:    defaultLimit / 2,
						Passed:   defaultLimit / 2,
						Rejected: 0,
					},
					SecondUser: User{
						Calls:    defaultLimit / 2,
						Passed:   defaultLimit / 2,
						Rejected: 0,
					},
				},
				{
					TestName: "Requests are close to the limit for each user",
					FirstUser: User{
						Calls:    defaultLimit - 4,
						Passed:   defaultLimit - 4,
						Rejected: 0,
					},
					SecondUser: User{
						Calls:    defaultLimit - 2,
						Passed:   defaultLimit - 2,
						Rejected: 0,
					},
				},
				{
					TestName: "Requests are exactly on the limit for each user",
					FirstUser: User{
						Calls:    defaultLimit,
						Passed:   defaultLimit,
						Rejected: 0,
					},
					SecondUser: User{
						Calls:    defaultLimit,
						Passed:   defaultLimit,
						Rejected: 0,
					},
				},
				{
					TestName: "Each user passes requests limit by 1",
					FirstUser: User{
						Calls:    defaultLimit + 1,
						Passed:   defaultLimit,
						Rejected: 1,
					},
					SecondUser: User{
						Calls:    defaultLimit + 1,
						Passed:   defaultLimit,
						Rejected: 1,
					},
				},
				{
					TestName: "Each user passes requests limit by a lot",
					FirstUser: User{
						Calls:    defaultLimit * 2,
						Passed:   defaultLimit,
						Rejected: defaultLimit,
					},
					SecondUser: User{
						Calls:    defaultLimit * 3,
						Passed:   defaultLimit,
						Rejected: defaultLimit * 2,
					},
				},
			}

			for _, tC := range testCases {
				t.Run(tC.TestName, func(t *testing.T) {
					var wg sync.WaitGroup
					firstUserPassed := 0
					firstUserRejected := 0
					secondUserPassed := 0
					secondUserRejected := 0

					ts := httptest.NewServer(NewRateLimiterMiddleware(defaultLimit, time.Second*40, keyFunc())(http.HandlerFunc(mockHandler)))
					defer ts.Close()

					wg.Add(2)
					go func() {
						defer wg.Done()
						for i := 0; i < tC.FirstUser.Calls; i++ {
							req, _ := http.NewRequest("GET", ts.URL, nil)
							req.Header.Add("User", firstUserKey)
							res, _ := ts.Client().Do(req)
							if res.StatusCode == http.StatusOK {
								firstUserPassed++
							} else {
								firstUserRejected++
							}
						}
					}()
					go func() {
						defer wg.Done()
						for i := 0; i < tC.SecondUser.Calls; i++ {
							req, _ := http.NewRequest("GET", ts.URL, nil)
							req.Header.Add("User", secondUserKey)
							res, _ := ts.Client().Do(req)
							if res.StatusCode == http.StatusOK {
								secondUserPassed++
							} else {
								secondUserRejected++
							}
						}
					}()
					wg.Wait()

					//Assertion
					assert.Equal(t, tC.FirstUser.Passed, firstUserPassed)
					assert.Equal(t, tC.FirstUser.Rejected, firstUserRejected)

					//Assertion
					assert.Equal(t, tC.SecondUser.Passed, secondUserPassed)
					assert.Equal(t, tC.SecondUser.Rejected, secondUserRejected)
				})
			}
		})
	})
	t.Run("Should fail on", func(t *testing.T) {
		t.Run("Function to get key is nil", func(t *testing.T) {
			assert.Panics(t, func() {
				NewRateLimiterMiddleware(20, time.Second*30, nil)
			})
		})
	})
}
