package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/scheduler-prototype/dto"
	requestDto "github.com/scheduler-prototype/dto/request"
	"github.com/scheduler-prototype/utility"
)

func (h *Handler) MGraphCalendarViewFirstSync(w http.ResponseWriter, r *http.Request) {
	// read the request body and insert into users table
	req := &requestDto.MGraphCalendarViewFirstSyncDto{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// insert into database if user does not exists
	newUserUuid, err := uuid.Parse(req.UserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = h.repo.GetUserByUserId(&newUserUuid)
	if err != nil {
		if err == utility.ErrNotFound {
			// create user if not exists
			newUser := &dto.UserDto{
				UserId:        newUserUuid,
				CurrentDelta:  nil,
				PreviousDelta: nil,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}

			err = h.repo.CreateUser(newUser)
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

	// once user is confirmed to be in database, grab the user and make first delta queries to Microsoft Graph
	userDto, err := h.repo.GetUserByUserId(&newUserUuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// make first delta queries to Microsoft Graph
	// Get the current time in the user's timezone
	now := time.Now().UTC().Add(time.Duration(time.Hour * -8))

	// Calculate the start of the current month
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	// Calculate the start of the next month
	nextMonth := now.AddDate(0, 1, 0)
	startOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, time.UTC)

	// Calculate the end of the next month
	endOfNextMonth := startOfNextMonth.Add(-time.Second)

	requestStart := time.Now()
	deltaLink, events, err := h.client.GetCalendarViewDelta(startOfMonth.Format(time.RFC3339), endOfNextMonth.Format(time.RFC3339), userDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	requestDuration := time.Since(requestStart)
	// Print the time the request took
	fmt.Printf("Graph Delta Request took: %s\n", requestDuration)

	for _, event := range *events {
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
					Type:            event.GetTypeEscaped().String(),
					IsRecurring:     event.GetSeriesMasterId() != nil,
					SeriesMasterId:  event.GetSeriesMasterId(),
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
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
	log.Println("completed processing events")

	// once events are inputted into the database, we need to store the token
	userDto.CurrentDelta = deltaLink
	err = h.repo.UpdateCurrentDeltaByUser(&userDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create subscription for the user
	userUuid := userDto.UserId.String()
	// create request to Microsoft Graph to create the event
	requestStart = time.Now()
	subscriptions, err := h.client.CreateCalendarViewSubscription(&userUuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	log.Printf("subscriptions: %s", subscriptions)
	requestDuration = time.Since(requestStart)
	fmt.Printf("Graph Subscription Request took: %s\n", requestDuration)

	userDto.SubscriptionId = subscriptions.GetId()
	userDto.SubscriptionExpiresAt = subscriptions.GetExpirationDateTime()
	log.Printf("userDto.SubscriptionId: %s", *userDto.SubscriptionId)
	log.Printf("userDto.SubscriptionExpiresAt: %s", *userDto.SubscriptionExpiresAt)
	err = h.repo.UpdateSubscriptionInfoByUser(&userDto)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Println("completed creating subscription")

	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Events successfully synced", "data": *deltaLink}
	json.NewEncoder(w).Encode(response)
}
