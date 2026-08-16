package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/tanzhijir-04/Envly/engine/internal/api"
	"github.com/tanzhijir-04/Envly/engine/internal/events"
	"github.com/tanzhijir-04/Envly/engine/internal/executor"
	"github.com/tanzhijir-04/Envly/engine/internal/state"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:17521", "listen address")
	webDir := flag.String("web-dir", "ui", "frontend static directory (relative to cwd)")
	dataDir := flag.String("data-dir", defaultDataDir(), "settings directory")
	flag.Parse()

	store := state.New(*dataDir)
	hub := events.NewHub()
	exec := executor.Simulated{Delay: 150 * time.Millisecond}
	srv := api.NewServer(store, hub, exec, *webDir, "0.1.0")

	log.Printf("Envly engine listening on http://%s (data: %s)", *addr, *dataDir)
	if err := http.ListenAndServe(*addr, srv.Handler()); err != nil {
		log.Fatal(err)
	}
}

func defaultDataDir() string {
	if d, err := os.UserConfigDir(); err == nil {
		return filepath.Join(d, "Envly")
	}
	return ".envly"
}
