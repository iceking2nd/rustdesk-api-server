package Models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint64 `json:"id" gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	GUID      string         `json:"guid" gorm:"uniqueIndex;type:varchar(64)"`
	Username  string         `json:"username" gorm:"uniqueIndex;type:varchar(64)"`
	Password  string
	Name      string `json:"name" gorm:"type:varchar(64)"`
	IsAdmin   bool   `json:"is_admin"`
	Status    int    `json:"status"`
	Tokens    []Token
	Note      string `json:"note"`
}
