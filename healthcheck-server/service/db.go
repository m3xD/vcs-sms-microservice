package service

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"healthcheck-server/util"
	"log"
)

func NewDatabase() *gorm.DB {
	config, err := util.LoadConfig("..")

	if err != nil {
		log.Fatal("cannot load config")
	}

	dsn := config.DBSource
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
