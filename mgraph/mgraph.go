package mgraph

import (
	"fmt"
	"os"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	_ "github.com/lib/pq"
	auth "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	requestDto "github.com/scheduler-prototype/dto/request"
)

type MGraphInterface interface {
	GetCalendarView(requestStartDateTime string, requestEndDateTime string) (graphmodels.EventCollectionResponseable, error)
	PostCreateEvent(requestDto *requestDto.MGraphCreateEventDto) (*graphmodels.Event, error)
}

type MGraph struct {
	adapter     *msgraphsdk.GraphRequestAdapter
	credentials *azidentity.ClientSecretCredential
	graphClient *msgraphsdk.GraphServiceClient
	scopes      []string
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

	authProvider, err := auth.NewAzureIdentityAuthenticationProviderWithScopes(cred, scopes)
	if err != nil {
		return nil, err
	}

	adapter, err := msgraphsdk.NewGraphRequestAdapter(authProvider)
	if err != nil {
		return nil, err
	}

	// Create client with credentials and required scopes
	// client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, scopes)
	client := msgraphsdk.NewGraphServiceClient(adapter)

	return &MGraph{adapter: adapter, credentials: cred, graphClient: client, scopes: scopes}, nil
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
