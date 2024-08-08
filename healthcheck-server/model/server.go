package model

import "time"

type Server struct {
	ID          string `gorm:"primary_key" gorm:"autoIncrement"`
	Name        string `gorm:"column:server_name" gorm:"unique"`
	Status      int
	IP          string
	CreatedTime time.Time
	LastUpdated time.Time
}
