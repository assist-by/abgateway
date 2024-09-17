package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	config "github.com/assist-by/abgateway/library"
	kafka "github.com/assist-by/abgateway/pkg/kafka"
	server "github.com/assist-by/abgateway/server"
	service "github.com/assist-by/abgateway/service"
)

func main() {
	cfg := config.Load()

	registrationWriter := kafka.NewWriter(cfg.KafkaBroker, cfg.RegistrationTopic)
	defer registrationWriter.Close()

	if err := service.RegisterService(registrationWriter, cfg.Host, cfg.Port); err != nil {
		log.Printf("Failed to register service: %v\n", err)
	}

	srv := server.NewServer(cfg)

	go func() {
		if err := srv.Run(); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	if err := srv.Shutdown(); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
