package interfaces

import (
	"bank-service/src/environments/client/resources/entity"
	"bank-service/src/libs/dto"
	"net/http"
)

/*
IMovementController methods to handle http requests
*/
type IMovementController interface {
	Index(response http.ResponseWriter, request *http.Request)
	Create(response http.ResponseWriter, request *http.Request)
}

/*
IMovementService methods to implement the bussiness logic
*/
type IMovementService interface {
	IndexByUserID(filterMovements *dto.FilterMovements, pagination *dto.Pagination) ([]entity.Movement, error)
	Create(createMovement *dto.CreateMovement) (*entity.Movement, error)
}

/*
IMovementRepository methods to interact with movement entity, independent of ORM
*/
type IMovementRepository interface {
	IndexByUserID(movementToFilter entity.Movement, pagination *dto.Pagination) ([]entity.Movement, error)
	Create(newMovement *entity.Movement) (*entity.Movement, error)
}
