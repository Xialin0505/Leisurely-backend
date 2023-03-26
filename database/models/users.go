package models

import (
	"time"
	//"github.com/lib/pq"
	_ "gorm.io/gorm"
)

type User struct {
	//gorm.Model
	// user ID
	UID         int          `gorm:"primaryKey"`
	Preferences []Preference `gorm:"foreignKey:UID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Plans       []Plan       `gorm:"foreignKey:UID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name        string
	Email       string    `gorm:"not null"`
	UserName    string    `gorm:"index;not null"`
	Password    string    `gorm:"not null"`
	Birthday    time.Time `gorm:"type:date;not null"`
	Gender      Gender    `gorm:"not null"`
	PhotoUrl    string
	Occupation  string
}
