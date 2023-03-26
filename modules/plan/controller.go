package plan

import (
	"fmt"
	"sync"
	"strings"
	Models "leisurely/database/models"
	Pref "leisurely/modules/preference"
	Event "leisurely/modules/event"
	Transportation "leisurely/modules/transportation"
	Recom "leisurely/recommendation"
	"leisurely/utils"
	"strconv"

	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type PlanResult struct {
	Schedule    Recom.Schedule
	Alternative []Recom.ScheduleEvents
}

func PlanRouter(app fiber.Router) {
	app.Post("/generatePlanFree/:uid", GeneratePlanFree)
	app.Post("/generatePlan/:uid", GeneratePlan)
	app.Post("/confirmPlan", ConfirmPlan)
	app.Post("/deletePlan/:planID", DeletePlanController)
}

func GeneratePlanFree(c *fiber.Ctx) error {
	log.Info().Msg("Getting generate plan request!")

	var userconstraint UserConstraintsDTO

	if err := c.BodyParser(&userconstraint); err != nil {
		log.Error().Stack().Err(err).Msg("PlanController: Failed to parse client request")
		return err
	}

	if res := utils.ValidateDTO(userconstraint); len(res) != 0 {
		return c.Status(400).JSON(res)
	}

	// get userConstraint from json
	uc := Recom.UserConstraints{
		StartTime: userconstraint.StartTime,
		EndTime:   userconstraint.EndTime,
		Date:      userconstraint.Date,
		// flags: willeat, hascar
		Location:  userconstraint.Location,
		Country:   userconstraint.Country,
		Transport: userconstraint.Transport,
		BudgetLevel: userconstraint.BudgetLevel,
	}
	
	var index = strings.Index(uc.Date, "0")
	if index != -1 {
		uc.Date = uc.Date[:index] + uc.Date[index+1:]
	}

	mustHave := 0

	for i := range userconstraint.Tags {
		if userconstraint.Tags[i].EventType == 1 {
			mustHave++
		}
	}
	uc.NumMustHave = mustHave

	uid, err := strconv.Atoi(c.Params("uid"))
	if err != nil {
		fmt.Printf("cannot fetch uid")
		return err
	}

	// fetch tags
	preference := Pref.GetUserPreferenceByID(uid)

	var initialTag []Recom.InitialTag
	var oneTag Recom.InitialTag
	for i := range preference {
		oneTag.Description = preference[i].Description
		oneTag.EventType = preference[i].TagIndicator
		oneTag.Count = preference[i].Count
		oneTag.SearchingTag = false
		initialTag = append(initialTag, oneTag)
	}

	for i := range userconstraint.Tags {
		oneTag.Count = userconstraint.Tags[i].Count
		oneTag.EventType = userconstraint.Tags[i].EventType
		oneTag.Description = userconstraint.Tags[i].Description
		oneTag.SearchingTag = true
		initialTag = append(initialTag, oneTag)
	}

	uc.Tags = initialTag

	var recEvent []Recom.ESTABLISHMENT
	var popEvent []Recom.EVENTTAGGED

	// threading
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		recEvent = Recom.RecEventGenerator(uc)
		wg.Done()
	}()

	go func() {
		popEvent = Recom.PopEventGenerator(uc)
		wg.Done()
	}()

	wg.Wait()

	// fmt.Println(recEvent)
	newSechdule, alternativeEvent := Recom.ItineraryManager(uc, popEvent, recEvent)

	var planresult PlanResult
	planresult.Schedule = newSechdule
	planresult.Alternative = alternativeEvent

	return c.JSON(fiber.Map{"success": true, "message": "Plan Generated", "data": planresult})
}

func GeneratePlanTest (c *fiber.Ctx) error {
	var planresult PlanResult

	jsonFile, err := os.Open("./API/planResponse.json")
	if err != nil {
		fmt.Println(err)
		return err
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &planresult)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.JSON(fiber.Map{"success": true, "message": "Plan Generated", "data": planresult})
}

func GeneratePlan(c *fiber.Ctx) error {
	// send through socket to call recommendation system
	var planresult PlanResult

	// num must have = number of tags passed in
	return c.JSON(fiber.Map{"success": true, "message": "Plan Generated", "data": planresult})
}

func ConfirmPlan(c *fiber.Ctx) error {
	var planForm PlanDTO

	if err := c.BodyParser(&planForm); err != nil {
		log.Error().Stack().Err(err).Msg("PlanController: Failed to parse client request")
		return err
	}

	if res := utils.ValidateDTO(planForm); len(res) != 0 {
		return c.Status(400).JSON(res)
	}

	resultPlan := Models.Plan{
		UID:           planForm.UID,
		Time:          planForm.Time,
		Country:			planForm.Country,
		Location:      planForm.Location,
	}

	resultPlan.PlanID,_ = GetCurrentPlanID()
	
	if planForm.PlanName != ""{
		resultPlan.PlanName = planForm.PlanName
	}
	
	if planForm.Photo != ""{
		resultPlan.Photo = planForm.Photo
	}

	planid, err := CreatePlan(&resultPlan)

	if (err != nil || planid == 0){
		return c.Status(409).JSON(fiber.Map{"success": false, "message": "Plan Fail to Insert Into Database"})
	}
	
	var eventsid []int
	tag := "default"

	for i := range(planForm.Events){
		insertErr, oneTag := Event.InsertEventDTO(planForm.UID, resultPlan.PlanID, planForm.Events[i])
		if insertErr != nil {
			// remove the plan
			return c.Status(409).JSON(fiber.Map{"success": false, "message": "Plan Event Fail to Insert Into Database"})
		}

		if tag == "default" && oneTag != "default"{
			tag = oneTag
		}

		currentEID, curErr := Event.GetCurrentEventID()
		if curErr != nil {
			return c.Status(409).JSON(fiber.Map{"success": false, "message": "Plan Event Fail to Insert Into Database"})
		}

		eventsid = append(eventsid, currentEID-1)
	}

	if len(eventsid) > 1 {
		Transportation.InsertTransportation(resultPlan.PlanID, planForm.Transportations, eventsid)
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Plan Confirmed", "tag": tag, "planID": planid})
}

func DeletePlanController(c *fiber.Ctx) error {
	planid, err := strconv.Atoi(c.Params("planID"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "cannot get plan id"})
	}

	plan, err := GetPlanByPlanID(planid)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "cannot get plan"})
	}

	DeletePlanByID(plan)
	return c.Status(200).JSON(fiber.Map{"success": true, "message": "deleted plan"})
}
