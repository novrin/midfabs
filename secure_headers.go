package midfabs

import (
	"net/http"
)

// SecureHeaders returns a middleware that sets a CSP policy and other default
// security headers consistent with OWASP guidance.
func SecureHeaders(csp string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Security-Policy", csp)
			// Omit referrer information in requests
			w.Header().Set("Referrer-Policy", "no-referrer")
			// Prevent MIME-type sniffing
			w.Header().Set("X-Content-Type-Options", "nosniff")
			// Prevent loading site contents in any frame
			w.Header().Set("X-Frame-Options", "deny")
			// Prevent loading over HTTP connections
			w.Header().Set("Strict-Transport-Security", "max-age=31536000")
			h.ServeHTTP(w, r)
		})
	}
}
