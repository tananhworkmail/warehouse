package realtime

import "sync"

type Hub struct {
	clients map[chan []byte]bool
	mu      sync.Mutex
}

var AlertHub = &Hub{
	clients: make(map[chan []byte]bool),
}

// gửi message cho tất cả client
func (h *Hub) Broadcast(msg []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for ch := range h.clients {
		select {
		case ch <- msg:
		default:
			// tránh bị block nếu client chết
		}
	}
}
func (h *Hub) AddClient(ch chan []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[ch] = true
}

func (h *Hub) RemoveClient(ch chan []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, ch)
	close(ch)
}