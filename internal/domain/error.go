package domain

import "errors"

var (
	// Paste errors
	ErrEmptyContent    = errors.New("content cannot be empty")
	ErrContentTooLarge = errors.New("content exceeds 1MB limit")
	ErrInvalidExpiry   = errors.New("expiry days cannot be negative")
	ErrPasteNotFound   = errors.New("paste not found")
	ErrPasteExpired    = errors.New("paste has expired")

	// URL errors
	ErrEmptyURL          = errors.New("URL cannot be empty")
	ErrInvalidURL        = errors.New("invalid URL format")
	ErrInvalidCustomCode = errors.New("custom code must be 3-20 characters (alphanumeric, dash, underscore)")
	ErrURLNotFound       = errors.New("short URL not found")
	ErrURLExpired        = errors.New("short URL has expired")
	ErrCustomCodeTaken   = errors.New("custom code already in use")

	// JSON bin errors
	ErrInvalidJSON  = errors.New("invalid JSON format")
	ErrJSONNotFound = errors.New("JSON bin not found")
	ErrJSONExpired  = errors.New("JSON bin has expired")
)
