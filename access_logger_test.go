package midfabs

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const errorString = "\nGot:\t%v\nWant:\t%v\n"

func TestWriteHeader(t *testing.T) {
	cases := map[string]struct {
		status int
		want   int
	}{
		"ok":  {status: http.StatusOK, want: http.StatusOK},
		"bad": {status: http.StatusBadRequest, want: http.StatusBadRequest},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := responseWriter{ResponseWriter: httptest.NewRecorder()}
			w.WriteHeader(c.status)
			if got := w.status; got != c.want {
				t.Fatalf(errorString, got, c.want)
			}
		})
	}

}

func TestAccessLogger(t *testing.T) {
	defaultSubstrings := []string{
		"status=200",
		"duration=",
		"method=GET",
		"proto=HTTP/1.1",
		"dest=/",
		"src=",
		"agent=",
	}
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	cases := map[string]struct {
		handler http.Handler
		want    string
	}{
		"no prefix": {
			handler: AccessLogger(logger, "")(base),
			want:    "msg=",
		},
		"with prefix": {
			handler: AccessLogger(logger, "foo")(base),
			want:    "msg=foo",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c.handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
			got := buf.String()
			for i, sub := range defaultSubstrings {
				if !strings.Contains(got, sub) {
					t.Fatalf("Did not find wanted substring %s", defaultSubstrings[i])
				}
				if !strings.Contains(got, c.want) {
					t.Fatalf("\nGot\t'%s'\n\tand could not find '%s'\n", got, c.want)
				}
			}
		})
	}
}
