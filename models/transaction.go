package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	UserId      uint `json:"user_id" gorm:"ForeignKey:UserId not null"`
	CategoryID  uint `json:"category_id" gorm:"ForeignKey:CategoryID not null"`
	MethodID    uint `json:"method_id" gorm:"ForeignKey:MethodID not null"`
	Amount      float64 `json:"amount" gorm:"not null"`
	Description string `json:"description" gorm:"not null"`
	Date        string `json:"date" gorm:"not null"`
}