package user

import (
	"bank-service/src/environments/client/resources/interfaces"
	controller "bank-service/src/libs/controllers/client"
	"bank-service/src/libs/dto"
	httpUtils "bank-service/src/libs/http"
	"bank-service/src/libs/i18n"
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

func (c *userController) WhoAmI(response http.ResponseWriter, request *http.Request) {
	userID, err := httpUtils.GetParamRequestInt(request, "id")
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	user, err := c.sUser.FindByID(userID)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakeSuccessResponse(response, user, http.StatusOK, i18n.T(i18n.Message{MessageID: "USER.WHO_AM_I"}))
}
