package http

import (
	"net/http"

	"github.com/aprksy/tinysvc/internal/usecase"
)

type IPHandler struct {
	ipService usecase.IPService
}

func NewIPHandler(ipService usecase.IPService) *IPHandler {
	return &IPHandler{
		ipService: ipService,
	}
}

// GetIP handles GET /ip
func (h *IPHandler) GetIP(w http.ResponseWriter, r *http.Request) {
	ip := h.ipService.GetClientIP(r)

	respondJSON(w, http.StatusOK, map[string]string{
		"ip": ip,
	})
}
