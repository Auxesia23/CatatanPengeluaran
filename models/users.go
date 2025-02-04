package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password  string `json:"password" gorm:"type:varchar(100);not null"`
	Superuser bool   `json:"is_superuser" gorm:"default:false"`
}
