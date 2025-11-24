package http

import (
	"encoding/json"
	"net/http"

	"github.com/aprksy/tinysvc/internal/domain"
	"github.com/aprksy/tinysvc/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	urlService usecase.URLService
}

func NewURLHandler(urlService usecase.URLService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}

// CreateShortURL handles POST /shorten
func (h *URLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	var req domain.URLCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	shortURL, err := h.urlService.CreateShortURL(r.Context(), req)
	if err != nil {
		switch err {
		case domain.ErrEmptyURL, domain.ErrInvalidURL, domain.ErrInvalidCustomCode, domain.ErrInvalidExpiry:
			respondError(w, http.StatusBadRequest, err.Error())
		case domain.ErrCustomCodeTaken:
			respondError(w, http.StatusConflict, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "Failed to create short URL")
		}
		return
	}

	respondJSON(w, http.StatusCreated, shortURL)
}

// GetShortURL handles GET /shorten/{code}
func (h *URLHandler) GetShortURL(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		respondError(w, http.StatusBadRequest, "Missing short code")
		return
	}

	shortURL, err := h.urlService.GetShortURL(r.Context(), code)
	if err != nil {
		switch err {
		case domain.ErrURLNotFound:
			respondError(w, http.StatusNotFound, "Short URL not found")
		case domain.ErrURLExpired:
			respondError(w, http.StatusGone, "Short URL has expired")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to retrieve short URL")
		}
		return
	}

	respondJSON(w, http.StatusOK, shortURL)
}

// RedirectShortURL handles GET /s/{code} - redirects to long URL
func (h *URLHandler) RedirectShortURL(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		http.Error(w, "Missing short code", http.StatusBadRequest)
		return
	}

	shortURL, err := h.urlService.GetShortURL(r.Context(), code)
	if err != nil {
		switch err {
		case domain.ErrURLNotFound:
			http.Error(w, "Short URL not found", http.StatusNotFound)
		case domain.ErrURLExpired:
			http.Error(w, "Short URL has expired", http.StatusGone)
		default:
			http.Error(w, "Failed to retrieve short URL", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, shortURL.LongURL, http.StatusFound)
}

// DeleteShortURL handles DELETE /shorten/{id}
func (h *URLHandler) DeleteShortURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing URL ID")
		return
	}

	if err := h.urlService.DeleteShortURL(r.Context(), id); err != nil {
		switch err {
		case domain.ErrURLNotFound:
			respondError(w, http.StatusNotFound, "Short URL not found")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to delete short URL")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
