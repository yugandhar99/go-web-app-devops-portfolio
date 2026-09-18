package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

var (
	startTime      = time.Now()
	requestTotal   atomic.Uint64
	healthRequests atomic.Uint64
)

// app builds a small production-ready HTTP router for the portfolio web app.
type app struct {
	staticDir string
}

func newApp(staticDir string) http.Handler {
	a := &app{staticDir: staticDir}
	mux := http.NewServeMux()

	mux.HandleFunc("/", a.rootHandler)
	mux.HandleFunc("/home", a.pageHandler("home.html"))
	mux.HandleFunc("/courses", a.pageHandler("courses.html"))
	mux.HandleFunc("/about", a.pageHandler("about.html"))
	mux.HandleFunc("/contact", a.pageHandler("contact.html"))
	mux.HandleFunc("/healthz", a.healthHandler)
	mux.HandleFunc("/readyz", a.readyHandler)
	mux.HandleFunc("/metrics", a.metricsHandler)

	return securityHeaders(loggingMiddleware(mux))
}

func (a *app) rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/home", http.StatusFound)
}

func (a *app) pageHandler(fileName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		filePath := filepath.Join(a.staticDir, fileName)
		w.Header().Set("Cache-Control", "public, max-age=300")
		http.ServeFile(w, r, filePath)
	}
}

func (a *app) healthHandler(w http.ResponseWriter, r *http.Request) {
	healthRequests.Add(1)
	writeJSON(w, http.StatusOK, `{"status":"ok","service":"go-web-app-devops"}`)
}

func (a *app) readyHandler(w http.ResponseWriter, r *http.Request) {
	// This app has no downstream dependency. A database/cache check can be added here later.
	writeJSON(w, http.StatusOK, `{"status":"ready"}`)
}

func (a *app) metricsHandler(w http.ResponseWriter, r *http.Request) {
	uptime := uint64(time.Since(startTime).Seconds())
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	fmt.Fprintf(w, "# HELP go_web_app_requests_total Total HTTP requests handled by the app.\n")
	fmt.Fprintf(w, "# TYPE go_web_app_requests_total counter\n")
	fmt.Fprintf(w, "go_web_app_requests_total %d\n", requestTotal.Load())
	fmt.Fprintf(w, "# HELP go_web_app_health_requests_total Total health check requests.\n")
	fmt.Fprintf(w, "# TYPE go_web_app_health_requests_total counter\n")
	fmt.Fprintf(w, "go_web_app_health_requests_total %d\n", healthRequests.Load())
	fmt.Fprintf(w, "# HELP go_web_app_uptime_seconds Application uptime in seconds.\n")
	fmt.Fprintf(w, "# TYPE go_web_app_uptime_seconds gauge\n")
	fmt.Fprintf(w, "go_web_app_uptime_seconds %d\n", uptime)
}

func writeJSON(w http.ResponseWriter, statusCode int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(body))
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data:; style-src 'self' 'unsafe-inline'")
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestTotal.Add(1)
		started := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("method=%s path=%s remote=%s duration_ms=%d", r.Method, r.URL.Path, r.RemoteAddr, time.Since(started).Milliseconds())
	})
}

func getPort() string {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		return "8080"
	}
	return port
}

func main() {
	addr := ":" + getPort()
	server := &http.Server{
		Addr:              addr,
		Handler:           newApp("static"),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("starting go-web-app-devops on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
