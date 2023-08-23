package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserDto struct {
	ID            uuid.UUID
	UserId        uuid.UUID
	CurrentDelta  *string
	PreviousDelta *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
