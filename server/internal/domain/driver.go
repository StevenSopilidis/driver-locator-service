package domain

import (
	"github.com/google/uuid"
)

type Driver struct {
	Id        uuid.UUID
	Latitude  float64
	Longitude float64
}
