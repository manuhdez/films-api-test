package rabbit

import (
	"fmt"

	"github.com/wagslane/go-rabbitmq"
)

type IncrementCountOnFilmCreatedConsumer struct {
	consumer *rabbitmq.Consumer
}

func NewIncrementCountOnFilmCreatedConsumer(conn *rabbitmq.Conn) (IncrementCountOnFilmCreatedConsumer, error) {
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

	return IncrementCountOnFilmCreatedConsumer{consumer}, nil
}

func (c IncrementCountOnFilmCreatedConsumer) Consume(errChan chan error) {
	err := c.consumer.Run(func(d rabbitmq.Delivery) (action rabbitmq.Action) {
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
