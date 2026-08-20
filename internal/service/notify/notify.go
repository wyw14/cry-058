package notify

import "sync"

type Message struct {
	UserID string
	Title  string
	Body   string
}
type Hub struct {
	mu    sync.Mutex
	items []Message
}

func New() *Hub               { return &Hub{items: []Message{}} }
func (h *Hub) Push(m Message) { h.mu.Lock(); defer h.mu.Unlock(); h.items = append(h.items, m) }
func (h *Hub) List(user string) []Message {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := []Message{}
	for _, m := range h.items {
		if user == "" || m.UserID == user {
			out = append(out, m)
		}
	}
	return out
}
func (h *Hub) Clear() { h.mu.Lock(); defer h.mu.Unlock(); h.items = nil }
