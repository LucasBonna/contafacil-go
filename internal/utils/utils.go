package utils

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lucasbonna/contafacil_api/internal/database"
	"github.com/lucasbonna/contafacil_api/internal/rabbitmq"
)

type MessageMetadata struct {
	RetryCount    int    `json:"retryCount"`
	RetryAt       int64  `json:"retryAt"`
	OriginalQueue string `json:"originalQueue"`
}

type TaskWithMetadata struct {
	Id       uuid.UUID
	Type     TaskType         `json:"type"`
	Payload  interface{}      `json:"payload"`
	Metadata *MessageMetadata `json:"messageMetadata,omitempty"`
}

type QueueHelper struct {
	rabbit *rabbitmq.RabbitMQ
}

func NewQueueHelper(rabbit *rabbitmq.RabbitMQ) *QueueHelper {
	return &QueueHelper{rabbit: rabbit}
}

func (qh *QueueHelper) EnqueueTask(taskType TaskType, payload interface{}) error {
	task := TaskWithMetadata{
		Id:      uuid.New(),
		Type:    taskType,
		Payload: payload,
	}

	body, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return qh.rabbit.Publish(string(taskType)+"Queue", body)
}

func GetUser(c *gin.Context) *database.User {
	userStored, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return nil
	}

	user, ok := userStored.(*database.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return nil
	}

	return user
}

func GetClient(c *gin.Context) *database.Client {
	clientStored, exists := c.Get("client")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return nil
	}

	client, ok := clientStored.(*database.Client)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return nil
	}

	return client
}
