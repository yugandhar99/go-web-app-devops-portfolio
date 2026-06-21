package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPageHandlers(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{name: "home", path: "/home"},
		{name: "courses", path: "/courses"},
		{name: "about", path: "/about"},
		{name: "contact", path: "/contact"},
	}

	handler := newApp("static")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %d", rr.Code)
			}
			if got := rr.Header().Get("Content-Type"); !strings.Contains(got, "text/html") {
				t.Fatalf("expected text/html content type, got %q", got)
			}
		})
	}
}

func TestRootRedirectsToHome(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	newApp("static").ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Fatalf("expected status 302, got %d", rr.Code)
	}
	if location := rr.Header().Get("Location"); location != "/home" {
		t.Fatalf("expected redirect to /home, got %q", location)
	}
}

func TestHealthReadinessAndMetrics(t *testing.T) {
	handler := newApp("static")

	for _, path := range []string{"/healthz", "/readyz"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("%s expected status 200, got %d", path, rr.Code)
		}
		if got := rr.Header().Get("Content-Type"); !strings.Contains(got, "application/json") {
			t.Fatalf("%s expected JSON content type, got %q", path, got)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("metrics expected status 200, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "go_web_app_requests_total") {
		t.Fatalf("metrics response missing request counter")
	}
}

func TestMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/home", nil)
	rr := httptest.NewRecorder()

	newApp("static").ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rr.Code)
	}
}
