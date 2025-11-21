package domain

import "errors"

var (
    ErrEmptyContent     = errors.New("content cannot be empty")
    ErrContentTooLarge  = errors.New("content exceeds 10MB limit")
    ErrInvalidExpiry    = errors.New("expiry days cannot be negative")
    ErrPasteNotFound    = errors.New("paste not found")
    ErrPasteExpired     = errors.New("paste has expired")
)