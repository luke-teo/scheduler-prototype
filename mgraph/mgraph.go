package mgraph

import (
	"fmt"
	"os"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	_ "github.com/lib/pq"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	requestDto "github.com/scheduler-prototype/dto/request"
)

type MGraphInterface interface {
	Init() (*msgraphsdk.GraphServiceClient, error)
	GetCalendarView(requestStartDateTime string, requestEndDateTime string) (graphmodels.EventCollectionResponseable, error)
	PostCreateEvent(requestDto *requestDto.MGraphCreateEventDto) (*graphmodels.Event, error)
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
