package tag

import (
	//"strings"
	"leisurely/database"
	"leisurely/database/models"

	_ "gorm.io/gorm"
)

func CreateTag(tag *models.Tags) error {
	return database.DB.Create(tag).Error
}

func CreateEventTag (eventtag *models.EventTag) error {
	return database.DB.Create(eventtag).Error
}

func UpdateTag(tag *models.Tags) error {
	if err := database.DB.Save(tag).Error; err != nil {
		return err
	}
	return nil
}

func DeleteTag(tag *models.Tags) error {
	if err := database.DB.Unscoped().Where("tag_id = ?", tag.TagID).Delete(tag).Error; err != nil {
		return err
	}
	return nil
}

func GetTag(tagid int) (*models.Tags, error) {
	var tag models.Tags
	if err := database.DB.Where("tag_id = ?", tagid).Find(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func GetTagByDescription(desc string) (models.Tags, error) {
	var tag models.Tags
	err := database.DB.Where("description = ?", desc).Find(&tag).Error

	if err != nil {
		return tag, err
	}

	if tag.TagID == 0 && tag.TagIndicator == 0 {
		return tag, nil
	}
	return tag, nil
}

func GetCurrentTagID() (int, error) {
	var tag models.Tags
	if err := database.DB.Order("tag_id desc").First(&tag).Error; err != nil {
		return 0, err
	}

	return tag.TagID + 1, nil
}

func GetTagByType(tagtype int) []models.Tags {
	var tag []models.Tags
	database.DB.Table("tags").Where("tag_indicator = ?", tagtype).Scan(&tag)
	return tag
}
