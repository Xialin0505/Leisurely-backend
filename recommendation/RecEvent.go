package recommendation

import (
	"math"
	"strconv"
	"strings"
	"time"
	"fmt"
	Tag "leisurely/modules/tag"
	"github.com/rs/zerolog/log"
)

var DateMapMonth = map[string]string{
	"Jan": "01",
	"Feb": "02",
	"Mar": "03",
	"Apr": "04",
	"May": "05",
	"Jun": "06",
	"Jul": "07",
	"Aug": "08",
	"Sep": "09",
	"Oct": "10",
	"Nov": "11",
	"Dec": "12",
}

var OutDoor = map[string]bool{
	// is out door
	"theme park": true,
	"park":       true,
	"farm":       true,
	"beach":      true,
}

func ParseDateTime(date string, starttime float64) string {
	var result string
	t := time.Now()
	year := t.Year()
	yearstring := strconv.Itoa(year)

	hour := math.Floor(starttime)
	min := (int)(starttime*100 - hour*100)

	hourStr := strconv.Itoa((int)(hour))
	if hour < 10 {
		hourStr = "0" + hourStr
	}

	minStr := strconv.Itoa(min)
	if min < 10 {
		minStr = "0" + minStr
	}

	datesplit := strings.Split(date, " ")

	var dateString string
	dateint, err := strconv.Atoi(datesplit[1])

	if err != nil {
		fmt.Println("Invalid date provided")
		return ""
	}

	if (dateint < 10){
		dateString = "0" + datesplit[1]
	}else {
		dateString = datesplit[1]
	}
	
	result = yearstring + "-" + DateMapMonth[datesplit[0]] + "-" + dateString +
		" " + hourStr + ":" + minStr + ":00"

	//fmt.Println(result)
	return result
}

func RecEventGenerator(uc UserConstraints) []ESTABLISHMENT {
	cityresult, _ := SearchCity(uc.Location, uc.Country)
	formattedDate := ParseDateTime(uc.Date, uc.StartTime)
	weather, err := CurlWeather(cityresult.Geometry.Location.Latitude, cityresult.Geometry.Location.Longitude, formattedDate)

	if (err != nil){
		fmt.Println("Cannot obtain weather information")
		var nullreturn []ESTABLISHMENT
		return nullreturn
	} 

	var newTags []InitialTag
	if strings.Contains(strings.ToLower(weather.Weather[0].Main), "rain") ||
		strings.Contains(strings.ToLower(weather.Weather[0].Main), "snow") {
		// use indoor only

		for i := range uc.Tags {
			if _, ok := OutDoor[uc.Tags[i].Description]; !ok {
				newTags = append(newTags, uc.Tags[i])
			}
		}

		uc.Tags = newTags
	}

	var allplaces []ESTABLISHMENT
	i := 0

	for; i < len(uc.Tags);i++ {
		if uc.Tags[i].EventType == 1 {
			allplaces = append(allplaces, getRecPlaces(uc, i)...)
		}
	}

	if len(allplaces) < 6 {
		rec := Tag.GetTagByType(1)
		for j := 0; j < len(rec); j++ {
			initTag := InitialTag {
				Description: rec[j].Description,
				EventType: 1,
				Count: 1,
				SearchingTag: true,
			}

			if _, ok := OutDoor[initTag.Description]; !ok {
				uc.Tags = append(uc.Tags, initTag)
				allplaces = append(allplaces, getRecPlaces(uc, i)...)
				i++
			}
		}
	}
	
	placesScored := ScoreRecEvent(uc, allplaces)
	loggingmsg := fmt.Sprintf("number of recurring event is: %d", len(placesScored))
	log.Info().Msg(loggingmsg)
	return placesScored
}

func ScoreRecEvent(uc UserConstraints, places []ESTABLISHMENT) []ESTABLISHMENT {
	var result []ESTABLISHMENT

	for i := range places {
		if places[i].MustPresent {
			places[i].Score += 10
		}

		for j := range uc.Tags {
			if uc.Tags[j].Description == places[i].Tag {
				places[i].Score = float64(uc.Tags[j].Count)
				break
			}
		}

		result = append(result, places[i])
	}
	
	return result
}

func getRecPlaces(uc UserConstraints, i int) []ESTABLISHMENT{
	var allplaces []ESTABLISHMENT
	if i >= len(uc.Tags){
		return allplaces
	}
	result, _ := SearchByKeyword(uc.Tags[i].Description, uc.Location+","+uc.Country)

	eventsAdded := 0
	for j := range result.Establishment {
		if eventsAdded >= 5 {
			break
		}
		result.Establishment[j].Tag = uc.Tags[i].Description
		result.Establishment[j].MustPresent = uc.Tags[i].SearchingTag
		placeDetails := GetRecDetail(result.Establishment[j].PlaceID)
		// If have time: get the actual day of the week instead of just monday
		result.Establishment[j].StartTime = uc.StartTime
		result.Establishment[j].EndTime = uc.EndTime

		if len(placeDetails.Hours.Period) > 0 {
			result.Establishment[j].StartTime, _ = strconv.ParseFloat(placeDetails.Hours.Period[0].Open.Time, 64)
			result.Establishment[j].StartTime /= 100
			result.Establishment[j].EndTime, _ = strconv.ParseFloat(placeDetails.Hours.Period[0].Close.Time, 64)
			result.Establishment[j].EndTime /= 100
		}
		
		if result.Establishment[j].Price_level <= uc.BudgetLevel {
			if (result.Establishment[j].PlaceID != ""){
				allplaces = append(allplaces, result.Establishment[j])
				eventsAdded += 1
			}	
		}
	}
	return allplaces
}