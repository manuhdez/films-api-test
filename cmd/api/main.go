package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/manuhdez/films-api-test/database"
	"github.com/manuhdez/films-api-test/internal/http/server"
	"github.com/manuhdez/films-api-test/internal/rabbit"
)

func main() {
	db, dbErr := database.NewPostgresConnection()
	if dbErr != nil {
		panic(dbErr)
	}

	eventBus := rabbit.NewRabbitEventBus(db)
	srv := server.New(db, eventBus)

	errChan := make(chan error)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// handle shutdown signal
	go func() {
		select {
		case <-errChan:
			{
				for _, consumer := range eventBus.Consumers() {
					consumer.Close()
				}
			}
		case <-sigChan:
			{
				eventBus.CloseConnection()
				srv.Stop()
			}
		}
	}()

	// start event consumers
	for _, consumer := range eventBus.Consumers() {
		fmt.Printf("starting consumer %T\n", consumer)
		go consumer.Consume(errChan)
	}

	log.Print("starting server...")
	srvErr := srv.Start()
	if srvErr != nil {
		log.Fatal(srvErr)
	}
}
