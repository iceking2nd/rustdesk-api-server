package Models

import (
	"time"
)

type Group struct {
	ID           uint64 `json:"id" gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	GUID         string `json:"guid" gorm:"uniqueIndex;type:varchar(36)"`
	TeamID       uint64 `json:"team_id" gorm:"index"`
	Team         Team
	Name         string   `json:"name" gorm:"uniqueIndex;type:varchar(255)"`
	Note         string   `json:"note" gorm:"type:varchar(255)"`
	Info         string   `json:"info" gorm:"type:text"`
	AccessTo     []*Group `json:"access_to" gorm:"many2many:group_access_to_groups;"`
	AccessedFrom []*Group `json:"accessed_from" gorm:"many2many:accessed_from_groups;"`
}
