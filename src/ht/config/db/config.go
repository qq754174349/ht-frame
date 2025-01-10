package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var gormConfig = &gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true,
	},
}
