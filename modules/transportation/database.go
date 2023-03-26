package transportation

import (
	"leisurely/database"
	"strings"
	Models "leisurely/database/models"
)

func CreateTransportation (trans *Models.Transportation) error {
	return database.DB.Create(trans).Error
}

func GetCurrentTransportationID() (int, error) {
	var transportation Models.Transportation
	err := database.DB.Order("transportation_id desc").First(&transportation).Error

	if err != nil && strings.Contains(err.Error(), "record not found") {
		return 1, nil
	}

	return transportation.TransportationID + 1, err
}