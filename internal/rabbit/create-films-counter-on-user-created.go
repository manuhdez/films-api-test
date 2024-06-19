package rabbit

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/wagslane/go-rabbitmq"

	"github.com/manuhdez/films-api-test/internal/domain/userfilms"
	"github.com/manuhdez/films-api-test/internal/service"
)

type CreateFilmsCounterOnUserCreated struct {
	consumer       *rabbitmq.Consumer
	counterCreator service.UserFilmsCounterCreator
}

func NewCreateFilmsCounterOnUserCreated(
	conn *rabbitmq.Conn,
	creator service.UserFilmsCounterCreator,
) (CreateFilmsCounterOnUserCreated, error) {
	consumer, err := rabbitmq.NewConsumer(
		conn,
		"events.create_films_counter_on_user_created",
		rabbitmq.WithConsumerOptionsLogging,
		rabbitmq.WithConsumerOptionsQueueDurable,
		rabbitmq.WithConsumerOptionsExchangeName(ExchangeName),
		rabbitmq.WithConsumerOptionsExchangeKind(ExchangeKind),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
		rabbitmq.WithConsumerOptionsExchangeDurable,
		rabbitmq.WithConsumerOptionsBinding(rabbitmq.Binding{
			RoutingKey:     "api.users.created",
			BindingOptions: rabbitmq.BindingOptions{Declare: true},
		}),
	)
	if err != nil {
		return CreateFilmsCounterOnUserCreated{}, err
	}

	return CreateFilmsCounterOnUserCreated{
		consumer:       consumer,
		counterCreator: creator,
	}, nil
}

func (c CreateFilmsCounterOnUserCreated) Consume(errChan chan error) {
	err := c.consumer.Run(func(d rabbitmq.Delivery) (action rabbitmq.Action) {
		var event service.UserCreatedEvent
		err := json.Unmarshal(d.Body, &event)
		if err != nil {
			log.Println("error unmarshalling message")
		}

		ctx := context.Background()
		id := uuid.MustParse(event.UserID)
		counter := userfilms.New(id, 0)
		err = c.counterCreator.Create(ctx, counter)
		if err != nil {
			log.Println("error creating films counter for user ", event.UserID)
		}

		return
	})
	if err != nil {
		log.Printf("error consuming user created event: %v", err)
		errChan <- err
	}
}

func (c CreateFilmsCounterOnUserCreated) Close() {
	c.consumer.Close()
}
