package entity

import (
	"time"

	"gorm.io/gorm"
)

type PasswordResetCode struct {
	ID        int            `json:"id" gorm:"primaryKey" groups:""`
	UserID    int            `json:"user_id" groups:""`
	CodeHash  string         `json:"code_hash" groups:""`
	ExpiresAt time.Time      `json:"expires_at" groups:""`
	Attempts  int            `json:"attempts" groups:""`
	UsedAt    *time.Time     `json:"used_at" groups:""`
	CreatedAt time.Time      `json:"created_at" groups:""`
	UpdatedAt time.Time      `json:"updated_at" groups:""`
	DeletedAt gorm.DeletedAt `json:"deleted_at" groups:""`

	User *User `json:"user,omitempty" groups:""`
}
