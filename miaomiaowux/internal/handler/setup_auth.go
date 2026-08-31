package handler

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	setupTokenHeader = "X-MMWX-Setup-Token"
	setupTokenCookie = "mmwx_setup"
)

// SetupProtector prevents an Internet-facing, uninitialized panel from being
// claimed by the first visitor. Loopback requests remain convenient; remote
// requests must first exchange the startup token for a short-lived HttpOnly
// cookie through /api/setup/authorize?token=...
type SetupProtector struct {
	token string
}

func NewSetupProtector() (*SetupProtector, error) {
	token := strings.TrimSpace(os.Getenv("MMWX_SETUP_TOKEN"))
	if token == "" {
		buf := make([]byte, 32)
		if _, err := rand.Read(buf); err != nil {
			return nil, err
		}
		token = base64.RawURLEncoding.EncodeToString(buf)
	}
	if len(token) < 24 {
		return nil, errors.New("MMWX_SETUP_TOKEN must contain at least 24 characters")
	}
	return &SetupProtector{token: token}, nil
}

func (p *SetupProtector) Token() string { return p.token }

func (p *SetupProtector) Protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isLoopbackRequest(r) && !p.authorized(r) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Cache-Control", "no-store")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"error":"setup authorization required"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (p *SetupProtector) AuthorizeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		candidate := strings.TrimSpace(r.URL.Query().Get("token"))
		if !constantTimeEqual(candidate, p.token) {
			http.Error(w, "invalid setup token", http.StatusForbidden)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     setupTokenCookie,
			Value:    p.token,
			Path:     "/api/setup/",
			MaxAge:   int((30 * time.Minute).Seconds()),
			HttpOnly: true,
			Secure:   requestIsHTTPS(r),
			SameSite: http.SameSiteStrictMode,
		})
		w.Header().Set("Cache-Control", "no-store")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
}

func (p *SetupProtector) authorized(r *http.Request) bool {
	candidate := strings.TrimSpace(r.Header.Get(setupTokenHeader))
	if candidate == "" {
		if cookie, err := r.Cookie(setupTokenCookie); err == nil {
			candidate = cookie.Value
		}
	}
	return constantTimeEqual(candidate, p.token)
}

func constantTimeEqual(candidate, expected string) bool {
	if candidate == "" || expected == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(candidate), []byte(expected)) == 1
}

func isLoopbackRequest(r *http.Request) bool {
	if r == nil {
		return false
	}
	remote := strings.TrimSpace(r.RemoteAddr)
	host, _, err := net.SplitHostPort(remote)
	if err != nil {
		host = strings.Trim(remote, "[]")
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		return false
	}
	// A local reverse proxy commonly connects from 127.0.0.1 while serving a
	// public hostname. Do not let that proxy turn the startup guard into a
	// public first-user registration endpoint. Loopback bypass is reserved for
	// direct localhost/loopback browser access without forwarded client data.
	if forwardedClientIP(r) != "" || hasForwardedHeaders(r) {
		return false
	}
	return isLocalHost(r.Host)
}

func hasForwardedHeaders(r *http.Request) bool {
	for _, name := range []string{"Forwarded", "X-Forwarded-For", "X-Forwarded-Host", "X-Forwarded-Proto", "X-Real-IP", "CF-Connecting-IP"} {
		if strings.TrimSpace(r.Header.Get(name)) != "" {
			return true
		}
	}
	return false
}

func isLocalHost(hostPort string) bool {
	host := strings.TrimSpace(hostPort)
	if parsedHost, _, err := net.SplitHostPort(host); err == nil {
		host = parsedHost
	}
	host = strings.Trim(strings.ToLower(host), "[]")
	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func requestIsHTTPS(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	// Only a trusted reverse proxy should be allowed to supply this header.
	return trustedProxyRequest(r) && strings.EqualFold(strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-Proto"), ",")[0]), "https")
}
