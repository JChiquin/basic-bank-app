package entity

import (
	"bank-service/src/libs/password"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID             int       `json:"id" gorm:"primaryKey" groups:"client"`
	FirstName      string    `json:"first_name" groups:"client"`
	LastName       string    `json:"last_name" groups:"client"`
	DocumentNumber string    `json:"document_number" groups:"client"`
	BirthDate      time.Time `json:"birth_date" groups:"client"`
	PhoneNumber    string    `json:"phone_number" groups:"client"`
	Email          string    `json:"email" groups:"client"`
	Password       string    `json:"password" groups:""`
	UserType       string    `json:"user_type" groups:""`

	CreatedAt time.Time      `json:"created_at" groups:""`
	UpdatedAt time.Time      `json:"updated_at" groups:""`
	DeletedAt gorm.DeletedAt `json:"deleted_at" groups:""`
}

//BeforeCreate is a database hook
func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := password.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}

func (u *User) CheckPassword(plainTextPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainTextPassword))
}
