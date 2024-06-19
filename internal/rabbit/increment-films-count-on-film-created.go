package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/wagslane/go-rabbitmq"

	"github.com/manuhdez/films-api-test/internal/service"
)

type IncrementCountOnFilmCreatedConsumer struct {
	consumer    *rabbitmq.Consumer
	incrementer service.UserFilmsCounterIncrementer
}

func NewIncrementCountOnFilmCreatedConsumer(conn *rabbitmq.Conn,
	incrementer service.UserFilmsCounterIncrementer) (IncrementCountOnFilmCreatedConsumer, error) {
	consumer, err := rabbitmq.NewConsumer(
		conn,
		"events.increase_films_count_on_film_created",
		rabbitmq.WithConsumerOptionsLogging,
		rabbitmq.WithConsumerOptionsQueueDurable,
		rabbitmq.WithConsumerOptionsExchangeName(ExchangeName),
		rabbitmq.WithConsumerOptionsExchangeKind(ExchangeKind),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
		rabbitmq.WithConsumerOptionsExchangeDurable,
		rabbitmq.WithConsumerOptionsBinding(rabbitmq.Binding{
			RoutingKey:     "api.films.created",
			BindingOptions: rabbitmq.BindingOptions{Declare: true},
		}),
	)
	if err != nil {
		return IncrementCountOnFilmCreatedConsumer{}, err
	}

	return IncrementCountOnFilmCreatedConsumer{
		consumer:    consumer,
		incrementer: incrementer,
	}, nil
}

func (c IncrementCountOnFilmCreatedConsumer) Consume(errChan chan error) {
	err := c.consumer.Run(func(d rabbitmq.Delivery) (action rabbitmq.Action) {
		var event service.FilmCreatedEvent
		err := json.Unmarshal(d.Body, &event)
		if err != nil {
			log.Println("error unmarshalling increment-films-count-on-film-created:", err)
			return
		}

		log.Printf("event received: %+v", event)

		user := uuid.MustParse(event.CreatedBy)
		err = c.incrementer.Increment(context.Background(), user)
		if err != nil {
			log.Printf("error incrementing film count for user %s: %e", event.CreatedBy, err)
			return
		}

		fmt.Println("film created event received", d.RoutingKey)
		return
	})
	if err != nil {
		fmt.Printf("error consuming film created event: %v", err)
		errChan <- err
	}
}

func (c IncrementCountOnFilmCreatedConsumer) Close() {
	c.consumer.Close()
}
