package workerpool

import (
	"context"
	"log"

	"github.com/StevenSopilidis/driver-locator-service/internal/domain"
	"github.com/StevenSopilidis/driver-locator-service/internal/interfaces"
)

type WorkerPool struct {
	repo interfaces.DriverRepo
}

func NewWorkerPool(repo interfaces.DriverRepo) *WorkerPool {
	return &WorkerPool{
		repo: repo,
	}
}

func (w *WorkerPool) Run(dataCh <-chan domain.Driver) {
	for data := range dataCh {
		go w.processData(data)
	}
}

func (w *WorkerPool) processData(data domain.Driver) {
	err := w.repo.CreateDriver(context.Background(), domain.Driver{
		Id:        data.Id,
		Longitude: data.Longitude,
		Latitude:  data.Latitude,
	})
	if err != nil {
		log.Printf("Failed to update driver location for driver: %s with error %s\n", data.Id.String(), err.Error())
		return
	}

	log.Printf("Updated location for driver: %s\n", data.Id.String())
}
