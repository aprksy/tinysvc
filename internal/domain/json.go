package domain

import (
	"encoding/json"
	"time"
)

// JSONBin represents a stored JSON document
type JSONBin struct {
	ID        string          `json:"id"`
	Content   json.RawMessage `json:"content"` // Store as raw JSON
	Name      string          `json:"name"`    // Optional name/description
	IsPublic  bool            `json:"is_public"`
	Views     int64           `json:"views"`
	ExpiresAt *time.Time      `json:"expires_at"`
	CreatedAt time.Time       `json:"created_at"`
}

// JSONBinCreateRequest represents input for creating a JSON bin
type JSONBinCreateRequest struct {
	Content    json.RawMessage `json:"content"`
	Name       string          `json:"name"`      // Optional
	IsPublic   bool            `json:"is_public"` // Default: true (no auth yet)
	ExpiryDays *int            `json:"expiry_days"`
}

// JSONBinUpdateRequest represents input for updating a JSON bin
type JSONBinUpdateRequest struct {
	Content json.RawMessage `json:"content"`
}

// Validate checks if the JSON bin creation request is valid
func (r *JSONBinCreateRequest) Validate() error {
	if len(r.Content) == 0 {
		return ErrEmptyContent
	}

	// Validate JSON syntax
	var js interface{}
	if err := json.Unmarshal(r.Content, &js); err != nil {
		return ErrInvalidJSON
	}

	// 10MB limit (same as paste)
	if len(r.Content) > 10*1024*1024 {
		return ErrContentTooLarge
	}

	if r.ExpiryDays != nil && *r.ExpiryDays < 0 && *r.ExpiryDays != -1 {
		return ErrInvalidExpiry
	}

	return nil
}

// Validate checks if the JSON bin update request is valid
func (r *JSONBinUpdateRequest) Validate() error {
	if len(r.Content) == 0 {
		return ErrEmptyContent
	}

	// Validate JSON syntax
	var js interface{}
	if err := json.Unmarshal(r.Content, &js); err != nil {
		return ErrInvalidJSON
	}

	if len(r.Content) > 10*1024*1024 {
		return ErrContentTooLarge
	}

	return nil
}

// IsExpired checks if the JSON bin has expired
func (j *JSONBin) IsExpired() bool {
	if j.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*j.ExpiresAt)
}
