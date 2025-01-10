package user

import "ht-crm/src/ht/model"

type BaseInfo struct {
	model.Model
	Nickname  string `gorm:"type:varchar(50)"`
	AvatarUrl string `gorm:"type:varchar(255)"`
}
