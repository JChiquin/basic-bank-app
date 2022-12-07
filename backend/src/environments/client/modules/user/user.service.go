package user

import (
	"bank-service/src/environments/client/resources/entity"
	"bank-service/src/environments/client/resources/interfaces"
	"bank-service/src/libs/dto"
	myErrors "bank-service/src/libs/errors"
	"bank-service/src/libs/logger"
	"bank-service/src/utils/validator"
	"errors"
)

/*
struct that implements IUserService
*/
type userService struct {
	rUser interfaces.IUserRepository
}

/*
NewUserService creates a new service, receives repository by dependency injection
and returns IUserService, so it needs to implement all its methods
*/
func NewUserService(rUser interfaces.IUserRepository) interfaces.IUserService {
	return &userService{rUser}
}

func (s *userService) Create(createUser *dto.CreateUser) (*entity.User, error) {
	if err := createUser.Validate(); err != nil {
		return nil, err
	}

	userToCreate := createUser.ParseToUser()

	_, err := s.rUser.FindByEmail(userToCreate.Email)
	if !errors.Is(err, myErrors.ErrNotFound) {
		logger.GetInstance().WithError(err).Warn("Error getting user by email or the user exists")
		return nil, myErrors.ErrUserExists
	}

	_, err = s.rUser.FindByDocumentNumber(userToCreate.DocumentNumber)
	if !errors.Is(err, myErrors.ErrNotFound) {
		logger.GetInstance().WithError(err).Warn("Error getting user by document number or the user exists")
		return nil, myErrors.ErrUserExists
	}

	return s.rUser.Create(userToCreate)
}

func (s *userService) Login(requestLogin *dto.RequestLogin) (*dto.ResponseLogin, error) {
	if err := requestLogin.Validate(); err != nil {
		return nil, err
	}

	userFound, err := s.rUser.FindByEmail(requestLogin.Email)
	if err != nil {
		return nil, myErrors.ErrUnauthorized
	}

	if userFound.Password != requestLogin.Password {
		return nil, myErrors.ErrUnauthorized
	}

	responseLogin := &dto.ResponseLogin{
		User: *userFound,
	}
	return responseLogin, nil
}

func (s *userService) FindByID(userID int) (*entity.User, error) {
	if err := validator.ValidateVar(userID, "user_id", "required,gte=1"); err != nil {
		return nil, err
	}
	return s.rUser.FindByID(userID)
}