package models

import (
	_ "gorm.io/gorm"
)

type Preference struct {
	//gorm.Model
	UID   int `gorm:"primaryKey"`
	TagID int `gorm:"primaryKey"`
	Count int
}
