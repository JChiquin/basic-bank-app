package movement

import (
	"bank-service/src/environments/client/resources/entity"
	"bank-service/src/environments/client/resources/interfaces"
	"bank-service/src/libs/dto"
	"bank-service/src/libs/errors"
)

/*
struct that implements IMovementService
*/
type movementService struct {
	rMovement interfaces.IMovementRepository
	rUser     interfaces.IUserRepository
}

/*
NewMovementService creates a new service, receives repository by dependency injection
and returns IMovementService, so it needs to implement all its methods
*/
func NewMovementService(rMovement interfaces.IMovementRepository, rUser interfaces.IUserRepository) interfaces.IMovementService {
	return &movementService{rMovement, rUser}
}

func (s *movementService) IndexByUserID(filterMovements *dto.FilterMovements, pagination *dto.Pagination) ([]entity.Movement, error) {
	if err := filterMovements.Validate(); err != nil {
		return nil, err
	}
	return s.rMovement.IndexByUserID(filterMovements.ParseToMovement(), pagination)
}

func (s *movementService) Create(createMovement *dto.CreateMovement) (*entity.Movement, error) {
	if err := createMovement.Validate(); err != nil {
		return nil, err
	}

	balance, err := s.rUser.GetBalance(createMovement.UserID)
	if err != nil {
		return nil, err
	}

	if createMovement.Amount > balance {
		return nil, errors.ErrInsufficientBalance
	}

	destinationUser, err := s.rUser.FindByAccountNumber(createMovement.AccountNumber)
	if err != nil {
		return nil, err
	}

	if destinationUser.ID == createMovement.UserID {
		return nil, errors.ErrSameWallet
	}

	sourceUser, err := s.rUser.FindByID(createMovement.UserID)
	if err != nil {
		return nil, err
	}

	debitMovement := createMovement.ParseToDebitMovement()
	newMovement, err := s.rMovement.Create(debitMovement)
	if err != nil {
		return nil, err
	}

	creditMovement := createMovement.ParseToCreditMovement(destinationUser.ID, sourceUser.AccountNumber)
	_, err = s.rMovement.Create(creditMovement)
	if err != nil {
		return nil, err
	}

	return newMovement, nil
}
