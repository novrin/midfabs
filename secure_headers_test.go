package midfabs

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	cases := map[string]struct {
		handler http.Handler
		want    []string
	}{
		"no csp": {
			handler: SecureHeaders("")(base),
			want: []string{
				"Referrer-Policy",
				"X-Content-Type-Options",
				"X-Frame-Options",
				"Strict-Transport-Security",
			},
		},
		"simple csp": {
			handler: SecureHeaders("defaut-src 'self';")(base),
			want: []string{
				"Content-Security-Policy",
				"Referrer-Policy",
				"X-Content-Type-Options",
				"X-Frame-Options",
				"Strict-Transport-Security",
			},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
			for i, key := range c.want {
				if got := w.Header().Get(key); got == "" {
					t.Fatalf(errorString, got, c.want[i])
				}
			}
		})
	}
}
