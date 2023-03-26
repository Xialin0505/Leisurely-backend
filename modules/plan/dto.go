package plan

import (
	"time"
	_"leisurely/modules/event"
	_ "leisurely/database/models"
	Event "leisurely/modules/event"
)

type UserConstraintsTagDTO struct {
	Description string `json:"description" form:"description" validate:"required"`
	// 1 is keyword
	// 2 is descriptive
	// 3 is pop-up
	EventType int `json:"eventType" form:"eventType" validate:"required"`
	Count     int `json:"count" form:"count" validate:"required"`
}

type UserConstraintsDTO struct {
	Country   string                  `json:"country" form:"country" validate:"required"`
	Location  string                  `json:"location" form:"location" validate:"required"`
	Date      string                  `json:"date" form:"date" validate:"required"`
	StartTime float64                 `json:"startTime" form:"startTime" validate:"required"`
	EndTime   float64                 `json:"endTime" form:"endTime" validate:"required"`
	Tags      []UserConstraintsTagDTO `json:"tags" form:"tags" validate:"omitempty"`
	Transport string                  `json:"transport" form:"transport" validate:"required"`
	BudgetLevel int						`json:"budgetLevel" form:"budgetLevel" validate:"omitempty`
}

type PlanDTO struct {
	PlanID         int 				`json:"planid" form:"planid" validate:"omitempty"`
	PlanName       string			`json:"planname" form:"planname" validate:"omitempty"`
	UID            int              `json:"uid" form:"uid" validate:"required"`
	Time           time.Time        `json:"time" form:"time" validate:"required"`
	Events         []Event.EventDTO         `json:"events" form:"events" validate:"required"`
	Transportations []float64 `json:"transport" form:"transport" validate:"omitempty"`
	Photo          string			`json:"photo" form:"photo" validate:"omitempty"`
	Country   string                  `json:"country" form:"country" validate:"required"`
	Location  string                  `json:"location" form:"location" validate:"required"`
}