package dto

import (
	"bank-service/src/environments/client/resources/entity"
	"bank-service/src/utils/validator"
	"fmt"
)

type FilterContacts struct {
	Alias  *string `json:"alias" schema:"alias" validate:"omitempty,max=20"`
	UserID int     `json:"-" validate:"required"` //Came from JWT
}

/*
Validate returns an error if the DTO doesn't pass any of its own validations
*/
func (dto *FilterContacts) Validate() error {
	//Checks for struct validations (see validate tag)
	if err := validator.ValidateStruct(dto); err != nil {
		return err
	}

	return nil
}

// LikeValue returns Alias to be used in like query
func (dto *FilterContacts) LikeValue() string {
	if dto.Alias != nil {
		return fmt.Sprint(`%`, *dto.Alias, `%`)
	}
	return "%"
}

type CreateContact struct {
	UserID        int     `json:"-" validate:"required"` //Came from JWT
	Alias         string  `json:"alias" validate:"required,max=20"`
	AccountNumber string  `json:"account_number" validate:"required,max=20"`
	Description   *string `json:"description" validate:"omitempty,max=100"`
}

/*
Validate returns an error if the DTO doesn't pass any of its own validations
*/
func (dto *CreateContact) Validate() error {
	//Checks for struct validations (see validate tag)
	if err := validator.ValidateStruct(dto); err != nil {
		return err
	}

	return nil
}

func (dto *CreateContact) ParseToEntity() *entity.Contact {
	return &entity.Contact{
		Alias:         dto.Alias,
		UserID:        dto.UserID,
		AccountNumber: dto.AccountNumber,
		Description:   dto.Description,
	}
}

type UpdateContact struct {
	UserID      int     `json:"-" validate:"required"` //Came from JWT
	ContactID   int     `json:"-" validate:"required"` //Came from path variable
	Alias       string  `json:"alias" validate:"required,max=20"`
	Description *string `json:"description" validate:"omitempty,max=100"`
}

/*
Validate returns an error if the DTO doesn't pass any of its own validations
*/
func (dto *UpdateContact) Validate() error {
	//Checks for struct validations (see validate tag)
	if err := validator.ValidateStruct(dto); err != nil {
		return err
	}

	return nil
}

func (dto *UpdateContact) ParseToEntity() *entity.Contact {
	return &entity.Contact{
		Alias:       dto.Alias,
		UserID:      dto.UserID,
		Description: dto.Description,
	}
}

type FilterOneContact struct {
	UserID    int `json:"-" validate:"required"` //Came from JWT
	ContactID int `json:"-" validate:"required"` //Came from path variable
}

/*
Validate returns an error if the DTO doesn't pass any of its own validations
*/
func (dto *FilterOneContact) Validate() error {
	//Checks for struct validations (see validate tag)
	if err := validator.ValidateStruct(dto); err != nil {
		return err
	}

	return nil
}
