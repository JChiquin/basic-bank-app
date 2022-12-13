package router

import (
	"bank-service/src/environments/client/modules/movement"
	"bank-service/src/environments/client/modules/user"
	"bank-service/src/libs/database"

	"github.com/gorilla/mux"
)

/*
SetupClientPrivateRoutes creates all instances for client private enviroment and calls each router
*/
func SetupClientPrivateRoutes(subRouter *mux.Router) {
	connection := database.GetBankGormConnection()
	rUser := user.NewUserGormRepo(connection)
	sUser := user.NewUserService(rUser)
	cUser := user.NewUserController(sUser)
	user.NewUserPrivateRouter(subRouter.PathPrefix("/user").Subrouter(), cUser)

	rMovement := movement.NewMovementGormRepo(connection)
	sMovement := movement.NewMovementService(rMovement)
	cMovement := movement.NewMovementController(sMovement)
	movement.NewMovementPrivateRouter(subRouter.PathPrefix("/movement").Subrouter(), cMovement)
}

/*
SetupClientPublicRoutes creates all instances for client public enviroment and calls each router
*/
func SetupClientPublicRoutes(subRouter *mux.Router) {
	connection := database.GetBankGormConnection()
	rUser := user.NewUserGormRepo(connection)
	sUser := user.NewUserService(rUser)
	cUser := user.NewUserController(sUser)
	user.NewUserPublicRouter(subRouter.PathPrefix("/user").Subrouter(), cUser)
}
