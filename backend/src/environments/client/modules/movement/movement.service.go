package movement

import (
	"bank-service/src/environments/client/resources/entity"
	"bank-service/src/environments/client/resources/interfaces"
	"bank-service/src/libs/dto"
)

/*
struct that implements IMovementService
*/
type movementService struct {
	rMovement interfaces.IMovementRepository
}

/*
NewMovementService creates a new service, receives repository by dependency injection
and returns IMovementService, so it needs to implement all its methods
*/
func NewMovementService(rMovement interfaces.IMovementRepository) interfaces.IMovementService {
	return &movementService{rMovement}
}

func (s *movementService) IndexByUserID(filterMovements *dto.FilterMovements, pagination *dto.Pagination) ([]entity.Movement, error) {
	if err := filterMovements.Validate(); err != nil {
		return nil, err
	}
	return s.rMovement.IndexByUserID(filterMovements.ParseToMovement(), pagination)
}
