package event

import (
	"github.com/gofiber/fiber/v2"
	"leisurely/utils"
	"github.com/rs/zerolog/log"
	"strings"
	"strconv"

	Models "leisurely/database/models"
	Tag "leisurely/modules/tag"
)

func EventRouter(app fiber.Router) {
	app.Post("/insertEvent/:uid", InsertEventAPI)
	app.Delete("/deleteEvent/:eventid", DeleteEventAPI)
}

func InsertEventAPI(c *fiber.Ctx) error{
	var eventdto EventDTO
	UID, UIDerr := strconv.Atoi(c.Params("uid"))

	if (UIDerr != nil){
		return UIDerr
	}

	if err := c.BodyParser(&eventdto); err != nil {
		log.Error().Stack().Err(err).Msg("EventController: Failed to parse client request")
		return err
	}

	if res := utils.ValidateDTO(eventdto); len(res) != 0 {
		return c.Status(400).JSON(res)
	}

	resultEvent := Models.Event{
		EventID: 	eventdto.EventID,     
		PlanID:   	1,
		Cost: 			 eventdto.Cost,
		StartTime:     eventdto.StartTime,
		EndTime:        eventdto.EndTime,
		Type:           eventdto.Type,     
		Address:          strings.Join(eventdto.Address,";"),    
	}

	if len(eventdto.EventTag) == 0{
		log.Error().Msg("EventController: Event has no tags")
	}


	if eventdto.Link != "" {
		resultEvent.Link = eventdto.Link
	}
	
	if eventdto.Venue != "" {
		resultEvent.Venue = eventdto.Venue
	}        

	if eventdto.Description != "" {
		resultEvent.Description = eventdto.Description
	}

	err := CreateEvent(&resultEvent)
	if err != nil{
		log.Error().Stack().Err(err).Msg("EventController: fail to insert event")
	}

	eventID, eventIDErr := GetCurrentEventID()

	if eventIDErr != nil{
		log.Error().Stack().Err(eventIDErr).Msg("EventController: Failed to fetch event ID.")
	}

	var tags []Models.Tags
	for i := 0; i < len(eventdto.EventTag); i++ {
		var onetag Models.Tags
		onetag.Description = eventdto.EventTag[i]
		onetag.TagIndicator = 2;

		if tagErr := Tag.InsertOneTag(eventdto.EventTag[i], eventID-1, UID); tagErr != nil {
			log.Error().Stack().Err(tagErr).Msg("EventController: Failed to insert event tag.")
		}

		tags = append(tags, onetag)
	}

	return nil;
}

func DeleteEventAPI(c *fiber.Ctx) error{
	return nil;
}

func InsertEventDTO(UID int, PlanID int, eventdto EventDTO) (error, string){

	resultEvent := Models.Event{
		EventID: 	eventdto.EventID,     
		Title:		eventdto.Title,
		PlanID:   	PlanID,
		Cost: 			 eventdto.Cost,
		StartTime:     eventdto.StartTime,
		EndTime:        eventdto.EndTime,
		Type:           eventdto.Type,     
		Address:          strings.Join(eventdto.Address,";"),    
	}

	if len(eventdto.EventTag) == 0{
		log.Error().Msg("EventController: Event has no tags")
		return nil, "default"
	}

	returntag := "default"
	if eventdto.Type == 1 {
		returntag = eventdto.EventTag[0]
	}
	

	if eventdto.Link != "" {
		resultEvent.Link = eventdto.Link
	}
	
	if eventdto.Venue != "" {
		resultEvent.Venue = eventdto.Venue
	}        

	if eventdto.Description != "" {
		resultEvent.Description = eventdto.Description
	}

	err := CreateEvent(&resultEvent)
	if err != nil{
		log.Error().Stack().Err(err).Msg("EventController: fail to insert event")
		return err, "default"
	}

	eventID, eventIDErr := GetCurrentEventID()

	if eventIDErr != nil{
		log.Error().Stack().Err(eventIDErr).Msg("EventController: Failed to fetch event ID.")
		return eventIDErr, "default"
	}

	var tags []Models.Tags
	for i := 0; i < len(eventdto.EventTag); i++ {
		var onetag Models.Tags
		onetag.Description = eventdto.EventTag[i]
		onetag.TagIndicator = 2;

		if tagErr := Tag.InsertOneTag(eventdto.EventTag[i], eventID-1, UID); tagErr != nil {
			log.Error().Stack().Err(tagErr).Msg("EventController: Failed to insert event tag.")
		}

		tags = append(tags, onetag)
	}

	return nil, returntag
}