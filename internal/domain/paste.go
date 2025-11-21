package domain

import (
	"time"
)

// Paste represents a text snippet with optional expiration
type Paste struct {
	ID         string     `json:"id"`
	Content    string     `json:"content"`
	IsMarkdown bool       `json:"is_markdown"`
	ExpiresAt  *time.Time `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

// PasteCreateRequest represents input for creating a paste
type PasteCreateRequest struct {
	Content    string `json:"content"`
	IsMarkdown bool   `json:"is_markdown"`
	ExpiryDays *int   `json:"expiry_days"` // nil = forever, 0 = default (30 days)
}

// Validate checks if the paste creation request is valid
func (r *PasteCreateRequest) Validate() error {
	if len(r.Content) == 0 {
		return ErrEmptyContent
	}
	if len(r.Content) > 10*1024*1024 { // 10MB
		return ErrContentTooLarge
	}
	if r.ExpiryDays != nil && *r.ExpiryDays < 0 && *r.ExpiryDays != -1 {
		return ErrInvalidExpiry
	}
	return nil
}

// IsExpired checks if the paste has expired
func (p *Paste) IsExpired() bool {
	if p.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*p.ExpiresAt)
}
