package requestDto

type MGraphCreateEventDto struct {
	Subject               string                          `json:"subject"`
	Content               string                          `json:"content"`
	StartTime             string                          `json:"start_time"`
	EndTime               string                          `json:"end_time"`
	TimeZone              string                          `json:"time_zone"`
	Attendees             []MGraphCreateEventAttendeeDto  `json:"attendees"`
	Locations             *[]MGraphCreateEventLocationDto `json:"locations"`
	IsRecurring           bool                            `json:"is_recurring"`
	PatternType           *string                         `json:"pattern_type"`
	PatternInterval       *int32                          `json:"pattern_interval"`
	PatternDaysOfWeek     *[]string                       `json:"pattern_days_of_week"`
	RecurrenceType        *string                         `json:"recurrence_type"`
	RecurrenceStart       *string                         `json:"recurrence_start"`
	RecurrenceEnd         *string                         `json:"recurrence_end"`
	IsOnlineMeeting       bool                            `json:"is_online_meeting"`
	OnlineMeetingProvider *string                         `json:"online_meeting_provider"`
}

// func TestType() bool {
// 	recurrenceRange := graphmodels.NewRecurrenceRange()
// 	recurrenceType := graphmodels.ENDDATE_RECURRENCERANGETYPE
// 	recurrenceRange.SetTypeEscaped(&recurrenceType)
// 	startDate, _ := time.Parse("2006-01-02", "2017-09-04")
// 	serializedStartDate := serialization.NewDateOnly(startDate)
// 	recurrenceRange.SetStartDate(serializedStartDate)
// 	endDate := "2017-12-31"
// 	recurrenceRange.SetEndDate(&endDate)
// 	recurrence.SetRange(recurrenceRange)
// 	requestBody.SetRecurrence(recurrence)
// 	location := graphmodels.NewLocation()
// 	displayName := "Harry's Bar"
// 	location.SetDisplayName(&displayName)
// 	requestBody.SetLocation(location)
//
// 	return true
// }

// use []graphmodels.attendeeable
type MGraphCreateEventAttendeeDto struct {
	EmailAddress string `json:"email_address"`
	Name         string `json:"name"`
	// attendee is of type *graphmodels.AttendeeType
	AttendeeType string `json:"attendee_type"`
}

// need to make sure we count the locations to determine if we need []graphmodels.Locationable
// if > 1 location, setLocations
// else setLocation
type MGraphCreateEventLocationDto struct {
	DisplayName string `json:"display_name"`
	// address is of type *graphmodels.PhysicalAddress
	Address         *MGraphCreateEventLocationAddressDto `json:"address"`
	DefaultLocation bool                                 `json:"default_location"`
}

type MGraphCreateEventLocationAddressDto struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
}
