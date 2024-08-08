package service

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func NewDatabase() *gorm.DB {
	dsn := os.Getenv("DB_SOURCE")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
