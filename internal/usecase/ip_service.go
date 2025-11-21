package usecase

import (
	"net"
	"net/http"
	"strings"
)

// IPService handles IP address detection
type IPService interface {
	GetClientIP(r *http.Request) string
}

type ipService struct{}

// NewIPService creates a new IP service
func NewIPService() IPService {
	return &ipService{}
}

func (s *ipService) GetClientIP(r *http.Request) string {
	// Check Cloudflare headers first
	if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
		return ip
	}

	// Check X-Forwarded-For
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
