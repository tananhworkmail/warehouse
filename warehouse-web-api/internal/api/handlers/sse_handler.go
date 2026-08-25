package handlers

import (
	"fmt"
	"net/http"
	"time" // thêm import

	"web-api/internal/realtime"

	"github.com/gin-gonic/gin"
)

func SSEHandler(c *gin.Context) {
	w := c.Writer
	r := c.Request

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // chống buffer nginx

	flusher, ok := w.(http.Flusher)
	if !ok {
		c.String(http.StatusInternalServerError, "Streaming unsupported")
		return
	}

	clientChan := make(chan []byte, 10)

	realtime.AlertHub.AddClient(clientChan)
	defer realtime.AlertHub.RemoveClient(clientChan)

	ctx := r.Context()

	// Gửi sự kiện "connected" ngay khi client vừa kết nối
	fmt.Fprintf(w, "event: connected\ndata: {}\n\n")
	flusher.Flush()

	// Heartbeat mỗi 30 giây
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg := <-clientChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()

		case <-ticker.C:
			// Gửi comment để giữ kết nối (không ảnh hưởng tới client)
			fmt.Fprintf(w, ": heartbeat\n\n")
			flusher.Flush()

		case <-ctx.Done():
			fmt.Println("Client disconnected")
			return
		}
	}
}