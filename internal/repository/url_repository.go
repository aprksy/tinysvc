package repository

import (
	"context"

	"github.com/aprksy/tinysvc/internal/domain"
)

// URLRepository defines storage operations for short URLs
type URLRepository interface {
	Create(ctx context.Context, url *domain.ShortURL) error
	GetByCode(ctx context.Context, code string) (*domain.ShortURL, error)
	GetByID(ctx context.Context, id string) (*domain.ShortURL, error)
	IncrementViews(ctx context.Context, code string) error
	Delete(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) (int64, error)
	CodeExists(ctx context.Context, code string) (bool, error)
}
