package store

import (
	"context"
	"what-to-do/internal/domain"

	_ "github.com/glebarez/go-sqlite"
)

type CommandStore interface {
	Create(ctx context.Context, c *domain.Command) (int64, error)
	Get(ctx context.Context, id int64) (*domain.Command, error)
	List(ctx context.Context) ([]*domain.Command, error)
	Update(ctx context.Context, c *domain.Command) error
	Delete(ctx context.Context, id int64) error
}
