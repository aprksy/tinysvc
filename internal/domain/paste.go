package domain

import (
    "time"
)

// Paste represents a text snippet with optional expiration
type Paste struct {
    ID         string
    Content    string
    IsMarkdown bool
    ExpiresAt  *time.Time // nil means never expires
    CreatedAt  time.Time
}

// PasteCreateRequest represents input for creating a paste
type PasteCreateRequest struct {
    Content       string
    IsMarkdown    bool
    ExpiryDays    *int // nil = forever, 0 = default (30 days)
}

// Validate checks if the paste creation request is valid
func (r *PasteCreateRequest) Validate() error {
    if len(r.Content) == 0 {
        return ErrEmptyContent
    }
    if len(r.Content) > 10*1024*1024 { // 10MB
        return ErrContentTooLarge
    }
    if r.ExpiryDays != nil && *r.ExpiryDays < 0 {
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