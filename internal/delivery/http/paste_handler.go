package http

import (
	"encoding/json"
	"net/http"

	"github.com/aprksy/tinysvc/internal/domain"
	"github.com/aprksy/tinysvc/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type PasteHandler struct {
	pasteService usecase.PasteService
}

func NewPasteHandler(pasteService usecase.PasteService) *PasteHandler {
	return &PasteHandler{
		pasteService: pasteService,
	}
}

// CreatePaste handles POST /paste
// func (h *PasteHandler) CreatePaste(w http.ResponseWriter, r *http.Request) {
// 	var req domain.PasteCreateRequest

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		respondError(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}

// 	paste, err := h.pasteService.CreatePaste(r.Context(), req)
// 	if err != nil {
// 		switch err {
// 		case domain.ErrEmptyContent, domain.ErrInvalidExpiry:
// 			respondError(w, http.StatusBadRequest, err.Error())
// 		case domain.ErrContentTooLarge:
// 			respondError(w, http.StatusRequestEntityTooLarge, err.Error())
// 		default:
// 			respondError(w, http.StatusInternalServerError, "Failed to create paste")
// 		}
// 		return
// 	}

//		respondJSON(w, http.StatusCreated, paste)
//	}
//
// CreatePaste handles POST /paste
func (h *PasteHandler) CreatePaste(w http.ResponseWriter, r *http.Request) {
	var req domain.PasteCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	paste, err := h.pasteService.CreatePaste(r.Context(), req)
	if err != nil {
		switch err {
		case domain.ErrEmptyContent, domain.ErrInvalidExpiry:
			respondError(w, http.StatusBadRequest, err.Error())
		case domain.ErrContentTooLarge:
			respondError(w, http.StatusRequestEntityTooLarge, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "Failed to create paste")
		}
		return
	}

	// Log the response for debugging
	responseJSON, _ := json.Marshal(paste)
	println("Created paste response:", string(responseJSON))

	respondJSON(w, http.StatusCreated, paste)
}

// GetPaste handles GET /paste/{id}
func (h *PasteHandler) GetPaste(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing paste ID")
		return
	}

	paste, err := h.pasteService.GetPaste(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrPasteNotFound:
			respondError(w, http.StatusNotFound, "Paste not found")
		case domain.ErrPasteExpired:
			respondError(w, http.StatusGone, "Paste has expired")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to retrieve paste")
		}
		return
	}

	respondJSON(w, http.StatusOK, paste)
}

// DeletePaste handles DELETE /paste/{id}
func (h *PasteHandler) DeletePaste(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing paste ID")
		return
	}

	if err := h.pasteService.DeletePaste(r.Context(), id); err != nil {
		switch err {
		case domain.ErrPasteNotFound:
			respondError(w, http.StatusNotFound, "Paste not found")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to delete paste")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
