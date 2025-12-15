package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aprksy/tinysvc/internal/domain"
	"github.com/aprksy/tinysvc/internal/repository"
)

// JSONService handles JSON bin business logic
type JSONService interface {
	CreateJSONBin(ctx context.Context, req domain.JSONBinCreateRequest) (*domain.JSONBin, error)
	GetJSONBin(ctx context.Context, id string) (*domain.JSONBin, error)
	UpdateJSONBin(ctx context.Context, id string, req domain.JSONBinUpdateRequest) (*domain.JSONBin, error)
	DeleteJSONBin(ctx context.Context, id string) error
	CleanupExpired(ctx context.Context) (int64, error)
}

type jsonService struct {
	repo repository.JSONRepository
}

// NewJSONService creates a new JSON service
func NewJSONService(repo repository.JSONRepository) JSONService {
	return &jsonService{
		repo: repo,
	}
}

func (s *jsonService) CreateJSONBin(ctx context.Context, req domain.JSONBinCreateRequest) (*domain.JSONBin, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	jsonBin := &domain.JSONBin{
		ID:        generateID(),
		Content:   req.Content,
		Name:      req.Name,
		IsPublic:  req.IsPublic,
		Views:     0,
		CreatedAt: time.Now().UTC(),
	}

	// Calculate expiry (same logic as paste/url)
	if req.ExpiryDays == nil {
		expiresAt := time.Now().UTC().AddDate(0, 0, DefaultExpiryDays)
		jsonBin.ExpiresAt = &expiresAt
	} else if *req.ExpiryDays == 0 {
		expiresAt := time.Now().UTC().AddDate(0, 0, DefaultExpiryDays)
		jsonBin.ExpiresAt = &expiresAt
	} else if *req.ExpiryDays > 0 {
		expiresAt := time.Now().UTC().AddDate(0, 0, *req.ExpiryDays)
		jsonBin.ExpiresAt = &expiresAt
	} else if *req.ExpiryDays == -1 {
		jsonBin.ExpiresAt = nil
	}

	if err := s.repo.Create(ctx, jsonBin); err != nil {
		return nil, err
	}

	return jsonBin, nil
}

func (s *jsonService) GetJSONBin(ctx context.Context, id string) (*domain.JSONBin, error) {
	jsonBin, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if jsonBin.IsExpired() {
		return nil, domain.ErrJSONExpired
	}

	// Increment views
	_ = s.repo.IncrementViews(ctx, id)

	return jsonBin, nil
}

func (s *jsonService) UpdateJSONBin(ctx context.Context, id string, req domain.JSONBinUpdateRequest) (*domain.JSONBin, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Check if exists and not expired
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existing.IsExpired() {
		return nil, domain.ErrJSONExpired
	}

	// Update content
	if err := s.repo.Update(ctx, id, req.Content); err != nil {
		return nil, err
	}

	// Get updated version
	return s.repo.GetByID(ctx, id)
}

func (s *jsonService) DeleteJSONBin(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *jsonService) CleanupExpired(ctx context.Context) (int64, error) {
	return s.repo.DeleteExpired(ctx)
}

// PrettyPrint formats JSON with indentation
func PrettyPrint(data json.RawMessage) (string, error) {
	var prettyJSON interface{}
	if err := json.Unmarshal(data, &prettyJSON); err != nil {
		return "", err
	}

	formatted, err := json.MarshalIndent(prettyJSON, "", "  ")
	if err != nil {
		return "", err
	}

	return string(formatted), nil
}
