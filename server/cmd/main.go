package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/StevenSopilidis/driver-locator-service/internal/infrastructure/config"
	"github.com/StevenSopilidis/driver-locator-service/internal/infrastructure/handlers"
	"github.com/StevenSopilidis/driver-locator-service/internal/infrastructure/repository"
)

func main() {
	config, err := config.NewConfig(".")
	if err != nil {
		log.Fatalf("Failed to read in config, %s", err.Error())
	}

	repo, err := repository.NewRedisRepo(config.RedisAddr, config.TTL)
	if err != nil {
		log.Fatalf("Could not create redis repo: %s", err.Error())
	}

	server, err := handlers.NewUDPServer(
		config.UdpAddr, config.UdpPort,
		config.MaxConcurrentRequests, repo)
	defer server.Shutdown()

	if err != nil {
		log.Fatalf("Failed to init udp server: %s", err.Error())
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go server.ListenAndServe(ctx)

	<-ctx.Done()
}
