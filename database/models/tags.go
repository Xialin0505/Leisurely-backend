package models

import (
	//"time"
	//"github.com/lib/pq"
	_ "gorm.io/gorm"
)

type Tags struct {
	//gorm.Model
	TagID        int          `gorm:"primaryKey"`
	Preferences  []Preference `gorm:"foreignKey:TagID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	EventTag     []EventTag   `gorm:"foreignKey:TagID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Description  string       `gorm:"not null;unique"`
	TagIndicator TagIndicator `gorm:"not null"`
}
