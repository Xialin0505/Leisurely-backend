package models

import (

	//"github.com/lib/pq"

	_ "gorm.io/gorm"
)

type Transportation struct {
	TransportationID int `gorm:"primaryKey"`
	PlanID           int
	Duration		float64`gorm:"not null"`
	StartEvent       int `gorm:"not null"`
	DestinationEvent int `gorm:"not null"`
}
