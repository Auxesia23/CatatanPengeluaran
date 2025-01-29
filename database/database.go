package database

import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"github.com/Auxesia23/CatatanPengeluaran/models"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("catatan_pengeluaran.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{}, &models.Method{}, &models.Category{}, &models.Transaction{})

	return db
}