package event

import (
	"leisurely/database"
	Models "leisurely/database/models"
)

//batch
func CreateEvent(event *Models.Event) error {
	return database.DB.Create(event).Error
}

func UpdateEvent(event *Models.Event) error {
	if err := database.DB.Save(event).Error; err != nil {
		return err
	}
	return nil
}

//batch
func DeleteEvent(event *Models.Event) error {
	err := database.DB.Unscoped().Where("event_id = ?", event.EventID).Delete(event).Error
	return err
}

//batch
func GetEvent(eventid int) (*Models.Event, error) {
	var event Models.Event
	err := database.DB.Where("event_id = ?", eventid).Find(&event).Error

	return &event, err
}

func GetCurrentEventID() (int, error){
	var event Models.Event
	if err := database.DB.Order("event_id desc").First(&event).Error; err != nil {
		return 0, err
	}

	return event.EventID + 1, nil
}