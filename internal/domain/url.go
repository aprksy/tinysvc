package domain

import (
	"net/url"
	"time"
)

// ShortURL represents a shortened URL
type ShortURL struct {
	ID        string     `json:"id"`
	LongURL   string     `json:"long_url"`
	ShortCode string     `json:"short_code"`
	Views     int64      `json:"views"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
}

// URLCreateRequest represents input for creating a short URL
type URLCreateRequest struct {
	LongURL    string `json:"long_url"`
	CustomCode string `json:"custom_code"` // Optional custom short code
	ExpiryDays *int   `json:"expiry_days"` // nil = forever, 0 = default (30 days)
}

// Validate checks if the URL creation request is valid
func (r *URLCreateRequest) Validate() error {
	if len(r.LongURL) == 0 {
		return ErrEmptyURL
	}

	// Validate URL format
	parsedURL, err := url.ParseRequestURI(r.LongURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return ErrInvalidURL
	}

	// Validate custom code if provided
	if r.CustomCode != "" {
		if len(r.CustomCode) < 3 || len(r.CustomCode) > 20 {
			return ErrInvalidCustomCode
		}
		// Only allow alphanumeric and dash
		for _, ch := range r.CustomCode {
			if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
				(ch >= '0' && ch <= '9') || ch == '-' || ch == '_') {
				return ErrInvalidCustomCode
			}
		}
	}

	if r.ExpiryDays != nil && *r.ExpiryDays < 0 && *r.ExpiryDays != -1 {
		return ErrInvalidExpiry
	}

	return nil
}

// IsExpired checks if the short URL has expired
func (u *ShortURL) IsExpired() bool {
	if u.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*u.ExpiresAt)
}
