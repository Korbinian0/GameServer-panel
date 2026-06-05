package rest

import (
	"encoding/json"
	"net/http"

	"github.com/Korbinian0/GameServer-panel/backend/internal/adapters/auth"
	websocketadapter "github.com/Korbinian0/GameServer-panel/backend/internal/adapters/websocket"
	"github.com/Korbinian0/GameServer-panel/backend/internal/app"
	"github.com/Korbinian0/GameServer-panel/backend/internal/domain"
	"github.com/go-chi/chi"
)

type server struct {
	appService *app.GatewayService
	auth       auth.JWTService
	hub        *websocketadapter.Hub
}

func RegisterRoutes(r *chi.Mux, appService *app.GatewayService, authService auth.JWTService) {
	s := &server{appService: appService, auth: authService}

	r.Post("/api/login", s.login)
	r.Get("/api/health", s.health)
}

func RegisterProtectedRoutes(r chi.Router, appService *app.GatewayService) {
	s := &server{appService: appService}
	r.Post("/users", s.registerUser)
	r.Get("/roles", s.getRoles)
}

func (s *server) login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	token, err := s.appService.Authenticate(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (s *server) registerUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	err := s.appService.RegisterUser(r.Context(), domain.User{ID: req.ID, Email: req.Email, Password: req.Password, Role: req.Role})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (s *server) getRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := s.appService.GetRoles(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(roles)
}

func (s *server) health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
