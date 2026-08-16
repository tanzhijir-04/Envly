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
	"github.com/tanzhijir-04/Envly/engine/internal/network"
	"github.com/tanzhijir-04/Envly/engine/internal/state"
	"github.com/tanzhijir-04/Envly/engine/internal/store"
	"github.com/tanzhijir-04/Envly/engine/internal/verify"
)

type RecordStore interface {
	Records() ([]store.Record, error)
	EnvOps() ([]store.EnvOp, error)
}

type Server struct {
	store   *state.Store
	records RecordStore
	hub     *events.Hub
	exec    executor.Executor
	ver     *verify.Verifier
	net     *network.Detector
	webDir  string
	version string

	runMu     sync.Mutex
	curRunID  string
	curCancel context.CancelFunc
	reportMu  sync.Mutex
	lastStatus string
}

func NewServer(settings *state.Store, records RecordStore, hub *events.Hub, exec executor.Executor, ver *verify.Verifier, net *network.Detector, webDir, version string) *Server {
	return &Server{store: settings, records: records, hub: hub, exec: exec, ver: ver, net: net, webDir: webDir, version: version}
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
	mux.HandleFunc("GET /api/report", s.handleReport)
	mux.HandleFunc("GET /api/network/status", s.handleNetworkStatus)
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
