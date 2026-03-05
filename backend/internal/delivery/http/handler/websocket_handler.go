package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	jwtpkg "github.com/n2pluto/cinema-booking-system/pkg/jwt"
	ws "github.com/n2pluto/cinema-booking-system/internal/delivery/http/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // CORS handled at the gin middleware level
	},
}

type WebSocketHandler struct {
	hub *ws.Hub
}

func NewWebSocketHandler(hub *ws.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

// GET /api/ws/seats/:cinemaId?token=<jwt>
func (h *WebSocketHandler) HandleSeatUpdates(c *gin.Context) {
	cinemaID := c.Param("cinemaId")
	if cinemaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing cinemaId"})
		return
	}

	// Authenticate via query-string token (browsers can't set WS headers).
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}
	if _, err := jwtpkg.Parse(token); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &ws.Client{
		Hub:      h.hub,
		CinemaID: cinemaID,
		Conn:     conn,
		Send:     make(chan []byte, 256),
	}

	h.hub.Register <- client
	go client.WritePump()
	client.ReadPump()
}
