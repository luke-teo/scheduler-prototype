package mgraph

import (
	"context"
	"errors"
	"log"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	graphusers "github.com/microsoftgraph/msgraph-sdk-go/users"
	"github.com/scheduler-prototype/dto"
)

func (m *MGraph) GetDeltaCalendarView(requestStartDateTime string, requestEndDateTime string, userDto dto.UserDto) (*[]graphmodels.Eventable, error) {
	// optional header config for restricting page size
	headers := abstractions.NewRequestHeaders()
	headers.Add("Prefer", "odata.maxpagesize=2")

	// construct config params for delta request
	requestParameters := &graphusers.ItemCalendarViewDeltaRequestBuilderGetQueryParameters{
		StartDateTime: &requestStartDateTime,
		EndDateTime:   &requestEndDateTime,
	}

	configuration := &graphusers.ItemCalendarViewDeltaRequestBuilderGetRequestConfiguration{
		Headers:         headers,
		QueryParameters: requestParameters,
	}

	delta, err := m.graphClient.Users().ByUserId(userDto.UserId.String()).CalendarView().Delta().Get(context.Background(), configuration)
	if err != nil {
		printOdataError(err)
		errorMessage := err.(*odataerrors.ODataError).GetErrorEscaped().GetMessage()
		return nil, errors.New(*errorMessage)
	}

	// instantiate data store
	var eventData []graphmodels.Eventable
	for _, event := range delta.GetValue() {
		// populate eventData
		eventData = append(eventData, event)
	}

	// insantiate initial tokens variable
	nextLink := delta.GetOdataNextLink()
	deltaLink := delta.GetOdataDeltaLink()

	// while nextLink is still available and deltaLink is not finalized, keep looping
	for nextLink != nil && deltaLink == nil {
		// craete new request for next page
		requestBuilder := graphusers.NewItemCalendarViewDeltaRequestBuilder(*nextLink, m.adapter)
		nextPage, err := requestBuilder.Get(context.Background(), nil)
		if err != nil {
			printOdataError(err)
			errorMessage := err.(*odataerrors.ODataError).GetErrorEscaped().GetMessage()
			return nil, errors.New(*errorMessage)
		}

		// populate event data
		for _, event := range nextPage.GetValue() {
			log.Printf("subject: %v\n", event.GetSubject())
			eventData = append(eventData, event)
		}

		// insantiate tokens variable
		newNextLink := nextPage.GetOdataNextLink()
		newDeltaLink := nextPage.GetOdataDeltaLink()

		// replace variable values with new values
		nextLink = newNextLink
		deltaLink = newDeltaLink

		log.Printf("nextLink:  %v\n", nextLink)
		log.Printf("deltaLink: %v\n", deltaLink)

	}

	return &eventData, nil
}
