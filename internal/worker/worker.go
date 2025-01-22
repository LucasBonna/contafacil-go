package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"

	"github.com/lucasbonna/contafacil_api/internal/config"
	"github.com/lucasbonna/contafacil_api/internal/rabbitmq"
	"github.com/lucasbonna/contafacil_api/internal/utils"
)

type Task struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func StartWorkers(rabbit *rabbitmq.RabbitMQ, dispatcher *Dispatcher) {
	queues := []string{"UploadFileQueue", "DownloadFileQueue", "IssueGNREQueue"}

	for _, queue := range queues {
		go func(qName string) {
			for {
				if !rabbit.IsConnectionOpen() {
					log.Printf("Connection lost for queue %s. Reconnecting...", qName)
					if err := rabbit.Reconnect(config.Env.RabbitMQUrl); err != nil {
						time.Sleep(5 * time.Second)
						continue
					}
				}

				deliveries, err := rabbit.Consume(qName, "")
				if err != nil {
					log.Printf("Error consuming queue %s: %v", qName, err)
					time.Sleep(5 * time.Second)
					continue
				}

				for msg := range deliveries {
					processMessage(rabbit, dispatcher, msg)
				}
			}
		}(queue)
	}
}

func processMessage(rabbit *rabbitmq.RabbitMQ, dispatcher *Dispatcher, msg amqp.Delivery) {
	var task utils.TaskWithMetadata
	if err := json.Unmarshal(msg.Body, &task); err != nil {
		log.Printf("Invalid message: %v", err)
		moveToExceptionQueue(rabbit, msg.Body, nil)
		msg.Reject(false)
		return
	}

	handler, err := dispatcher.GetHandler(task.Type)
	if err != nil {
		log.Printf("No handler for task type %s: %v", task.Type, err)
		moveToExceptionQueue(rabbit, msg.Body, task.Metadata)
		msg.Reject(false)
		return
	}

	ctx := context.Background()
	if err := handler.Handle(ctx, task.Payload); err != nil {
		handleFailure(rabbit, task, msg)
		return
	}

	msg.Ack(false)
	log.Printf("task %s processed successfully", task.Type)
}

func handleFailure(rabbitmq *rabbitmq.RabbitMQ, task utils.TaskWithMetadata, msg amqp.Delivery) {
	if task.Metadata == nil {
		task.Metadata = &utils.MessageMetadata{
			RetryCount:    0,
			OriginalQueue: string(task.Type) + "Queue",
		}
	}

	log.Printf("task %v failed, retry count: %v/3", msg.Type, task.Metadata.RetryCount)

	if task.Metadata.RetryCount >= 3 {
		moveToExceptionQueue(rabbitmq, msg.Body, task.Metadata)
		msg.Reject(false)
		return
	}

	task.Metadata.RetryCount++
	task.Metadata.RetryAt = time.Now().Add(getBackoffDuration(task.Metadata.RetryCount)).Unix()

	body, _ := json.Marshal(task)
	rabbitmq.Publish("RetryQueue", body)
	msg.Ack(false)
}

func getBackoffDuration(retryCount int) time.Duration {
	switch retryCount {
	case 1:
		return time.Minute
	case 2:
		return 3 * time.Minute
	case 3:
		return 5 * time.Minute
	default:
		return time.Minute
	}
}

func moveToExceptionQueue(rabbit *rabbitmq.RabbitMQ, body []byte, metadata *utils.MessageMetadata) {
	if metadata == nil {
		metadata = &utils.MessageMetadata{RetryCount: 3}
	}

	task := utils.TaskWithMetadata{
		Metadata: metadata,
	}

	json.Unmarshal(body, &task)
	newBody, _ := json.Marshal(task)
	rabbit.Publish("ExceptionQueue", newBody)
}

func StartRetryWorker(rabbit *rabbitmq.RabbitMQ) {
	go func() {
		for {
			deliveries, err := rabbit.Consume("RetryQueue", "cause error")
			if err != nil {
				time.Sleep(5 * time.Second)
				continue
			}

			for msg := range deliveries {
				var task utils.TaskWithMetadata
				if err := json.Unmarshal(msg.Body, &task); err != nil {
					msg.Reject(false)
					continue
				}

				if time.Now().Unix() < task.Metadata.RetryAt {
					// Ainda nao e hora de retentar
					rabbit.Publish("RetryQueue", msg.Body)
					msg.Ack(false)
					continue
				}

				rabbit.Publish(task.Metadata.OriginalQueue, msg.Body)
				msg.Ack(false)
			}
		}
	}()
}
