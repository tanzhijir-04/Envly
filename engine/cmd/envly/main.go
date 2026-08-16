package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/tanzhijir-04/Envly/engine/internal/api"
	"github.com/tanzhijir-04/Envly/engine/internal/env"
	"github.com/tanzhijir-04/Envly/engine/internal/events"
	"github.com/tanzhijir-04/Envly/engine/internal/executor"
	"github.com/tanzhijir-04/Envly/engine/internal/installer"
	"github.com/tanzhijir-04/Envly/engine/internal/network"
	"github.com/tanzhijir-04/Envly/engine/internal/runner"
	"github.com/tanzhijir-04/Envly/engine/internal/state"
	"github.com/tanzhijir-04/Envly/engine/internal/store"
	"github.com/tanzhijir-04/Envly/engine/internal/verify"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:17521", "listen address")
	webDir := flag.String("web-dir", "ui", "frontend static directory (relative to cwd)")
	dataDir := flag.String("data-dir", defaultDataDir(), "settings and records directory")
	simulate := flag.Bool("simulate", false, "use simulated executor (no real installs)")
	launchUI := flag.String("launch-ui", "", "path to Pake UI exe to launch after server start")
	flag.Parse()

	run := runner.OS{}
	storeDB := store.New(filepath.Join(*dataDir, "store"))
	ver := verify.New(run)
	inst := installer.New(run)
	applier := env.NewApplier(run, storeDB)
	detector := network.NewDetector(httpProbe)

	settingsStore := state.New(*dataDir)
	hub := events.NewHub()

	var exec executor.Executor
	if *simulate {
		exec = executor.Simulated{Delay: 150 * time.Millisecond}
	} else {
		exec = executor.NewReal(inst, applier, ver, storeDB,
			func() string {
				st, _ := settingsStore.Load()
				return st.Region
			},
			func() string { return "win" }, // M2 目标平台为 Windows；macOS/Linux 在 M4 引入 runtime.GOOS 分支
		)
	}

	srv := api.NewServer(settingsStore, storeDB, applier, hub, exec, ver, detector, *webDir, "0.3.0")
	log.Printf("Envly engine listening on http://%s (data: %s, simulate: %v)", *addr, *dataDir, *simulate)
	if *launchUI != "" {
		go launchAndWatchUI(*launchUI)
	}
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

func httpProbe(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return &httpError{code: resp.StatusCode}
	}
	_, err = io.Copy(io.Discard, io.LimitReader(resp.Body, 1024))
	return err
}

type httpError struct{ code int }

func (e *httpError) Error() string { return http.StatusText(e.code) }

func launchAndWatchUI(path string) {
	cmd := exec.Command(path)
	if err := cmd.Start(); err != nil {
		log.Printf("failed to launch UI: %v", err)
		return
	}
	_ = cmd.Wait()
	log.Printf("UI window closed, shutting down engine")
	os.Exit(0)
}
