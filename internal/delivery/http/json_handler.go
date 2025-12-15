package http

import (
	"encoding/json"
	"net/http"

	"github.com/aprksy/tinysvc/internal/domain"
	"github.com/aprksy/tinysvc/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type JSONHandler struct {
	jsonService usecase.JSONService
}

func NewJSONHandler(jsonService usecase.JSONService) *JSONHandler {
	return &JSONHandler{
		jsonService: jsonService,
	}
}

// CreateJSONBin handles POST /json
func (h *JSONHandler) CreateJSONBin(w http.ResponseWriter, r *http.Request) {
	var req domain.JSONBinCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	jsonBin, err := h.jsonService.CreateJSONBin(r.Context(), req)
	if err != nil {
		switch err {
		case domain.ErrEmptyContent, domain.ErrInvalidJSON, domain.ErrInvalidExpiry:
			respondError(w, http.StatusBadRequest, err.Error())
		case domain.ErrContentTooLarge:
			respondError(w, http.StatusRequestEntityTooLarge, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "Failed to create JSON bin")
		}
		return
	}

	respondJSON(w, http.StatusCreated, jsonBin)
}

// GetJSONBin handles GET /json/{id}
func (h *JSONHandler) GetJSONBin(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing JSON bin ID")
		return
	}

	jsonBin, err := h.jsonService.GetJSONBin(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrJSONNotFound:
			respondError(w, http.StatusNotFound, "JSON bin not found")
		case domain.ErrJSONExpired:
			respondError(w, http.StatusGone, "JSON bin has expired")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to retrieve JSON bin")
		}
		return
	}

	respondJSON(w, http.StatusOK, jsonBin)
}

// GetJSONBinRaw handles GET /json/{id}/raw - returns just the JSON content
func (h *JSONHandler) GetJSONBinRaw(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing JSON bin ID")
		return
	}

	jsonBin, err := h.jsonService.GetJSONBin(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrJSONNotFound:
			respondError(w, http.StatusNotFound, "JSON bin not found")
		case domain.ErrJSONExpired:
			respondError(w, http.StatusGone, "JSON bin has expired")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to retrieve JSON bin")
		}
		return
	}

	// Return raw JSON content
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBin.Content)
}

// UpdateJSONBin handles PUT /json/{id}
func (h *JSONHandler) UpdateJSONBin(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing JSON bin ID")
		return
	}

	var req domain.JSONBinUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	jsonBin, err := h.jsonService.UpdateJSONBin(r.Context(), id, req)
	if err != nil {
		switch err {
		case domain.ErrJSONNotFound:
			respondError(w, http.StatusNotFound, "JSON bin not found")
		case domain.ErrJSONExpired:
			respondError(w, http.StatusGone, "JSON bin has expired")
		case domain.ErrEmptyContent, domain.ErrInvalidJSON:
			respondError(w, http.StatusBadRequest, err.Error())
		case domain.ErrContentTooLarge:
			respondError(w, http.StatusRequestEntityTooLarge, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "Failed to update JSON bin")
		}
		return
	}

	respondJSON(w, http.StatusOK, jsonBin)
}

// DeleteJSONBin handles DELETE /json/{id}
func (h *JSONHandler) DeleteJSONBin(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing JSON bin ID")
		return
	}

	if err := h.jsonService.DeleteJSONBin(r.Context(), id); err != nil {
		switch err {
		case domain.ErrJSONNotFound:
			respondError(w, http.StatusNotFound, "JSON bin not found")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to delete JSON bin")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
