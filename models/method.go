package models

import "gorm.io/gorm"

type Method struct {
	gorm.Model
	Name string `json:"name" gorm:"not null;uniqueIndex"`
}