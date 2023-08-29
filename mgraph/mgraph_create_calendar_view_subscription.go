package mgraph

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

func (m *MGraph) CreateCalendarViewSubscription(userId *string) (graphmodels.Subscriptionable, error) {
	// create a subscription for specific user
	// -- should optimally be called on first sync so that any changes can be tracked
	requestBody := graphmodels.NewSubscription()
	changeType := "created,updated,deleted"
	requestBody.SetChangeType(&changeType)

	// the below url should be the domain + routes to handle the notification response
	// in local environment it will be the ngrok url
	// otherwise ensure that it is set in the env value or if there are better solutions
	notificationUrl := os.Getenv("AZURE_NOTIFICATION_URL")
	log.Printf("notificationUrl: %s", notificationUrl)
	requestBody.SetNotificationUrl(&notificationUrl)
	lifecycleNotificationUrl := os.Getenv("AZURE_LIFECYCLE_NOTIFICATION_URL")
	log.Printf("lifecycleNotificationUrl: %s", lifecycleNotificationUrl)
	requestBody.SetLifecycleNotificationUrl(&lifecycleNotificationUrl)

	resource := fmt.Sprintf("/users/%s/events", *userId)
	requestBody.SetResource(&resource)

	// Get the current time in the user's timezone
	now := time.Now().UTC()

	// Calculate the end date time 2 days from now
	endDateTime := now.Add(time.Duration(2) * 24 * time.Hour)

	// Format the end date time as a string
	endDateTimeString := endDateTime.Format(time.RFC3339)
	expirationDateTime, err := time.Parse(time.RFC3339, endDateTimeString)
	requestBody.SetExpirationDateTime(&expirationDateTime)

	// client state is like a signature, make it unique to know that info is actually coming from microsoft
	clientState := os.Getenv("AZURE_CLIENT_STATE_SECRET")
	log.Printf("clientState: %s", clientState)
	requestBody.SetClientState(&clientState)

	subscriptions, err := m.graphClient.Subscriptions().Post(context.Background(), requestBody, nil)
	if err != nil {
		printOdataError(err)
		errorMessage := err.(*odataerrors.ODataError).GetErrorEscaped().GetMessage()
		return nil, errors.New(*errorMessage)
	}

	return subscriptions, nil
}
