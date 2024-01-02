package midfabs

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

// PanicBroker returns a Middleware that uses slog to write the debug stack to
// stderr and run broker's ServeHTTP if a panic originates after the handler
// calls ServeHTTP. Use it implement graceful recovery.
func PanicBroker(broker http.Handler) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					trace := fmt.Sprintf("%s\n%s", err, debug.Stack())
					slog.Error(trace, slog.String("by", "midway.PanicBroker"))
					broker.ServeHTTP(w, r)
				}
			}()
			h.ServeHTTP(w, r)
		})
	}
}
