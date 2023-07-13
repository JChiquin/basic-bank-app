package dto

import (
	"bank-service/src/environments/client/resources/entity"
	"bank-service/src/utils/constant"
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

type CreateMovement struct {
	UserID        int     `json:"-" schema:"-"` //From JWT
	Description   string  `json:"description" validate:"required,max=100"`
	AccountNumber string  `json:"account_number" validate:"required,numeric,len=20"`
	Amount        float64 `json:"amount" validate:"gt=0"`
}

/*
Validate returns an error if the DTO doesn't pass any of its own validations
*/
func (dto *CreateMovement) Validate() error {
	//Checks for struct validations (see validate tag)
	if err := validator.ValidateStruct(dto); err != nil {
		return err
	}

	return nil
}

func (dto *CreateMovement) ParseToDebitMovement() *entity.Movement {
	movementToFilter := &entity.Movement{
		UserID:        dto.UserID,
		Amount:        dto.Amount,
		AccountNumber: dto.AccountNumber,
		Description:   dto.Description,
		Multiplier:    constant.DebitMultiplier,
	}
	return movementToFilter
}

func (dto *CreateMovement) ParseToCreditMovement(userId int, accountNumber string) *entity.Movement {
	movementToFilter := &entity.Movement{
		Description:   dto.Description,
		Amount:        dto.Amount,
		UserID:        userId,
		AccountNumber: accountNumber,
		Multiplier:    constant.CreditMultiplier,
	}
	return movementToFilter
}
