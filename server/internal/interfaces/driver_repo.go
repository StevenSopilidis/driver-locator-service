package interfaces

import (
	"context"

	"github.com/StevenSopilidis/driver-locator-service/internal/domain"
)

type DriverRepo interface {
	CreateDriver(ctx context.Context, driver domain.Driver) error
	GetDriversWithingRadius(ctx context.Context, lat float64, long float64, radiusKm float64, count int) ([]domain.Driver, error)
}
