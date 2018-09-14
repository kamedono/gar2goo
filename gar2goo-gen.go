package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"google.golang.org/api/calendar/v3"
)

const (
	basePath = "https://cybozu.com/g/api/v1/schedule/events"
	header   = "X-Cybozu-Authorization"
)

type EventsCall struct {
	urlParams_ url.Values
	host       string
	headerVal  string
}

// TODO: fuck!
func (c *EventsCall) Init() *EventsCall {
	c.urlParams_ = make(url.Values)
	return c
}

func (c *EventsCall) Limit(limit int) *EventsCall {
	c.urlParams_.Set("limit", fmt.Sprint(limit))
	return c
}

func (c *EventsCall) Offset(offset int) *EventsCall {
	c.urlParams_.Set("offset", fmt.Sprint(offset))
	return c
}

func (c *EventsCall) Field(field string) *EventsCall {
	c.urlParams_.Set("offset", field)
	return c
}

func (c *EventsCall) OrderBy(orderBy string) *EventsCall {
	c.urlParams_.Set("orderBy", orderBy)
	return c
}

func (c *EventsCall) RangeStart(rangeStart string) *EventsCall {
	c.urlParams_.Set("rangeStart", rangeStart)
	return c
}

func (c *EventsCall) RangeEnd(rangeEnd string) *EventsCall {
	c.urlParams_.Set("rangeEnd", rangeEnd)
	return c
}

func (c *EventsCall) Target(target string) *EventsCall {
	c.urlParams_.Set("target", target)
	return c
}

func (c *EventsCall) TargetType(targetType string) *EventsCall {
	c.urlParams_.Set("targetType", targetType)
	return c
}

func (c *EventsCall) Keyword(keyword string) *EventsCall {
	c.urlParams_.Set("keyword", keyword)
	return c
}

func (c *EventsCall) ExcludeFromSearch(excludeFromSearch string) *EventsCall {
	c.urlParams_.Set("excludeFromSearch", excludeFromSearch)
	return c
}

func (c *EventsCall) HeaderValue(value string) *EventsCall {
	c.headerVal = value
	return c
}

func (c *EventsCall) Host(host string) *EventsCall {
	c.host = host
	return c
}

func (c *EventsCall) Do() (*Events, error) {
	u, err := url.Parse(basePath)
	if err != nil {
		return nil, err
	}
	u.Host = c.host + "." + u.Hostname()
	u.RawQuery = c.urlParams_.Encode()

	req, _ := http.NewRequest("GET", u.ResolveReference(u).String(), nil)
	fmt.Println(u.ResolveReference(u).String())

	req.Header.Set(header, c.headerVal)

	client := new(http.Client)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	target := &Events{}
	e := DecodeResponse(target, resp)
	return target, e
}

// convert event from Garoon to Google
func (e *Events) toGoogleEvents() []*calendar.Event {
	eventList := []*calendar.Event{}

	for _, v := range e.Events {
		eventList = append(eventList, v.toGoogleEvent())
	}
	return eventList
}

// convert event from Garoon to Google
func (e *Event) toGoogleEvent() *calendar.Event {

	return &calendar.Event{
		Created:     e.CreatedAt.String(),
		Summary:     e.Subject,
		Description: e.Notes,
		Location:    e.toGoogleEventLocation(),
		Creator:     e.toGoogleEventCreator(),
		Start:       e.toGoogleEventStart(),
		End:         e.toGoogleEventEnd(),
	}
}

// convert creator from Garoon to Google
func (e *Event) toGoogleEventCreator() *calendar.EventCreator {

	return &calendar.EventCreator{
		DisplayName: e.Creator.Name,
	}

}

// convert location from Garoon to Google
func (e *Event) toGoogleEventLocation() string {
	s := []string{}

	for _, v := range e.Facilities {
		s = append(s, v.Name)
	}
	return strings.Join(s, "&")
}

// convert start DateTime from Garoon to Google
func (e *Event) toGoogleEventStart() *calendar.EventDateTime {
	dt := &calendar.EventDateTime{}

	if e.IsAllDay == "" {
		dt.Date = e.Start.DateTime.String()
	} else {
		dt.DateTime = e.Start.DateTime.String()
	}

	dt.TimeZone = e.Start.TimeZone

	return dt
}

// convert end DateTime from Garoon to Google
func (e *Event) toGoogleEventEnd() *calendar.EventDateTime {
	dt := &calendar.EventDateTime{}

	if e.IsAllDay == "" {
		dt.Date = e.End.DateTime.String()
	} else {
		dt.DateTime = e.End.DateTime.String()
	}

	dt.TimeZone = e.End.TimeZone

	return dt
}

// Decode json
func DecodeResponse(target interface{}, res *http.Response) error {
	if res.StatusCode == http.StatusNoContent {
		return nil
	}
	return json.NewDecoder(res.Body).Decode(target)
}

type Events struct {
	Events  []Event `json:"events"`
	HasNext bool    `json:"hasNext"`
}

type User struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type DateTime struct {
	DateTime time.Time `json:"dateTime"`
	TimeZone string    `json:"timeZone"`
}

type CompanyInfo struct {
	Name      string `json:"name"`
	ZipCode   string `json:"zipCode"`
	Address   string `json:"address"`
	Route     string `json:"route"`
	RouteTime string `json:"routeTime"`
	RouteFare string `json:"routeFare"`
	Phone     string `json:"phone"`
}

type Attachment struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ContentType string `json:"contentType"`
	Size        string `json:"size"`
}

type Attendee struct {
	ID                 int                 `json:"id"`
	Code               string              `json:"code"`
	Name               string              `json:"name"`
	Type               string              `json:"type"`
	AttendanceResponse *AttendanceResponse `json:"attendanceResponse"`
}

type AttendanceResponse struct {
	Status  string `json:"status"`
	Comment string `json:"comment"`
}

type Watcher struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Facilitie struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type FacilityReservationInfo struct {
	AdditionalProp1 *AdditionalProp `json:"additionalProp1"`
	AdditionalProp2 *AdditionalProp `json:"additionalProp2"`
	AdditionalProp3 *AdditionalProp `json:"additionalProp3"`
}

type AdditionalProp struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type FacilityUsageRequest struct {
	Status           string      `json:"status"`
	Facility         *Facility   `json:"facility"`
	ApprovedBy       *ApprovedBy `json:"approvedBy"`
	ApprovedDateTime time.Time   `json:"approvedDateTime"`
}

type Facility struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type ApprovedBy struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type RepeatInfo struct {
	Type               string               `json:"type"`
	Period             *Period              `json:"period"`
	Time               *Time                `json:"time"`
	TimeZone           string               `json:"timeZone"`
	IsAllDay           bool                 `json:"isAllDay"`
	IsStartOnly        bool                 `json:"isStartOnly"`
	DayOfWeek          string               `json:"dayOfWeek"`
	DayOfMonth         string               `json:"dayOfMonth"`
	ExclusiveDateTimes []*ExclusiveDateTime `json:"exclusiveDateTimes"`
}

type Period struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type Time struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type ExclusiveDateTime struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type TemporaryEventCandidate struct {
	End      *DateTime `json:"end"`
	Start    *DateTime `json:"start"`
	Facility *Facility `json:"facility"`
}

type AdditionalItems struct {
	Item *Item `json:"item"`
}

type Item struct {
	Value string `json:"value"`
}

type Event struct {
	ID                       string                     `json:"id"`
	Creator                  *User                      `json:"creator"`
	CreatedAt                time.Time                  `json:"createdAt"`
	Updater                  *User                      `json:"updater"`
	UpdatedAt                time.Time                  `json:"updatedAt"`
	EventType                string                     `json:"eventType"`
	EventMenu                string                     `json:"eventMenu"`
	Subject                  string                     `json:"subject"`
	Notes                    string                     `json:"notes"`
	VisibilityType           string                     `json:"visibilityType"`
	UseAttendanceCheck       bool                       `json:"useAttendanceCheck"`
	CompanyInfo              *CompanyInfo               `json:"companyInfo"`
	Attachments              []*Attachment              `json:"attachments"`
	Start                    *DateTime                  `json:"start"`
	End                      *DateTime                  `json:"end"`
	IsAllDay                 string                     `json:"isAllDay"`
	IsStartOnly              string                     `json:"isStartOnly"`
	OriginalStartTimeZone    string                     `json:"originalStartTimeZone"`
	OriginalEndTimeZone      string                     `json:"originalEndTimeZone"`
	Attendees                []Attendee                 `json:"attendees"`
	Watchers                 []*Watcher                 `json:"watchers"`
	Facilities               []*Facilitie               `json:"facilities"`
	FacilityUsingPurpose     string                     `json:"facilityUsingPurpose"`
	FacilityReservationInfo  *FacilityReservationInfo   `json:"facilityReservationInfo"`
	FacilityUsageRequests    []*FacilityUsageRequest    `json:"facilityUsageRequests"`
	RepeatInfo               *RepeatInfo                `json:"repeatInfo"`
	TemporaryEventCandidates []*TemporaryEventCandidate `json:"temporaryEventCandidates"`
	AdditionalItems          *AdditionalItems           `json:"additionalItems"`
	RepeatID                 string                     `json:"repeatId"`
}
