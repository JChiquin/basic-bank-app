package entity

import (
	"time"

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
