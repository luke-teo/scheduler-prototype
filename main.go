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
	"github.com/scheduler-prototype/mgraph"
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

	// chi router
	r := chi.NewRouter()

	subRouter := chi.NewRouter()
	subRouter.Get("/calendarview", mGraphGetCalendarView(client))

	r.Mount("/mgraph", subRouter)

	http.ListenAndServe(":8080", r)
}

func mGraphGetCalendarView(client *mgraph.MGraph) http.HandlerFunc {
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
		w.Write([]byte("{"))
		w.Write(eventsJson)
		w.Write([]byte("}"))

		// Itterating over events
		for _, event := range events.GetValue() {
			fmt.Println("Event found:")
			fmt.Printf("Event: %s\n", *event.GetSubject())
			fmt.Printf("Start: %s\n", *event.GetStart().GetDateTime())
			fmt.Printf("End: %s\n", *event.GetEnd().GetDateTime())
			fmt.Printf("Location: %s\n", *event.GetLocation().GetDisplayName())
			fmt.Printf("Attendees: ")
			attendees := event.GetAttendees()
			for _, attendee := range attendees {
				fmt.Printf("%s ", *attendee.GetEmailAddress().GetAddress())
			}
			fmt.Printf("\n\n")
		}

		// Print the time the request took
		fmt.Printf("Graph Request took: %s\n", requestDuration)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(eventsJson)
	}
}
