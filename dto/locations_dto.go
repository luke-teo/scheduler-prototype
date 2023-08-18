package dto

import (
	"time"

	"github.com/google/uuid"
)

type MGraphLocationDto struct {
	ID          uuid.UUID
	ICalUid     string
	DisplayName string
	LocationUri *string
	Address     *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
