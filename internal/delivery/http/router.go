package http

import (
	"net/http"
	"time"

	"github.com/aprksy/tinysvc/internal/delivery/http/middleware"
	"github.com/aprksy/tinysvc/internal/usecase"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"golang.org/x/time/rate"
)

// Router holds all HTTP handlers
type Router struct {
	pasteHandler *PasteHandler
	ipHandler    *IPHandler
	rateLimiter  *middleware.RateLimiter
}

// NewRouter creates a new HTTP router with all handlers
func NewRouter(
	pasteService usecase.PasteService,
	ipService usecase.IPService,
) *Router {
	return &Router{
		pasteHandler: NewPasteHandler(pasteService),
		ipHandler:    NewIPHandler(ipService),
		rateLimiter:  middleware.NewRateLimiter(rate.Limit(10), 20), // 10 req/s, burst 20
	}
}

// SetupRoutes configures all routes
func (rt *Router) SetupRoutes() http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))
	r.Use(rt.rateLimiter.Limit)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// IP service
		r.Get("/ip", rt.ipHandler.GetIP)

		// Paste service
		r.Post("/paste", rt.pasteHandler.CreatePaste)
		r.Get("/paste/{id}", rt.pasteHandler.GetPaste)
		r.Delete("/paste/{id}", rt.pasteHandler.DeletePaste)
	})

	return r
}
