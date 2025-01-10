package model

import "time"

type Model struct {
	ID         uint      `gorm:"primarykey"`
	CreateTime time.Time `gorm:"autoCreateTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime"`
}
