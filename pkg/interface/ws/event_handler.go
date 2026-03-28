// @AI_GENERATED
package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/make-bin/groundhog/pkg/application/eventbus"
	"github.com/make-bin/groundhog/pkg/infrastructure/service"
	"github.com/make-bin/groundhog/pkg/interface/http/response"
	"github.com/make-bin/groundhog/pkg/utils/bcode"
)

// WSMessage is the message format sent to WebSocket clients.
type WSMessage struct {
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

// wsClient represents a connected WebSocket client using SSE-style long polling.
// TODO: Replace with gorilla/websocket for full WebSocket support in production.
type wsClient struct {
	mu     sync.Mutex
	ch     chan WSMessage
	closed bool
}

func newWSClient() *wsClient {
	return &wsClient{ch: make(chan WSMessage, 64)}
}

func (c *wsClient) send(msg WSMessage) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.closed {
		select {
		case c.ch <- msg:
		default:
			// drop if buffer full
		}
	}
}

func (c *wsClient) close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.closed {
		c.closed = true
		close(c.ch)
	}
}

// WSEventHandler handles WebSocket connections for real-time event streaming.
type WSEventHandler interface {
	Handle(c *gin.Context)
}

type wsEventHandler struct {
	EventBus   eventbus.EventBus   `inject:""`
	JWTService *service.JWTService `inject:""`

	mu      sync.RWMutex
	clients map[string]*wsClient
}

// NewWSEventHandler creates a new WSEventHandler.
func NewWSEventHandler() WSEventHandler {
	h := &wsEventHandler{
		clients: make(map[string]*wsClient),
	}
	return h
}

// Handle upgrades the HTTP connection to a streaming endpoint and pushes events.
// JWT is validated from the "token" query param or Authorization header.
//
// TODO: Replace the SSE-style streaming stub with gorilla/websocket for full
// WebSocket support (RFC 6455) in production.
func (h *wsEventHandler) Handle(c *gin.Context) {
	// --- JWT validation ---
	token := c.Query("token")
	if token == "" {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		}
	}
	if token == "" {
		response.Error(c, bcode.ErrUnauthorized, fmt.Errorf("missing token"))
		return
	}
	claims, err := h.JWTService.ValidateToken(token)
	if err != nil {
		response.Error(c, bcode.ErrUnauthorized, err)
		return
	}

	principalID := claims.PrincipalID

	// --- Register client ---
	client := newWSClient()
	h.mu.Lock()
	h.clients[principalID] = client
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, principalID)
		h.mu.Unlock()
		client.close()
	}()

	// Subscribe to EventBus events and forward to this client.
	eventTypes := []string{
		"agent.turn.started",
		"agent.turn.streaming",
		"agent.turn.completed",
		"agent.tool.executing",
		"agent.tool.completed",
		"agent.tool.approval",
	}
	for _, et := range eventTypes {
		eventType := et // capture loop var
		h.EventBus.Subscribe(eventType, func(event interface{}) {
			payload, _ := json.Marshal(event)
			client.send(WSMessage{
				Type:      eventType,
				Payload:   payload,
				Timestamp: time.Now().UTC(),
			})
		})
	}

	// --- Stream events using chunked HTTP (SSE-style stub) ---
	// In production this would be replaced by a proper WebSocket upgrade.
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)

	// Heartbeat ticker for keep-alive.
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	flusher, canFlush := c.Writer.(http.Flusher)

	for {
		select {
		case msg, ok := <-client.ch:
			if !ok {
				return
			}
			data, _ := json.Marshal(msg)
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			if canFlush {
				flusher.Flush()
			}
		case <-ticker.C:
			// heartbeat ping
			fmt.Fprintf(c.Writer, ": ping\n\n")
			if canFlush {
				flusher.Flush()
			}
		case <-c.Request.Context().Done():
			return
		}
	}
}

// @AI_GENERATED: end
