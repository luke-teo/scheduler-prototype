package dto

import (
	"time"

	"github.com/google/uuid"
)

type MGraphAttendeeDto struct {
	ID           uuid.UUID
	UserId       string
	Name         string
	EmailAddress string
	ICalUid      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
