package producerconsumer

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	ch   *amqp091.Channel
	conn *amqp091.Connection
}

func NewProducer(url string) (*Producer, error) {
	var err error
	consume := Producer{}
	consume.conn, err = amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	consume.ch, err = consume.conn.Channel()
	if err != nil {
		return nil, err
	}
	return &consume, nil
}

func (p *Producer) PublishCtx(ctx context.Context, queueName string, body []byte) error {
	_, err := p.ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	return p.ch.PublishWithContext(ctx, "", queueName, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func (p *Producer) Close() {
	p.ch.Close()
	p.conn.Close()
}
