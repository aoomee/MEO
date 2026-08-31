package handler

import (
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

var (
	trustedProxyOnce sync.Once
	trustedProxyNets []*net.IPNet
)

// trustedProxyRequest reports whether forwarding headers may be trusted.
// The default trusts loopback only. Reverse-proxy deployments must explicitly
// set MMWX_TRUSTED_PROXIES to comma-separated IPs/CIDRs (for example
// 172.18.0.0/16 for a dedicated Docker network).
func trustedProxyRequest(r *http.Request) bool {
	if r == nil {
		return false
	}
	host := remoteHost(r.RemoteAddr)
	ip := net.ParseIP(host)
	return ip != nil && isTrustedProxyIP(ip)
}

func isTrustedProxyIP(ip net.IP) bool {
	trustedProxyOnce.Do(loadTrustedProxyNets)
	for _, network := range trustedProxyNets {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

func loadTrustedProxyNets() {
	values := []string{"127.0.0.0/8", "::1/128"}
	for _, value := range strings.Split(os.Getenv("MMWX_TRUSTED_PROXIES"), ",") {
		if value = strings.TrimSpace(value); value != "" {
			values = append(values, value)
		}
	}
	for _, value := range values {
		if ip := net.ParseIP(value); ip != nil {
			bits := 128
			if ip.To4() != nil {
				bits = 32
			}
			trustedProxyNets = append(trustedProxyNets, &net.IPNet{IP: ip, Mask: net.CIDRMask(bits, bits)})
			continue
		}
		if _, network, err := net.ParseCIDR(value); err == nil {
			trustedProxyNets = append(trustedProxyNets, network)
		}
	}
}

func remoteHost(remoteAddr string) string {
	remote := strings.TrimSpace(remoteAddr)
	if host, _, err := net.SplitHostPort(remote); err == nil {
		return strings.Trim(host, "[]")
	}
	return strings.Trim(remote, "[]")
}

func forwardedClientIP(r *http.Request) string {
	if !trustedProxyRequest(r) {
		return ""
	}
	if candidate := validIPString(r.Header.Get("CF-Connecting-IP")); candidate != "" {
		return candidate
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		for i := len(parts) - 1; i >= 0; i-- {
			candidate := validIPString(parts[i])
			if candidate == "" {
				continue
			}
			ip := net.ParseIP(candidate)
			if !isTrustedProxyIP(ip) || i == 0 {
				return candidate
			}
		}
	}
	return validIPString(r.Header.Get("X-Real-IP"))
}

func validIPString(value string) string {
	value = strings.Trim(strings.TrimSpace(value), "[]")
	if ip := net.ParseIP(value); ip != nil {
		return ip.String()
	}
	return ""
}

func sameOriginOrNoOrigin(r *http.Request) bool {
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		// Native Agent/tester clients do not send a browser Origin header.
		return true
	}
	u, err := url.Parse(origin)
	if err != nil || u.Host == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" {
		return false
	}
	if !strings.EqualFold(u.Host, r.Host) {
		return false
	}
	requestScheme := "http"
	if requestIsHTTPS(r) {
		requestScheme = "https"
	}
	return strings.EqualFold(u.Scheme, requestScheme)
}

func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Content-Security-Policy", "default-src 'self'; base-uri 'self'; object-src 'none'; frame-ancestors 'none'; form-action 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; font-src 'self' data:; img-src 'self' data: blob: https:; connect-src 'self' ws: wss: https:")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=(), usb=()")
		if requestIsHTTPS(r) {
			h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		next.ServeHTTP(w, r)
	})
}
