package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/tanzhijir-04/Envly/engine/internal/config"
	"github.com/tanzhijir-04/Envly/engine/internal/events"
	"github.com/tanzhijir-04/Envly/engine/internal/executor"
	"github.com/tanzhijir-04/Envly/engine/internal/state"
)

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "version": s.version})
}

func (s *Server) handleCatalog(w http.ResponseWriter, r *http.Request) {
	platform := r.URL.Query().Get("platform")
	if platform == "" {
		platform = platformName()
	}
	type toolView struct {
		ID        string `json:"id"`
		NameZh    string `json:"name_zh"`
		NameEn    string `json:"name_en"`
		DescZh    string `json:"desc_zh"`
		DescEn    string `json:"desc_en"`
		GroupID   string `json:"group_id"`
		Optional  bool   `json:"optional"`
		Method    string `json:"method"`
		Package   string `json:"package"`
		VerifyCmd string `json:"verify_cmd,omitempty"`
	}
	tools := make([]toolView, 0, len(config.Tools))
	for _, t := range config.Tools {
		spec, ok := t.Install[platform]
		if !ok {
			continue
		}
		tools = append(tools, toolView{t.ID, t.NameZh, t.NameEn, t.DescZh, t.DescEn, t.GroupID, t.Optional, spec.Method, spec.Package, spec.VerifyCmd})
	}
	writeJSON(w, http.StatusOK, map[string]any{"groups": config.Groups, "tools": tools, "platform": platform})
}

func (s *Server) handleTemplates(w http.ResponseWriter, r *http.Request) {
	type templateView struct {
		ID      string   `json:"id"`
		NameZh  string   `json:"name_zh"`
		NameEn  string   `json:"name_en"`
		DescZh  string   `json:"desc_zh"`
		DescEn  string   `json:"desc_en"`
		Count   int      `json:"count"`
		ToolIDs []string `json:"tool_ids"`
	}
	views := make([]templateView, 0, len(config.Templates))
	for _, t := range config.Templates {
		ids := config.TemplateToolIDs(t)
		views = append(views, templateView{t.ID, t.NameZh, t.NameEn, t.DescZh, t.DescEn, len(ids), ids})
	}
	writeJSON(w, http.StatusOK, views)
}

type planRequest struct {
	ToolIDs []string `json:"tool_ids"`
}

func (s *Server) handlePlan(w http.ResponseWriter, r *http.Request) {
	var req planRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	type item struct {
		ToolID string `json:"tool_id"`
		NameZh string `json:"name_zh"`
		NameEn string `json:"name_en"`
		Method string `json:"method"`
		Status string `json:"status"`
	}
	items := make([]item, 0, len(req.ToolIDs))
	for _, id := range req.ToolIDs {
		t, ok := config.ToolByID(id)
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "unknown tool: " + id})
			return
		}
		spec, ok := t.Install[platformName()]
		method := "unsupported"
		if ok {
			method = spec.Method
		}
		items = append(items, item{id, t.NameZh, t.NameEn, method, "pending"})
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

type runRequest struct {
	ToolIDs []string `json:"tool_ids"`
}

func (s *Server) handleRun(w http.ResponseWriter, r *http.Request) {
	var req runRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	if len(req.ToolIDs) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "empty tool_ids"})
		return
	}
	runID := newRunID()
	ctx, cancel := context.WithCancel(context.Background())
	s.runMu.Lock()
	if s.curRunID != "" {
		s.runMu.Unlock()
		cancel()
		writeJSON(w, http.StatusConflict, map[string]string{"error": "run already active"})
		return
	}
	s.curRunID = runID
	s.curCancel = cancel
	s.runMu.Unlock()

	go func() {
		err := s.exec.Run(ctx, req.ToolIDs, func(p executor.Progress) {
			s.hub.Publish(events.Event{Type: p.Status, RunID: runID, ToolID: p.ToolID, MessageKey: p.MsgKey, Params: p.Params})
		})
		status := "success"
		if err != nil {
			if errors.Is(err, context.Canceled) {
				status = "cancelled"
			} else {
				status = "failed"
			}
		}
		s.hub.Publish(events.Event{Type: "run_done", RunID: runID, Status: status, MessageKey: "run.done"})
		s.runMu.Lock()
		if s.curRunID == runID {
			s.curRunID = ""
			s.curCancel = nil
		}
		s.runMu.Unlock()
	}()
	writeJSON(w, http.StatusAccepted, map[string]any{"run_id": runID})
}

func (s *Server) handleCancel(w http.ResponseWriter, r *http.Request) {
	s.runMu.Lock()
	cancel := s.curCancel
	s.runMu.Unlock()
	if cancel == nil {
		writeJSON(w, http.StatusOK, map[string]any{"ok": false, "reason": "no active run"})
		return
	}
	cancel()
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "streaming unsupported"})
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch, unsubscribe := s.hub.Subscribe()
	defer unsubscribe()

	ctx := r.Context()
	runID := r.URL.Query().Get("run_id")
	for {
		select {
		case <-ctx.Done():
			return
		case e, open := <-ch:
			if !open {
				return
			}
			if runID != "" && e.RunID != runID {
				continue
			}
			data, _ := json.Marshal(e)
			_, _ = w.Write([]byte("data: " + string(data) + "\n\n"))
			flusher.Flush()
		}
	}
}

func (s *Server) handleGetSettings(w http.ResponseWriter, r *http.Request) {
	st, _ := s.store.Load()
	writeJSON(w, http.StatusOK, st)
}

func (s *Server) handlePostSettings(w http.ResponseWriter, r *http.Request) {
	var st state.Settings
	if err := json.NewDecoder(r.Body).Decode(&st); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	if st.Language != "zh" && st.Language != "en" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "language must be zh or en"})
		return
	}
	if st.Region != "auto" && st.Region != "cn" && st.Region != "global" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "region must be auto, cn or global"})
		return
	}
	if err := s.store.Save(st); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, st)
}

func newRunID() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
