package user

import (
	"bank-service/src/environments/client/resources/interfaces"
	httpUtils "bank-service/src/libs/http"
	"net/http"

	"github.com/gorilla/mux"
)

type userRouter struct {
	cUser interfaces.IUserController
}

/*
NewUserPrivateRouter receives controller, creates the router and calls function to setup all endpoints
*/
func NewUserPrivateRouter(subRouter *mux.Router, cUser interfaces.IUserController) {
	routerUser := userRouter{
		cUser: cUser,
	}
	routerUser.privateRoutes(subRouter)
}

/*
privateRoutes assigns controller function for routes
*/
func (r *userRouter) privateRoutes(subRouter *mux.Router) {
	subRouter.
		Path("/whoami").
		Handler(httpUtils.Middleware(
			http.HandlerFunc(r.cUser.WhoAmI),
		)).
		Methods(http.MethodGet)

	subRouter.
		Path("/account/{account_number}").
		Handler(httpUtils.Middleware(
			http.HandlerFunc(r.cUser.FindByAccountNumber),
		)).
		Methods(http.MethodGet)

	subRouter.
		Path("/balance").
		Handler(httpUtils.Middleware(
			http.HandlerFunc(r.cUser.GetBalance),
		)).
		Methods(http.MethodGet)

	subRouter.
		Path("/password").
		Handler(httpUtils.Middleware(
			http.HandlerFunc(r.cUser.UpdatePassword),
		)).
		Methods(http.MethodPatch)
}

/*
NewUserPublicRouter receives controller, creates the router and calls function to setup all endpoints
*/
func NewUserPublicRouter(subRouter *mux.Router, cUser interfaces.IUserController) {
	routerUser := userRouter{
		cUser: cUser,
	}
	routerUser.publicRoutes(subRouter)
}

/*
publicRoutes assigns controller function for routes
*/
func (r *userRouter) publicRoutes(subRouter *mux.Router) {
	subRouter.
		Path("/register").
		Handler(httpUtils.Middleware(
			http.HandlerFunc(r.cUser.Create),
		)).
		Methods(http.MethodPost)

	subRouter.
		Path("/login").
		Handler(httpUtils.Middleware(
			http.HandlerFunc(r.cUser.Login),
		)).
		Methods(http.MethodPost)
}
