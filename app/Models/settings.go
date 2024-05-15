package Models

import (
	"gorm.io/gorm"
	"time"
)

type Settings struct {
	ID        uint64 `json:"id" gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Key       string         `json:"key" gorm:"uniqueIndex;type:varchar(64)"`
	Value     string         `json:"value"`
}
