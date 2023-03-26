package user

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

const (
	YYYYMMDD = "2006-01-02"
)

func IsUsernameTaken(username string) bool {
	_, err := GetUserProfileByName(username)

	if errors.Is(err, gorm.ErrRecordNotFound) == true {
		return false
	}
	return true
}

func getDBErr(err error) error {
	if strings.Contains(err.Error(), "idx_users_username") {
		err = gorm.ErrInvalidValue
	}
	return err
}
