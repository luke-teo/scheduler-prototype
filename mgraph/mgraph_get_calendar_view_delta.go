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

func (m *MGraph) GetCalendarViewDelta(requestStartDateTime string, requestEndDateTime string, userDto dto.UserDto) (*string, *[]graphmodels.Eventable, error) {
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
		return nil, nil, errors.New(*errorMessage)
	}

	// instantiate data store
	var eventData []graphmodels.Eventable
	for _, event := range delta.GetValue() {
		// check for event type, if series master, get the instance, loop and add
		eventType := event.GetTypeEscaped()
		if *eventType == graphmodels.OCCURRENCE_EVENTTYPE || *eventType == graphmodels.EXCEPTION_EVENTTYPE {
			// skip occurrence & exception type as they only reference back to the series master
			// will get it through series master instance below
			continue
		} else if *eventType == graphmodels.SERIESMASTER_EVENTTYPE {
			instances, err := m.GetEventSeriesMasterInstance(requestStartDateTime, requestEndDateTime, userDto.UserId.String(), *event.GetId())
			if err != nil {
				printOdataError(err)
				errorMessage := err.(*odataerrors.ODataError).GetErrorEscaped().GetMessage()
				return nil, nil, errors.New(*errorMessage)
			}

			// for each instance (occurence or exception) add to eventData
			for _, instance := range instances.GetValue() {
				eventData = append(eventData, instance)
			}
		}

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
			return nil, nil, errors.New(*errorMessage)
		}

		// populate event data
		for _, event := range delta.GetValue() {
			// check for event type, if series master, get the instance, loop and add
			eventType := event.GetTypeEscaped()
			if *eventType == graphmodels.OCCURRENCE_EVENTTYPE || *eventType == graphmodels.EXCEPTION_EVENTTYPE {
				// skip occurrence & exception type as they only reference back to the series master
				// will get it through series master instance below
				continue
			} else if *eventType == graphmodels.SERIESMASTER_EVENTTYPE {
				log.Println("-- series master --")
				instances, err := m.GetEventSeriesMasterInstance(requestStartDateTime, requestEndDateTime, userDto.UserId.String(), *event.GetId())
				if err != nil {
					printOdataError(err)
					errorMessage := err.(*odataerrors.ODataError).GetErrorEscaped().GetMessage()
					return nil, nil, errors.New(*errorMessage)
				}

				// for each instance (occurence or exception) add to eventData
				for _, instance := range instances.GetValue() {
					log.Println("-- instance --")
					log.Printf("instance-subject: %v\n", *instance.GetSubject())
					eventData = append(eventData, instance)
				}
			}

			// populate eventData
			log.Printf("event-type: %v\n", *event.GetTypeEscaped())
			log.Printf("event-subject: %v\n", *event.GetSubject())
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

	return deltaLink, &eventData, nil
}
