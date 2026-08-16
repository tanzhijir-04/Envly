package api

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/tanzhijir-04/Envly/engine/internal/events"
	"github.com/tanzhijir-04/Envly/engine/internal/executor"
	"github.com/tanzhijir-04/Envly/engine/internal/network"
	"github.com/tanzhijir-04/Envly/engine/internal/state"
	"github.com/tanzhijir-04/Envly/engine/internal/store"
	"github.com/tanzhijir-04/Envly/engine/internal/verify"
)

func newTestServer(t *testing.T, exec executor.Executor) *Server {
	return NewServer(state.New(t.TempDir()), store.New(t.TempDir()), nil, events.NewHub(), exec, nil, nil, "", "test")
}

func doReq(t *testing.T, srv *Server, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	var reader *bytes.Reader
	if body == "" {
		reader = bytes.NewReader(nil)
	} else {
		reader = bytes.NewReader([]byte(body))
	}
	req := httptest.NewRequest(method, path, reader)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)
	return rec
}

func TestHealth(t *testing.T) {
	rec := doReq(t, newTestServer(t, executor.Simulated{}), http.MethodGet, "/api/health", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"ok":true`) {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestCatalogFiltersByPlatform(t *testing.T) {
	srv := newTestServer(t, executor.Simulated{})
	rec := doReq(t, srv, http.MethodGet, "/api/catalog?platform=mac", "")
	if strings.Contains(rec.Body.String(), "mingw") {
		t.Fatal("mingw should not appear on macOS")
	}
	rec = doReq(t, srv, http.MethodGet, "/api/catalog?platform=win", "")
	if !strings.Contains(rec.Body.String(), "mingw") {
		t.Fatal("mingw should appear on Windows")
	}
}

func TestTemplatesReturnsFour(t *testing.T) {
	rec := doReq(t, newTestServer(t, executor.Simulated{}), http.MethodGet, "/api/templates", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	var templates []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &templates); err != nil {
		t.Fatal(err)
	}
	if len(templates) != 4 {
		t.Fatalf("expected 4 templates, got %d", len(templates))
	}
}

func TestPlanRejectsUnknownTool(t *testing.T) {
	rec := doReq(t, newTestServer(t, executor.Simulated{}), http.MethodPost, "/api/plan", `{"tool_ids":["nope"]}`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status %d", rec.Code)
	}
}

func TestRunEmitsRunDoneAndCancelWorks(t *testing.T) {
	srv := newTestServer(t, executor.Simulated{Delay: time.Millisecond})
	ch, unsubscribe := srv.hub.Subscribe()
	defer unsubscribe()
	rec := doReq(t, srv, http.MethodPost, "/api/run", `{"tool_ids":["nodejs"]}`)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("status %d, body %s", rec.Code, rec.Body.String())
	}
	var runResp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &runResp); err != nil {
		t.Fatal(err)
	}
	if runResp["run_id"] == "" {
		t.Fatal("expected run_id")
	}
	waitRunDone(t, ch, "success")

	rec = doReq(t, srv, http.MethodPost, "/api/run", `{"tool_ids":["nodejs","git"]}`)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("status %d", rec.Code)
	}
	rec = doReq(t, srv, http.MethodPost, "/api/cancel", "")
	if !strings.Contains(rec.Body.String(), `"ok":true`) {
		t.Fatalf("cancel failed: %s", rec.Body.String())
	}
	waitRunDone(t, ch, "cancelled")
}

func TestEventsSSESendsHistory(t *testing.T) {
	srv := newTestServer(t, executor.Simulated{})
	srv.hub.Publish(events.Event{Type: "running", RunID: "r1", ToolID: "nodejs", MessageKey: "tool.start"})
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+"/api/events", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	if !scanner.Scan() {
		t.Fatal("expected an SSE line")
	}
	line := scanner.Text()
	if !strings.Contains(line, "tool.start") {
		t.Fatalf("unexpected SSE line: %s", line)
	}
}

func waitRunDone(t *testing.T, ch <-chan events.Event, wantStatus string) {
	t.Helper()
	deadline := time.After(3 * time.Second)
	for {
		select {
		case e := <-ch:
			if e.Type == "run_done" {
				if e.Status != wantStatus {
					t.Fatalf("run_done status %q, want %q", e.Status, wantStatus)
				}
				return
			}
		case <-deadline:
			t.Fatalf("timed out waiting for run_done %q", wantStatus)
		}
	}
}

type fakeVerifyRunner struct {
	ok      bool
	version string
}

func (f *fakeVerifyRunner) Run(context.Context, string, ...string) (string, error) {
	if f.ok {
		return f.version, nil
	}
	return "", errors.New("not found")
}

func TestPlanMarksInstalled(t *testing.T) {
	run := &fakeVerifyRunner{ok: true, version: "22.14.0"}
	ver := verify.New(run)
	srv := NewServer(state.New(t.TempDir()), store.New(t.TempDir()), nil, events.NewHub(), executor.Simulated{}, ver, nil, "", "test")
	rec := doReq(t, srv, http.MethodPost, "/api/plan", `{"tool_ids":["nodejs"]}`)
	if !strings.Contains(rec.Body.String(), `"status":"installed"`) {
		t.Fatalf("expected installed status, got %s", rec.Body.String())
	}
}

func TestRunRejectsUnknownTool(t *testing.T) {
	srv := newTestServer(t, executor.Simulated{})
	rec := doReq(t, srv, http.MethodPost, "/api/run", `{"tool_ids":["ghost"]}`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestNetworkStatusReturnsRegion(t *testing.T) {
	detector := network.NewDetector(func(context.Context, string) error { return nil })
	srv := NewServer(state.New(t.TempDir()), store.New(t.TempDir()), nil, events.NewHub(), executor.Simulated{}, nil, detector, "", "test")
	rec := doReq(t, srv, http.MethodGet, "/api/network/status", "")
	if !strings.Contains(rec.Body.String(), `"region":"global"`) {
		t.Fatalf("expected global region, got %s", rec.Body.String())
	}
}

func TestReportReturnsLastStatus(t *testing.T) {
	srv := newTestServer(t, executor.Simulated{Delay: time.Millisecond})
	ch, unsubscribe := srv.hub.Subscribe()
	defer unsubscribe()
	rec := doReq(t, srv, http.MethodPost, "/api/run", `{"tool_ids":["nodejs"]}`)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("status %d", rec.Code)
	}
	waitRunDone(t, ch, "success")
	rec = doReq(t, srv, http.MethodGet, "/api/report", "")
	if !strings.Contains(rec.Body.String(), `"status":"success"`) {
		t.Fatalf("expected success status, got %s", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"records"`) {
		t.Fatalf("expected records field, got %s", rec.Body.String())
	}
}

type fakeRestorer struct {
	called bool
	err    error
}

func (f *fakeRestorer) Restore(context.Context) error {
	f.called = true
	return f.err
}

func TestRestoreEnvCallsRestorer(t *testing.T) {
	restorer := &fakeRestorer{}
	srv := NewServer(state.New(t.TempDir()), store.New(t.TempDir()), restorer, events.NewHub(), executor.Simulated{}, nil, nil, "", "test")
	rec := doReq(t, srv, http.MethodPost, "/api/settings/restore-env", "")
	if rec.Code != http.StatusOK || !restorer.called {
		t.Fatalf("status %d, called %v, body %s", rec.Code, restorer.called, rec.Body.String())
	}
}

type fakeWindowController struct {
	action string
	err    error
}

func (f *fakeWindowController) Action(action string) error {
	f.action = action
	return f.err
}

func TestWindowActionCallsController(t *testing.T) {
	fc := &fakeWindowController{}
	srv := newTestServer(t, executor.Simulated{})
	srv.SetWindowController(fc)
	rec := doReq(t, srv, http.MethodPost, "/api/window/action", `{"action":"minimize"}`)
	if rec.Code != http.StatusOK || fc.action != "minimize" {
		t.Fatalf("status %d action %q body %s", rec.Code, fc.action, rec.Body.String())
	}
}
