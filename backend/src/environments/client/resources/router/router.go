package router

import (
	"bank-service/src/environments/client/modules/contact"
	"bank-service/src/environments/client/modules/jwt"
	"bank-service/src/environments/client/modules/movement"
	"bank-service/src/environments/client/modules/user"
	"bank-service/src/libs/database"
	jwtMiddleware "bank-service/src/libs/middleware/jwt"
	"bank-service/src/utils/constant"

	"github.com/gorilla/mux"
)

/*
SetupClientPrivateRoutes creates all instances for client private enviroment and calls each router
*/
func SetupClientPrivateRoutes(subRouter *mux.Router) {
	subRouter.Use(jwtMiddleware.NewJWTMiddleware(constant.ClientUserTypes).Handler)

	connection := database.GetBankGormConnection()
	rUser := user.NewUserGormRepo(connection)
	sJWT := jwt.NewJwtService()
	sUser := user.NewUserService(rUser, sJWT)
	cUser := user.NewUserController(sUser)
	user.NewUserPrivateRouter(subRouter.PathPrefix("/user").Subrouter(), cUser)

	rMovement := movement.NewMovementGormRepo(connection)
	sMovement := movement.NewMovementService(rMovement, rUser)
	cMovement := movement.NewMovementController(sMovement)
	movement.NewMovementPrivateRouter(subRouter.PathPrefix("/movement").Subrouter(), cMovement)

	rContact := contact.NewContactGormRepo(connection)
	sContact := contact.NewContactService(rContact, rUser)
	cContact := contact.NewContactController(sContact)
	contact.NewContactPrivateRouter(subRouter.PathPrefix("/contact").Subrouter(), cContact)
}

/*
SetupClientPublicRoutes creates all instances for client public enviroment and calls each router
*/
func SetupClientPublicRoutes(subRouter *mux.Router) {
	connection := database.GetBankGormConnection()
	rUser := user.NewUserGormRepo(connection)
	sJWT := jwt.NewJwtService()
	sUser := user.NewUserService(rUser, sJWT)
	cUser := user.NewUserController(sUser)
	user.NewUserPublicRouter(subRouter.PathPrefix("/user").Subrouter(), cUser)
}
