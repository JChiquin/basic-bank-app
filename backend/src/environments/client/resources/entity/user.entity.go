package entity

import (
	"bank-service/src/libs/password"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID             int       `json:"id" gorm:"primaryKey" groups:""`
	FirstName      string    `json:"first_name" groups:"client"`
	LastName       string    `json:"last_name" groups:"client"`
	DocumentNumber string    `json:"document_number" groups:"client"`
	BirthDate      time.Time `json:"birth_date" groups:"client"`
	PhoneNumber    string    `json:"phone_number" groups:"client"`
	Email          string    `json:"email" groups:"client"`
	Password       string    `json:"password" groups:""`
	UserType       string    `json:"user_type" groups:""`
	AccountNumber  string    `json:"account_number" groups:"client"`

	CreatedAt time.Time      `json:"created_at" groups:""`
	UpdatedAt time.Time      `json:"updated_at" groups:""`
	DeletedAt gorm.DeletedAt `json:"deleted_at" groups:""`
}

// BeforeCreate is a database hook
func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := password.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword

	//Generating a random account number
	u.AccountNumber = fmt.Sprintf("%020d", time.Now().UnixNano())[:20]

	return nil
}

// BeforeUpdate is a database hook
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	hashedPassword, err := password.HashPassword(tx.Statement.Dest.(User).Password)
	if err != nil {
		return err
	}
	tx.Statement.SetColumn("Password", hashedPassword)
	return nil
}

func (u *User) CheckPassword(plainTextPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainTextPassword))
}
