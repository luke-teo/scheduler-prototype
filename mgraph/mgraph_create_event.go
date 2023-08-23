package mgraph

import (
	"context"
	"errors"
	"log"

	"github.com/microsoft/kiota-abstractions-go/serialization"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	requestDto "github.com/scheduler-prototype/dto/request"
)

// sample json request body
// {
//     "subject": "test recurring",
//     "content": "test content",
//     "start_time": "2023-08-24T00:00:00Z",
//     "end_time": "2023-08-24T00:03:00Z",
//     "time_zone": "UTC",
//     "attendees" : [
//         {
//             "email_address": "luke.teo@isao.cloud",
//             "name": "Luke Teo",
//             "attendee_type": "required"
//         }
//     ],
//     "locations": [
//         {
//             "display_name": "Floris Room",
//             "address": {
//                 "street": "street",
//                 "city": "city",
//                 "state": "state",
//                 "country": "country",
//                 "postal_code": "postal code"
//             },
//             "default_location": true
//         }
//     ],
//     "is_recurring": true,
//     "pattern_type": "weekly",
//     "pattern_interval": 1,
//     "pattern_days_of_week": ["wednesday", "tuesday"],
//     "recurrence_type": "endDate",
//     "recurrence_start": "2023-08-23",
//     "recurrence_end": "2023-09-01",
//     "is_online_meeting": true,
//     "online_meeting_provider": "teamsForBusiness"
// }

func (m *MGraph) PostCreateEvent(requestDto *requestDto.MGraphCreateEventDto) (*graphmodels.Eventable, error) {
	// Create an event
	requestBody := graphmodels.NewEvent()

	// Set the event subject
	requestBody.SetSubject(&requestDto.Subject)

	// Set the event body
	contentBody := graphmodels.NewItemBody()
	contentBodyType := graphmodels.HTML_BODYTYPE
	contentBody.SetContentType(&contentBodyType)
	contentBody.SetContent(&requestDto.Content)
	requestBody.SetBody(contentBody)

	// Set the start time
	start := graphmodels.NewDateTimeTimeZone()
	start.SetDateTime(&requestDto.StartTime)
	start.SetTimeZone(&requestDto.TimeZone)
	requestBody.SetStart(start)

	// Set the end time
	end := graphmodels.NewDateTimeTimeZone()
	end.SetDateTime(&requestDto.EndTime)
	end.SetTimeZone(&requestDto.TimeZone)
	requestBody.SetEnd(end)

	// Set recurrence
	if requestDto.IsRecurring == true {
		recurrenceObj := graphmodels.NewPatternedRecurrence()

		// set recurrence pattern
		recurrencePattern := graphmodels.NewRecurrencePattern()

		// -- set pattern type
		if &requestDto.PatternType != nil {
			patternType, err := graphmodels.ParseRecurrencePatternType(*requestDto.PatternType)
			if err != nil {
				printOdataError(err)
				errorMessage := err.(*odataerrors.ODataError).GetErrorEscaped().GetMessage()
				return nil, errors.New(*errorMessage)
			}

			if pt, ok := (patternType).(*graphmodels.RecurrencePatternType); ok {
				recurrencePattern.SetTypeEscaped(pt)
			}
		}

		// -- set pattern interval
		patternInterval := requestDto.PatternInterval
		recurrencePattern.SetInterval(patternInterval)

		// -- set pattern for days of week
		if *requestDto.PatternDaysOfWeek != nil {
			if len(*requestDto.PatternDaysOfWeek) > 0 {
				patternDaysOfWeek := []graphmodels.DayOfWeek{}
				for _, dayOfWeek := range *requestDto.PatternDaysOfWeek {
					dayOfWeek, err := graphmodels.ParseDayOfWeek(dayOfWeek)
					if err != nil {
						printOdataError(err)
						errorMessage := err.(*odataerrors.ODataError).GetErrorEscaped().GetMessage()
						return nil, errors.New(*errorMessage)
					}

					if dow, ok := (dayOfWeek).(*graphmodels.DayOfWeek); ok {
						patternDaysOfWeek = append(patternDaysOfWeek, *dow)
					}
				}
				recurrencePattern.SetDaysOfWeek(patternDaysOfWeek)
			}
		}

		// -- TODO: set pattern for days of month

		recurrenceObj.SetPattern(recurrencePattern)

		// -- set recurrence range
		recurrenceRange := graphmodels.NewRecurrenceRange()

		if &requestDto.RecurrenceType != nil {
			rangeType, err := graphmodels.ParseRecurrenceRangeType(*requestDto.RecurrenceType)
			if err != nil {
				printOdataError(err)
				errorMessage := err.(*odataerrors.ODataError).GetErrorEscaped().GetMessage()
				return nil, errors.New(*errorMessage)
			}

			if rt, ok := (rangeType).(*graphmodels.RecurrenceRangeType); ok {
				recurrenceRange.SetTypeEscaped(rt)
			}
		}

		serializedRangeStart, err := serialization.ParseDateOnly(*requestDto.RecurrenceStart)
		if err != nil {
			printOdataError(err)
			errorMessage := err.(*odataerrors.ODataError).GetErrorEscaped().GetMessage()
			return nil, errors.New(*errorMessage)
		}
		recurrenceRange.SetStartDate(serializedRangeStart)

		serializedRangeEnd, err := serialization.ParseDateOnly(*requestDto.RecurrenceEnd)
		if err != nil {
			printOdataError(err)
			errorMessage := err.(*odataerrors.ODataError).GetErrorEscaped().GetMessage()
			return nil, errors.New(*errorMessage)
		}
		recurrenceRange.SetEndDate(serializedRangeEnd)

		recurrenceObj.SetRangeEscaped(recurrenceRange)

		requestBody.SetRecurrence(recurrenceObj)
	}

	// Set attendees
	attendees := []graphmodels.Attendeeable{}
	for _, attendee := range requestDto.Attendees {
		attendeeObj := graphmodels.NewAttendee()

		// Set attendee email information
		emailObj := graphmodels.NewEmailAddress()
		emailObj.SetAddress(&attendee.EmailAddress)
		emailObj.SetName(&attendee.Name)
		attendeeObj.SetEmailAddress(emailObj)

		// Set attendee type
		attendeeType, err := graphmodels.ParseAttendeeType(attendee.AttendeeType)
		if err != nil {
			printOdataError(err)
			errorMessage := err.(*odataerrors.ODataError).GetErrorEscaped().GetMessage()
			return nil, errors.New(*errorMessage)
		}

		if at, ok := attendeeType.(graphmodels.AttendeeType); ok {
			attendeeObj.SetTypeEscaped(&at)
		}

		attendees = append(attendees, attendeeObj)
	}
	requestBody.SetAttendees(attendees)

	// Set location
	if &requestDto.Locations != nil && len(*requestDto.Locations) > 1 {
		// multiple locations
		for _, location := range *requestDto.Locations {
			locationObj := graphmodels.NewLocation()
			locationObj.SetDisplayName(&location.DisplayName)

			// set address if there is one
			if &location.Address != nil {
				addressObj := graphmodels.NewPhysicalAddress()
				addressObj.SetStreet(&location.Address.Street)
				addressObj.SetCity(&location.Address.City)
				addressObj.SetState(&location.Address.State)
				addressObj.SetCountryOrRegion(&location.Address.Country)
				addressObj.SetPostalCode(&location.Address.PostalCode)
				locationObj.SetAddress(addressObj)
			}

			// check for default location
			if location.DefaultLocation == true {
				defaultLocation := graphmodels.DEFAULTESCAPED_LOCATIONTYPE
				locationObj.SetLocationType(&defaultLocation)
			}

			requestBody.SetLocation(locationObj)
		}
	} else if *requestDto.Locations != nil && len(*requestDto.Locations) == 1 {
		// single location
		location := (*requestDto.Locations)[0] // access first element
		locationObj := graphmodels.NewLocation()
		locationObj.SetDisplayName(&location.DisplayName)

		// set address if there is one
		if &location.Address != nil {
			addressObj := graphmodels.NewPhysicalAddress()
			addressObj.SetStreet(&location.Address.Street)
			addressObj.SetCity(&location.Address.City)
			addressObj.SetState(&location.Address.State)
			addressObj.SetCountryOrRegion(&location.Address.Country)
			addressObj.SetPostalCode(&location.Address.PostalCode)
			locationObj.SetAddress(addressObj)
		}

		requestBody.SetLocation(locationObj)
	}

	// Set Online Meeting
	if requestDto.IsOnlineMeeting == true {
		requestBody.SetIsOnlineMeeting(&requestDto.IsOnlineMeeting)
		if requestDto.OnlineMeetingProvider != nil {
			onlineMeetingProvider, err := graphmodels.ParseOnlineMeetingProviderType(*requestDto.OnlineMeetingProvider)
			if err != nil {
				printOdataError(err)
				errorMessage := err.(*odataerrors.ODataError).GetErrorEscaped().GetMessage()
				return nil, errors.New(*errorMessage)
			}

			if omp, ok := (onlineMeetingProvider).(*graphmodels.OnlineMeetingProviderType); ok {
				requestBody.SetOnlineMeetingProvider(omp)
			}
		}
	}

	event, err := m.graphClient.Users().ByUserId("24dc94f1-08bf-4d47-850b-5690533b8236").Events().Post(context.Background(), requestBody, nil)
	if err != nil {
		printOdataError(err)
		errorMessage := err.(*odataerrors.ODataError).GetErrorEscaped().GetMessage()
		return nil, errors.New(*errorMessage)
	}
	// better way to handle creation of events and sync? should we wait for delta? or just return the event?
	log.Println(event.GetBody().GetContent())
	return &event, nil
}
