package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"bank-service/src/libs/dto"
	"bank-service/src/libs/errors"
	"bank-service/src/libs/logger"
	"bank-service/src/utils/constant"

	"github.com/gorilla/mux"
)

/*
MakePaginateResponse Set Data array, pagination headers and Errors to empty array
*/
func MakePaginateResponse(response http.ResponseWriter, data interface{}, statusCode int, pagination *dto.Pagination) {
	response.Header().Set("X-pagination-total-count", strconv.FormatInt(pagination.TotalCount, 10))
	response.Header().Set("X-pagination-page-count", strconv.Itoa(pagination.PageCount()))
	response.Header().Set("X-pagination-current-page", strconv.Itoa(pagination.Page))
	response.Header().Set("X-pagination-page-size", strconv.Itoa(pagination.PageSize))
	body := dto.NewBodyResponse("Success", make([]map[string]string, 0), data)
	makeResponse(response, body, statusCode)
}

/*
MakeSuccessResponse Set Message, Data object and Errors to empty array
*/
func MakeSuccessResponse(response http.ResponseWriter, data interface{}, statusCode int, message string) {
	body := dto.NewBodyResponse(message, make([]map[string]string, 0), data)
	makeResponse(response, body, statusCode)
}

/*
MakeErrorResponse Set Message, Errors to an Array of objects (JSON) and Data to null
*/
func MakeErrorResponse(response http.ResponseWriter, err error) {
	errorMessage := errors.GetErrorMessage(err)
	statusCode := errors.GetStatusCode(err)
	keyBody := errors.GetKeyBody(err)
	logger.GetInstance().Warningf("Response to client error \n client error: %s \n internal error: %s \n status: %d", errorMessage, err, statusCode)
	errors := []map[string]string{{keyBody: errorMessage}}
	body := dto.NewBodyResponse(errorMessage, errors, nil)
	SetActionNeeded(err, response)
	makeResponse(response, body, statusCode)
}

/*
makeResponse Serialize and send the JSON body to client. Above methods end here
*/
func makeResponse(response http.ResponseWriter, body *dto.BodyResponse, statusCode int) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	json.NewEncoder(response).Encode(&body)
}

/*
GetBodyRequest parses request body to received variable
*/
func GetBodyRequest(req *http.Request, data interface{}) error {
	bodyBytes, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	err := json.Unmarshal(bodyBytes, data)
	if err != nil {
		return handleUnmarshalErr(err)
	}
	return nil
}

func handleUnmarshalErr(err error) error {
	switch t := err.(type) {
	case *json.UnmarshalTypeError:
		return errors.ErrUnsupportedFieldValue(t.Field)
	default:
		return errors.ErrInvalidRequestBody
	}
}

/*
GetBodyBankResponse parses body to received variable
*/
func GetBodyBankResponse(res *http.Response, data interface{}) (*dto.BodyResponse, error) {
	bodyResponse := dto.NewBodyResponse("", nil, data)
	err := GetBodyResponse(res, bodyResponse)
	if err != nil {
		return nil, err
	}
	return bodyResponse, nil
}

/*
GetBodyResponse parses response body to received variable and closes the body
*/
func GetBodyResponse(res *http.Response, data interface{}) error {
	body := res.Body
	defer body.Close()
	if err := json.NewDecoder(body).Decode(data); err != nil {
		return err
	}
	return nil
}

/*
SetActionNeeded receives error and response, then calls GetAction from errors package
finally it sets on header the action needed
*/
func SetActionNeeded(err error, response http.ResponseWriter) {
	action := errors.GetAction(err)
	if action != nil {
		response.Header().Set(constant.HeaderNeedsAction, *action)
	}
}

/*
GetParamRequest retrieves requested param from request
*/
func GetParamRequest(request *http.Request, name string) (string, error) {
	param, ok := mux.Vars(request)[name]
	if !ok {
		return "", errors.ErrUnsupportedFieldValue(name)
	}
	return param, nil
}

/*
GetParamRequestInt partial application of GetParamRequest for Int types
*/
func GetParamRequestInt(request *http.Request, name string) (int, error) {
	paramString, err := GetParamRequest(request, name)
	if err != nil {
		return 0, err
	}
	param, err := strconv.Atoi(paramString)
	if err != nil {
		return 0, errors.ErrUnsupportedFieldValue(name)
	}
	return param, nil
}

/*
Middleware (this function) makes adding more than one layer of middleware easy
by specifying them as a list. It will run the first specified middleware first.
*/
func Middleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	length := len(middlewares)
	for i := range middlewares {
		handler = middlewares[length-1-i](handler)
	}
	return handler
}
