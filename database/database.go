package database

import (
	"os"

	"github.com/Auxesia23/CatatanPengeluaran/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := os.Getenv("DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{}, &models.Method{}, &models.Category{}, &models.Transaction{})

	return db
}
