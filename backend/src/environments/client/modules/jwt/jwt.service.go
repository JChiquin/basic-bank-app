package jwt

import (
	"bank-service/src/environments/client/resources/entity"
	"bank-service/src/environments/client/resources/interfaces"
	"bank-service/src/libs/env"

	"github.com/golang-jwt/jwt/v4"
)

/*
struct that implements IJWTService
*/
type jwtService struct {
}

/*
NewJwtService creates a new service and returns IJWTService, so it needs to implement all its methods
*/
func NewJwtService() interfaces.IJWTService {
	return &jwtService{}
}

func (s *jwtService) Create(user *entity.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"user_type": user.UserType,
	})
	tokenString, err := token.SignedString([]byte(env.JWTSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
