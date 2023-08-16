package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
	db, err := sql.Open("postgres", dbConnStr)
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

	token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: scopes,
	})
	if err != nil { 
		log.Fatal(err)
	}

	fmt.Println(token.Token)

	// chi router 
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// // get the current time in the user's timezone
		// now := time.Now().UTC().Add(time.Duration(time.Hour * -8))

		// // calculate the start of this week
		// startOfWeek := now.Truncate(time.Hour * 24).Add(time.Duration(time.Hour * 24 * time.Duration(int(now.Weekday())-1) * -1))

		// // calculate the end of next week
		// endOfNextWeek := startOfWeek.Add(time.Duration(time.Hour * 24 * 7 * 2)).Add(time.Duration(time.Hour * -1))

		// format the dates as strings
		requestStartDateTime := "2023-08-16T08:33:50.415Z"
		requestEndDateTime := "2023-08-23T08:33:50.415Z"

		requestParameters := &graphusers.ItemCalendarCalendarViewRequestBuilderGetQueryParameters{
			StartDateTime: &requestStartDateTime,
			EndDateTime: &requestEndDateTime,
		}
		configuration := &graphusers.ItemCalendarCalendarViewRequestBuilderGetRequestConfiguration{
			QueryParameters: requestParameters,
		}

		events, err := client.Users().ByUserId("24dc94f1-08bf-4d47-850b-5690533b8236").Calendar().CalendarView().Get(context.Background(), configuration)	
		if err != nil {
			printOdataError(err)
		}

		eventsJSON, err := json.Marshal(events)
		if err != nil {
			log.Fatal(err)
		}

		// users, err := client.Users().ByUserId("24dc94f1-08bf-4d47-850b-5690533b8236").Get(context.Background(), nil)
		// if err != nil {
		// 	printOdataError(err)
		// }

		log.Println(eventsJSON)

		w.Write([]byte("welcome"))
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