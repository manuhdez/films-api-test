package rabbit

import (
	"context"
	"log"

	"github.com/wagslane/go-rabbitmq"

	"github.com/manuhdez/films-api-test/internal/domain"
)

type ApiPublisher struct {
	publisher *rabbitmq.Publisher
}

func NewApiPublisher(conn *rabbitmq.Conn) (ApiPublisher, error) {
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(ExchangeName),
		rabbitmq.WithPublisherOptionsExchangeKind(ExchangeKind),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
		rabbitmq.WithPublisherOptionsExchangeDurable,
	)

	if err != nil {
		return ApiPublisher{}, err
	}

	return ApiPublisher{publisher: publisher}, nil
}

func (p ApiPublisher) Publish(c context.Context, e domain.Event) error {
	err := p.publisher.PublishWithContext(c,
		e.Data(),
		[]string{e.Key()},
		rabbitmq.WithPublishOptionsExchange(ExchangeName),
		rabbitmq.WithPublishOptionsContentType(ExchangeContentType),
	)
	if err != nil {
		log.Printf("failed to publish event: %v", err)
	}
	return err
}
