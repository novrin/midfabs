// Package midfabs provides prefabrications for common HTTP middleware.
package midfabs

import "net/http"

// Middleware wraps an http.Handler. Use it to insert code before or after a
// given handler calls ServeHTTP.
type Middleware func(http.Handler) http.Handler
