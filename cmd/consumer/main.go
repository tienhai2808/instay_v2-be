package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/InstaySystem/is_v2-be/internal/container"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/background/consumer"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctn := container.NewContainer(cfg)
	if err := ctn.InitConsumer(); err != nil {
		log.Fatal(err)
	}
	defer ctn.Cleanup()

	csm := consumer.NewConsumer(ctn.Log, ctn.MQPro, ctn.SMTPPro)
	csm.Start()

	log.Println("Consumer is running")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Consumer stopped successfully")
}
