package contact

import (
	"bank-service/src/environments/client/resources/interfaces"
	httpUtils "bank-service/src/libs/http"
	"net/http"

	"github.com/gorilla/mux"
)

type contactRouter struct {
	cContact interfaces.IContactController
}

/*
NewContactPrivateRouter receives controller, creates the router and calls function to setup all endpoints
*/
func NewContactPrivateRouter(subRouter *mux.Router, cContact interfaces.IContactController) {
	routerContact := contactRouter{
		cContact: cContact,
	}
	routerContact.privateRoutes(subRouter)
}

/*
privateRoutes assigns controller function for routes
*/
func (r *contactRouter) privateRoutes(subRouter *mux.Router) {
	subRouter.
		Path("").
		Handler(httpUtils.Middleware(
			http.HandlerFunc(r.cContact.Index),
		)).
		Methods(http.MethodGet)
	subRouter.
		Path("").
		Handler(httpUtils.Middleware(
			http.HandlerFunc(r.cContact.Create),
		)).
		Methods(http.MethodPost)
	subRouter.
		Path("/{id}").
		Handler(httpUtils.Middleware(
			http.HandlerFunc(r.cContact.Update),
		)).
		Methods(http.MethodPatch)
	subRouter.
		Path("/{id}").
		Handler(httpUtils.Middleware(
			http.HandlerFunc(r.cContact.Delete),
		)).
		Methods(http.MethodDelete)
	subRouter.
		Path("/{id}").
		Handler(httpUtils.Middleware(
			http.HandlerFunc(r.cContact.GetOne),
		)).
		Methods(http.MethodGet)
}
