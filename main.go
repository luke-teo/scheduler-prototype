package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	msjson "github.com/microsoft/kiota-serialization-json-go"
	"github.com/scheduler-prototype/dto"
	"github.com/scheduler-prototype/mgraph"
	"github.com/scheduler-prototype/repository"
	"github.com/scheduler-prototype/utility"
)

func main() {
	// loading env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// db connection
	dbConnStr := os.Getenv("DB_CONN_STR")
	dbDriver := os.Getenv("DB_DRIVER")
	db, err := sql.Open(dbDriver, dbConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to database!")

	// initialize msgraph client
	client, err := mgraph.NewMGraphClient()
	if err != nil {
		log.Fatal(err)
	}

	// initialize reposotories
	repo := repository.NewRepository(db)

	// chi router
	r := chi.NewRouter()

	subRouter := chi.NewRouter()
	subRouter.Get("/calendarview", mGraphGetCalendarView(client, repo))

	r.Mount("/mgraph", subRouter)

	http.ListenAndServe(":8080", r)
}

func mGraphGetCalendarView(client *mgraph.MGraph, repo *repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		events, err := client.GetCalendarView(requestStartDateTime, requestEndDateTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		requestDuration := time.Since(requestStart)

		// Converting the model into a JSON object
		serializer := msjson.NewJsonSerializationWriter()
		events.Serialize(serializer)
		eventsJson, _ := serializer.GetSerializedContent()

		// Write JSON to the Chi response
		w.Header().Set("Content-Type", "application/json")

		// Iterating over events
		for _, event := range events.GetValue() {
			log.Println("entered iteration")
			iCalUid := *event.GetICalUId()

			// Check if event already exists in DB
			_, err := repo.GetEventByICalUid(iCalUid)
			if err != nil {
				if err == utility.ErrNotFound {
					// Create DTO for insertion
					eventDto := &dto.MGraphEventDto{
						UserId:          "1",
						ICalId:          *event.GetICalUId(),
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
						MeetingUrl:      *event.GetOnlineMeetingUrl(),
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					}

					err = repo.CreateEvent(eventDto)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
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
		json.NewEncoder(w).Encode(eventsJson)
	}
}
