package Models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID          uint64 `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Username    string         `json:"username" gorm:"uniqueIndex;type:varchar(64)"`
	Password    string
	Name        string `json:"name" gorm:"type:varchar(64)"`
	IsValidated bool   `json:"is-validated"`
	Tokens      []Token
}
