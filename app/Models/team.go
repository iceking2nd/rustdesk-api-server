package Models

import (
	"time"
)

type Team struct {
	ID        uint64 `json:"id" gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	GUID      string `json:"guid" gorm:"uniqueIndex;type:varchar(36)"`
	Name      string `json:"name" gorm:"uniqueIndex;type:varchar(255)"`
	EMail     string `json:"email" gorm:"type:varchar(255)"`
	Note      string `json:"note"`
	Info      string `json:"info" gorm:"type:text"`
}
