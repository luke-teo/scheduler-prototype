package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	msjson "github.com/microsoft/kiota-serialization-json-go"
	"github.com/scheduler-prototype/dto"
	"github.com/scheduler-prototype/mgraph"
	"github.com/scheduler-prototype/repository"
	"github.com/scheduler-prototype/utility"
)

type Handler struct {
	client *mgraph.MGraph
	repo   *repository.Repository
}

func NewHandler(client *mgraph.MGraph, repo *repository.Repository) *Handler {
	return &Handler{
		client: client,
		repo:   repo,
	}
}

func (h *Handler) MGraphGetCalendarView(w http.ResponseWriter, r *http.Request) {
	// get the current time in the user's timezone
	now := time.Now().UTC().Add(time.Duration(time.Hour * -8))

	// calculate the start of this week
	startOfWeek := now.Truncate(time.Hour * 24).Add(time.Duration(time.Hour * 24 * time.Duration(int(now.Weekday())-1) * -1))

	// calculate the end of next week
	endOfNextWeek := startOfWeek.Add(time.Duration(time.Hour * 24 * 7 * 2)).Add(time.Duration(time.Hour * -1))

	requestStart := time.Now()

	// format the dates as strings
	requestStartDateTime := startOfWeek.Format(time.RFC3339)
	requestEndDateTime := endOfNextWeek.Format(time.RFC3339)

	events, err := h.client.GetCalendarView(requestStartDateTime, requestEndDateTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	requestDuration := time.Since(requestStart)

	// Converting the model into a JSON object
	serializer := msjson.NewJsonSerializationWriter()
	events.Serialize(serializer)
	eventsJson, _ := serializer.GetSerializedContent()
	w.Header().Set("Content-Type", "application/json")

	// Iterating over events
	for _, event := range events.GetValue() {
		iCalUid := *event.GetICalUId()

		// Check if event already exists in DB
		_, err := h.repo.GetEventByICalUid(iCalUid)
		if err != nil {
			if err == utility.ErrNotFound {
				meetingUrl := ""
				if event.GetOnlineMeeting() != nil {
					meetingUrl = *event.GetOnlineMeeting().GetJoinUrl()
				}
				// Event creation
				eventDto := &dto.MGraphEventDto{
					UserId:          "1",
					ICalUid:         *event.GetICalUId(),
					EventId:         *event.GetId(),
					Title:           *event.GetSubject(),
					Description:     *event.GetBody().GetContent(),
					LocationsCount:  len(event.GetLocations()),
					StartTime:       *event.GetStart().GetDateTime(),
					EndTime:         *event.GetEnd().GetDateTime(),
					IsOnline:        *event.GetIsOnlineMeeting(),
					IsAllDay:        *event.GetIsAllDay(),
					IsCancelled:     *event.GetIsCancelled(),
					OrganizerUserId: "1",
					CreatedTime:     *event.GetCreatedDateTime(),
					UpdatedTime:     *event.GetLastModifiedDateTime(),
					Timezone:        *event.GetStart().GetTimeZone(),
					PlatformUrl:     *event.GetWebLink(),
					MeetingUrl:      meetingUrl,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
					IsRecurring:     event.GetSeriesMasterId() != nil,
					SeriesMasterId:  event.GetSeriesMasterId(),
				}

				err = h.repo.CreateEvent(eventDto)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					response := map[string]string{"error": err.Error()}
					json.NewEncoder(w).Encode(response)
					return
				}

				// Attendees creation
				// TODO: Check for duplicates
				for _, attendee := range event.GetAttendees() {
					// emailAddress := *attendee.GetEmailAddress().GetAddress()
					// iCalUid := *event.GetICalUId()

					attendeeDto := &dto.MGraphAttendeeDto{
						UserId:       "1",
						Name:         *attendee.GetEmailAddress().GetName(),
						EmailAddress: *attendee.GetEmailAddress().GetAddress(),
						ICalUid:      event.GetICalUId(),
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
					}
					err := h.repo.CreateAttendee(attendeeDto)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						response := map[string]string{"error": err.Error()}
						json.NewEncoder(w).Encode(response)
						return
					}
				}

				// Location creation
				// TODO: Check for duplicates
				for _, location := range event.GetLocations() {
					// combine address props to create a single string
					street := *location.GetAddress().GetStreet()
					city := *location.GetAddress().GetCity()
					state := *location.GetAddress().GetState()
					postalCode := *location.GetAddress().GetPostalCode()
					country := *location.GetAddress().GetCountryOrRegion()

					address := street + ", " + city + ", " + state + ", " + postalCode + ", " + country
					if address == ",,,," {
						address = ""
					}

					locationDto := &dto.MGraphLocationDto{
						ICalUid:     event.GetICalUId(),
						DisplayName: location.GetDisplayName(),
						LocationUri: location.GetLocationUri(),
						Address:     address,
					}

					err := h.repo.CreateLocation(locationDto)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						response := map[string]string{"error": err.Error()}
						json.NewEncoder(w).Encode(response)
						return
					}
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				response := map[string]string{"error": err.Error()}
				json.NewEncoder(w).Encode(response)
				return
			}
		}
		continue
	}

	// Print the time the request took
	fmt.Printf("Graph Request took: %s\n", requestDuration)

	w.WriteHeader(http.StatusOK)
	w.Write(eventsJson)
}
