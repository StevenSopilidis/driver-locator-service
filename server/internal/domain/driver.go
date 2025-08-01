package domain

import (
	"time"

	"github.com/google/uuid"
)

type Driver struct {
	Id        uuid.UUID
	Last_seen time.Time
	Latitude  float64
	Longitude float64
}
