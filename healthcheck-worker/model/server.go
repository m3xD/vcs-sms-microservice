package model

import "time"

type Server struct {
	ID          string `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"column:name" gorm:"unique"`
	Status      int
	IP          string    `gorm:"column:ipv4" gorm:"unique"`
	CreatedTime time.Time `gorm:"column:created_at"`
	LastUpdated time.Time `gorm:"column:update_at"`
}
