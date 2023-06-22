package Models

import (
	"gorm.io/gorm"
	"time"
)

type Token struct {
	ID          uint64 `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	UserID      uint64         `json:"user_id"`
	User        User
	ClientID    string `json:"client_id" gorm:"index;type:varchar(16)"`
	AccessToken string `json:"access_token" gorm:"uniqueIndex"`
}
