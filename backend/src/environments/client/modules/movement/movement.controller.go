package movement

import (
	"bank-service/src/environments/client/resources/interfaces"
	controller "bank-service/src/libs/controllers/client"
	"bank-service/src/libs/dto"
	"bank-service/src/utils/constant"
	"bank-service/src/utils/helpers"
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

	movements, err := c.sMovement.IndexByUserID(jwtContext.UserID)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakeSuccessResponse(response, movements, http.StatusOK, "success")
}
