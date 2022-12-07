package interfaces

import "bank-service/src/environments/client/resources/entity"

/*
IJWTService methods to implement the bussiness logic
*/
type IJWTService interface {
	Create(user *entity.User) (string, error)
}
