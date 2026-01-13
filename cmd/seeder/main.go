package main

import (
	"log"

	"github.com/InstaySystem/is_v2-be/internal/container"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/background/seeder"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctn := container.NewContainer(cfg)
	if err := ctn.InitSeed(); err != nil {
		log.Fatal(err)
	}
	defer ctn.Cleanup()

	sd := seeder.NewSeeder(cfg.SuperUser, ctn.Log, ctn.DB.Gorm, ctn.IDGen, ctn.UserRepo)
	if err = sd.Start(); err != nil {
		log.Fatal(err)
	}
}
