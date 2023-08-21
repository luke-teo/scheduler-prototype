package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	msjson "github.com/microsoft/kiota-serialization-json-go"
	"github.com/scheduler-prototype/dto"
	"github.com/scheduler-prototype/utility"
)

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
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
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
		iCalUid := event.GetICalUId()

		// Check if event already exists in DB
		_, err := h.repo.GetEventByICalUid(*iCalUid)
		if err != nil {
			if err == utility.ErrNotFound {
				var meetingUrl *string
				if event.GetOnlineMeeting() != nil {
					meetingUrl = event.GetOnlineMeeting().GetJoinUrl()
				} else {
					meetingUrl = nil
				}
				// Event creation
				eventDto := &dto.MGraphEventDto{
					UserId:          "1",
					ICalUid:         *iCalUid,
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
				for _, attendee := range event.GetAttendees() {
					emailAddress := *attendee.GetEmailAddress().GetAddress()
					_, err := h.repo.GetAttendeeByICalUidAndEmailAddress(*iCalUid, emailAddress)
					if err != nil {
						if err == utility.ErrNotFound {
							attendeeDto := &dto.MGraphAttendeeDto{
								UserId:       "1",
								Name:         *attendee.GetEmailAddress().GetName(),
								EmailAddress: *attendee.GetEmailAddress().GetAddress(),
								ICalUid:      *iCalUid,
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
						} else {
							w.WriteHeader(http.StatusInternalServerError)
							response := map[string]string{"error": err.Error()}
							json.NewEncoder(w).Encode(response)
							return
						}
					}
				}

				// Location creation
				if len(event.GetLocations()) > 0 {
					for _, location := range event.GetLocations() {
						displayName := location.GetDisplayName()
						_, err := h.repo.GetLocationByICalUidAndDisplayName(*iCalUid, *displayName)
						if err != nil {
							if err == utility.ErrNotFound {
								var address *string
								// combine address props to create a single string
								street := *location.GetAddress().GetStreet()
								city := *location.GetAddress().GetCity()
								state := *location.GetAddress().GetState()
								postalCode := *location.GetAddress().GetPostalCode()
								country := *location.GetAddress().GetCountryOrRegion()

								fullAddress := street + ", " + city + ", " + state + ", " + postalCode + ", " + country
								if fullAddress == ", , , , " {
									address = nil
								} else {
									address = &fullAddress
								}

								locationDto := &dto.MGraphLocationDto{
									ICalUid:     *iCalUid,
									DisplayName: *displayName,
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
							} else {
								w.WriteHeader(http.StatusInternalServerError)
								response := map[string]string{"error": err.Error()}
								json.NewEncoder(w).Encode(response)
								return
							}
						}

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
