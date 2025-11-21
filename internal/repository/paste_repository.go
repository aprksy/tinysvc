package repository

import (
	"context"

	"github.com/aprksy/tinysvc/internal/domain"
)

// PasteRepository defines storage operations for pastes
type PasteRepository interface {
	Create(ctx context.Context, paste *domain.Paste) error
	GetByID(ctx context.Context, id string) (*domain.Paste, error)
	Delete(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) (int64, error) // Returns count of deleted
}
