package src

import (
	"bank-service/src/libs/middleware"
	"net/http"
	"time"

	clientRouter "bank-service/src/environments/client/resources/router"
	utilsCors "bank-service/src/libs/cors"
	utilsHttp "bank-service/src/libs/http"
	utilNotFound "bank-service/src/libs/middleware/notFound"
	"bank-service/src/libs/middleware/xss"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

/*
SetupHandler returns the handler with all routes and middlewares using mux
*/
func SetupHandler() *http.Handler {
	muxRouter := mux.NewRouter()

	muxRouter.Use(middleware.LanguageMiddleware)
	settingRoutes(muxRouter)
	utilNotFound.CustomNotFoundHandler(muxRouter)
	handler := handlers.RecoveryHandler()(muxRouter)
	handler = utilsCors.SetCors(handler)

	return &handler
}

/*
settingRoutes takes a pointer to Router and call all environment routers passing its prefix
*/
func settingRoutes(muxRouter *mux.Router) {
	muxRouter.Use(xss.NewXSSMiddleware().Handler)
	pingEndpoint(muxRouter)
	clientRouter.SetupClientPublicRoutes(muxRouter.PathPrefix("/v1/public/client").Subrouter())
	clientRouter.SetupClientPrivateRoutes(muxRouter.PathPrefix("/v1/client").Subrouter())
}

//pingEndpoint is a public endpoint to check the status of this running instance
func pingEndpoint(muxRouter *mux.Router) {
	muxRouter.HandleFunc("/ping", func(response http.ResponseWriter, request *http.Request) {
		bodyResponse := struct {
			Service  string    `json:"service"`
			Datetime time.Time `json:"datetime"`
		}{
			Service:  "bank-service is online",
			Datetime: time.Now(),
		}
		utilsHttp.MakeSuccessResponse(response, bodyResponse, http.StatusOK, http.StatusText(http.StatusOK))
	})
}
