package user

import (
	_ "fmt"
	"time"
	"leisurely/database"
	"leisurely/database/models"
)

type DBErr string

type TransportationResult struct {
	Duration 		float64 `gorm:"column:duration"`
	StartEvent      string `gorm:"column:originevent"`
	EndEvent        string `gorm:"column:destevent"`
}

type EventTagDes struct {
	EventID 	int `gorm:"column:event_id"`
	Description string  `gorm:"column:description"`
}

type FetchEvent struct {
	EventID         int       
	PlanID          int
	Link			string 
	Title			string
	StartTime       float64 
	EndTime         float64
	Type            models.EventType
	Description     string
	Cost            float64
	Address         string /* deliminator ";" */
	Venue           string
	// constraint
	EventTagDes        []EventTagDes 
}

type PlanResult struct {
	PlanID         int 
	PlanName       string
	UID            int           
	Time           time.Time 
	Events         []FetchEvent 
	EventNum		int
	Transportation []TransportationResult
	Photo          string
	Country			string
	Location		string
}

const (
	DUPLICATE_USER_NAME DBErr = "duplicated user name"
)

func CreateUserProfile(user *models.User) error {
	return database.DB.Create(user).Error
}

func UpdateUserProfile(user *models.User) error {
	if err := database.DB.Save(user).Error; err != nil {
		return getDBErr(err)
	}
	return nil
}

func GetUserProfileByID(uid int) (*models.User, error) {
	var user models.User
	err := database.DB.Where("UID = ?", uid).First(&user).Error
	return &user, err
}

func GetUserProfileByName(username string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("user_name = ?", username).First(&user).Error

	return &user, err
}

func DeleteUserProfile(user *models.User) error {
	err := database.DB.Unscoped().Where("UID = ?", user.UID).Delete(user).Error
	return err
}

func GetUserLogin(email string, password string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("email = ? AND password = ?", email, password).First(&user).Error

	return &user, err
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error

	return &user, err
}

func LoadINitialData() {

}

func GetCurrentUserID() (int, error) {
	var user models.User
	err := database.DB.Order("UID desc").First(&user).Error

	return user.UID + 1, err
}

func GetPlanByUID(uid int) ([]PlanResult, error) {
	var planresult []PlanResult
	var plans []models.Plan
	err := database.DB.Where("UID = ?", uid).Find(&plans).Error

	if err != nil {
		return []PlanResult{}, err
	}

	for i:=0; i < len(plans); i++{
		// fill in fields
		eachplan := PlanResult {
			PlanID: plans[i].PlanID,
			PlanName: plans[i].PlanName,
			UID: plans[i].UID,
			Time: plans[i].Time,
			Photo: plans[i].Photo,
			Country: plans[i].Country,
			Location: plans[i].Location,
		}

		// fetch event
		var events []models.Event
		err = database.DB.Where("plan_id = ?", plans[i].PlanID).Find(&events).Error
		if err != nil {
			eachplan.Events = []FetchEvent{}
		}else{
			
			// fetch eventtag
			for j:=0; j < len(events); j++ {
				var tags []EventTagDes

				eachEvent := FetchEvent {
					EventID:      events[j].EventID,  
					PlanID:          events[j].EventID, 
					Link:			events[j].Link,  
					Title:			events[j].Title,
					StartTime:       events[j].StartTime, 
					EndTime:         events[j].EndTime,
					Type:            events[j].Type, 
					Description:     events[j].Description,
					Cost:            events[j].Cost,
					Address:         events[j].Address,
					Venue:           events[j].Venue,
				}

				err = database.DB.Table("event_tags").
				Select("description, event_id").
				Joins("join tags ON event_tags.tag_id = tags.tag_id").
				Where("event_id = ?", events[j].EventID).Find(&tags).Error

				if err != nil {
					eachEvent.EventTagDes = []EventTagDes{}
				}else{
					eachEvent.EventTagDes = tags
				}
				eachplan.Events = append(eachplan.Events, eachEvent)
			}
		}

		// fetch transportation
		var trans []TransportationResult
		database.DB.Table("transportations").
		Select("duration, startevent.title as originevent, destinationevent.title as destevent").
		Joins("join events startevent ON transportations.start_event = startevent.event_id").
		Joins("join events destinationevent ON transportations.destination_event = destinationevent.event_id").
		Where("transportations.plan_id = ?", plans[i].PlanID).
		Scan(&trans)

		eachplan.Transportation = trans
		eachplan.EventNum = len(events)
		planresult = append(planresult, eachplan)

	}

	return planresult, err
}