package sse

import (
	"sync"

	"github.com/google/uuid"

	"github.com/lucasbonna/contafacil_api/internal/schemas"
)

type SSEManager struct {
	clients map[uuid.UUID]chan schemas.SSEMessage
	mu      sync.Mutex
}

func NewSSEManager() *SSEManager {
	return &SSEManager{
		clients: make(map[uuid.UUID]chan schemas.SSEMessage),
	}
}

func (m *SSEManager) AddClient(userId uuid.UUID, clientChan chan schemas.SSEMessage) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[userId] = clientChan
}

func (m *SSEManager) RemoveClient(userId uuid.UUID) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if clientChan, ok := m.clients[userId]; ok {
		close(clientChan)
		delete(m.clients, userId)
	}
}

func (m *SSEManager) SendToClient(userId uuid.UUID, message schemas.SSEMessage) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	clientChan, ok := m.clients[userId]
	if !ok {
		return false
	}

	// BLOCKING send. Will wait until the SSE route receives from the channel.
	clientChan <- message
	return true
}
