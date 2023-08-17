package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	msjson "github.com/microsoft/kiota-serialization-json-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	graphusers "github.com/microsoftgraph/msgraph-sdk-go/users"
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
		panic(err.Error())
	}

	fmt.Println("Successfully connected to database!")

	// getting azure credentials object
	tenantId := os.Getenv("AZURE_TENANT_ID")
	clientId := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	scopes := []string{"https://graph.microsoft.com/.default"}

	cred, _ := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)

	// getting request adapter
	client, _ := msgraphsdk.NewGraphServiceClientWithCredentials(cred, scopes)

	// chi router

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// get the current time in the user's timezone
		now := time.Now().UTC().Add(time.Duration(time.Hour * -8))

		// calculate the start of this week
		startOfWeek := now.Truncate(time.Hour * 24).Add(time.Duration(time.Hour * 24 * time.Duration(int(now.Weekday())-1) * -1))

		// calculate the end of next week
		endOfNextWeek := startOfWeek.Add(time.Duration(time.Hour * 24 * 7 * 2)).Add(time.Duration(time.Hour * -1))

		// format the dates as strings
		requestStartDateTime := startOfWeek.Format(time.RFC3339)
		requestEndDateTime := endOfNextWeek.Format(time.RFC3339)

		requestParameters := &graphusers.ItemCalendarCalendarViewRequestBuilderGetQueryParameters{
			StartDateTime: &requestStartDateTime,
			EndDateTime:   &requestEndDateTime,
		}
		configuration := &graphusers.ItemCalendarCalendarViewRequestBuilderGetRequestConfiguration{
			QueryParameters: requestParameters,
		}

		// Calculate time the request takes
		requestStart := time.Now()

		events, err := client.Users().ByUserId("24dc94f1-08bf-4d47-850b-5690533b8236").Calendar().CalendarView().Get(context.Background(), configuration)
		if err != nil {
			printOdataError(err)
		}

		requestDuration := time.Since(requestStart)

		// Converting the model into a JSON object
		serializer := msjson.NewJsonSerializationWriter()
		events.Serialize(serializer)
		json, _ := serializer.GetSerializedContent()
		fmt.Println("JSON output:")
		fmt.Println(string(json)) // Print JSON to console

		// Write JSON to the Chi response
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{"))
		w.Write(json)
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
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}

func printOdataError(err error) {
	switch err.(type) {
	case *odataerrors.ODataError:
		typed := err.(*odataerrors.ODataError)
		fmt.Printf("error: %s", typed.Error())
		if terr := typed.GetErrorEscaped(); terr != nil {
			fmt.Printf("code: %s", *terr.GetCode())
			fmt.Printf("msg: %s", *terr.GetMessage())
		}
	default:
		fmt.Printf("%T > error: %#v", err, err)
	}
}
