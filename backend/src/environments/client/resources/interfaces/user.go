package interfaces

import (
	"bank-service/src/environments/client/resources/entity"
	"bank-service/src/libs/dto"
	"net/http"
)

/*
IUserController methods to handle http requests
*/
type IUserController interface {
	Create(response http.ResponseWriter, request *http.Request)
	Login(response http.ResponseWriter, request *http.Request)
	WhoAmI(response http.ResponseWriter, request *http.Request)
	FindByAccountNumber(response http.ResponseWriter, request *http.Request)
	GetBalance(response http.ResponseWriter, request *http.Request)
	UpdatePassword(response http.ResponseWriter, request *http.Request)
}

/*
IUserService methods to implement the bussiness logic
*/
type IUserService interface {
	Create(createUser *dto.CreateUser) (*entity.User, error)
	Login(requestLogin *dto.RequestLogin) (*dto.ResponseLogin, error)
	FindByID(userID int) (*entity.User, error)
	FindByAccountNumber(accountNumber string) (*entity.User, error)
	GetBalance(userId int) (*dto.LastBalance, error)
	UpdatePassword(updatePassword *dto.UpdatePassord) error
}

/*
IUserRepository methods to interact with user entity, independent of ORM
*/
type IUserRepository interface {
	Create(user *entity.User) (*entity.User, error)
	FindByDocumentNumber(documentNumber string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindByID(userID int) (*entity.User, error)
	FindByAccountNumber(accountNumber string) (*entity.User, error)
	GetBalance(userID int) (float64, error)
	UpdatePassword(userID int, newPlainPassword string) error
}
