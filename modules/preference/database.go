package preference

import (
	"leisurely/database"
	"leisurely/database/models"

	"gorm.io/gorm"
)

type PreferenceResult struct {
	UID          int
	TagID        int
	Count        int
	Description  string
	TagIndicator int
}

func CreateUserPreference(preference *models.Preference) error {
	return database.DB.Create(preference).Error
}

func GetUserPreferenceByID(uid int) []PreferenceResult {
	var prefs []PreferenceResult
	database.DB.Table("preferences").
		Select("preferences.uid, preferences.tag_id, preferences.count, tags.description, tags.tag_indicator").
		Where("tags.tag_indicator = 1 AND uid = ?", uid).
		Joins("join tags ON preferences.tag_id = tags.tag_id").
		Order("preferences.count desc").
		Limit(5).
		Scan(&prefs)

	var prefsDes []PreferenceResult
	database.DB.Table("preferences").
		Select("preferences.uid, preferences.tag_id, preferences.count, tags.description, tags.tag_indicator").
		Where("tags.tag_indicator != 1 AND uid = ?", uid).
		Joins("join tags ON preferences.tag_id = tags.tag_id").
		Scan(&prefsDes)

	if prefs == nil && prefsDes == nil {
		return []PreferenceResult{}
	}

	for i := range prefsDes {
		prefs = append(prefs, prefsDes[i])
	}

	return prefs
}

func UpdateUserPreferences(prefs []models.Preference) error {
	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(prefs); i++ {
			err := tx.Save(&prefs[i]).Error

			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func DeleteUserPreferenceByID(uid int, tagid int) {

}

func GetUserPreferenceByIDTag (uid int, tagid int) (models.Preference, error) {
	var pref models.Preference

	err := database.DB.Where("uid = ? AND tag_id = ?", uid, tagid).Find(&pref).Error
	if err != nil{
		return pref, err
	}

	return pref, nil
}

func UpdateUserPreferenceByIDTag (uid int, tagid int) error {
	pref, err := GetUserPreferenceByIDTag(uid, tagid)
	if err != nil {
		return err
	}
	
	pref.Count = pref.Count + 1
	
	if err = database.DB.Save(pref).Error; err != nil {
		return err
	}
	
	return nil
}