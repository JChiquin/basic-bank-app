package entity

import (
	"time"

	"gorm.io/gorm"
)

type Contact struct {
	ID            int     `json:"id" gorm:"primaryKey" groups:"client"`
	Alias         string  `json:"alias" groups:"client"`
	UserID        int     `json:"user_id" groups:""`
	Description   *string `json:"description" groups:"client"`
	AccountNumber string  `json:"account_number" groups:"client"`

	//Timestampt fields
	CreatedAt time.Time      `json:"created_at" groups:"client"`
	UpdatedAt time.Time      `json:"updated_at" groups:"client"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" groups:""`

	//Associations
	User *User `json:"user,omitempty" gorm:"foreignKey:AccountNumber;references:AccountNumber" groups:"client"`
}
