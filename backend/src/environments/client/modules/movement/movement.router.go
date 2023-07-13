package movement

import (
	"bank-service/src/environments/client/resources/interfaces"
	httpUtils "bank-service/src/libs/http"
	"net/http"

	"github.com/gorilla/mux"
)

type movementRouter struct {
	cMovement interfaces.IMovementController
}

/*
NewMovementPrivateRouter receives controller, creates the router and calls function to setup all endpoints
*/
func NewMovementPrivateRouter(subRouter *mux.Router, cMovement interfaces.IMovementController) {
	routerMovement := movementRouter{
		cMovement: cMovement,
	}
	routerMovement.privateRoutes(subRouter)
}

/*
privateRoutes assigns controller function for routes
*/
func (r *movementRouter) privateRoutes(subRouter *mux.Router) {
	subRouter.
		Path("").
		Handler(httpUtils.Middleware(
			http.HandlerFunc(r.cMovement.Index),
		)).
		Methods(http.MethodGet)
	subRouter.
		Path("").
		Handler(httpUtils.Middleware(
			http.HandlerFunc(r.cMovement.Create),
		)).
		Methods(http.MethodPost)
}
