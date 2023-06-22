package Models

import (
	"time"
)

type Address struct {
	ID        uint64 `json:"id" gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint64 `json:"user_id" gorm:"uniqueIndex:unique_address"`
	User      User
	Username  string `json:"username" gorm:"type:varchar(128)"`
	ClientID  string `json:"client_id" gorm:"type:varchar(16);uniqueIndex:unique_address"`
	Hostname  string `json:"hostname" gorm:"type:varchar(128)"`
	Platform  string `json:"platform" gorm:"type:varchar(20)"`
	Alias     string `json:"alias" gorm:"type:varchar(20)"`
	Tags      []Tag  `json:"tags" gorm:"many2many:address_tags;"`
}
