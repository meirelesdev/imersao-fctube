package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/streadway/amqp"
)

type RabbitClientInterface interface {
	ConsumeMessages(exchange, routingKey, queueName string) (<-chan amqp.Delivery, error)
	PublishMessage(exchange, routingKey, queueName string, message []byte) error
	Close() error
	IsClosed() bool
}

type RabbitClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

func newConnection(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open a channel: %v", err)
	}
	return conn, ch, nil
}

func NewRabbitClient(ctx context.Context, connectionURL string) (*RabbitClient, error) {
	conn, ch, err := newConnection(connectionURL)
	if err != nil {
		return nil, err
	}
	return &RabbitClient{conn, ch, connectionURL}, nil
}

func (c *RabbitClient) ConsumeMessages(exchange, routingKey, queueName string) (<-chan amqp.Delivery, error) {
	err := c.channel.ExchangeDeclare(
		exchange,
		"direct",
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange %s: %v", exchange, err)
	}

	queue, err := c.channel.QueueDeclare(
		queueName,
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue %s: %v", queueName, err)
	}

	err = c.channel.QueueBind(
		queue.Name,
		routingKey,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue %s to exchange %s with routing key %s: %v", queueName, exchange, routingKey, err)
	}

	msgs, err := c.channel.Consume(
		queueName,
		"conversor_go",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages from queue %s: %v", queueName, err)
	}
	return msgs, nil
}

func (c *RabbitClient) PublishMessage(exchange, routingKey, queueName string, message []byte) error {
	err := c.channel.ExchangeDeclare(
		exchange,
		"direct",
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange %s: %v", exchange, err)
	}

	queue, err := c.channel.QueueDeclare(
		queueName,
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %v", queueName, err)
	}

	err = c.channel.QueueBind(
		queue.Name,
		routingKey,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue %s to exchange %s with routing key %s: %v", queueName, exchange, routingKey, err)
	}

	err = c.channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message to queue %s: %v", queueName, err)
	}
	return nil
}

func (c *RabbitClient) Close() {
	c.channel.Close()
	c.conn.Close()
}

func (client *RabbitClient) Reconnect(ctx context.Context) error {
	var err error
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled while trying to reconnect")
		default:
			slog.Info("Attempting to reconnect to RabbitMQ...")
			client.conn, client.channel, err = newConnection(client.url)
			if err == nil {
				slog.Info("Reconnected to RabbitMQ successfully")
				return nil
			}
			slog.Error("Failed to reconnect to RabbitMQ", slog.String("error", err.Error()))
			time.Sleep(5 * time.Second)
		}
	}
}
