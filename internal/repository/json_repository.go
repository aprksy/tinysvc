package repository

import (
    "context"
    "github.com/aprksy/tinysvc/internal/domain"
)

// JSONRepository defines storage operations for JSON bins
type JSONRepository interface {
    Create(ctx context.Context, jsonBin *domain.JSONBin) error
    GetByID(ctx context.Context, id string) (*domain.JSONBin, error)
    Update(ctx context.Context, id string, content []byte) error
    IncrementViews(ctx context.Context, id string) error
    Delete(ctx context.Context, id string) error
    DeleteExpired(ctx context.Context) (int64, error)
}