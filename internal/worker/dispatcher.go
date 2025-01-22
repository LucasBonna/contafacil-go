package worker

import (
	"errors"
	"sync"

	"github.com/lucasbonna/contafacil_api/internal/utils"
)

// Dispatcher gerencia os handlers para diferentes tipos de tarefas.
type Dispatcher struct {
	handlers map[string]TaskHandler
	mu       sync.RWMutex
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string]TaskHandler),
	}
}

func (d *Dispatcher) RegisterHandler(taskType string, handler TaskHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[taskType] = handler
}

func (d *Dispatcher) GetHandler(taskType utils.TaskType) (TaskHandler, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	handler, exists := d.handlers[string(taskType)]
	if !exists {
		return nil, errors.New("no handler registered for task type: " + string(taskType))
	}
	return handler, nil
}

