package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/StevenSopilidis/driver-locator-service/internal/infrastructure/handlers"
)

func main() {
	server, err := handlers.NewUDPServer("0.0.0.0", 8080, 100_000)
	defer server.Shutdown()

	if err != nil {
		log.Fatalf("Failed to init udp server: %s", err.Error())
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go server.ListenAndServe(ctx)

	<-ctx.Done()
}
