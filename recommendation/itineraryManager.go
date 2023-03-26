package recommendation

// we have forsaken restaurants in general
// recommend 2 restaurants since their time is flexible
// sort recurring event by rating

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"
)

// enum for fixed time of recurring events

type ScheduleEvents struct {
	EventInfoRecur ESTABLISHMENT
	EventInfoPop   EVENTTAGGED
	// 1 for recur, 2 for pop
	Type             int
	Position_x       float64
	Position_y       float64
	InitialStartTime float64
	InitialEndTime   float64
}

type Schedule struct {
	MaxEvents int
	// the index of the popup event, -1 if none
	FixedSlot int
	// this array should be ordered
	Events  []ScheduleEvents
	TMatrix TransportationMatrix
}

type InitialTag struct {
	Description string
	// 1 is keyword
	// 2 is descriptive
	// 3 is pop-up
	EventType    int
	Count        int
	SearchingTag bool
}

type UserConstraints struct {
	StartTime float64
	EndTime   float64
	Date      string
	// flags: willeat, hascar
	Location    string
	Country     string
	// State		string
	Tags        []InitialTag
	NumMustHave int
	Transport   string
	BudgetLevel int
}

func removeDuplicateStr(events []ESTABLISHMENT) []ESTABLISHMENT {
    allKeys := make(map[string]bool)
    list := []ESTABLISHMENT{}
    for _, item := range events {
		if _, value := allKeys[item.PlaceID]; !value {
			allKeys[item.PlaceID] = true
			list = append(list, item)
		}
    }
    return list
}

func RandomDuration(lowerB float64, upperB float64) float64 {
	return math.Round((rand.Float64()*(upperB-lowerB))+lowerB/float64(0.05)) * float64(0.05)
}

// struct for schedule

// priority: time and preference
//
func ItineraryManager(userConstraints UserConstraints, popEvents []EVENTTAGGED, recurringEvents []ESTABLISHMENT) (Schedule, []ScheduleEvents) {
	recurringEvents = removeDuplicateStr(recurringEvents)
	var finalSchedule Schedule
	var alternatives []ScheduleEvents
	//MAXITERATION := 100
	// iteration
	var popEvent []ScheduleEvents
	var recEvent []ScheduleEvents

	for i := range popEvents {
		oneEvent := ScheduleEvents{
			EventInfoPop: popEvents[i],
			Type:         2,
			Position_x:   float64(popEvents[i].Coordinate.Latitude),
			Position_y:   float64(popEvents[i].Coordinate.Longitude),
		}

		popEvent = append(popEvent, oneEvent)
	}

	for i := range recurringEvents {
		var addr []string
		addr = append(addr, recurringEvents[i].FormattedAddress)
		recurringEvents[i].Address = addr

		oneEvent := ScheduleEvents{
			EventInfoRecur: recurringEvents[i],
			Type:           1,
			Position_x:     float64(recurringEvents[i].Geometry.Location.Latitude),
			Position_y:     float64(recurringEvents[i].Geometry.Location.Longitude),
		}

		recEvent = append(recEvent, oneEvent)
	}

	// fmt.Printf("%+v\n", recEvent)
	// fmt.Printf("%+v\n", popEvent)

	if len(popEvent) != 0 {
		sort.SliceStable(popEvent, func(i, j int) bool {
			return popEvent[i].EventInfoPop.Score > popEvent[j].EventInfoPop.Score
		})
	}

	sort.SliceStable(recEvent, func(i, j int) bool {
		if recEvent[i].EventInfoRecur.Score == recEvent[j].EventInfoRecur.Score {
			return recEvent[i].EventInfoRecur.Rating > recEvent[j].EventInfoRecur.Rating
		}
		return recEvent[i].EventInfoRecur.Score > recEvent[j].EventInfoRecur.Score
	})

	//popScheduleEvent := PopParseTime(userConstraints, popEvent)
	//recurScheduleEvent := RecurParseTime(recurringEvent)
	initialSchedule, alternatives := GenerateInitial(popEvent, recEvent, userConstraints)

	// for i := 0; i < MAXITERATION; i++ {
	// 	// tabu search
	// }

	finalSchedule = SwapRec(initialSchedule)
	//fmt.Printf("%+v", finalSchedule)

	finalSchedule.TMatrix = GetTransportation(finalSchedule, userConstraints, alternatives)
	finalSchedule = AssignTime(finalSchedule, userConstraints)

	return finalSchedule, alternatives
}

func addTime(x float64, y_inSec float64) float64 {
	y_inMin := math.Floor(y_inSec / 60)
	hour_x := math.Floor(x)
	x_inMin := hour_x*60 + (x-hour_x)*100
	hour_res := math.Floor((x_inMin + y_inMin) / 60)
	min_res := ((x_inMin + y_inMin) - hour_res*60) / 100
	return hour_res + min_res
}

// check violation function
func CheckViolation() bool {
	// since there is max one pop up event per schedule, no need ot check time violation

	return true
}

func AssignTime(finalSchedule Schedule, userConstraints UserConstraints) Schedule {
	if finalSchedule.MaxEvents == 0 {
		return Schedule{}
	}
	
	// fmt.Printf("%+v", finalSchedule)
	duration := userConstraints.EndTime - userConstraints.StartTime
	duration_inMin := (math.Floor(duration)*60 + (duration-math.Floor(duration))*100)
	interval_inSec := math.Round(duration_inMin/float64(finalSchedule.MaxEvents)/5) * 5 * 60

	simToRec := true
	if (finalSchedule.FixedSlot != -1) && (finalSchedule.Events[finalSchedule.FixedSlot].EventInfoPop.EndTime-
		finalSchedule.Events[finalSchedule.FixedSlot].EventInfoPop.StartTime < 3.3) {
		simToRec = false
	}

	if simToRec {

		for i := 0; i < len(finalSchedule.Events); i++ {

			if i == 0 {
				finalSchedule.Events[i].InitialStartTime = userConstraints.StartTime
				finalSchedule.Events[i].InitialEndTime = addTime(finalSchedule.Events[i].InitialStartTime, interval_inSec)
			} else if i == len(finalSchedule.Events)-1 {
				// tTime is in minutes
				tTime := (float64)(finalSchedule.TMatrix.TransportationROWS[i-1].Elements[i].Duration.Value)

				finalSchedule.Events[i].InitialStartTime = addTime(finalSchedule.Events[i-1].InitialEndTime, tTime)
				finalSchedule.Events[i].InitialEndTime = userConstraints.EndTime
			} else {
				tTime := (float64)(finalSchedule.TMatrix.TransportationROWS[i-1].Elements[i].Duration.Value)
				finalSchedule.Events[i].InitialStartTime = addTime(finalSchedule.Events[i-1].InitialEndTime, tTime)
				finalSchedule.Events[i].InitialEndTime = addTime(finalSchedule.Events[i].InitialStartTime, interval_inSec)
			}
		}

	} else {
		finalSchedule.Events[finalSchedule.FixedSlot].InitialStartTime = finalSchedule.Events[finalSchedule.FixedSlot].EventInfoPop.StartTime
		finalSchedule.Events[finalSchedule.FixedSlot].InitialEndTime = finalSchedule.Events[finalSchedule.FixedSlot].EventInfoPop.EndTime

		for i := finalSchedule.FixedSlot - 1; i >= 0; i-- {
			tTime := (float64)((finalSchedule.TMatrix.TransportationROWS[i].Elements[i+1].Duration.Value))
			if i == 0 {
				finalSchedule.Events[i].InitialStartTime = userConstraints.StartTime
				finalSchedule.Events[i].InitialEndTime = addTime(finalSchedule.Events[i+1].InitialStartTime, -1*tTime)
			} else {
				finalSchedule.Events[i].InitialEndTime = addTime(finalSchedule.Events[i+1].InitialStartTime, -1*tTime)
				finalSchedule.Events[i].InitialStartTime = addTime(finalSchedule.Events[i].InitialEndTime, -1*interval_inSec)
			}
		}

		for i := finalSchedule.FixedSlot + 1; i < len(finalSchedule.Events); i++ {
			tTime := (float64)(finalSchedule.TMatrix.TransportationROWS[i-1].Elements[i].Duration.Value)
			if i == len(finalSchedule.Events)-1 {
				finalSchedule.Events[i].InitialEndTime = userConstraints.EndTime
				finalSchedule.Events[i].InitialStartTime = addTime(finalSchedule.Events[i-1].InitialEndTime, tTime)
				
				if finalSchedule.Events[i].InitialEndTime - finalSchedule.Events[i].InitialStartTime > 3 {
					finalSchedule.Events[i].InitialStartTime = finalSchedule.Events[i].InitialEndTime - 3
				}
			} else {
				finalSchedule.Events[i].InitialStartTime = addTime(finalSchedule.Events[i-1].InitialEndTime, tTime)
				finalSchedule.Events[i].InitialEndTime = addTime(finalSchedule.Events[i].InitialStartTime, interval_inSec)
			}
		}

		if finalSchedule.FixedSlot == 0 && finalSchedule.Events[finalSchedule.FixedSlot].InitialStartTime < userConstraints.StartTime{
			finalSchedule.Events[finalSchedule.FixedSlot].InitialStartTime = userConstraints.StartTime
		}else if finalSchedule.FixedSlot == finalSchedule.MaxEvents - 1 && finalSchedule.Events[finalSchedule.FixedSlot].InitialEndTime > userConstraints.EndTime{
			finalSchedule.Events[finalSchedule.FixedSlot].InitialEndTime = userConstraints.EndTime
		}

	}
	
	return finalSchedule
}

// cost function
func CalculateCost(schedule Schedule) float64 {
	// assume schedules that violat time order will not be inputs
	// Travelling salesman on distance
	cost := 0.0

	for i := 0; i < schedule.MaxEvents; i++ {
		score := 0.0
		if schedule.Events[0].Type == 1 {
			score = (float64)(schedule.Events[i].EventInfoRecur.Score)
		} else {
			score = (float64)(schedule.Events[i].EventInfoPop.Score)
		}

		if schedule.MaxEvents == 1 {
			return score
		}

		// minimizing cost
		if i == schedule.MaxEvents-1 {
			cost += 1 / score
		} else {
			cost += math.Abs(schedule.Events[i].Position_x-schedule.Events[i+1].Position_x) +
				math.Abs(schedule.Events[i].Position_y-schedule.Events[i+1].Position_y) + 1/score
		}
	}
	// cost: minimize distance and maximize score
	// calculate distance using (x,y) of each event
	// for i := 0; i < schedule.NumEvents; i++ {
	// 	cost += int(schedule.Events[i].EventInfo.Score)
	// }

	return cost
}

// swap/generate neighbor function
func SwapRec(schedule Schedule) Schedule {
	var newSchedule Schedule
	newSchedule = schedule
	//newSchedule.MaxEvents = schedule.MaxEvents
	//newSchedule.FixedSlot = schedule.FixedSlot
	//indicator := rand.Intn(2)
	bestCost := math.Inf(0)

	for i := range schedule.Events {
		for j := range schedule.Events {
			if j != schedule.FixedSlot && i != schedule.FixedSlot {
				temEvent := schedule.Events[j]
				curS := schedule
				curS.Events[j] = curS.Events[i]
				curS.Events[i] = temEvent
				cost := CalculateCost(curS)

				if cost < bestCost {
					bestCost = cost
					newSchedule = curS
				}
			}
		}
	}

	//fmt.Printf("%+v\n", newSchedule)

	return newSchedule
}

// swap popup to another popup
func SwapPop(schedule Schedule, alternatives []ScheduleEvents) Schedule {
	var newSchedule Schedule
	var thisPop ScheduleEvents
	var index int

	for i := range schedule.Events {
		if schedule.Events[i].Type == 2 {
			thisPop = schedule.Events[i]
			index = i
			break
		}
	}

	for i := range alternatives {
		if alternatives[i].Type == 2 {
			if (alternatives[i].EventInfoPop.StartTime*0.85) <= thisPop.EventInfoPop.StartTime ||
				thisPop.EventInfoPop.StartTime <= (alternatives[i].EventInfoPop.StartTime*1.15) {
				newSchedule.Events[index] = alternatives[i]
			}
		}
	}

	return schedule
}

func AllRecurring(scheduleRecEvents []ScheduleEvents, numEvent int, uc UserConstraints) (Schedule, []ScheduleEvents) {
	if len(scheduleRecEvents) == 0 {
		return Schedule{}, []ScheduleEvents{}
	}

	var schedule Schedule
	schedule.Events = make([]ScheduleEvents, numEvent)
	schedule.MaxEvents = numEvent
	var alternatives []ScheduleEvents

	var prevTag string
	//numNotMustHave := 0
	var indicesMustHave []int
	var indicesOther []int

	for i := range scheduleRecEvents {
		if scheduleRecEvents[i].EventInfoRecur.Tag != prevTag &&
			scheduleRecEvents[i].EventInfoRecur.MustPresent {

			indicesMustHave = append(indicesMustHave, i)
			prevTag = scheduleRecEvents[i].EventInfoRecur.Tag
		} else if scheduleRecEvents[i].EventInfoRecur.Tag != prevTag {

			indicesOther = append(indicesOther, i)
			prevTag = scheduleRecEvents[i].EventInfoRecur.Tag
		}
	}

	if uc.NumMustHave < numEvent {

		currentIndex := 0
		// for i := range indicesOther {
		// 	fmt.Printf("%d\n", indicesOther[i])
		// }

		selected := make(map[int]int)

		for i := 0; i < uc.NumMustHave; i++ {
			if currentIndex < numEvent && i < len(indicesMustHave){
				schedule.Events[currentIndex] = scheduleRecEvents[indicesMustHave[i]]
				selected[indicesMustHave[i]] = 10
				currentIndex++
			} else {
				break
			}
		}

		for i := 0; i < len(indicesOther); i++ {
			if currentIndex < numEvent {
				schedule.Events[currentIndex] = scheduleRecEvents[indicesOther[i]]
				selected[indicesOther[i]] = 10
				currentIndex++
			} else {
				break
			}
		}

		for i := currentIndex; i < numEvent; i++ {

			if len(indicesMustHave) != 0 {
				schedule.Events[i] = scheduleRecEvents[indicesMustHave[0]+i]
				selected[indicesMustHave[0]+i] = 10
			} else {
				schedule.Events[i] = scheduleRecEvents[indicesOther[0]+i]
				selected[indicesOther[0]+i] = 10
			}

		}

		for i := 0; i < len(scheduleRecEvents); i += 2 {
			if selected[i] != 10 {
				alternatives = append(alternatives, scheduleRecEvents[i])
			}
		}

	} else {
		selected := make(map[int]int)
		currentIndex := 0

		for i := range scheduleRecEvents {
			if scheduleRecEvents[i].EventInfoRecur.MustPresent &&
				scheduleRecEvents[i].EventInfoRecur.Tag != prevTag {

				schedule.Events[currentIndex] = scheduleRecEvents[i]
				prevTag = schedule.Events[currentIndex].EventInfoRecur.Tag
				currentIndex ++ 
				selected[i] = 10

			}
		}

		for i := 0; i < len(scheduleRecEvents); i ++ {
			if selected[i] != 10 {
				alternatives = append(alternatives, scheduleRecEvents[i])
			}
		}
	}

	var num_alt = (int)(math.Min(10, (float64)((int)((schedule.MaxEvents + len(alternatives))/2)*2))) - schedule.MaxEvents
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(alternatives), func(i, j int) { alternatives[i], alternatives[j] = alternatives[j], alternatives[i] })
	var alternativesreduced []ScheduleEvents
	// var num_alt = 10 - schedule.MaxEvents
	for i := 0; i < num_alt; i++ {
		alternativesreduced = append(alternativesreduced, alternatives[0])
		alternatives = alternatives[1:]
	}
	schedule.FixedSlot = -1

	for i := 0; i < len(schedule.Events); i++{
		if schedule.Events[i].EventInfoRecur.Name == "" && schedule.Events[i].EventInfoPop.Title == ""{
			if len(alternatives) > 0{
				schedule.Events[i] = alternatives[0]
				alternatives = alternatives[1:]
			}	
		}
	}

	return schedule, alternativesreduced
}

// generate initial solution
// input: list of tagged events
func GenerateInitial(schedulePopEvents []ScheduleEvents, scheduleRecEvents []ScheduleEvents, userConstraints UserConstraints) (Schedule, []ScheduleEvents) {
	if len(scheduleRecEvents) == 0 {
		return Schedule{}, []ScheduleEvents{}
	}

	var maxEvent = 0
	var alternativesRaw []ScheduleEvents
	var alternatives []ScheduleEvents

	duration := userConstraints.EndTime - userConstraints.StartTime

	if duration >= 12 {
		maxEvent = 4
	} else if duration >= 6 {
		maxEvent = 3
	} else if duration >= 3 {
		maxEvent = 2
	} else {
		maxEvent = 1
	}

	if len(schedulePopEvents) == 0 {
		return AllRecurring(scheduleRecEvents, maxEvent, userConstraints)
	} else {

		bestPop := schedulePopEvents[0]
		index := getPopEventIndex(bestPop, userConstraints, maxEvent)

		var schedule Schedule
		schedule.Events = make([]ScheduleEvents, maxEvent)
		schedule.MaxEvents = maxEvent
		schedule.Events[index] = bestPop
		schedule.FixedSlot = index

		var indicesMustHave []int
		var indicesOther []int

		var prevTag string

		for i := range scheduleRecEvents {
			if scheduleRecEvents[i].EventInfoRecur.Tag != prevTag &&
				scheduleRecEvents[i].EventInfoRecur.MustPresent {

				indicesMustHave = append(indicesMustHave, i)
				prevTag = scheduleRecEvents[i].EventInfoRecur.Tag
			} else if scheduleRecEvents[i].EventInfoRecur.Tag != prevTag {

				indicesOther = append(indicesOther, i)
				prevTag = scheduleRecEvents[i].EventInfoRecur.Tag
			}
		}

		currentIndex := 0
		// for i := range indicesOther {
		// 	fmt.Printf("%d\n", indicesOther[i])
		// }

		selected := make(map[int]int)

		for i := 0; i < userConstraints.NumMustHave; i++ {
			if currentIndex != index && currentIndex < maxEvent && i < len(indicesMustHave){
				schedule.Events[currentIndex] = scheduleRecEvents[indicesMustHave[i]]
				selected[indicesMustHave[i]] = 10
				currentIndex++
			} else if currentIndex == index {
				currentIndex++
			} else {
				break
			}
		}

		for i := 0; i < len(indicesOther); i++ {
			if currentIndex != index && currentIndex < maxEvent {
				schedule.Events[currentIndex] = scheduleRecEvents[indicesOther[i]]
				selected[indicesOther[i]] = 10
				currentIndex++
			} else if currentIndex == index {
				currentIndex++
			} else {
				break
			}
		}

		for i := currentIndex; i < maxEvent; i++ {
			if currentIndex == index {
				continue
			} else {
				if len(indicesMustHave) != 0 && indicesMustHave[0]+i < len(indicesMustHave){
					schedule.Events[i] = scheduleRecEvents[indicesMustHave[0]+i]
					selected[indicesMustHave[0]+i] = 10
				} else if indicesOther[0]+i < len(indicesOther){
					schedule.Events[i] = scheduleRecEvents[indicesOther[0]+i]
					selected[indicesOther[0]+i] = 10
				}
			}
		}

		for i := 0; i < len(scheduleRecEvents); i++{
			if (selected[i] != 10){
				alternativesRaw = append(alternativesRaw, scheduleRecEvents[i])
			}
		}

		var num_alt = (int)(math.Min(10, (float64)((int)((schedule.MaxEvents + len(alternativesRaw))/2)*2))) - schedule.MaxEvents
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(alternativesRaw), func(i, j int) { alternativesRaw[i], alternativesRaw[j] = alternativesRaw[j], alternativesRaw[i] })
		// var num_alt = 10 - schedule.MaxEvents

		for i := 0; i < num_alt; i++ {
			alternatives = append(alternatives, alternativesRaw[0])
			alternativesRaw = alternativesRaw[1:]
		}

		allrec := 0

		for i := 0; i < schedule.MaxEvents; i++{
			if len(schedule.Events[i].EventInfoRecur.Name) > 0 {
				allrec++
			}

			if len(schedule.Events[i].EventInfoRecur.Name) <= 0 && len(schedule.Events[i].EventInfoPop.Title) <= 0 {
				if i == schedule.FixedSlot {
					schedule.FixedSlot = -1
				}
				if len(alternativesRaw) > 0{
					schedule.Events[i] = alternativesRaw[0]
					alternativesRaw = alternativesRaw[1:]
				}	
			}
		}

		if allrec == schedule.MaxEvents {
			schedule.FixedSlot = -1
		}

		return schedule, alternatives
	}

	// order events by start time
	// sort.SliceStable(scheduleEvents, func(i, j int) bool {
	// 	return scheduleEvents[i].StartTime < scheduleEvents[j].StartTime
	// })

	// schedule.Events = append(schedule.Events, scheduleEvents[0])
	// schedule.NumEvents = 1

	// for i := 1; i <= len(scheduleEvents); i++ {
	// 	if scheduleEvents[i].StartTime > schedule.Events[schedule.NumEvents-1].EndTime {
	// 		for j := range scheduleEvents[i].EventInfo.Tag {
	// 			if scheduleEvents[i].EventInfo.Tag[j] == "restaurant" {
	// 				isRestaurant = true
	// 				break
	// 			}
	// 		}

	// 		if !isRestaurant {
	// 			schedule.Events = append(schedule.Events, scheduleEvents[0])
	// 			schedule.NumEvents += 1
	// 		}
	// 	}
	// }
	// Note: don't consider recurring events for initial solution since too hard
}

func getPopEventIndex(event ScheduleEvents, userConstraints UserConstraints, maxEvent int) int {
	duration := userConstraints.EndTime - userConstraints.StartTime
	interval := duration / (float64)(maxEvent)

	index := (int)((event.EventInfoPop.StartTime - userConstraints.StartTime) / interval)
	if index < 1 {
		index = 0
	} else if index > (maxEvent - 1) {
		index = maxEvent - 1
	}

	return index
}

func GetTmpEvent() []EVENTTAGGED {
	jsonFile, err := os.Open("./API/eventTag.json")
	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var eventsTaged EVENTARRAY

	err = json.Unmarshal(byteValue, &eventsTaged)
	if err != nil {
		fmt.Printf(err.Error())
	}

	return eventsTaged.EVENTS
}

// func RecurParseTime(eventsTagged []ESTABLISHMENT) []ScheduleEvents {
// 	var scheduleEvents []ScheduleEvents
// 	return scheduleEvents
// }
