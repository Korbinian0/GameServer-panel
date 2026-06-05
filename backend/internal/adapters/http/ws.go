package rest

import (
	"net/http"

	"github.com/Korbinian0/GameServer-panel/backend/internal/adapters/auth"
	websocketadapter "github.com/Korbinian0/GameServer-panel/backend/internal/adapters/websocket"
	"github.com/go-chi/chi"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func RegisterWebsocketRoutes(r *chi.Mux, hub *websocketadapter.Hub, authService auth.JWTService) {
	s := &server{auth: authService}
	r.With(s.auth.Middleware).Get("/ws/events", s.serveWebsocket(hub))
}

func (s *server) serveWebsocket(hub *websocketadapter.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "websocket upgrade failed", http.StatusBadRequest)
			return
		}

		client := &websocketadapter.Client{Conn: conn, Send: make(chan []byte, 256)}
		hub.RegisterClient(client)

		defer func() {
			hub.UnregisterClient(client)
			conn.Close()
		}()

		go func() {
			for msg := range client.Send {
				if err := conn.WriteMessage(ws.TextMessage, msg); err != nil {
					return
				}
			}
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return
			}
			hub.BroadcastMessage(message)
		}
	}
}
