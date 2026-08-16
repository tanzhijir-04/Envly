package api

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/tanzhijir-04/Envly/engine/internal/events"
	"github.com/tanzhijir-04/Envly/engine/internal/executor"
	"github.com/tanzhijir-04/Envly/engine/internal/state"
)

type Server struct {
	store   *state.Store
	hub     *events.Hub
	exec    executor.Executor
	webDir  string
	version string

	runMu     sync.Mutex
	curRunID  string
	curCancel context.CancelFunc
}

func NewServer(store *state.Store, hub *events.Hub, exec executor.Executor, webDir, version string) *Server {
	return &Server{store: store, hub: hub, exec: exec, webDir: webDir, version: version}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("GET /api/catalog", s.handleCatalog)
	mux.HandleFunc("GET /api/templates", s.handleTemplates)
	mux.HandleFunc("POST /api/plan", s.handlePlan)
	mux.HandleFunc("POST /api/run", s.handleRun)
	mux.HandleFunc("POST /api/cancel", s.handleCancel)
	mux.HandleFunc("GET /api/events", s.handleEvents)
	mux.HandleFunc("GET /api/settings", s.handleGetSettings)
	mux.HandleFunc("POST /api/settings", s.handlePostSettings)
	mux.Handle("/", s.staticHandler())
	return mux
}

func (s *Server) staticHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		clean := filepath.Clean(strings.TrimPrefix(r.URL.Path, "/"))
		p := filepath.Join(s.webDir, clean)
		if r.URL.Path == "/" {
			p = filepath.Join(s.webDir, "index.html")
		}
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			http.ServeFile(w, r, p)
			return
		}
		http.ServeFile(w, r, filepath.Join(s.webDir, "index.html"))
	})
}

func platformName() string {
	switch runtime.GOOS {
	case "darwin":
		return "mac"
	case "linux":
		return "linux"
	default:
		return "win"
	}
}
