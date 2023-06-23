package Models

type ActivateToken struct {
	Token    string `json:"token" gorm:"primarykey;type:varchar(40)"`
	Username string `json:"username" gorm:"uniqueIndex;type:varchar(64);not null"`
	UserID   uint64 `json:"userID"`
	User     User
}
