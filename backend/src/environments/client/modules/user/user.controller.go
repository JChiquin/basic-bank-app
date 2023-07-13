package user

import (
	"bank-service/src/environments/client/resources/interfaces"
	controller "bank-service/src/libs/controllers/client"
	"bank-service/src/libs/dto"
	httpUtils "bank-service/src/libs/http"
	"bank-service/src/libs/i18n"
	"bank-service/src/utils/constant"
	"bank-service/src/utils/helpers"
	"net/http"
)

/*
struct that implements IUserController
*/
type userController struct {
	controller.ClientController
	sUser interfaces.IUserService
}

/*
NewUserController creates a new controller, receives service by dependency injection
and returns IUserController, so it needs to implement all its methods
*/
func NewUserController(sUser interfaces.IUserService) interfaces.IUserController {
	cUser := userController{sUser: sUser}
	return &cUser
}

/*
Create extracts request body and calls Create service
*/
func (c *userController) Create(response http.ResponseWriter, request *http.Request) {
	createUser := &dto.CreateUser{}
	err := httpUtils.GetBodyRequest(request, createUser)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}

	user, err := c.sUser.Create(createUser)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakeSuccessResponse(response, user, http.StatusCreated, i18n.T(i18n.Message{MessageID: "USER.CREATED"}))
}

/*
Login extracts request body and calls Login service
*/
func (c *userController) Login(response http.ResponseWriter, request *http.Request) {
	requestLogin := &dto.RequestLogin{}
	err := httpUtils.GetBodyRequest(request, requestLogin)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}

	user, err := c.sUser.Login(requestLogin)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakeSuccessResponse(response, user, http.StatusOK, i18n.T(i18n.Message{MessageID: "USER.LOGIN"}))
}

/*
WhoAmI extracts userID from context (came from JWT) and calls service to find user by that value
*/
func (c *userController) WhoAmI(response http.ResponseWriter, request *http.Request) {
	jwtContext := request.Context().Value(helpers.ContextKey(constant.JWTContext)).(*dto.JWTContext)

	user, err := c.sUser.FindByID(jwtContext.UserID)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakeSuccessResponse(response, user, http.StatusOK, i18n.T(i18n.Message{MessageID: "USER.WHO_AM_I"}))
}

/*
GetBalance extracts userID from context (came from JWT) and calls service to the user last balance
*/
func (c *userController) GetBalance(response http.ResponseWriter, request *http.Request) {
	jwtContext := request.Context().Value(helpers.ContextKey(constant.JWTContext)).(*dto.JWTContext)

	balance, err := c.sUser.GetBalance(jwtContext.UserID)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakeSuccessResponse(response, balance, http.StatusOK, i18n.T(i18n.Message{MessageID: "USER.BALANCE"}))
}

/*
FindByAccountNumber extracts account_number from path variable and calls service to get that user
*/
func (c *userController) FindByAccountNumber(response http.ResponseWriter, request *http.Request) {
	accountNumber, err := httpUtils.GetParamRequest(request, "account_number")
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}

	user, err := c.sUser.FindByAccountNumber(accountNumber)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}

	c.MakeSuccessResponse(response, user, http.StatusOK, i18n.T(i18n.Message{MessageID: "USER.FOUND"}))
}

/*
UpdatePassword extracts request body and user ID from JWT, then calls UpdatePassword service
*/
func (c *userController) UpdatePassword(response http.ResponseWriter, request *http.Request) {
	jwtContext := request.Context().Value(helpers.ContextKey(constant.JWTContext)).(*dto.JWTContext)

	updatePassword := &dto.UpdatePassord{}
	err := httpUtils.GetBodyRequest(request, updatePassword)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}

	updatePassword.UserID = jwtContext.UserID

	err = c.sUser.UpdatePassword(updatePassword)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakeSuccessResponse(response, nil, http.StatusOK, i18n.T(i18n.Message{MessageID: "USER.PASSWORD_UPDATED"}))
}
