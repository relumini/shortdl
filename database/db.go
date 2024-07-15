package database

import (
	"github.com/relumini/shortdl/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=12345678 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	DB = db

	err = DB.AutoMigrate(&models.ChecksumData{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
