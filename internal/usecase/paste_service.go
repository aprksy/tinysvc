package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/aprksy/tinysvc/internal/domain"
	"github.com/aprksy/tinysvc/internal/repository"
)

const (
	DefaultExpiryDays = 30
	IDLength          = 8
)

// PasteService handles paste business logic
type PasteService interface {
	CreatePaste(ctx context.Context, req domain.PasteCreateRequest) (*domain.Paste, error)
	GetPaste(ctx context.Context, id string) (*domain.Paste, error)
	DeletePaste(ctx context.Context, id string) error
	CleanupExpired(ctx context.Context) (int64, error)
}

type pasteService struct {
	repo repository.PasteRepository
}

// NewPasteService creates a new paste service
func NewPasteService(repo repository.PasteRepository) PasteService {
	return &pasteService{
		repo: repo,
	}
}

func (s *pasteService) CreatePaste(ctx context.Context, req domain.PasteCreateRequest) (*domain.Paste, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	paste := &domain.Paste{
		ID:         generateID(),
		Content:    req.Content,
		IsMarkdown: req.IsMarkdown,
		CreatedAt:  time.Now(),
	}

	// Calculate expiry
	expiryDays := DefaultExpiryDays
	if req.ExpiryDays != nil {
		if *req.ExpiryDays == 0 {
			expiryDays = DefaultExpiryDays
		} else if *req.ExpiryDays > 0 {
			expiryDays = *req.ExpiryDays
		} else {
			// Negative means never expire
			paste.ExpiresAt = nil
		}
	}

	if paste.ExpiresAt == nil && expiryDays > 0 {
		expiresAt := time.Now().AddDate(0, 0, expiryDays)
		paste.ExpiresAt = &expiresAt
	}

	if err := s.repo.Create(ctx, paste); err != nil {
		return nil, err
	}

	return paste, nil
}

func (s *pasteService) GetPaste(ctx context.Context, id string) (*domain.Paste, error) {
	paste, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if paste.IsExpired() {
		return nil, domain.ErrPasteExpired
	}

	return paste, nil
}

func (s *pasteService) DeletePaste(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *pasteService) CleanupExpired(ctx context.Context) (int64, error) {
	return s.repo.DeleteExpired(ctx)
}

// generateID creates a random URL-safe ID
func generateID() string {
	b := make([]byte, IDLength)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:IDLength]
}
