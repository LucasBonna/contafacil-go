package rabbitmq

import (
	"fmt"
	"log"
	"time"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(dsn string) (*RabbitMQ, error) {
	conn, err := amqp.DialConfig(dsn, amqp.Config{
		Heartbeat: 10 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	rabbit := &RabbitMQ{
		conn:    conn,
		channel: channel,
	}

	if err := SetupQueues(rabbit); err != nil {
		rabbit.Close()
		return nil, err
	}

	go func() {
		closeChan := make(chan *amqp.Error)
		rabbit.conn.NotifyClose(closeChan)
		
		for {
			err := <-closeChan
			if err != nil {
				log.Printf("Connection closed, attempting to reconnect: %v", err)
				for {
					if err := rabbit.Reconnect(dsn); err == nil {
						break
					}
					time.Sleep(5 * time.Second)
				}
			}
		}
	}()

	return rabbit, nil
}

func (r *RabbitMQ) Consume(queueName, consumerName string) (<-chan amqp.Delivery, error) {
	if r.channel == nil || r.conn.IsClosed() {
		return nil, fmt.Errorf("channel/connection is not open")
	}

	return r.channel.Consume(
		queueName,
		consumerName,
		false,
		false,
		false,
		false,
		nil,
	)
}

func (r *RabbitMQ) Publish(queueName string, body []byte) error {
	if r.channel == nil || r.conn.IsClosed() {
		return fmt.Errorf("channel/connection is not open")
	}

	return r.channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

func (r *RabbitMQ) DeclareQueue(queueName string) error {
	if r.channel == nil || r.conn.IsClosed() {
		return fmt.Errorf("channel/connection is not open")
	}

	log.Printf("creating queue: %v", queueName)
	_, err := r.channel.QueueDeclare(
		queueName,
		true,  // Durable
		false, // Auto-delete
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)

	if err != nil {
		return fmt.Errorf("error declaring queue %s: %v", queueName, err)
	}

	log.Printf("queue %v created", queueName)
	return nil
}

func SetupQueues(rabbit *RabbitMQ) error {
	queues := []string{"UploadFileQueue", "DownloadFileQueue", "IssueGNREQueue", "ExceptionQueue", "RetryQueue"}
	for _, queueName := range queues {
		if err := rabbit.DeclareQueue(queueName); err != nil {
			return err
		}
	}
	return nil
}

func (r *RabbitMQ) IsConnectionOpen() bool {
	return r.conn != nil && !r.conn.IsClosed()
}

func (r *RabbitMQ) Reconnect(dsn string) error {
	log.Println("Attempting to reconnect to RabbitMQ...")
	r.Close()

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return fmt.Errorf("failed to reconnect: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to create channel: %v", err)
	}

	r.conn = conn
	r.channel = channel

	if err := SetupQueues(r); err != nil {
		r.Close()
		return fmt.Errorf("failed to setup queues: %v", err)
	}

	log.Println("Successfully reconnected to RabbitMQ")
	return nil
}
