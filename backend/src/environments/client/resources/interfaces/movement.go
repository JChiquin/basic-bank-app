package interfaces

import (
	"bank-service/src/environments/client/resources/entity"
	"net/http"
)

/*
IMovementController methods to handle http requests
*/
type IMovementController interface {
	Index(response http.ResponseWriter, request *http.Request)
}

/*
IMovementService methods to implement the bussiness logic
*/
type IMovementService interface {
	IndexByUserID(userID int) ([]entity.Movement, error)
}

/*
IMovementRepository methods to interact with movement entity, independent of ORM
*/
type IMovementRepository interface {
	IndexByUserID(userID int) ([]entity.Movement, error)
}
