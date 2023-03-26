package event

import (
	Models "leisurely/database/models"
)

type EventDTO struct {
	EventID         int  
	Title			string `json:"title" form:"title" validate:"omitempty"`      
	Link			string `json:"link" form:"link" validate:"omitempty"`
	EventTag        []string `json:"eventtags" form:"eventtags" validate:"omitempty"`
	StartTime       float64 `json:"starttime" form:"startime" validate:"required"`
	EndTime         float64 `json:"endtime" form:"endtime" validate:"required"`
	Type            Models.EventType `json:"type" form:"type" validate:"required"`
	Description     string		`json:"description" form:"description" validate:"omitempty"`
	Cost            float64			`json:"cost" form:"cost" validate:"omitempty" default:"-1"`
	Address         []string /* deliminator ";" */ `json:"address" form:"address" validate:"required"`
	Venue           string		`json:"venue" form:"venue" validate:"omitempty"`
}