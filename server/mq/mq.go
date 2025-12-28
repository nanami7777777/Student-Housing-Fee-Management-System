package mq

import (
	"context"
	"errors"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn *amqp.Connection
	ch   *amqp.Channel
	mu   sync.Mutex
)

func Init(url string) error {
	mu.Lock()
	defer mu.Unlock()
	if conn != nil {
		return nil
	}
	c, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	channel, err := c.Channel()
	if err != nil {
		c.Close()
		return err
	}
	conn = c
	ch = channel
	return nil
}

func Publish(queue string, body []byte) error {
	mu.Lock()
	defer mu.Unlock()
	if ch == nil {
		return errors.New("mq not initialized")
	}
	_, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return ch.PublishWithContext(ctx, "", queue, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
}

func Close() error {
	mu.Lock()
	defer mu.Unlock()
	if ch != nil {
		if err := ch.Close(); err != nil {
			return err
		}
		ch = nil
	}
	if conn != nil {
		if err := conn.Close(); err != nil {
			return err
		}
		conn = nil
	}
	return nil
}

