package workerpool

import (
	"log"

	"github.com/StevenSopilidis/driver-locator-service/internal/domain"
)

type WorkerPool struct {
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{}
}

func (w *WorkerPool) Run(dataCh <-chan *domain.Driver) {
	for data := range dataCh {
		go w.processData(data)
	}
}

func (w *WorkerPool) processData(data *domain.Driver) {
	log.Printf("Worker processing data for: %s", data.Id.String())
}
