package Models

import "time"

type Tag struct {
	ID        uint64 `json:"id" gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint64 `json:"user_id" gorm:"uniqueIndex:unique_tag"`
	Name      string `json:"name" gorm:"uniqueIndex:unique_tag;type:varchar(64)"`
}
