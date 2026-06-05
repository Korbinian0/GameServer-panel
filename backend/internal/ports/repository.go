package ports

import (
	"context"
	"github.com/Korbinian0/GameServer-panel/backend/internal/domain"
)

type Repository interface {
	CreateUser(ctx context.Context, user domain.User) error
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	ListRoles(ctx context.Context) ([]domain.Role, error)
	CreateNode(ctx context.Context, node domain.Node) error
	GetNode(ctx context.Context, nodeID string) (domain.Node, error)
	ListNodes(ctx context.Context) ([]domain.Node, error)
}
