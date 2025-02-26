package sse

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Message struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type Manager struct {
	redis   *redis.Client
	clients map[uuid.UUID]chan Message
	mu      sync.RWMutex
}

func NewManager(redis *redis.Client) *Manager {
	return &Manager{
		redis:   redis,
		clients: make(map[uuid.UUID]chan Message),
	}
}

func (m *Manager) Register(userID uuid.UUID) chan Message {
	log.Println("registering: ", userID)
	m.mu.Lock()
	defer m.mu.Unlock()

	ch := make(chan Message, 10)
	m.clients[userID] = ch

	// Atualizar presença no Redis
	m.redis.SetEX(context.Background(),
		fmt.Sprintf("sse:presence:%s", userID),
		"online",
		30*time.Second,
	)

	return ch
}

func (m *Manager) Unregister(userID uuid.UUID) {
	log.Println("unregistering: ", userID)
	m.mu.Lock()
	defer m.mu.Unlock()

	if ch, ok := m.clients[userID]; ok {
		close(ch)
		delete(m.clients, userID)
	}

	m.redis.Del(context.Background(),
		fmt.Sprintf("sse:presence:%s", userID),
	)
}

func (m *Manager) IsConnected(userID uuid.UUID) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	log.Println("checking...")

	m.ListClients()

	_, memOk := m.clients[userID]
	redisVal, redisErr := m.redis.Get(context.Background(),
		fmt.Sprintf("sse:presence:%s", userID)).Result()

	log.Printf("Connection check for %s - Memory: %t, Redis: %s",
		userID, memOk, redisVal)

	return memOk && redisErr == nil && redisVal == "online"
}

func (m *Manager) ListClients() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	log.Println("Active SSE connections:")
	for userID := range m.clients {
		log.Printf("- UserID: %s", userID)
	}
}

func (m *Manager) Send(userID uuid.UUID, msg Message) bool {
	log.Println("sending to:", userID)
	m.mu.RLock()
	defer m.mu.RUnlock()

	if ch, ok := m.clients[userID]; ok {
		select {
		case ch <- msg:
			return true
		default:
			return false
		}
	}
	return false
}
