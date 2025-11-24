package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"

	"github.com/aprksy/tinysvc/internal/domain"
	"github.com/aprksy/tinysvc/internal/repository"
)

const (
	ShortCodeLength = 6
)

// URLService handles URL shortening business logic
type URLService interface {
	CreateShortURL(ctx context.Context, req domain.URLCreateRequest) (*domain.ShortURL, error)
	GetShortURL(ctx context.Context, code string) (*domain.ShortURL, error)
	GetShortURLByID(ctx context.Context, id string) (*domain.ShortURL, error)
	DeleteShortURL(ctx context.Context, id string) error
	CleanupExpired(ctx context.Context) (int64, error)
}

type urlService struct {
	repo repository.URLRepository
}

// NewURLService creates a new URL service
func NewURLService(repo repository.URLRepository) URLService {
	return &urlService{
		repo: repo,
	}
}

func (s *urlService) CreateShortURL(ctx context.Context, req domain.URLCreateRequest) (*domain.ShortURL, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	shortURL := &domain.ShortURL{
		ID:        generateID(),
		LongURL:   req.LongURL,
		Views:     0,
		CreatedAt: time.Now().UTC(),
	}

	// Handle custom code or generate random
	if req.CustomCode != "" {
		// Check if custom code already exists
		exists, err := s.repo.CodeExists(ctx, req.CustomCode)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, domain.ErrCustomCodeTaken
		}
		shortURL.ShortCode = req.CustomCode
	} else {
		// Generate random short code
		for i := 0; i < 5; i++ { // Try 5 times
			code := generateShortCode()
			exists, err := s.repo.CodeExists(ctx, code)
			if err != nil {
				return nil, err
			}
			if !exists {
				shortURL.ShortCode = code
				break
			}
		}
		if shortURL.ShortCode == "" {
			return nil, domain.ErrInvalidCustomCode // Couldn't generate unique code
		}
	}

	// Calculate expiry (same logic as paste)
	if req.ExpiryDays == nil {
		expiresAt := time.Now().UTC().AddDate(0, 0, DefaultExpiryDays)
		shortURL.ExpiresAt = &expiresAt
	} else if *req.ExpiryDays == 0 {
		expiresAt := time.Now().UTC().AddDate(0, 0, DefaultExpiryDays)
		shortURL.ExpiresAt = &expiresAt
	} else if *req.ExpiryDays > 0 {
		expiresAt := time.Now().UTC().AddDate(0, 0, *req.ExpiryDays)
		shortURL.ExpiresAt = &expiresAt
	} else if *req.ExpiryDays == -1 {
		shortURL.ExpiresAt = nil
	}

	if err := s.repo.Create(ctx, shortURL); err != nil {
		return nil, err
	}

	return shortURL, nil
}

func (s *urlService) GetShortURL(ctx context.Context, code string) (*domain.ShortURL, error) {
	url, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if url.IsExpired() {
		return nil, domain.ErrURLExpired
	}

	// Increment views
	_ = s.repo.IncrementViews(ctx, code)

	return url, nil
}

func (s *urlService) GetShortURLByID(ctx context.Context, id string) (*domain.ShortURL, error) {
	url, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if url.IsExpired() {
		return nil, domain.ErrURLExpired
	}

	return url, nil
}

func (s *urlService) DeleteShortURL(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *urlService) CleanupExpired(ctx context.Context) (int64, error) {
	return s.repo.DeleteExpired(ctx)
}

// generateShortCode creates a random URL-safe short code
func generateShortCode() string {
	b := make([]byte, ShortCodeLength)
	rand.Read(b)
	code := base64.URLEncoding.EncodeToString(b)[:ShortCodeLength]
	// Make it more URL-friendly (remove special chars)
	code = strings.ReplaceAll(code, "-", "")
	code = strings.ReplaceAll(code, "_", "")
	if len(code) < ShortCodeLength {
		code += "x" // Pad if needed
	}
	return code[:ShortCodeLength]
}
