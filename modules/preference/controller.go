package preference

import (
	"github.com/rs/zerolog/log"
	Models "leisurely/database/models"
)

func UpdateUserPreference(uid int, tagid int) error{
	pref, Geterr := GetUserPreferenceByIDTag(uid, tagid)

	if Geterr != nil {
		log.Error().Stack().Err(Geterr).Msg("PreferenceController: Failed to fetch preference.")
		return Geterr
	}

	if pref.UID == 0{
		var newpref Models.Preference
		newpref.UID = uid
		newpref.TagID = tagid
		newpref.Count = 1
		createErr := CreateUserPreference(&newpref)
		if createErr != nil {
			log.Error().Stack().Err(createErr).Msg("PreferenceController: Failed to create preference.")
			return createErr
		}
		return nil
	} 
	
	err := UpdateUserPreferenceByIDTag(uid, tagid)
	
	if err != nil{
		log.Error().Stack().Err(err).Msg("PreferenceController: Failed to update preference.")
		return err
	}

	return nil
}