package models

import "gorm.io/gorm"

type Method struct {
	gorm.Model
	Name string `json:"name" gorm:"type:varchar(100);not null;uniqueIndex"`
}