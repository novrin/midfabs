package midfabs

import (
	"log/slog"
	"net/http"
	"time"
)

// responseWriter encapsulates http.ResponseWriter along with an additional
// status int field. It is used in AccessLogger to capture response status
// before the encapsulated ResponseWriter calls WriteHeader.
type responseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader sets rw's status before the encapsulated ResponseWriter calls
// WriteHeader.
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// AccessLogger returns a middleware that uses slog logger to track request and
// response details.
func AccessLogger(logger *slog.Logger, prefix string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w}
			h.ServeHTTP(rw, r)
			logger.Info(
				prefix,
				slog.Int("status", rw.status),
				slog.Duration("duration", time.Since(start)),
				slog.String("method", r.Method),
				slog.String("proto", r.Proto),
				slog.String("dest", r.URL.RequestURI()),
				slog.String("src", r.RemoteAddr),
				slog.String("agent", r.Header.Get("User-Agent")),
			)
		})
	}
}
