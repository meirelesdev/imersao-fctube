package rabbitmq

import (
  "github.com/streadway/amqp"
	"fmt"
)
type RabbitClient struct {
	conn *amqp.Connection
	channel *amqp.Channel
	url string
}

func newConnection(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
  if err != nil {
    return nil, nil, fmt.Errorf("Failed to connect to RabbitMQ: %v", err)
  }

  ch, err := conn.Channel()
  if err != nil {
    return nil, nil, fmt.Errorf("Failed to open a channel: %v", err)
  }
  return conn, ch, nil
}

func NewRabbitClient(connectionURL string) (*RabbitClient, error) {
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
    return nil, fmt.Errorf("Failed to declare exchange %s: %v", exchange, err)
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
    return nil, fmt.Errorf("Failed to declare queue %s: %v", queueName, err)
  }

	err = c.channel.QueueBind(
    queue.Name,
    routingKey,
    exchange,
    false,
		nil,
	)
	if err!= nil {
    return nil, fmt.Errorf("Failed to bind queue %s to exchange %s with routing key %s: %v", queueName, exchange, routingKey, err)
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
    return nil, fmt.Errorf("Failed to consume messages from queue %s: %v", queueName, err)
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
    return fmt.Errorf("Failed to declare exchange %s: %v", exchange, err)
  }

  queue, err := c.channel.QueueDeclare(
    queueName,
    true,
    true,
    false,
    false,
    nil,
  )
  if err!= nil {
    return fmt.Errorf("Failed to declare queue %s: %v", queueName, err)
  }

  err = c.channel.QueueBind(
    queue.Name,
    routingKey,
    exchange,
    false,
    nil,
  )
  if err!= nil {
    return fmt.Errorf("Failed to bind queue %s to exchange %s with routing key %s: %v", queueName, exchange, routingKey, err)
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
  if err!= nil {
    return fmt.Errorf("Failed to publish message to queue %s: %v", queueName, err)
  }
  return nil
}

func (c *RabbitClient) Close() {
  c.channel.Close()
  c.conn.Close()
}