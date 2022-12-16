package dto

import (
	"bank-service/src/environments/client/resources/entity"
	"bank-service/src/utils/validator"
)

type FilterMovements struct {
	UserID     int  `schema:"-"` //From JWT
	Multiplier *int `schema:"multiplier" validate:"omitempty,oneof=1 -1"`
}

/*
Validate returns an error if the DTO doesn't pass any of its own validations
*/
func (dto *FilterMovements) Validate() error {
	//Checks for struct validations (see validate tag)
	if err := validator.ValidateStruct(dto); err != nil {
		return err
	}

	return nil
}

func (dto *FilterMovements) ParseToMovement() entity.Movement {
	movementToFilter := entity.Movement{
		UserID: dto.UserID,
	}
	if dto.Multiplier != nil {
		movementToFilter.Multiplier = *dto.Multiplier
	}
	return movementToFilter
}
