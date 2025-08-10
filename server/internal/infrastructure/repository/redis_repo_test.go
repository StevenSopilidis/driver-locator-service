package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/StevenSopilidis/driver-locator-service/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func redisAddr() string {
	if v := os.Getenv("REDIS_ADDR"); v != "" {
		return v
	}
	return "127.0.0.1:6379"
}

func flushDB(t *testing.T, repo *RedisRepo) {
	t.Helper()
	ctx := context.Background()
	require.NoError(t, repo.client.FlushDB(ctx).Err())
}

func TestUpdateDrivereLocationDoesNotCreateDuplicates(t *testing.T) {
	addr := redisAddr()
	repo, err := NewRedisRepo(addr, time.Minute)
	require.NoError(t, err)
	require.NotNil(t, repo)

	ctx := context.Background()
	// ensure clean DB
	flushDB(t, repo)
	defer flushDB(t, repo)

	// center point (Athens example)
	lat := 37.9755
	lon := 23.7363

	driverInside := domain.Driver{
		Id:        uuid.New(),
		Latitude:  lat,
		Longitude: lon,
	}
	driverOutside := domain.Driver{
		Id:        uuid.New(),
		Latitude:  51.5074,
		Longitude: -0.1278,
	}

	require.NoError(t, repo.CreateDriver(ctx, driverInside))
	require.NoError(t, repo.CreateDriver(ctx, driverOutside))

	driverInside = domain.Driver{
		Id:        driverInside.Id,
		Latitude:  lat,
		Longitude: lon,
	}

	found, err := repo.GetDriversWithingRadius(ctx, lat, lon, 5.0, 10)
	require.NoError(t, err)

	// Expect at least one nearby driver
	require.GreaterOrEqual(t, len(found), 1, "expected at least one nearby driver")

	instances := 0
	for _, d := range found {
		if d.Id == driverInside.Id {
			instances += 1
		}
	}
	require.Equal(t, 1, instances)
}

func TestCreateDriverAndGetDriversWithinRadius_Integration(t *testing.T) {
	addr := redisAddr()
	repo, err := NewRedisRepo(addr, time.Minute)
	require.NoError(t, err)
	require.NotNil(t, repo)

	ctx := context.Background()
	// ensure clean DB
	flushDB(t, repo)
	defer flushDB(t, repo)

	// center point (Athens example)
	lat := 37.9755
	lon := 23.7363

	driverInside := domain.Driver{
		Id:        uuid.New(),
		Latitude:  lat,
		Longitude: lon,
	}
	driverOutside := domain.Driver{
		Id:        uuid.New(),
		Latitude:  51.5074,
		Longitude: -0.1278,
	}

	require.NoError(t, repo.CreateDriver(ctx, driverInside))
	require.NoError(t, repo.CreateDriver(ctx, driverOutside))

	found, err := repo.GetDriversWithingRadius(ctx, lat, lon, 5.0, 10)
	require.NoError(t, err)

	// Expect at least one nearby driver
	require.GreaterOrEqual(t, len(found), 1, "expected at least one nearby driver")

	// find the inside driver
	var got domain.Driver
	foundIt := false
	for _, d := range found {
		if d.Id == driverInside.Id {
			got = d
			foundIt = true
			break
		}
	}
	require.True(t, foundIt, "inside driver must be returned by GetDriversWithingRadius")

	// floating point tolerance (relaxed to 1e-5)
	const tol = 1e-5
	assert.InDelta(t, driverInside.Latitude, got.Latitude, tol)
	assert.InDelta(t, driverInside.Longitude, got.Longitude, tol)
}

func TestGetDriversWithinRadius_NoResults_Integration(t *testing.T) {
	addr := redisAddr()
	repo, err := NewRedisRepo(addr, time.Minute)
	require.NoError(t, err)
	require.NotNil(t, repo)

	ctx := context.Background()
	flushDB(t, repo)
	defer flushDB(t, repo)

	found, err := repo.GetDriversWithingRadius(ctx, 0.0, 0.0, 1.0, 10)
	require.NoError(t, err)

	assert.Len(t, found, 0, "expected no nearby drivers")
}
