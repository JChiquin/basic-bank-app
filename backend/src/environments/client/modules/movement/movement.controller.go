package movement

import (
	"bank-service/src/environments/client/resources/interfaces"
	controller "bank-service/src/libs/controllers/client"
	httpUtils "bank-service/src/libs/http"
	"bank-service/src/utils/pagination"
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
	userID, err := httpUtils.GetParamRequestInt(request, "id")
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}

	pagination, err := pagination.GetPaginationFromQuery(request.URL.Query())
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}

	movements, err := c.sMovement.IndexByUserID(userID, pagination)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakePaginateResponse(response, movements, http.StatusOK, pagination)
}
