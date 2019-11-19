package rabbit

import (
	"context"
	"encoding/json"
	"github.com/koind/banner-rotation/api/internal/domain/repository"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Publisher interface
type PublisherInterface interface {
	// Send message to the queue
	Publish(ctx context.Context, statistic repository.Statistic) error
}

// Publisher
type Publisher struct {
	conn         amqp.Connection
	exchangeName string
	queueName    string
}

// Returns new publisher
func NewPublisher(conn amqp.Connection, exchangeName string, queueName string) *Publisher {
	return &Publisher{
		conn:         conn,
		exchangeName: exchangeName,
		queueName:    queueName,
	}
}

// Send message to the queue
func (p *Publisher) Publish(ctx context.Context, statistic repository.Statistic) error {
	if ctx.Err() == context.Canceled {
		return errors.New("sending statistics was aborted due to context cancellation")
	}

	ch, err := p.conn.Channel()
	if err != nil {
		return errors.Wrap(err, "unable to connect to channel")
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		p.queueName,
		true,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "failed to declare queue")
	}

	err = ch.QueueBind(
		q.Name,
		"",
		p.exchangeName,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "failed to bind the queue to the exchanger")
	}

	data, err := json.Marshal(statistic)
	if err != nil {
		return errors.Wrap(err, "failed to encode in json")
	}

	err = ch.Publish(
		p.exchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
	if err != nil {
		return errors.Wrap(err, "unable to publish message")
	}

	return nil
}
