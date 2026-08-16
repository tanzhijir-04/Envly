package events

import (
	"sync"
	"time"
)

type Event struct {
	Type       string         `json:"type"` // running | success | failed | run_done
	RunID      string         `json:"run_id"`
	ToolID     string         `json:"tool_id,omitempty"`
	Status     string         `json:"status,omitempty"`
	MessageKey string         `json:"message_key,omitempty"`
	Params     map[string]any `json:"params,omitempty"`
	TS         int64          `json:"ts"`
}

type Hub struct {
	mu           sync.Mutex
	subs         map[chan Event]struct{}
	history      []Event
	historyLimit int
}

func NewHub() *Hub {
	return &Hub{subs: map[chan Event]struct{}{}, historyLimit: 200}
}

func (h *Hub) Subscribe() (<-chan Event, func()) {
	ch := make(chan Event, 32)
	h.mu.Lock()
	for _, e := range h.history {
		select {
		case ch <- e:
		default:
		}
	}
	h.subs[ch] = struct{}{}
	h.mu.Unlock()
	return ch, func() {
		h.mu.Lock()
		delete(h.subs, ch)
		close(ch)
		h.mu.Unlock()
	}
}

func (h *Hub) Publish(e Event) {
	e.TS = time.Now().UnixMilli()
	h.mu.Lock()
	h.history = append(h.history, e)
	if len(h.history) > h.historyLimit {
		h.history = h.history[len(h.history)-h.historyLimit:]
	}
	for ch := range h.subs {
		select {
		case ch <- e:
		default:
		}
	}
	h.mu.Unlock()
}
