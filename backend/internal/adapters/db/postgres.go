package db

import (
	"context"

	"github.com/Korbinian0/GameServer-panel/backend/internal/domain"
	"github.com/Korbinian0/GameServer-panel/backend/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(ctx context.Context, dsn string) (ports.Repository, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{pool: pool}, nil
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user domain.User) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO users (id, email, password, role, created_at) VALUES ($1, $2, $3, $4, $5)`, user.ID, user.Email, user.Password, user.Role, user.CreatedAt)
	return err
}

func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	row := r.pool.QueryRow(ctx, `SELECT id, email, password, role, created_at FROM users WHERE email = $1`, email)
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	return user, err
}

func (r *PostgresRepository) ListRoles(ctx context.Context) ([]domain.Role, error) {
	rows, err := r.pool.Query(ctx, `SELECT name, permissions FROM roles`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]domain.Role, 0)
	for rows.Next() {
		var role domain.Role
		var permissions []string
		if err := rows.Scan(&role.Name, &permissions); err != nil {
			return nil, err
		}
		role.Permissions = permissions
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *PostgresRepository) CreateNode(ctx context.Context, node domain.Node) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO nodes (id, platform, hostname, ip_address, capabilities, version, last_seen) VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (id) DO UPDATE SET platform = EXCLUDED.platform, hostname = EXCLUDED.hostname, ip_address = EXCLUDED.ip_address, capabilities = EXCLUDED.capabilities, version = EXCLUDED.version, last_seen = EXCLUDED.last_seen`, node.ID, node.Platform, node.Hostname, node.IPAddress, node.Capabilities, node.Version, node.LastSeen)
	return err
}

func (r *PostgresRepository) GetNode(ctx context.Context, nodeID string) (domain.Node, error) {
	var node domain.Node
	row := r.pool.QueryRow(ctx, `SELECT id, platform, hostname, ip_address, capabilities, version, last_seen FROM nodes WHERE id = $1`, nodeID)
	err := row.Scan(&node.ID, &node.Platform, &node.Hostname, &node.IPAddress, &node.Capabilities, &node.Version, &node.LastSeen)
	return node, err
}

func (r *PostgresRepository) ListNodes(ctx context.Context) ([]domain.Node, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, platform, hostname, ip_address, capabilities, version, last_seen FROM nodes ORDER BY last_seen DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodes := make([]domain.Node, 0)
	for rows.Next() {
		var node domain.Node
		if err := rows.Scan(&node.ID, &node.Platform, &node.Hostname, &node.IPAddress, &node.Capabilities, &node.Version, &node.LastSeen); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (r *PostgresRepository) Close() {
	r.pool.Close()
}
