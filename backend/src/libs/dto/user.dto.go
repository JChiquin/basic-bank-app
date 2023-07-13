package dto

import (
	"bank-service/src/environments/client/resources/entity"
	"bank-service/src/utils/constant"
	"bank-service/src/utils/validator"
	"time"
)

type CreateUser struct {
	FirstName      string    `json:"first_name" validate:"required,max=40"`
	LastName       string    `json:"last_name" validate:"required,max=40"`
	DocumentNumber string    `json:"document_number" validate:"required,max=20"`
	BirthDate      time.Time `json:"birth_date" validate:"required"`
	PhoneNumber    string    `json:"phone_number" validate:"required,max=20"`
	Email          string    `json:"email" validate:"required,email"`
	Password       string    `json:"password" validate:"required,min=8,max=16"`
}

/*
Validate returns an error if the DTO doesn't pass any of its own validations
*/
func (dto *CreateUser) Validate() error {
	//Checks for struct validations (see validate tag)
	if err := validator.ValidateStruct(dto); err != nil {
		return err
	}

	return nil
}

func (dto *CreateUser) ParseToUser() *entity.User {
	return &entity.User{
		FirstName:      dto.FirstName,
		LastName:       dto.LastName,
		DocumentNumber: dto.DocumentNumber,
		BirthDate:      dto.BirthDate,
		PhoneNumber:    dto.PhoneNumber,
		Email:          dto.Email,
		Password:       dto.Password,
		UserType:       constant.UserTypeClient,
	}
}

type RequestLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=16"`
}

/*
Validate returns an error if the DTO doesn't pass any of its own validations
*/
func (dto *RequestLogin) Validate() error {
	//Checks for struct validations (see validate tag)
	if err := validator.ValidateStruct(dto); err != nil {
		return err
	}

	return nil
}

type ResponseLogin struct {
	entity.User
	JWT string `json:"jwt" groups:"client,admin"`
}

type JWTContext struct {
	UserID   int    `json:"-"`
	UserType string `json:"-"`
}

type LastBalance struct {
	Balance  float64   `json:"balance" groups:"client"`
	LastTime time.Time `json:"last_time" groups:"client"`
}

type UpdatePassord struct {
	UserID      int    `json:"-" validate:"required"` //came from JWT
	Password    string `json:"password" validate:"required,min=8,max=16"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=16"`
}

/*
Validate returns an error if the DTO doesn't pass any of its own validations
*/
func (dto *UpdatePassord) Validate() error {
	//Checks for struct validations (see validate tag)
	if err := validator.ValidateStruct(dto); err != nil {
		return err
	}

	return nil
}
