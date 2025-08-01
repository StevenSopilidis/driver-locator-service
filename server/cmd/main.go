package main

import (
	"log"

	"github.com/StevenSopilidis/driver-locator-service/internal/infrastructure/handlers"
)

func main() {
	server, err := handlers.NewUDPServer("0.0.0.0", 8080, 1000)
	defer server.Shutdown()

	if err != nil {
		log.Fatalf("Failed to init udp server: %s", err.Error())
	}

	server.ListenAndServe()
}
