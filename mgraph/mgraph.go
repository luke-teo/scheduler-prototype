package mgraph

import (
	"context"
	"errors"
	"fmt"
	"os"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	_ "github.com/lib/pq"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	graphusers "github.com/microsoftgraph/msgraph-sdk-go/users"
)

type MGraphInterface interface {
	Init() (*msgraphsdk.GraphServiceClient, error)
	GetCalendarView(client *msgraphsdk.GraphServiceClient, requestStartDateTime string, requestEndDateTime string) (models.EventCollectionResponseable, error)
}

type MGraph struct {
	credentials *azidentity.ClientSecretCredential
	graphClient *msgraphsdk.GraphServiceClient
}

func NewMGraphClient() (*MGraph, error) {
	tenantId := os.Getenv("AZURE_TENANT_ID")
	clientId := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	scopes := []string{"https://graph.microsoft.com/.default"}

	// Create an OAuth client with the credential.
	cred, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
	if err != nil {
		return nil, err
	}

	// Create client with credentials and required scopes
	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, scopes)

	return &MGraph{credentials: cred, graphClient: client}, nil
}

func (m *MGraph) Init() (*msgraphsdk.GraphServiceClient, error) {
	tenantId := os.Getenv("AZURE_TENANT_ID")
	clientId := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	scopes := []string{"https://graph.microsoft.com/.default"}

	// Create an OAuth client with the credential.
	cred, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
	if err != nil {
		return nil, err
	}

	// Create client with credentials and required scopes
	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, scopes)

	return client, nil
}

func (m *MGraph) GetCalendarView(requestStartDateTime string, requestEndDateTime string) (models.EventCollectionResponseable, error) {
	requestParameters := &graphusers.ItemCalendarCalendarViewRequestBuilderGetQueryParameters{
		StartDateTime: &requestStartDateTime,
		EndDateTime:   &requestEndDateTime,
	}
	configuration := &graphusers.ItemCalendarCalendarViewRequestBuilderGetRequestConfiguration{
		QueryParameters: requestParameters,
	}
	// Get the events
	events, err := m.graphClient.Users().ByUserId("24dc94f1-08bf-4d47-850b-5690533b8236").Calendar().CalendarView().Get(context.Background(), configuration)
	if err != nil {
		printOdataError(err)
		errorMessage := err.(*odataerrors.ODataError).GetErrorEscaped().GetMessage()
		return nil, errors.New(*errorMessage)
	}

	return events, nil
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
