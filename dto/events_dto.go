package dto

import (
	"time"

	"github.com/google/uuid"
)

type MGraphEventDto struct {
	ID              uuid.UUID
	UserId          string
	ICalUid         string
	EventId         string
	Title           string
	Description     string
	LocationsCount  int
	StartTime       string
	EndTime         string
	IsOnline        bool
	IsAllDay        bool
	IsCancelled     bool
	OrganizerUserId string
	CreatedTime     time.Time
	UpdatedTime     time.Time
	Timezone        string
	PlatformUrl     string
	MeetingUrl      *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	IsRecurring     bool
	SeriesMasterId  *string
}
