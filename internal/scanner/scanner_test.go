package scanner

import (
	"context"
	"github.com/azinamouto/kuro/internal/scope"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScan(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
	}))
	defer s.Close()
	p := scope.Policy{Version: 1, Origin: s.URL, MaxRequests: 10}
	r, err := Scan(context.Background(), p, s.URL)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Findings) != 3 {
		t.Fatalf("expected 3 observations, got %d", len(r.Findings))
	}
}
func TestRedirectNotFollowed(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "http://example.org", 302) }))
	defer s.Close()
	p := scope.Policy{Version: 1, Origin: s.URL, MaxRequests: 10}
	r, err := Scan(context.Background(), p, s.URL)
	if err != nil {
		t.Fatal(err)
	}
	if r.Status != 302 {
		t.Fatalf("redirect followed: %d", r.Status)
	}
}

