package service

import (
	"github.com/StevenSopilidis/driver-locator-service/internal/domain"
	"github.com/google/uuid"
)

type DriverService interface {
	CreateDriver(driver domain.Driver) error
	RemoveDriver(id uuid.UUID) error
	GetDriver(id uuid.UUID) (domain.Driver, error)
}
