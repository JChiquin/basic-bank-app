package controller

import (
	"net/http"
	"testing"

	"bank-service/src/libs/dto"
	myErrors "bank-service/src/libs/errors"
	utilsHttp "bank-service/src/libs/http"
	utilsMocks "bank-service/src/utils/test/mock"

	"github.com/stretchr/testify/assert"
)

type DData struct {
	Client  string `groups:"client"`
	Admin   string `groups:"admin"`
	Console string `groups:"console"`
}

var defaultController ClientController = ClientController{}
var defaultMessage string = "Test"
var defaultData = DData{Client: "Client", Admin: "Admin", Console: "Console"}

func TestBaseController_MakePaginateResponse(t *testing.T) {
	t.Run("Should Succeed", func(t *testing.T) {
		t.Run("Correct header data", func(t *testing.T) {
			// Fixture
			page := &dto.Pagination{}

			// Run Foo inside request
			response := utilsMocks.MHTTPHandle("GET", "/",
				func(response http.ResponseWriter, request *http.Request) {
					assert.NotPanics(t, func() {
						defaultController.MakePaginateResponse(
							response,
							defaultData,
							http.StatusOK,
							page,
						)
					})
				}, "", nil, nil)
			defer response.Body.Close()
			body, _ := utilsHttp.GetBodyBankResponse(response, &DData{})

			// Assert Data
			assert.Equal(t, defaultData.Client, body.Data.(*DData).Client)
			assert.Zero(t, body.Data.(*DData).Admin)
			assert.Zero(t, body.Data.(*DData).Console)
		})
	})
}

func TestBaseController_MakeSuccessResponse(t *testing.T) {
	t.Run("Should Succeed", func(t *testing.T) {
		t.Run("Correct OK response", func(t *testing.T) {
			// Run Foo inside request
			response := utilsMocks.MHTTPHandle("GET", "/",
				func(response http.ResponseWriter, request *http.Request) {
					assert.NotPanics(t, func() {
						defaultController.MakeSuccessResponse(
							response,
							defaultData,
							http.StatusOK,
							defaultMessage,
						)
					})
				}, "", nil, nil)
			defer response.Body.Close()
			body, _ := utilsHttp.GetBodyBankResponse(response, &DData{})

			// Assert Data
			assert.Equal(t, defaultData.Client, body.Data.(*DData).Client)
			assert.Zero(t, body.Data.(*DData).Admin)
			assert.Zero(t, body.Data.(*DData).Console)
		})
	})
}

func TestBaseController_MakeErrorResponse(t *testing.T) {
	t.Run("Should Succeed", func(t *testing.T) {
		t.Run("Correct Error response", func(t *testing.T) {
			// Run Foo inside request
			response := utilsMocks.MHTTPHandle("GET", "/",
				func(response http.ResponseWriter, request *http.Request) {
					assert.NotPanics(t, func() {
						defaultController.MakeErrorResponse(
							response,
							myErrors.ErrNotFound,
						)
					})
				}, "", nil, nil)
			defer response.Body.Close()
			body, _ := utilsHttp.GetBodyBankResponse(response, &DData{})

			// Assert Data
			assert.Nil(t, body.Data)
		})
	})
}
