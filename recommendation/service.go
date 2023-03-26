package recommendation

import (
	"errors"
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"leisurely/database"

	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"net/http"

	g "github.com/serpapi/google-search-results-golang"
	"github.com/rs/zerolog/log"
)

/* Google Map Finding */
var GoogleMapAPIKey = "AIzaSyCYWn0AbmKUqxjapYKLWViWOyawKxRS_D8"
var GoogleURLBase = "https://maps.googleapis.com/maps/api/place/textsearch/json?"
var QueryParameter = "query="
var KeyParameter = "key="
var RadiusParameter = "radius="
var ResultLenParameter = ""
var LocationParameter = "location="
var OpenNowParameter = "opennow="
var TypeParameter = "type="
var GoogleMethod = "GET"

type Coordinate struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lng"`
}

type Geometry struct {
	Location Coordinate `json:"location"`
}

type PeriodDetail struct {
	Day  int    `json:"day"`
	Time string `json:"time"`
}

type Periods struct {
	Open  PeriodDetail `json:"open"`
	Close PeriodDetail `json:"close"`
}

type Hour struct {
	IsOpen  bool     `json:"open_now"`
	Period  []Periods  `json:"periods"`
	Weekday []string `json:"weekday_text"`
}

type Response struct {
	Address   string   `json:"formatted_address"`
	Geometry  Geometry `json:"geometry"`
	Name      string   `json:"name"`
	Hours     Hour     `json:"opening_hours"`
	PlaceID   string   `json:"place_id"`
	Price     int      `json:"price_level"`
	Rating    float32  `json:"rating"`
	PlaceType string   `json:"type"`
}

type Responses struct {
	Response []Response `json:"results"`
}

type ViewPort struct {
	NE Coordinate `json:"northeast"`
	SW Coordinate `json:"southwest"`
}

type CityGeometry struct {
	Location Coordinate `json:"location"`
	ViewPort ViewPort   `json:"viewport"`
}

type CITYINFO struct {
	FormattedAddress string       `json:"formatted_address"`
	Geometry         CityGeometry `json:"geometry"`
	Name             string       `json:"name"`
	PlaceID          string       `json:"place_id"`
}

type ALLCITY struct {
	CityInfo []CITYINFO `json:"results"`
}

func SearchCity(city string, country string) (CITYINFO, error) {
	query := city + "," + country
	URL := GoogleURLBase + QueryParameter + query + "&" + KeyParameter + GoogleMapAPIKey

	client := &http.Client{}
	req, err := http.NewRequest(GoogleMethod, URL, nil)

	if err != nil {
		return CITYINFO{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return CITYINFO{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return CITYINFO{}, err
	}

	//////////////////////////////////////
	// jsonFile, err := os.Open("./API/city.json")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// body, _ := ioutil.ReadAll(jsonFile)
	/////////////////////////////////////////

	var cityinfo ALLCITY
	json.Unmarshal(body, &cityinfo)

	if err != nil {
		println("error when fetching city")
		
	}

	return cityinfo.CityInfo[0], nil

}

type OpeningHours struct {
	OpenNow bool `json:"open_now"`
}

type ESTABLISHMENT struct {
	BusinessStatus   string       `json:"business_status"`
	FormattedAddress string       `json:"formatted_address"`
	Address			[]string
	Geometry         CityGeometry `json:"geometry"`
	Name             string       `json:"name"`
	OpeningHours     OpeningHours `json:"opening_hours"`
	PlaceID          string       `json:"place_id"`
	Price_level      int      `json:"price_level"`
	Rating           float64      `json:"rating"`
	UserRatingTotal  int          `json:"user_ratings_total"`
	Types            []string
	StartTime        float64
	EndTime          float64
	Score            float64
	Tag              string
	MustPresent      bool
}

type ALLESTABLISHMENT struct {
	Establishment []ESTABLISHMENT `json:"results"`
}

func SearchByKeyword(keyword string, location string) (ALLESTABLISHMENT, error) {

	keywordnew := strings.ReplaceAll(keyword, " ", "%20")
	query := keywordnew + "%20in%20" + location

	// Location within 20 Km
	URL := GoogleURLBase + QueryParameter + query + "&" + RadiusParameter + "20000" + "&" + KeyParameter + GoogleMapAPIKey

	client := &http.Client{}
	req, err := http.NewRequest(GoogleMethod, URL, nil)

	if err != nil {
		return ALLESTABLISHMENT{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return ALLESTABLISHMENT{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return ALLESTABLISHMENT{}, err
	}

	var places ALLESTABLISHMENT
	json.Unmarshal(body, &places)

	return places, nil
}

func SearchByCoordinate(keyword string, longitude float32, latitude float32) (string, error) {

	coordinate := fmt.Sprintf("%.2f,%.2f", latitude, longitude)
	URL := GoogleURLBase + QueryParameter + keyword + "&" + LocationParameter + coordinate + "&" + KeyParameter + GoogleMapAPIKey

	client := &http.Client{}
	req, err := http.NewRequest(GoogleMethod, URL, nil)

	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}
	return string(body), nil
}

func ParseResult(response string) (Response, error) {

	var result Response

	err := json.Unmarshal([]byte(response), &result)

	if err != nil {
		fmt.Println(err)
		return result, err
	}

	return result, nil
}

func Transportation(start Geometry, dest Geometry) error {
	return nil
}

/*
	query places detail api for the exact business hours
*/
var placeDetailURL = "https://maps.googleapis.com/maps/api/place/details/json?"
var placeIDParameter = "place_id="
var fieldParameter = "fields=opening_hours"

func GetRecDetail(placeID string) Response{
	URL := placeDetailURL + placeIDParameter + placeID + "&" + fieldParameter + "&" + KeyParameter + GoogleMapAPIKey
	// fmt.Println(URL)
	client := &http.Client{}
	req, err := http.NewRequest(GoogleMethod, URL, nil)
	if err != nil {
		return Response{}
	}

	res, err := client.Do(req)
	if err != nil {
		return Response{}
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return Response{}
	}

	var details Response
	json.Unmarshal(body, &details)

	//fmt.Println(fmt.Sprintf("%+v", details))
	return details
}

/* TMDE - movie showtime */
var MovieMethod = "GET"
var MovieURLBase = ""
var MovieAPIKey = ""

/* serpapi */
var TravelMethod = "GET"
var TravelURLBase = ""
var TravelAPIKey = "8a5d86b78eef21cc590c8fcee4820dc3a3e99f080c8b9eca777b9d29a420bf07"

/* Event - serpapi */
var EventMethod = "GET"
var EventAPIKey = "8a5d86b78eef21cc590c8fcee4820dc3a3e99f080c8b9eca777b9d29a420bf07"

type EVENTDATE struct {
	StartDate string
	WHEN      string
}

type EVENTINFO struct {
	Title       string
	Link        string
	Description string
	Address     []string
	Venue       string
	Date        EVENTDATE
	Image       string
	StartTime   float64
	EndTime     float64
}

type EVENTTAGGED struct {
	Title       string    `json:"Title"`
	Link        string    `json:"Link"`
	Description string    `json:"Description"`
	Address     []string  `json:"Address"`
	Venue       string    `json:"Venue"`
	Date        EVENTDATE `json:"Date"`
	Image       string    `json:"Image"`
	Tag         []string  `json:"EventTags"`
	Score       float32   `json:"Score"`
	StartTime   float64   `json:"StartTime"`
	EndTime     float64   `json:"EndTime"`
	Coordinate  Coordinate
}

type EVENTARRAY struct {
	EVENTS []EVENTTAGGED `json:"Events"`
}

// parse start and end time of events
func PopParseTime(userConstraints UserConstraints, events []EVENTINFO) []EVENTINFO {
	var viableEvents []EVENTINFO
	var startTime float64
	var endTime float64

	for i := range events {
		eventTime := strings.SplitAfter(events[i].Date.WHEN, ", ")
		splitTime := strings.Split(eventTime[len(eventTime)-1], " ")

		startTime = converTime(splitTime[0])
		if startTime == -1 {
			startTime = userConstraints.StartTime
		}

		endTime = userConstraints.EndTime

		// start time, no end time
		if len(splitTime) >= 2 {
			if splitTime[1] == "AM" && startTime == 12 {
				startTime = 0
			}

			if splitTime[1] == "PM" && startTime != 12 {
				startTime += 12
			}
		}

		// possibly have endtime
		if len(splitTime) > 3 {
			if splitTime[1] == "\u2013" {
				endTime = converTime(splitTime[2])
			}
			if splitTime[2] == "\u2013" {
				endTime = converTime(splitTime[3])
			}

			if splitTime[3] == "AM" {
				if startTime == 12 {
					startTime = 0
				}
				if endTime == 12 {
					endTime = 0
				}
			}

			if splitTime[3] == "PM" {
				if startTime != 12 {
					startTime += 12
				}
				if endTime != 12 {
					endTime += 12
				}
			}

			if len(splitTime) >= 5 {
				if splitTime[4] == "AM" && endTime == 12 {
					endTime = 0
				}
				if splitTime[4] == "PM" && endTime != 12 {
					endTime += 12
				}
			}
		}

		startTime = math.Floor(startTime*100) / 100
		endTime = math.Floor(endTime*100) / 100

		if startTime >= userConstraints.StartTime && startTime < userConstraints.EndTime && endTime <= userConstraints.EndTime {
			if endTime < startTime {
				events[i].StartTime = endTime
				events[i].EndTime = startTime
			} else {
				events[i].StartTime = startTime
				events[i].EndTime = endTime
			}

			viableEvents = append(viableEvents, events[i])

		}
	}

	return viableEvents
}

func converTime(time string) float64 {
	min := (float64)(0)
	hour := (float64)(0)

	splitTime := strings.Split(time, ":")

	hourRaw, err := strconv.ParseFloat(splitTime[0], 64)
	if err != nil {
		return -1
	}
	hour = hourRaw
	if len(splitTime) > 1 {
		minRaw, _ := strconv.ParseFloat(splitTime[1], 64)
		min = minRaw * 0.01
	}
	
	return hour + min
}

func CurlEvent(keyword string, Location string, planTime string, index string) ([]EVENTINFO, error) {

	query := strings.ToLower(keyword) + " in " + Location + " " + planTime

	parameter := map[string]string{
		"q":       query, /* example "events+in+Austin" */
		"start":   index,
		"engine":  "google_events",
		"gl":      "us",
		"hl":      "en",
		"api_key": EventAPIKey,
	}

	search := g.NewGoogleSearch(parameter, EventAPIKey)
	results, err := search.GetJSON()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	eventsResults := results["events_results"].([]interface{})
	var events []EVENTINFO

	for _, element := range eventsResults {

		var event EVENTINFO
		c := element.(map[string]interface{})
		event.Title = c["title"].(string)

		if c["link"] != nil {
			event.Link = c["link"].(string)
		} else {
			event.Link = ""
		}

		if c["description"] != nil {
			event.Description = c["description"].(string)
		} else {
			event.Description = ""
		}

		if c["venue"] != nil {
			v := c["venue"].(map[string]interface{})
			event.Venue = v["name"].(string)
		} else {
			event.Venue = ""
		}

		if c["thumbnail"] != nil {
			event.Image = c["thumbnail"].(string)
		} else {
			event.Image = ""
		}

		a := c["address"].([]interface{})
		var address []string
		for _, addr := range a {
			address = append(address, addr.(string))
		}
		event.Address = address

		d := c["date"].(map[string]interface{})
		event.Date.StartDate = d["start_date"].(string)
		event.Date.WHEN = d["when"].(string)

		event.StartTime = -1
		event.EndTime = -1

		events = append(events, event)

	}

	return events, nil
}

func CurlEventAll(tag string, userConstraints UserConstraints) []EVENTINFO {
	var eventsraw []EVENTINFO
	var events []EVENTINFO

	for i := 0; i < 2; i++ {
		firstNevents, err := CurlEvent(tag, userConstraints.Location, userConstraints.Date, strconv.Itoa(i*10))

		if err == nil {
			for i := range firstNevents {
				eventsraw = append(eventsraw, firstNevents[i])
			}
		} else {
			break
		}
	}

	if len(eventsraw) == 0 {
		return []EVENTINFO{}
	}

	for i := 0; i < len(eventsraw); i++ {
		// change this for #pop
		if (i > 8) {
			break
		}
		events = append(events, eventsraw[i])

	}

	formattedEvents := PopParseTime(userConstraints, events)

	return formattedEvents
}

func SearchByAddress(address string) (Coordinate, error) {

	addressNew := strings.ReplaceAll(address, " ", "%20")
	URL := GoogleURLBase + QueryParameter + addressNew + "&" + KeyParameter + GoogleMapAPIKey

	client := &http.Client{}
	req, err := http.NewRequest(GoogleMethod, URL, nil)

	if err != nil {
		return Coordinate{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return Coordinate{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return Coordinate{}, err
	}

	var places ALLESTABLISHMENT
	json.Unmarshal(body, &places)

	if len(places.Establishment) == 0 {
		return Coordinate{}, errors.New("cannot find the address")
	} else {
		return places.Establishment[0].Geometry.Location, err
	}
}

func PopEventGenerator(userConstraints UserConstraints) []EVENTTAGGED {
	var events []EVENTINFO

	oneEvent := CurlEventAll("Events ", userConstraints)

	if len(oneEvent) == 0{
		return []EVENTTAGGED{}
	}

	for j := range oneEvent {
		if j >= 15 {
			break
		}
		
		events = append(events, oneEvent[j])
	}
	
	servAddr := database.ENVSettings.BackendRecom_Name + ":4000"
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, _ := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		println("Dial failed:", err.Error())
		return []EVENTTAGGED{}
	}

	eventsdata, err := json.Marshal(events)
	if err != nil {
		println("Marshal events fail:", err.Error())
	}

	prefjson, err := json.Marshal(userConstraints.Tags)

	datasize := uint32(len(prefjson))
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, datasize)
	_, err = conn.Write(bs)
	_, err = conn.Write(prefjson)

	datasize = uint32(len(eventsdata))
	bs = make([]byte, 4)
	binary.BigEndian.PutUint32(bs, datasize)
	_, err = conn.Write(bs)
	_, err = conn.Write(eventsdata)

	if err != nil {
		println("Write to server failed:", err.Error())
		return []EVENTTAGGED{}
	}

	// message_size_buf := make([]byte, 8)
	// reply_len, _ := conn.Read(message_size_buf)
	//fmt.Println(message_size_buf)
	// message_size := binary.BigEndian.Uint64(message_size_buf)
	// fmt.Println(message_size)

	reply := make([]byte, 100000000)
	reply_len, err := conn.Read(reply)

	if err != nil {
		fmt.Print("receive error")
	}

	var eventsTagged EVENTARRAY
	var validEvent []EVENTTAGGED

	///////////////////////////////////////
	// jsonFile, err := os.Open("./API/event.json")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// byteValue, _ := ioutil.ReadAll(jsonFile)
	// json.Unmarshal(byteValue, &eventsTagged)

	// if err != nil {
	// 	println("error when fetching pop-up event")

	// }
	////////////////////////////////////////
		
	err = json.Unmarshal(reply[0:reply_len], &eventsTagged)
	if err != nil {
		fmt.Printf(err.Error())
	}

	conn.Close()

	for i := range eventsTagged.EVENTS {
		eventsTagged.EVENTS[i].Coordinate, err = SearchByAddress(eventsTagged.EVENTS[i].Address[0])
		if err == nil {
			validEvent = append(validEvent, eventsTagged.EVENTS[i])
		}
	}

	return validEvent
}

/* Weather Finding */
var WeatherMethod = "GET"
var WeatherURLBase = "https://api.openweathermap.org/data/3.0/onecall?"
var WeatherAPIKey = "088891ece2ac7c1b953a4af019851d2f"
var Weatherlatitude = "lat="
var Weatherlongitude = "lon="
var WeatherAppID = "appid="
var WeatherExclude = "exclude="
var TimeLayout = "2006-01-02 15:04:05"

type FELLSLIKE struct {
	Day   float64 `json:"day"`
	Night float64 `json:"night"`
	Eve   float64 `json:"eve"`
	Morn  float64 `json:"morn"`
}

type WEATHER struct {
	Main        string `json:"main"`
	Description string `json:"description"`
}

type DAILY struct {
	Feels_like FELLSLIKE `json:"feels_like"`
	Weather    []WEATHER `json:"weather"`
}

type ALLWEATHER struct {
	Daily []DAILY `json:"daily"`
}

func CurlWeather(latitude float32, longitude float32, date string) (DAILY, error) {
	var curWeatherlatitude = Weatherlatitude + fmt.Sprintf("%.2f", latitude)
	var curWeatherlongitude = Weatherlongitude + fmt.Sprintf("%.2f", longitude)
	var curWeatherAppID = WeatherAppID + WeatherAPIKey

	currentTime := time.Now()
	futureDate, err := time.Parse(TimeLayout, date)

	if err != nil {
		fmt.Println("Invalid date provided")
		return DAILY{}, err
	}

	diff := futureDate.Sub(currentTime)
	days := int(diff.Hours() / 24)
	dayIndex := days + 1

	var exclude string
	exclude = WeatherExclude + "current,minutely,hourly"

	WeatherURL := WeatherURLBase + curWeatherlatitude + "&" + curWeatherlongitude + "&" + exclude + "&" + curWeatherAppID

	client := &http.Client{}
	req, err := http.NewRequest(WeatherMethod, WeatherURL, nil)

	if err != nil {
		log.Error().Msg("Fail to setup HTTP Request for WeatherAPI")
		return DAILY{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		log.Error().Msg("Cannot curl weather")
		return DAILY{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Error().Msg("Curl weather response error")
		return DAILY{}, err
	}

	var weatherinfo ALLWEATHER
	json.Unmarshal(body, &weatherinfo)

	if err != nil {
		log.Error().Msg("error when fetching weather")
		return DAILY{}, nil
	}

	if (dayIndex > 7) {
		return DAILY{}, nil
	}

	if (dayIndex < 0) {
		dayIndex = 0
	}
	
	return weatherinfo.Daily[dayIndex], nil
}

type Plan struct {
	PlanResult PlanResult `json:"data"`
}

type PlanResult struct {
	Schedule    Schedule
	Alternative []ScheduleEvents
}

var TransportationURLBASE = "https://maps.googleapis.com/maps/api/distancematrix/json?"
var ORIGIN = "origins="
var DESTINATION = "destinations="
var travelMode = "mode="
var DrivingMode = "driving"
var PublicMode = "transit"

func GetTmpPlaces() (PlanResult, error) {
	var places Plan

	jsonFile, err := os.Open("./API/planResponse.json")
	if err != nil {
		fmt.Println(err)
		return PlanResult{}, err
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &places)

	if err != nil {
		fmt.Println(err.Error())
	}

	return places.PlanResult, nil
}

type Distance struct {
	Text  string `json:"text"`
	Value int    `json:"value"`
}

type Duration struct {
	Text  string `json:"text"`
	Value int    `json:"value"`
}

type TransportationElement struct {
	Distance Distance `json:"distance"`
	Duration Duration `json:"duration"`
}

type TransportationROW struct {
	Elements []TransportationElement `json:"elements"`
}

type TransportationMatrix struct {
	Origins            []string            `json:"origin_addresses"`
	Destinations       []string            `json:"destination_addresses"`
	TransportationROWS []TransportationROW `json:"rows"`
}

func GetTransportation(plan Schedule, uc UserConstraints, alternatives []ScheduleEvents) TransportationMatrix {

	var origin string
	var dest string
	var result TransportationMatrix

	// only getting diagonal, need to send 4 at a time for 49 separate calls
	allevents := append(plan.Events, alternatives...)
	var coordinatecol string
	// var coordinaterow string

	// if len(allevents)%5 != 0 {
	// 	for i := (len(allevents)/5)*5; i < len(allevents); i++ {
	// 		coordinate = fmt.Sprintf("%.2f,%.2f", allevents[i].Position_x, allevents[i].Position_y)
	// 		if (i != len(allevents)){
	// 			origindest = origindest + coordinate + "|"
	// 		}else {
	// 			origindest = origindest + coordinate
	// 		}
	// 	}
		
	// 	trans := GetTransportationByString(origindest, uc)
	// 	for row := range trans.TransportationROWS {
	// 		for col := range trans.TransportationROWS[row].Elements {
	// 			result.TransportationROWS[len(allevents)/5)*5 + i].Elements[j - 4 + col] = trans.TransportationROWS[row].Elements[col];
	// 		}
	// 	}
		
	// 	return trans
	// }

	// initialize result
	initialDist := Distance {
		Text: "0.0 km",
		Value: 0,
	}

	initialDur := Duration {
		Text: "0 mins",
		Value: 0,
	}
	
	initialTrans := TransportationElement {
		Distance: initialDist,
		Duration: initialDur,
	}
	
	result.TransportationROWS = make([]TransportationROW, len(allevents)) 
	for i := range result.TransportationROWS {
		result.TransportationROWS[i].Elements = make([]TransportationElement, len(allevents))
		for j := range result.TransportationROWS[i].Elements {
			result.TransportationROWS[i].Elements[j] = initialTrans
		}
	}

	half := (int)(len(allevents)/2)
	loggingmsg := fmt.Sprintf("Length of allevents: %d", len(allevents))
	log.Info().Msg(loggingmsg)

	for i := 0; i < int(math.Ceil(float64(len(allevents))/float64(half))); i++{
		
		for j := 0; j < len(allevents); j++ {
			coordinatecol = fmt.Sprintf("%.5f,%.5f", allevents[j].Position_x, allevents[j].Position_y)
			
			if (j != 0){
				origin = origin + "|" + coordinatecol
			}else {
				origin = origin + coordinatecol
			}
			
			if ((j % half) == (half-1)){
				for k := i*half; k < (i+1)*half; k++{
					if (k == (i+1)*half-1){
						dest += fmt.Sprintf("%.5f,%.5f", allevents[k].Position_x, allevents[k].Position_y)
					}else{
						dest += fmt.Sprintf("%.5f,%.5f", allevents[k].Position_x, allevents[k].Position_y) + "|"
					}
				}

				// behavior of dest and origin are opposite, so reversing the order here instead of changing all the names
				trans := GetTransportationByString(dest, origin, uc)
				for row := range trans.TransportationROWS {
					for col := range trans.TransportationROWS[row].Elements {
						if trans.TransportationROWS[row].Elements[col].Distance.Text != ""{
							result.TransportationROWS[i*half + row].Elements[j - (half - 1) + col] = trans.TransportationROWS[row].Elements[col];
						}
					}
				}
				
				coordinatecol = ""
				origin = ""
				dest = ""
			// case: number of events is not divisible by 5
			// no need to worry about divisible by 5 since it will be taken care of by the previous case
			}
		}
	}

	for i:=0; i < len(allevents); i++{
		result.Origins = append(result.Origins, fmt.Sprintf("%f,%f", allevents[i].Position_x, allevents[i].Position_y))
		result.Destinations = append(result.Destinations, fmt.Sprintf("%f,%f", allevents[i].Position_x, allevents[i].Position_y))
	}

	return result

}

func GetTransportationByString(origin string, dest string, uc UserConstraints) TransportationMatrix{
	url := TransportationURLBASE + ORIGIN + origin + "&" + DESTINATION + dest + "&" + travelMode + DrivingMode +
		"&" + KeyParameter + GoogleMapAPIKey

	if uc.Transport == "Transit" {
		url = TransportationURLBASE + ORIGIN + origin + "&" + DESTINATION + dest + "&" + travelMode + PublicMode +
			"&" + KeyParameter + GoogleMapAPIKey
	}

	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return TransportationMatrix{}
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return TransportationMatrix{}
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return TransportationMatrix{}
	}
	
	var trans TransportationMatrix
	json.Unmarshal(body, &trans)
	return trans
}