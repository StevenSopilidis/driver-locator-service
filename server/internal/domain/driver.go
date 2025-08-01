package domain

import (
	"time"

	"github.com/google/uuid"
)

type Driver struct {
	id        uuid.UUID
	username  string
	last_seen time.Time
	latitude  float64
	longitude float64
}
