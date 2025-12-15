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
	urlHandler   *URLHandler
	jsonHandler  *JSONHandler // ADD THIS
	rateLimiter  *middleware.RateLimiter
}

// NewRouter creates a new HTTP router with all handlers
func NewRouter(
	pasteService usecase.PasteService,
	ipService usecase.IPService,
	urlService usecase.URLService,
	jsonService usecase.JSONService, // ADD THIS
) *Router {
	return &Router{
		pasteHandler: NewPasteHandler(pasteService),
		ipHandler:    NewIPHandler(ipService),
		urlHandler:   NewURLHandler(urlService),
		jsonHandler:  NewJSONHandler(jsonService), // ADD THIS
		rateLimiter:  middleware.NewRateLimiter(rate.Limit(10), 20),
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

	// Serve static files (CSS, JS)
	fileServer := http.FileServer(http.Dir("./web"))
	r.Handle("/css/*", fileServer)
	r.Handle("/js/*", fileServer)

	// Serve HTML pages
	r.Get("/", serveHTML("./web/index.html"))
	r.Get("/ip.html", serveHTML("./web/ip.html"))
	r.Get("/paste.html", serveHTML("./web/paste.html"))
	r.Get("/shorten.html", serveHTML("./web/shorten.html"))
	r.Get("/json.html", serveHTML("./web/json.html")) // ADD THIS

	// Short URL redirect
	r.Get("/s/{code}", rt.urlHandler.RedirectShortURL)

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

		// URL shortener service
		r.Post("/shorten", rt.urlHandler.CreateShortURL)
		r.Get("/shorten/{code}", rt.urlHandler.GetShortURL)
		r.Delete("/shorten/{id}", rt.urlHandler.DeleteShortURL)

		// JSON bin service - ADD THIS BLOCK
		r.Post("/json", rt.jsonHandler.CreateJSONBin)
		r.Get("/json/{id}", rt.jsonHandler.GetJSONBin)
		r.Get("/json/{id}/raw", rt.jsonHandler.GetJSONBinRaw)
		r.Put("/json/{id}", rt.jsonHandler.UpdateJSONBin)
		r.Delete("/json/{id}", rt.jsonHandler.DeleteJSONBin)
	})

	return r
}

// serveHTML returns a handler that serves an HTML file
func serveHTML(filepath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath)
	}
}
