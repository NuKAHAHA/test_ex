package logger

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"sync"
	"time"
)

type Event struct {
	Time   time.Time      `json:"time"`
	Action string         `json:"action"`
	TaskID string         `json:"task_id,omitempty"`
	Meta   map[string]any `json:"meta,omitempty"`
}

type AsyncLogger struct {
	ch     chan Event
	done   chan struct{}
	wg     sync.WaitGroup
	closed bool
	mu     sync.Mutex
}

func NewAsyncLogger(out io.Writer, bufSize int) *AsyncLogger {
	l := &AsyncLogger{
		ch:   make(chan Event, bufSize),
		done: make(chan struct{}),
	}
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		enc := json.NewEncoder(out)
		for e := range l.ch {
			if err := enc.Encode(e); err != nil {
				log.Printf("Logger encode error: %v", err)
			}
		}
		close(l.done)
	}()
	return l
}

func (l *AsyncLogger) Log(e Event) {
	select {
	case l.ch <- e:
	default:
		log.Printf("Logger buffer full, dropping event: %v", e)
	}
}

func (l *AsyncLogger) Close(ctx context.Context) error {
	l.mu.Lock()
	if l.closed {
		l.mu.Unlock()
		return nil
	}
	l.closed = true
	close(l.ch)
	l.mu.Unlock()

	doneCh := make(chan struct{})
	go func() {
		l.wg.Wait()
		close(doneCh)
	}()

	select {
	case <-doneCh:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
