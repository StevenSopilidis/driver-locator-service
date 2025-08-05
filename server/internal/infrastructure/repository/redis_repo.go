package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/StevenSopilidis/driver-locator-service/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisRepo(addr string, ttl time.Duration) (*RedisRepo, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("could not connect to redis: %s", err.Error())
	}

	return &RedisRepo{
		client: client,
		ttl:    ttl,
	}, nil
}

func (r *RedisRepo) CreateDriver(ctx context.Context, driver domain.Driver) error {
	_, err := r.client.GeoAdd(ctx, "drivers:locations", &redis.GeoLocation{
		Name:      driver.Id.String(),
		Latitude:  driver.Latitude,
		Longitude: driver.Longitude,
	}).Result()

	if err != nil {
		return fmt.Errorf("failed to set driver to repo: %s", err.Error())
	}

	return nil
}

func (r *RedisRepo) GetDriversWithingRadius(ctx context.Context, lat float64, long float64, radiusKm float64, count int) ([]domain.Driver, error) {
	drivers, err := r.client.GeoSearchLocation(ctx, "drivers:locations", &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Longitude:  long,
			Latitude:   lat,
			Radius:     radiusKm,
			RadiusUnit: "km",
			Sort:       "ASC",
			Count:      count,
		},
		WithCoord: true,
		WithDist:  true,
	}).Result()

	if err != nil {
		return nil, fmt.Errorf("could not get nearby users %s", err.Error())
	}

	result := make([]domain.Driver, len(drivers))
	for _, driver := range drivers {
		id, err := uuid.Parse(driver.Name)
		if err != nil {
			continue
		}

		result = append(result, domain.Driver{
			Id:        id,
			Longitude: driver.Longitude,
			Latitude:  driver.Latitude,
		})
	}

	return result, nil
}
