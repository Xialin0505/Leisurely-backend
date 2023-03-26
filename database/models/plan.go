package models

import (
	"time"

	_ "github.com/lib/pq"
	_ "gorm.io/gorm"
)

type Plan struct {
	//gorm.Model
	//metadata of plan
	PlanID         int `gorm:"primaryKey"`
	PlanName       string
	UID            int              `gorm:"not null;foreignKey:UID"`
	Time           time.Time        `gorm:"type:time"`
	Events         []Event          `gorm:"foreignKey:PlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Transportations []Transportation `gorm:"foreignKey:PlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Photo          string
	Country			string
	Location			string			
}
