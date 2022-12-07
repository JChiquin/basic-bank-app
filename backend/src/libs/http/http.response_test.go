package http

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bank-service/src/libs/dto"
	myErrors "bank-service/src/libs/errors"
	"bank-service/src/utils/constant"
	customMock "bank-service/src/utils/test/mock"

	"github.com/stretchr/testify/assert"
)

type DData struct {
	ID   string `groups:"client"`
	Name string `groups:"admin"`
}

var defaultMessage string = "Test"
var defaultErrors []map[string]string = []map[string]string{
	{"test01": "Test01", "test02": "Test02"},
	{"test03": "Test03", "test04": "Test04"},
}
var defaultData = DData{ID: "1", Name: "Name"}

func TestMakePaginateResponse(t *testing.T) {
	t.Run("Should Succeed", func(t *testing.T) {
		t.Run("Correct header data", func(t *testing.T) {
			// Fixture
			page := &dto.Pagination{
				Page:       1,
				PageSize:   10,
				TotalCount: 20,
			}

			// Run Foo inside request
			response := customMock.MHTTPHandle("GET", "/",
				func(response http.ResponseWriter, request *http.Request) {
					assert.NotPanics(t, func() {
						MakePaginateResponse(
							response,
							dto.NewBodyResponse(defaultMessage, defaultErrors, defaultData),
							http.StatusOK,
							page,
						)
					})
				}, "", nil, nil)

			// Assert Data
			assert.Equal(t, "1", response.Header.Get("X-pagination-current-page"))
			assert.Equal(t, "10", response.Header.Get("X-pagination-page-size"))
			assert.Equal(t, "2", response.Header.Get("X-pagination-page-count"))
			assert.Equal(t, "20", response.Header.Get("X-pagination-total-count"))
		})
	})
}

func TestMakeSuccessResponse(t *testing.T) {
	t.Run("Should Succeed", func(t *testing.T) {
		t.Run("Correct OK response", func(t *testing.T) {
			// Run Foo inside request
			response := customMock.MHTTPHandle("GET", "/",
				func(response http.ResponseWriter, request *http.Request) {
					assert.NotPanics(t, func() {
						MakeSuccessResponse(
							response,
							defaultData,
							http.StatusOK,
							defaultData.Name,
						)
					})
				}, "", nil, nil)
			defer response.Body.Close()
			body, _ := GetBodyBankResponse(response, &DData{})

			// Assert Data
			assert.Equal(t, http.StatusOK, response.StatusCode)
			assert.Equal(t, defaultData.Name, body.Message)
			assert.Equal(t, []map[string]string{}, body.Errors)
			assert.Equal(t, defaultData.ID, body.Data.(*DData).ID)
		})
	})
}

func TestMakeErrorResponse(t *testing.T) {
	t.Run("Should Succeed", func(t *testing.T) {
		t.Run("Correct Error response", func(t *testing.T) {
			// Fixture
			expError := []map[string]string{
				{"error": myErrors.ErrNotFound.Error()},
			}

			// Run Foo inside request
			response := customMock.MHTTPHandle("GET", "/",
				func(response http.ResponseWriter, request *http.Request) {
					assert.NotPanics(t, func() {
						MakeErrorResponse(
							response,
							myErrors.ErrNotFound,
						)
					})
				}, "", nil, nil)
			defer response.Body.Close()
			body, _ := GetBodyBankResponse(response, &DData{})

			// Assert Data
			assert.Equal(t, http.StatusNotFound, response.StatusCode)
			assert.Equal(t, myErrors.ErrNotFound.Error(), body.Message)
			assert.Equal(t, expError, body.Errors)
			assert.Nil(t, body.Data)
		})
	})
}

func TestGetBodyRequest(t *testing.T) {
	t.Run("Should Succeed", func(t *testing.T) {
		t.Run("Valid JSON struct", func(t *testing.T) {
			// Fixture
			type Body struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}
			reqBody := []byte(`{
					"name":"Name","age":20
				}`)

			// Prepare Request
			req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(reqBody))
			body := Body{}
			err := GetBodyRequest(req, &body)

			// Response Assert
			assert.NoError(t, err)
			assert.Equal(t, "Name", body.Name)
			assert.Equal(t, 20, body.Age)
		})
	})
	t.Run("Should Fail", func(t *testing.T) {
		//Fixture
		type Body struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		t.Run("Invalid JSON syntax", func(t *testing.T) {
			// Fixture
			reqBody := []byte(`{
					"name":"Nam
				}`)

			// Prepare Request
			req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(reqBody))
			body := Body{}
			err := GetBodyRequest(req, &body)

			// Response Assert
			assert.EqualError(t, err, myErrors.ErrInvalidRequestBody.Error())
			assert.Zero(t, body)
		})
		t.Run("Invalid JSON value for struct field", func(t *testing.T) {
			// Fixture
			expectedErr := myErrors.ErrUnsupportedFieldValue("name")
			reqBody := []byte(`{
					"name": false
				}`)

			// Prepare Request
			req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(reqBody))
			body := Body{}
			err := GetBodyRequest(req, &body)

			// Response Assert
			assert.EqualError(t, err, expectedErr.Error())
			assert.Zero(t, body)
		})
	})
}

func TestGetBodyBankResponse(t *testing.T) {
	t.Run("Should Succeed", func(t *testing.T) {
		// Fixture
		type Body struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		reqBody := []byte(`{
				"Message":"Message",
				"Errors":[{"a":"aa","b":"bb"}],
				"Data":{"name":"Name","age":20}
			}`)

		// Prepare request
		response := customMock.MHTTPHandle("GET", "/", func(response http.ResponseWriter, request *http.Request) {
			response.WriteHeader(http.StatusOK)
			response.Write(reqBody)
		}, "", nil, nil)

		// Get response
		data := &Body{}
		bodyResponse, err := GetBodyBankResponse(response, data)

		// Data assertion
		assert.NoError(t, err)
		assert.Equal(t, "Message", bodyResponse.Message)
		assert.Equal(t, []map[string]string{{"a": "aa", "b": "bb"}}, bodyResponse.Errors)
		assert.Equal(t, "Name", data.Name)
		assert.Equal(t, 20, data.Age)
	})
	t.Run("Should Fail", func(t *testing.T) {
		// Fixture
		type Body struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		reqBody := []byte(`{
				"Message":"Messa
				"Errors":[{"a":"aa"
				"Data":{"name":"Na
			}`)

		// Prepare request
		response := customMock.MHTTPHandle("GET", "/", func(response http.ResponseWriter, request *http.Request) {
			response.WriteHeader(http.StatusOK)
			response.Write(reqBody)
		}, "", nil, nil)

		// Get response
		data := &Body{}
		bodyResponse, err := GetBodyBankResponse(response, data)

		// Data assertion
		assert.Error(t, err)
		assert.Nil(t, bodyResponse)
	})
}

func TestGetBodyResponse(t *testing.T) {

	// Fixture
	type Body struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	reqBody := []byte(`{
		"name": "Name",
		"age": 20
	}`)

	t.Run("Should Succeed", func(t *testing.T) {

		// Prepare request
		response := customMock.MHTTPHandle("GET", "/", func(response http.ResponseWriter, request *http.Request) {
			response.WriteHeader(http.StatusOK)
			response.Write(reqBody)
		}, "", nil, nil)

		// Get response
		data := &Body{}
		err := GetBodyResponse(response, data)

		// Data assertion
		assert.NoError(t, err)
		assert.Equal(t, "Name", data.Name)
		assert.Equal(t, 20, data.Age)
	})
	t.Run("Should Fail", func(t *testing.T) {
		t.Run("Bad JSON", func(t *testing.T) {
			// Fixture
			reqBody := []byte(`{"name": "Name",`)

			// Prepare request
			response := customMock.MHTTPHandle("GET", "/", func(response http.ResponseWriter, request *http.Request) {
				response.WriteHeader(http.StatusOK)
				response.Write(reqBody)
			}, "", nil, nil)

			// Get response
			data := &Body{}
			err := GetBodyResponse(response, data)

			// Data assertion
			assert.Error(t, err)
			assert.Empty(t, data)
		})
		t.Run("Invalid destination variable type", func(t *testing.T) {

			// Prepare request
			response := customMock.MHTTPHandle("GET", "/", func(response http.ResponseWriter, request *http.Request) {
				response.WriteHeader(http.StatusOK)
				response.Write(reqBody)
			}, "", nil, nil)

			// Get response
			data := make(chan bool)
			err := GetBodyResponse(response, data)

			// Data assertion
			assert.Error(t, err)
		})
	})
}

func TestSetAction(t *testing.T) {
	testCases := []struct {
		testName             string
		err                  error
		expectedHeaderAction string
	}{
		{
			testName:             "Custom error without action",
			err:                  myErrors.ErrNotFound,
			expectedHeaderAction: "", //Empty
		},
		{
			testName:             "Custom error with action",
			err:                  myErrors.ErrNotFound.SetAction("Hola"),
			expectedHeaderAction: "Hola",
		},
		{
			testName:             "Generic error",
			err:                  errors.New("Generic error"),
			expectedHeaderAction: "", //Empty
		},
	}
	for _, tC := range testCases {
		t.Run(tC.testName, func(t *testing.T) {
			writer := httptest.NewRecorder()
			SetActionNeeded(tC.err, writer)
			headerAction := writer.Header().Get(constant.HeaderNeedsAction)

			//Data Assertion
			assert.Equal(t, tC.expectedHeaderAction, headerAction)
		})
	}
}

func TestGetParamRequest(t *testing.T) {
	validParamName := "param"
	validParamValue := "2"
	validMethod := http.MethodGet
	t.Run("GetParamRequest", func(t *testing.T) {
		t.Run("Should Succeed", func(t *testing.T) {
			validHandler := func(response http.ResponseWriter, request *http.Request) {
				// Action
				result, err := GetParamRequest(request, validParamName)

				// Assert data
				assert.Nil(t, err)
				assert.Equal(t, validParamValue, result)

				response.WriteHeader(http.StatusOK)
			}

			customMock.MHTTPHandle(validMethod, fmt.Sprintf("/{%s}", validParamName), validHandler, "/"+validParamValue, nil, nil)
		})
		t.Run("Should Fail", func(t *testing.T) {
			validHandler := func(response http.ResponseWriter, request *http.Request) {
				// Action
				result, err := GetParamRequest(request, validParamName)

				// Assert data
				assert.EqualError(t, err, myErrors.ErrUnsupportedFieldValue(validParamName).Error())
				assert.Zero(t, result)

				response.WriteHeader(http.StatusOK)
			}

			customMock.MHTTPHandle(validMethod, "/{invalid}", validHandler, "/asd", nil, nil)
		})
	})
	t.Run("GetParamRequestInt", func(t *testing.T) {
		t.Run("Should Succeed", func(t *testing.T) {
			validHandler := func(response http.ResponseWriter, request *http.Request) {
				// Action
				result, err := GetParamRequestInt(request, validParamName)

				// Assert data
				assert.Nil(t, err)
				assert.Equal(t, 2, result)

				response.WriteHeader(http.StatusOK)
			}

			customMock.MHTTPHandle(validMethod, fmt.Sprintf("/{%s}", validParamName), validHandler, "/"+validParamValue, nil, nil)
		})
		t.Run("Should Fail", func(t *testing.T) {
			testCases := []struct {
				desc, path, value string
			}{
				{
					desc:  "invalid value",
					path:  fmt.Sprintf("/{%s}", validParamName),
					value: "/asd",
				},
				{
					desc:  "missing param",
					path:  "/{invalid}",
					value: "/asd",
				},
			}
			for _, tC := range testCases {
				t.Run(tC.desc, func(t *testing.T) {
					validHandler := func(response http.ResponseWriter, request *http.Request) {
						// Action
						result, err := GetParamRequestInt(request, validParamName)

						// Assert data
						assert.EqualError(t, err, myErrors.ErrUnsupportedFieldValue(validParamName).Error())
						assert.Zero(t, result)

						response.WriteHeader(http.StatusOK)
					}

					customMock.MHTTPHandle(validMethod, tC.path, validHandler, tC.value, nil, nil)
				})
			}
		})
	})
}

func TestMiddleware(t *testing.T) {
	callStack := make([]string, 0)
	//Fixture
	handler := func(response http.ResponseWriter, request *http.Request) {
		callStack = append(callStack, "handler")
		response.WriteHeader(http.StatusTeapot)
	}
	firstMiddleware := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callStack = append(callStack, "firstMiddleware")
			h.ServeHTTP(w, r)
		})
	}
	secondMiddleware := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callStack = append(callStack, "secondMiddleware")
			h.ServeHTTP(w, r)
		})
	}

	//This line do the actual work
	result := Middleware(
		http.HandlerFunc(handler),
		firstMiddleware,
		secondMiddleware,
	)
	ts := httptest.NewServer(result)
	defer ts.Close()
	req, _ := http.NewRequest("GET", ts.URL, nil)
	res, _ := ts.Client().Do(req)

	//Data Assertion
	assert.Equal(t, http.StatusTeapot, res.StatusCode)
	assert.Equal(t, "firstMiddleware", callStack[0])
	assert.Equal(t, "secondMiddleware", callStack[1])
	assert.Equal(t, "handler", callStack[2])
}
