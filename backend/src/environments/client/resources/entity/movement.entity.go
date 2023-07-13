package entity

import (
	"time"

	"gorm.io/gorm"
)

type Movement struct {
	ID            int     `json:"id" gorm:"primaryKey" groups:"client"`
	UserID        int     `json:"user_id" groups:""`
	Amount        float64 `json:"amount" groups:"client"`
	Balance       float64 `json:"balance" groups:"client"`
	Multiplier    int     `json:"multiplier" groups:"client"` //-1 debit, 0 neutral and 1 credit
	Description   string  `json:"description" groups:"client"`
	AccountNumber string  `json:"account_number" groups:"client"`

	//Timestampt fields
	CreatedAt time.Time      `json:"created_at" groups:"client"`
	UpdatedAt time.Time      `json:"updated_at" groups:"client"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" groups:""`

	//Associations
	User *User `json:"user,omitempty" groups:"client"`
}
