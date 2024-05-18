package Models

import (
	"time"
)

type User struct {
	ID        uint64 `json:"id" gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	GUID      string `json:"guid" gorm:"uniqueIndex;type:varchar(36)"`
	Username  string `json:"username" gorm:"uniqueIndex;type:varchar(64)"`
	Password  string
	Name      string `json:"name" gorm:"type:varchar(64)"`
	GroupID   uint64 `json:"group_id"`
	Group     Group
	IsAdmin   bool `json:"is_admin"`
	Status    int  `json:"status"`
	Tokens    []Token
	Note      string `json:"note"`
	Info      string `json:"info"`
}
