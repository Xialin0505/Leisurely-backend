package models

import (
	_ "gorm.io/gorm"
)

type EventTag struct {
	EventID int `gorm:"primaryKey"`
	TagID   int `gorm:"primaryKey"`
}

type Event struct {
	EventID         int        `gorm:"primaryKey"`
	PlanID          int
	Link			string 
	Title			string
	StartTime       float64 
	EndTime         float64
	Type            EventType
	Description     string
	Cost            float64
	Address         string /* deliminator ";" */
	Venue           string
	// constraint
	EventTag        []EventTag `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TransportationS Transportation `gorm:"foreignKey:StartEvent;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TransportationD Transportation `gorm:"foreignKey:DestinationEvent;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
