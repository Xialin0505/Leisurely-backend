package tag

import (
	"github.com/rs/zerolog/log"

	Models "leisurely/database/models"
	Prefs "leisurely/modules/preference"
)

func InsertOneTag(description string, eventID int, UID int) error {
	var newtag Models.Tags
	
	existingTag, err := GetTagByDescription(description)
	if err != nil {
		log.Error().Stack().Err(err).Msg("TagController: get tag ID by description error")
	}

	var et Models.EventTag
	et.EventID = eventID
	et.TagID = existingTag.TagID
	
	if existingTag.TagID == 0 {
		tagid, tagDatabaseErr := GetCurrentTagID()
		if tagDatabaseErr != nil {
			log.Error().Stack().Err(tagDatabaseErr).Msg("TagController: get tag ID error")
		}

		newtag.TagID = tagid
		newtag.Description = description
		newtag.TagIndicator = 2
		err = CreateTag(&newtag)
		if err != nil {
			log.Error().Stack().Err(err).Msg("TagController: cannot create new tag")
			return err
		}
		et.TagID = tagid
	}

	err = CreateEventTag(&et)

	if err != nil{
		log.Error().Stack().Err(err).Msg("TagController: cannot create event tag")
		return err
	}
	
	err = Prefs.UpdateUserPreference(UID, et.TagID)
	if err != nil {
		log.Error().Stack().Err(err).Msg("TagController: cannot create user preference")
		return err
	}

	return nil
}