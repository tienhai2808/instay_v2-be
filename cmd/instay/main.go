package main

import (
	"log"

	"github.com/InstayPMS/backend/internal/infrastructure/api"
	"github.com/InstayPMS/backend/internal/infrastructure/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	sv, err := api.NewServer(cfg)
	if err != nil {
		log.Fatalf("Server initialization failed: %v", err)
	}

	ch := make(chan error, 1)
	go func() {
		if err := sv.Start(); err != nil {
			ch <- err
		}
	}()

	log.Printf("Server is running at: http://localhost:%d", cfg.Server.Port)

	sv.GracefulShutdown(ch)
}
