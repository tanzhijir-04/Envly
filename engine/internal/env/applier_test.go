package env

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tanzhijir-04/Envly/engine/internal/store"
)

type stubRunner struct {
	results map[string]string
	calls   []string
}

func (s *stubRunner) Run(_ context.Context, name string, args ...string) (string, error) {
	key := name + " " + strings.Join(args, " ")
	s.calls = append(s.calls, key)
	if v, ok := s.results[key]; ok {
		return v, nil
	}
	return "", nil
}

func TestApplyMirrorsRecordsEnvOp(t *testing.T) {
	run := &stubRunner{results: map[string]string{
		"npm config get registry": "https://registry.npmjs.org",
	}}
	db := store.New(t.TempDir())
	a := NewApplier(run, db)
	if err := a.ApplyMirrors(context.Background(), "cn"); err != nil {
		t.Fatal(err)
	}
	if len(run.calls) < 2 {
		t.Fatalf("expected npm set/get calls, got %v", run.calls)
	}
	ops, _ := db.EnvOps()
	if len(ops) != 1 || ops[0].Key != "npm_registry" || ops[0].After == "" {
		t.Fatalf("unexpected env ops: %+v", ops)
	}
}

func TestApplyProfileWritesOnce(t *testing.T) {
	profile := filepath.Join(t.TempDir(), "profile.ps1")
	run := &stubRunner{results: map[string]string{
		"powershell -NoProfile -Command $PROFILE": profile,
	}}
	a := NewApplier(run, store.New(t.TempDir()))
	if err := a.ApplyProfile(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := a.ApplyProfile(context.Background()); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(profile)
	if strings.Count(string(b), "# >>> Envly >>>") != 1 {
		t.Fatalf("expected marker once, got:\n%s", string(b))
	}
}
