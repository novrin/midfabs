package midfabs

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestPanicBroker(t *testing.T) {
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sid := r.URL.Query().Get("id")
		if _, err := strconv.Atoi(sid); err != nil {
			panic("reached panic")
		}
		w.Write([]byte("ok"))
	})
	broker := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("broker caught panic"))
	})
	handler := PanicBroker(broker)(base)
	cases := map[string]struct {
		url        string
		wantBody   string
		wantStatus int
	}{
		"no panic": {
			url:        "/?id=10",
			wantBody:   "ok",
			wantStatus: http.StatusOK,
		},
		"panic": {
			url:        "/?id=foo",
			wantBody:   "broker caught panic",
			wantStatus: http.StatusInternalServerError,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, c.url, nil))
			if gotS := w.Code; gotS != c.wantStatus {
				t.Fatalf(errorString, gotS, c.wantStatus)
			}
			if gotB := w.Body.String(); gotB != c.wantBody {
				t.Fatalf(errorString, gotB, c.wantBody)
			}
		})
	}
}
