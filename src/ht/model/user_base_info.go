package model

type BaseInfo struct {
	Model
	Nickname  string `gorm:"type:varchar(50)"`
	AvatarUrl string `gorm:"type:varchar(255)"`
}
