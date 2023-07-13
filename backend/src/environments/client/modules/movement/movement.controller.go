package movement

import (
	"bank-service/src/environments/client/resources/interfaces"
	controller "bank-service/src/libs/controllers/client"
	"bank-service/src/libs/dto"
	httpUtils "bank-service/src/libs/http"
	"bank-service/src/libs/i18n"
	"bank-service/src/utils/constant"
	"bank-service/src/utils/helpers"
	"bank-service/src/utils/pagination"
	"bank-service/src/utils/querystring"
	"net/http"
)

/*
struct that implements IMovementController
*/
type movementController struct {
	controller.ClientController
	sMovement interfaces.IMovementService
}

/*
NewMovementController creates a new controller, receives service by dependency injection
and returns IMovementController, so it needs to implement all its methods
*/
func NewMovementController(sMovement interfaces.IMovementService) interfaces.IMovementController {
	cMovement := movementController{sMovement: sMovement}
	return &cMovement
}

func (c *movementController) Index(response http.ResponseWriter, request *http.Request) {
	jwtContext := request.Context().Value(helpers.ContextKey(constant.JWTContext)).(*dto.JWTContext)
	pagination, err := pagination.GetPaginationFromQuery(request.URL.Query())
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}

	filterMovements := &dto.FilterMovements{}
	if err := querystring.Decode(filterMovements, request.URL.Query()); err != nil {
		c.MakeErrorResponse(response, err)
		return
	}

	filterMovements.UserID = jwtContext.UserID

	movements, err := c.sMovement.IndexByUserID(filterMovements, pagination)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakePaginateResponse(response, movements, http.StatusOK, pagination)
}

func (c *movementController) Create(response http.ResponseWriter, request *http.Request) {
	jwtContext := request.Context().Value(helpers.ContextKey(constant.JWTContext)).(*dto.JWTContext)
	createMovement := &dto.CreateMovement{}
	err := httpUtils.GetBodyRequest(request, createMovement)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}

	createMovement.UserID = jwtContext.UserID

	movement, err := c.sMovement.Create(createMovement)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakeSuccessResponse(response, movement, http.StatusCreated, i18n.T(i18n.Message{MessageID: "MOVEMENT.CREATED"}))
}
