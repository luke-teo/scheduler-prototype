package mgraph

import (
	"context"
	"errors"

	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	graphusers "github.com/microsoftgraph/msgraph-sdk-go/users"
)

func (m *MGraph) GetEventSeriesMasterInstance(requestStartDateTime string, requestEndDateTime string, userId string, eventId string) (graphmodels.EventCollectionResponseable, error) {
	requestParameters := &graphusers.ItemEventsItemInstancesRequestBuilderGetQueryParameters{
		StartDateTime: &requestStartDateTime,
		EndDateTime:   &requestEndDateTime,
	}
	configuration := &graphusers.ItemEventsItemInstancesRequestBuilderGetRequestConfiguration{
		QueryParameters: requestParameters,
	}

	// Get the events instances
	events, err := m.graphClient.Users().ByUserId("24dc94f1-08bf-4d47-850b-5690533b8236").Events().ByEventId(eventId).Instances().Get(context.Background(), configuration)
	if err != nil {
		printOdataError(err)
		errorMessage := err.(*odataerrors.ODataError).GetErrorEscaped().GetMessage()
		return nil, errors.New(*errorMessage)
	}

	return events, nil
}
