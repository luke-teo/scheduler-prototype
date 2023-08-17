package dto

import "time"

type MGraphEventDto struct {
	ID              int64
	UserId          string
	ICalId          string
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
	MeetingUrl      string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
