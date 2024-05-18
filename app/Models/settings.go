package Models

import (
	"time"
)

type Settings struct {
	ID        uint64 `json:"id" gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Key       string `json:"key" gorm:"uniqueIndex;type:varchar(64)"`
	Value     string `json:"value"`
}
