package shutdown

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type ShutdownFunc func(ctx context.Context) error

type Manager struct {
	mu       sync.Mutex
	handlers []ShutdownFunc
	timeout  time.Duration
}

func New(timeout time.Duration) *Manager {
	return &Manager{
		timeout: timeout,
	}
}

func (m *Manager) Register(fn ShutdownFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers = append(m.handlers, fn)
}

func (m *Manager) Wait() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Shutdown signal received, closing resources...")

	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	// Закрываем в обратном порядке регистрации (LIFO)
	m.mu.Lock()
	for i := len(m.handlers) - 1; i >= 0; i-- {
		if err := m.handlers[i](ctx); err != nil {
			log.Printf("Shutdown error: %v", err)
		}
	}
	m.mu.Unlock()

	log.Println("All resources closed. Bye!")
}
