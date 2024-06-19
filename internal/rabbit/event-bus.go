package rabbit

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/wagslane/go-rabbitmq"

	"github.com/manuhdez/films-api-test/internal/domain"
	"github.com/manuhdez/films-api-test/internal/infra"
	"github.com/manuhdez/films-api-test/internal/service"
)

type EventBus struct {
	connection *rabbitmq.Conn

	publisher Publisher
	consumers []Consumer
}

func NewRabbitEventBus(db *sql.DB) EventBus {
	conn, err := rabbitmq.NewConn("amqp://guest:guest@rabbit:5672/", rabbitmq.WithConnectionOptionsLogging)
	if err != nil {
		log.Panicf("Failed to connect to RabbitMQ: %v", err)
	}

	// instantiate use cases
	userFilmRepository := infra.NewUserFilmsPostgresRepository(db)
	counterCreator := service.NewUserFilmsCounterCreator(userFilmRepository)
	counterIncrementer := service.NewUserFilmsCounterIncrementer(userFilmRepository)

	publisher, err := setupPublisher(conn)
	if err != nil {
		log.Panicf("Failed to setup publisher: %v", err)
	}

	consumers, err := setupConsumers(conn, counterCreator, counterIncrementer)
	if err != nil {
		log.Panicf("Failed to setup consumers: %v", err)
	}

	return EventBus{
		connection: conn,
		publisher:  publisher,
		consumers:  consumers,
	}
}

func setupPublisher(conn *rabbitmq.Conn) (Publisher, error) {
	return NewApiPublisher(conn)
}

func setupConsumers(
	conn *rabbitmq.Conn,
	counterCreator service.UserFilmsCounterCreator,
	counterIncrementer service.UserFilmsCounterIncrementer,
) ([]Consumer, error) {
	var consumers []Consumer

	incrementCountOnFilmCreatedConsumer, err := NewIncrementCountOnFilmCreatedConsumer(conn, counterIncrementer)
	if err != nil {
		log.Printf("Failed to create increment count on film created consumer: %v", err)
	} else {
		consumers = append(consumers, incrementCountOnFilmCreatedConsumer)
	}

	createCounterOnUserCreated, err := NewCreateFilmsCounterOnUserCreated(conn, counterCreator)
	if err != nil {
		log.Printf("Failed to create films counter on user created consumer: %v", err)
	} else {
		consumers = append(consumers, createCounterOnUserCreated)
	}

	return consumers, nil
}

func (bus EventBus) Publish(ctx context.Context, event domain.Event) error {
	return bus.publisher.Publish(ctx, event)
}

func (bus EventBus) Consumers() []Consumer {
	return bus.consumers
}

func (bus EventBus) CloseConnection() {
	err := bus.connection.Close()
	if err != nil {
		fmt.Println("Error closing connection", err)
	}
}
