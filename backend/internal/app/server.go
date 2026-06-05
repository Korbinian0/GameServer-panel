package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Korbinian0/GameServer-panel/backend/internal/adapters/auth"
	"github.com/Korbinian0/GameServer-panel/backend/internal/adapters/grpc"
	"github.com/Korbinian0/GameServer-panel/backend/internal/domain"
	"github.com/Korbinian0/GameServer-panel/backend/internal/ports"
)

type GatewayService struct {
	repo        ports.Repository
	authService auth.JWTService
	registry    *grpcclient.BackendRegistry
}

func NewGatewayService(repo ports.Repository, authService auth.JWTService, registry *grpcclient.BackendRegistry) *GatewayService {
	return &GatewayService{repo: repo, authService: authService, registry: registry}
}

func (s *GatewayService) Authenticate(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user.Password != password {
		return "", errors.New("invalid credentials")
	}
	token, err := s.authService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *GatewayService) RegisterUser(ctx context.Context, user domain.User) error {
	user.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	return s.repo.CreateUser(ctx, user)
}

func (s *GatewayService) RegisterNode(ctx context.Context, node domain.Node) error {
	node.LastSeen = time.Now().UTC().Format(time.RFC3339)
	if err := s.repo.CreateNode(ctx, node); err != nil {
		return err
	}
	if s.registry != nil {
		target := fmt.Sprintf("%s:50051", node.IPAddress)
		_ = s.registry.AddBackend(node.ID, target)
	}
	return nil
}

func (s *GatewayService) GetNode(ctx context.Context, nodeID string) (domain.Node, error) {
	return s.repo.GetNode(ctx, nodeID)
}

func (s *GatewayService) ListNodes(ctx context.Context) ([]domain.Node, error) {
	return s.repo.ListNodes(ctx)
}

func (s *GatewayService) RecordHeartbeat(ctx context.Context, nodeID string) error {
	node, err := s.repo.GetNode(ctx, nodeID)
	if err != nil {
		return err
	}
	node.LastSeen = time.Now().UTC().Format(time.RFC3339)
	return s.repo.CreateNode(ctx, node)
}

func (s *GatewayService) GetRoles(ctx context.Context) ([]domain.Role, error) {
	return s.repo.ListRoles(ctx)
}
