package models

import (

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID      uint      `json:"user_id" gorm:"not null"`
	CategoryID  uint      `json:"category_id" gorm:"not null"`
	MethodID    uint      `json:"method_id" gorm:"not null"`
	Amount      float64   `json:"amount" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	Date        string    `json:"date" gorm:"type:date;not null"`
}
