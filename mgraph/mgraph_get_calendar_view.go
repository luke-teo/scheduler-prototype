package mgraph

import (
	"context"
	"errors"

	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	graphusers "github.com/microsoftgraph/msgraph-sdk-go/users"
)

func (m *MGraph) GetCalendarView(requestStartDateTime string, requestEndDateTime string) (graphmodels.EventCollectionResponseable, error) {
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
