package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Store interface {
	CreateServer(ctx context.Context, arg CreateServerParams) (Server, error)
	GetServer(ctx context.Context, arg GetServerParams) ([]Server, error)
	UpdateServer(ctx context.Context, arg UpdateServerParams) (Server, error)
	DeleteServer(ctx context.Context, id int64) error
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetUser(ctx context.Context, username string) (User, error)
	GetAllServers(ctx context.Context) ([]Server, error)
	UpdateRole(ctx context.Context, arg UpdateRoleParams) (User, error)
	CreateScope(ctx context.Context, arg CreateScopeParams) (Scope, error)
	DeleteScope(ctx context.Context, id int64) error
	GetScope(ctx context.Context, id int64) ([]string, error)
	UpdateScope(ctx context.Context, arg UpdateScopeParams) (Scope, error)
}

type SQLStore struct {
	db *pgx.Conn
	*Queries
}

func NewStore(db *pgx.Conn) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}
