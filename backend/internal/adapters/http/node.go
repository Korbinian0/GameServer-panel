package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Korbinian0/GameServer-panel/backend/internal/adapters/websocket"
	"github.com/Korbinian0/GameServer-panel/backend/internal/app"
	"github.com/Korbinian0/GameServer-panel/backend/internal/domain"
	"github.com/go-chi/chi"
)

type nodeRegistrationRequest struct {
	NodeID       string   `json:"nodeId"`
	Platform     string   `json:"platform"`
	Hostname     string   `json:"hostname"`
	IPAddress    string   `json:"ipAddress"`
	Capabilities []string `json:"capabilities"`
	Version      string   `json:"version"`
}

type nodeEvent struct {
	Type      string      `json:"type"`
	Node      domain.Node `json:"node"`
	Timestamp string      `json:"timestamp"`
}

func RegisterNodeRoutes(r chi.Router, appService *app.GatewayService, hub *websocket.Hub) {
	s := &server{appService: appService, hub: hub}
	r.Get("/nodes", s.listNodes)
	r.Post("/nodes", s.registerNode)
	r.Post("/nodes/{nodeId}/heartbeat", s.heartbeatNode)
	r.Get("/nodes/{nodeId}", s.getNode)
}

func (s *server) registerNode(w http.ResponseWriter, r *http.Request) {
	var req nodeRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	node := domain.Node{
		ID:           req.NodeID,
		Platform:     req.Platform,
		Hostname:     req.Hostname,
		IPAddress:    req.IPAddress,
		Capabilities: req.Capabilities,
		Version:      req.Version,
	}
	if err := s.appService.RegisterNode(r.Context(), node); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.broadcastNodeEvent("node.registered", node)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (s *server) heartbeatNode(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "nodeId")
	if err := s.appService.RecordHeartbeat(r.Context(), nodeID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	node, err := s.appService.GetNode(r.Context(), nodeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	s.broadcastNodeEvent("node.heartbeat", node)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (s *server) getNode(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "nodeId")
	node, err := s.appService.GetNode(r.Context(), nodeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(node)
}

func (s *server) listNodes(w http.ResponseWriter, r *http.Request) {
	nodes, err := s.appService.ListNodes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(nodes)
}

func (s *server) broadcastNodeEvent(eventType string, node domain.Node) {
	if s.hub == nil {
		return
	}
	message := nodeEvent{
		Type:      eventType,
		Node:      node,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	payload, err := json.Marshal(message)
	if err != nil {
		return
	}
	s.hub.BroadcastMessage(payload)
}
