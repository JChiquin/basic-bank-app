package user

import (
	"bank-service/src/environments/client/resources/entity"
	"bank-service/src/environments/client/resources/interfaces"
	"bank-service/src/libs/dto"
	myErrors "bank-service/src/libs/errors"
	"bank-service/src/libs/logger"
	"bank-service/src/libs/mailer"
	"bank-service/src/libs/password"
	"bank-service/src/utils/validator"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"
)

const (
	passwordResetCodeDigits      = 6
	passwordResetCodeTTL         = 15 * time.Minute
	passwordResetCodeMaxAttempts = 5
)

/*
struct that implements IUserService
*/
type userService struct {
	rUser interfaces.IUserRepository
	sJWT  interfaces.IJWTService
}

/*
NewUserService creates a new service, receives repository by dependency injection
and returns IUserService, so it needs to implement all its methods
*/
func NewUserService(rUser interfaces.IUserRepository, sJWT interfaces.IJWTService) interfaces.IUserService {
	return &userService{rUser, sJWT}
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

	if err := userFound.CheckPassword(requestLogin.Password); err != nil {
		return nil, myErrors.ErrUnauthorized
	}

	jwt, err := s.sJWT.Create(userFound)
	if err != nil {
		return nil, err
	}

	responseLogin := &dto.ResponseLogin{
		User: *userFound,
		JWT:  jwt,
	}
	return responseLogin, nil
}

func (s *userService) RequestPasswordReset(requestPasswordReset *dto.RequestPasswordReset) error {
	if err := requestPasswordReset.Validate(); err != nil {
		return err
	}

	if !mailer.IsConfigured() {
		return myErrors.ErrEmailNotConfigured
	}

	userFound, err := s.rUser.FindByEmail(requestPasswordReset.Email)
	if errors.Is(err, myErrors.ErrNotFound) {
		return nil
	} else if err != nil {
		return err
	}

	code, err := generatePasswordResetCode()
	if err != nil {
		return err
	}

	codeHash, err := password.HashPassword(code)
	if err != nil {
		return err
	}

	passwordResetCode := &entity.PasswordResetCode{
		UserID:    userFound.ID,
		CodeHash:  codeHash,
		ExpiresAt: time.Now().Add(passwordResetCodeTTL),
	}
	if err := s.rUser.CreatePasswordResetCode(passwordResetCode); err != nil {
		return err
	}

	return mailer.SendRecoveryCode(userFound.Email, code)
}

func (s *userService) ResetPassword(resetPassword *dto.ResetPassword) error {
	if err := resetPassword.Validate(); err != nil {
		return err
	}

	userFound, err := s.rUser.FindByEmail(resetPassword.Email)
	if err != nil {
		return myErrors.ErrInvalidPasswordResetCode
	}

	passwordResetCode, err := s.rUser.FindValidPasswordResetCode(userFound.ID)
	if err != nil {
		return myErrors.ErrInvalidPasswordResetCode
	}

	if passwordResetCode.Attempts >= passwordResetCodeMaxAttempts {
		_ = s.rUser.MarkPasswordResetCodeUsed(passwordResetCode.ID)
		return myErrors.ErrPasswordResetCodeAttemptsExceeded
	}

	if err := password.CheckPassword(passwordResetCode.CodeHash, resetPassword.Code); err != nil {
		_ = s.rUser.IncrementPasswordResetCodeAttempts(passwordResetCode.ID)
		if passwordResetCode.Attempts+1 >= passwordResetCodeMaxAttempts {
			_ = s.rUser.MarkPasswordResetCodeUsed(passwordResetCode.ID)
			return myErrors.ErrPasswordResetCodeAttemptsExceeded
		}
		return myErrors.ErrInvalidPasswordResetCode
	}

	if err := s.rUser.UpdatePassword(userFound.ID, resetPassword.NewPassword); err != nil {
		return err
	}

	return s.rUser.MarkPasswordResetCodeUsed(passwordResetCode.ID)
}

func (s *userService) FindByID(userID int) (*entity.User, error) {
	if err := validator.ValidateVar(userID, "user_id", "required,gte=1"); err != nil {
		return nil, err
	}
	return s.rUser.FindByID(userID)
}

func (s *userService) GetBalance(userID int) (*dto.LastBalance, error) {
	if err := validator.ValidateVar(userID, "user_id", "required,gte=1"); err != nil {
		return nil, err
	}

	balance, err := s.rUser.GetBalance(userID)
	if err != nil {
		return nil, err
	}

	lastBalance := &dto.LastBalance{
		Balance:  balance,
		LastTime: time.Now(),
	}

	return lastBalance, nil
}

func (s *userService) FindByAccountNumber(accountNumber string) (*entity.User, error) {
	if err := validator.ValidateVar(accountNumber, "account_number", "required,len=20"); err != nil {
		return nil, err
	}
	return s.rUser.FindByAccountNumber(accountNumber)
}

func (s *userService) UpdatePassword(updatePassword *dto.UpdatePassord) error {
	if err := updatePassword.Validate(); err != nil {
		return err
	}

	userFound, err := s.rUser.FindByID(updatePassword.UserID)
	if err != nil {
		return err
	}

	if err := userFound.CheckPassword(updatePassword.Password); err != nil {
		return myErrors.ErrUnauthorized
	}

	return s.rUser.UpdatePassword(updatePassword.UserID, updatePassword.NewPassword)
}

func generatePasswordResetCode() (string, error) {
	max := big.NewInt(1000000)
	number, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%0*d", passwordResetCodeDigits, number.Int64()), nil
}
