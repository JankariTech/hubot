package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/badkaktus/gorocket"
	"github.com/go-resty/resty/v2"
	"github.com/go-yaml/yaml"
)

const fileName = "config.yml"

var ErrLogin = errors.New("the login was not successfull most likely due to invalid credentails")

const logFile = "latest.log"

const eventsTrackerFile = "events_tracker.json"

// Global logger
var logger *log.Logger

// Wrapper struct to keep track of notified events
// Also will be useful while resuming
type EventsForDay struct {
	Day      string                 `json:"day"`
	EventIDs []EventIDWithStartTime `json:"event_ids"`
}

type EventIDWithStartTime struct {
	EventID   string `json:"event_id"`
	StartTime string `json:"start_dt"`
}
type AllEvents struct {
	DayEvents []EventsForDay `json:"all_events"`
}

// Create events tracking json file if it does not exist
func createJSONFile() {
	_, err := os.Stat(eventsTrackerFile)

	// TODO: Improvement
	// Add better logic
	// Not just check if file exists
	// but also if content is correct
	// If not , write down an empty struct
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.OpenFile(eventsTrackerFile, os.O_WRONLY|os.O_CREATE, 0644) // 0644 = user can read and write, other and groups can read
		if err != nil {
			logger.Fatal(err.Error())
		}
		defer f.Close()

		allEvents := AllEvents{DayEvents: []EventsForDay{}}
		writeData, err := json.Marshal(allEvents)
		if err != nil {
			logger.Fatal(err.Error())
		}
		_, err = f.Write(writeData)
		if err != nil {
			logger.Fatal(err.Error())
		}
	}
}

// reads the events that were already notifed for the given day
func readFromJSONFile(day string) EventsForDay {
	var dayEvents *EventsForDay
	data, err := os.ReadFile(eventsTrackerFile)
	if err != nil {
		logger.Fatal(err.Error())
	}
	var allEvents AllEvents
	// Parse the json
	err = json.Unmarshal(data, &allEvents)
	if err != nil {
		logger.Fatal(err.Error())
	}

	for i := 0; i < len(allEvents.DayEvents); i++ {
		if (allEvents.DayEvents[i]).Day == day {
			dayEvents = &allEvents.DayEvents[i]
			break
		}
	}
	if dayEvents == nil {
		return EventsForDay{}
	}

	return *dayEvents
}

// writes to the events tracker json file
// inside the provided day's object
func writeToJSONFile(day, eventID, startTime string) {

	data, err := os.ReadFile(eventsTrackerFile)
	if err != nil {
		logger.Fatal(err.Error())
	}
	var allEvents AllEvents
	// Parse the json
	err = json.Unmarshal(data, &allEvents)
	if err != nil {
		logger.Fatal(err.Error())
	}

	var dayEvents *EventsForDay
	for i := 0; i < len(allEvents.DayEvents); i++ {
		if (allEvents.DayEvents[i]).Day == day {
			dayEvents = &allEvents.DayEvents[i]
			break
		}
	}

	otherDayEvents := EventsForDay{Day: day, EventIDs: []EventIDWithStartTime{}}

	// Means no events have been notified yet
	if dayEvents == nil {
		dayEvents = &otherDayEvents
		// Append the event id to the slice
		dayEvents.EventIDs = append(dayEvents.EventIDs, EventIDWithStartTime{EventID: eventID, StartTime: startTime})

		allEvents.DayEvents = append(allEvents.DayEvents, *dayEvents)
	} else {
		// Append the event id to the slice
		dayEvents.EventIDs = append(dayEvents.EventIDs, EventIDWithStartTime{EventID: eventID, StartTime: startTime})

	}

	f, err := os.OpenFile(eventsTrackerFile, os.O_WRONLY|os.O_CREATE, 0644) // 0644 = user can read and write, other and groups can read
	if err != nil {
		logger.Fatal(err.Error())
	}

	writeData, err := json.Marshal(allEvents)

	if err != nil {
		logger.Fatal(err.Error())
	}
	_, err = f.Write(writeData)
	if err != nil {
		logger.Fatal(err.Error())
	}
	f.Close()
}

// Configuration holds different options
// required to run the bot
type Configuration struct {
	URL          string `yaml:"URL"`
	Username     string `yaml:"USERNAME"`
	Password     string `yaml:"PASSWORD"`
	UseSSL       bool   `yaml:"USE_SSL"`
	Room         string `yaml:"ROOM"`
	MeetingsCode string `yaml:"MEETINGS_TEAMUP"`
	TeamupToken  string `yaml:"TOKEN_TEAMUP"`
	RepeatIn     int    `yaml:"REPEAT_IN"`
}

type TeamupEvent struct {
	ID             string        `json:"id"`
	SeriesID       interface{}   `json:"series_id"`
	RemoteID       string        `json:"remote_id"`
	SubcalendarID  int           `json:"subcalendar_id"`
	SubcalendarIds []int         `json:"subcalendar_ids"`
	AllDay         bool          `json:"all_day"`
	Rrule          string        `json:"rrule"`
	Title          string        `json:"title"`
	Who            string        `json:"who"`
	Location       string        `json:"location"`
	Notes          string        `json:"notes"`
	Version        string        `json:"version"`
	Readonly       bool          `json:"readonly"`
	Tz             interface{}   `json:"tz"`
	Attachments    []interface{} `json:"attachments"`
	StartDt        string        `json:"start_dt"`
	EndDt          string        `json:"end_dt"`
	RistartDt      interface{}   `json:"ristart_dt"`
	RsstartDt      interface{}   `json:"rsstart_dt"`
	CreationDt     time.Time     `json:"creation_dt"`
	UpdateDt       interface{}   `json:"update_dt"`
	DeleteDt       interface{}   `json:"delete_dt"`
}

type TeamupEvents struct {
	Events    []TeamupEvent `json:"events"`
	Timestamp int           `json:"timestamp"`
}

func checkForMeetings(config *Configuration, chatClient *gorocket.Client) {

	locale, _ := time.LoadLocation("Asia/Kathmandu")
	now := time.Now().In(locale)

	fmt.Println("Trying to check for new meetings at Nepal time:", now)

	events, err := fetchMeetingEvents(config)
	if err != nil {
		logger.Println("Error while fetching events. Please check the error\n", err.Error())
	}
	today := time.Now().Format("2006-01-02")

	if len(events.Events) == 0 {
		logger.Println("No meetings found today!!!")
	} else {

		dayEvents := readFromJSONFile(today)

		fmt.Println(dayEvents)
		var futureEvents []TeamupEvent
		if len(dayEvents.EventIDs) > 0 {
			futureEvents = getFutureEvents(events, dayEvents.EventIDs)
		} else {
			futureEvents = getFutureEvents(events, []EventIDWithStartTime{})

		}

		// add logic to send notice for  modified events
		// with start time later and eventsIDs stored in json
		toNotifyEventsIds := []EventIDWithStartTime{}
		toSendMsgs := []string{}
		for _, event := range futureEvents {
			diff := timeDiffWithNow(event.StartDt)
			if diff > 10 && diff < 21 {
				toNotifyEventsIds = append(toNotifyEventsIds, EventIDWithStartTime{event.ID, event.StartDt})
				toSendMsgs = append(toSendMsgs, prepareMeetingMsg(event))
			}
		}

		finalMsg := strings.Join(toSendMsgs, ("\n" + strings.Repeat("-", 100) + "\n"))

		if len(finalMsg) > 0 {
			// See what msg will be sent
			fmt.Println("Trying to send the following messge currently:\n", finalMsg)
			msgSent, err := chatClient.PostMessage(&gorocket.Message{Channel: config.Room, Text: finalMsg})
			if err != nil {
				logger.Println(err)
			}
			// For checking message sent status
			fmt.Println("Message send status: ", msgSent.Success, msgSent.Error)
		}

		for _, val := range toNotifyEventsIds {
			writeToJSONFile(today, val.EventID, val.StartTime)
		}

	}
}

func main() {

	// Setup logger
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644) // 0644 = user can read and write, other and groups can read
	if err != nil {
		logger.Fatal(err.Error())
	}

	//defer to close when you're done with it
	defer f.Close()

	// set this logger to the global variable
	logger = log.New(f, "TEAMUP-ROCKETCHAT-BOT ", log.Flags())

	config, err := readConfig()

	// The app will exit if any errors while reading configuration file
	if err != nil {
		logger.Fatalln("Error while reading configuration file. Please check the error.\n", err.Error())
	}

	logger.Println("Read the following configuration", config)
	// Create json file
	createJSONFile()

	// login to rocketchat
	//chatClient, err := loginRocketChat(config)

	loginResp, err := UpadatedLogin(config, "api/v1")

	if err != nil {
		logger.Fatalln("Error while trying to login. Please check the error.\n", err.Error())
	}

	userIDOpt := gorocket.WithUserID(loginResp.Data.UserID)
	xTokenOpt := gorocket.WithXToken(loginResp.Data.AuthToken)

	chatClient := gorocket.NewWithOptions(config.URL, userIDOpt, xTokenOpt)

	// run this once before cron job
	checkForMeetings(config, chatClient)

	// Ticker will send a signal at each specified time period
	ticker := time.NewTicker(time.Duration(config.RepeatIn) * time.Minute)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	// Run a separate go routine
	go func(wg *sync.WaitGroup) {

		// Range over the ticker
		// Will run each time the ticker ticks, .i.e each 10 mins
		for range ticker.C {
			// check for meeting
			checkForMeetings(config, chatClient)

		}
		wg.Done()

	}(wg)

	// Wait for the goroutine to complete
	wg.Wait()

}

// readConfig reads config.yml files
// and returns config or
// corresponding error if any
func readConfig() (*Configuration, error) {
	config := Configuration{}

	yamlData, err := os.ReadFile(fileName)

	// Check for errors
	if err != nil {
		return nil, err
	}

	// err = yaml.Unmarshal(yamlData, &config)
	err = yaml.UnmarshalStrict(yamlData, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}

func UpadatedLogin(config *Configuration, apiVersion string) (*UpdatedLoginResponse, error) {

	loginPayload := gorocket.LoginPayload{
		User:     config.Username,
		Password: config.Password,
	}
	url := fmt.Sprintf("%s/%s/login", config.URL, apiVersion)
	client := resty.New()
	var errMsg interface{}
	result := UpdatedLoginResponse{}
	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(loginPayload).
		SetResult(&result). // set the response to result with required json parsing
		SetError(&errMsg).
		Post(url)

	if err != nil {
		return &result, err
	} else if res.StatusCode() != 200 {
		return &result, fmt.Errorf("%v", fmt.Errorf("recieved status code %v with body\n%v", res.StatusCode(), string(res.Body())))
	} else if result.Status != "success" {
		return &result, fmt.Errorf("%v", fmt.Errorf("recieved status code %v with body\n%v", res.StatusCode(), string(res.Body())))
	}

	// success auth
	return &result, nil
}

// fetchMeetingEvents fetches events for the day
// for the meetings sub-calendar
func fetchMeetingEvents(config *Configuration) (*TeamupEvents, error) {
	startDate := time.Now().Format("2006-01-02")

	url := fmt.Sprintf("https://api.teamup.com/%s/events?startDate=%v&endDate=%v", config.MeetingsCode, startDate, startDate)
	client := resty.New()
	result := TeamupEvents{} // empty struct
	var errMsg interface{}

	res, err := client.R().
		SetHeader("Accept", "application/json").
		SetHeader("Teamup-Token", config.TeamupToken).
		SetResult(&result). // set the response to result with required json parsing
		SetError(&errMsg).
		Get(url)

	if err != nil {
		return &result, err
	} else if res.StatusCode() != 200 {
		return &result, fmt.Errorf("%v", fmt.Errorf("recieved status code %v with body\n%v", res.StatusCode(), string(res.Body())))
	}

	return &result, nil
}

// timeDiffWithNow returns the time difference (in minutes) between provided time
// and the current time based on unix timestamp
// Example: dateTime = "2023-01-04T07:00:00+05:45"
// Return -1 if any errors
func timeDiffWithNow(dateTime string) int64 {
	epochNow := time.Now().Unix()
	timeGiven, err := time.Parse(time.RFC3339, dateTime)
	if err != nil {
		fmt.Println("Error while converting time stamps")
		return -1
	}
	epochGiven := timeGiven.Unix()
	diff := epochGiven - epochNow
	mins := (diff / 60)
	return mins
}

// returns time difference (in mins) between two time stamps
// Example: dateTime = "2023-01-04T07:00:00+05:45"
func timeDiffBetween(less, more string) int64 {
	timeLess, err := time.Parse(time.RFC3339, less)
	if err != nil {
		fmt.Println("Error while converting time stamps")
		return -1
	}
	timeMore, err := time.Parse(time.RFC3339, more)
	if err != nil {
		fmt.Println("Error while converting time stamps")
		return -1
	}

	epochLess := timeLess.Unix()
	epochMore := timeMore.Unix()

	diff := epochMore - epochLess
	mins := (diff / 60)
	return mins
}

// getFutureEvents returns a slice of future events
// from the current time period
func getFutureEvents(events *TeamupEvents, alreadyNotified []EventIDWithStartTime) []TeamupEvent {
	futureEvents := []TeamupEvent{}
	for _, ev := range events.Events {
		diff := timeDiffWithNow(ev.StartDt)
		// Events yet to come is added to slice
		if diff > 0 {
			shouldAdd := true

			// Loop through noticeIDs to check if it was already notified
			for _, val := range alreadyNotified {
				d := timeDiffBetween(val.StartTime, ev.StartDt)
				if val.EventID == ev.ID && d == 0 {
					shouldAdd = false
					break
				}
			}
			// Add if the flag is true, i.e the id is not listed in already notified
			if shouldAdd {
				futureEvents = append(futureEvents, ev)
			}
		}
	}

	return futureEvents
}

// prepareMeetingMsg prepares message to be
// sent for the upcoming event
// with following format
//
// @mentions *Event title* will start soon
// Start-time: 10:00
// let's meet in:
// Some description
// of the meeting
func prepareMeetingMsg(event TeamupEvent) string {
	title := event.Title
	who := event.Who
	locale, _ := time.LoadLocation("Asia/Kathmandu")
	startTime, _ := time.Parse(time.RFC3339, event.StartDt) // Parse the time in locale time
	location := event.Location
	notes := event.Notes
	description := strings.Split(notes, "\n")
	descSlice := make([]string, 0)

	// Catches the html tags present
	removeTag := regexp.MustCompile(`<\/?[^>]+(>|$)`) // Create a global func for this
	// Catches the whole line containing Who: @mention1 @mention2
	removeMention := regexp.MustCompile(`^(Who|who):\s*@.*$`) // Create a global func for this

	// Trim the white space, remove html tags, remove "Who: @"
	for i := 0; i < len(description); i++ {
		str := strings.TrimSpace(description[i])
		str = removeTag.ReplaceAllString(str, "")
		str = removeMention.ReplaceAllString(str, "")
		// Only add the non-empty lines
		if len(str) > 0 {
			descSlice = append(descSlice, str)
		}
	}
	finalDesc := strings.Join(descSlice, "\n")

	msg := fmt.Sprintf("%s *%s* will start soon\nStart-time: *%s*\nlet's meet in: %s\n%s",
		who,
		title,
		strings.Split(startTime.In(locale).String(), "+")[0],
		location,
		finalDesc,
	)
	return msg
}

// Some struct definitions

type UpdatedMe struct {
	ID                    string            `json:"_id"`
	Services              gorocket.Services `json:"services"`
	Emails                []gorocket.Email  `json:"emails"`
	Status                string            `json:"status"`
	Active                bool              `json:"active"`
	UpdatedAt             time.Time         `json:"_updatedAt"`
	Roles                 []string          `json:"roles"`
	Name                  string            `json:"name"`
	StatusConnection      string            `json:"statusConnection"`
	Username              string            `json:"username"`
	UtcOffset             float64           `json:"utcOffset"`
	StatusText            string            `json:"statusText"`
	Settings              gorocket.Settings `json:"settings"`
	AvatarOrigin          string            `json:"avatarOrigin"`
	RequirePasswordChange bool              `json:"requirePasswordChange"`
	Language              string            `json:"language"`
	Email                 string            `json:"email"`
	AvatarURL             string            `json:"avatarUrl"`
}

type UpdatedDataLogin struct {
	UserID    string    `json:"userId"`
	AuthToken string    `json:"authToken"`
	Me        UpdatedMe `json:"me"`
}

type UpdatedLoginResponse struct {
	Status  string           `json:"status"`
	Data    UpdatedDataLogin `json:"data"`
	Message string           `json:"message,omitempty"`
}
