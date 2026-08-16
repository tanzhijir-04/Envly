package installer

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/tanzhijir-04/Envly/engine/internal/config"
)

type recordingRunner struct {
	calls [][]string
	err   error
}

func (r *recordingRunner) Run(_ context.Context, name string, args ...string) (string, error) {
	r.calls = append(r.calls, append([]string{name}, args...))
	return "", r.err
}

func TestInstallDispatchesByMethod(t *testing.T) {
	run := &recordingRunner{}
	inst := New(run)
	cases := []struct {
		method string
		pkg    string
		want   []string
	}{
		{"winget", "Git.Git", []string{"winget", "install", "--id", "Git.Git", "--accept-package-agreements", "--accept-source-agreements", "--silent"}},
		{"npm", "typescript", []string{"npm", "install", "-g", "typescript"}},
		{"pip", "notebook", []string{"pip", "install", "notebook"}},
		{"rustup", "stable", []string{"rustup", "default", "stable"}},
	}
	for _, c := range cases {
		run.calls = nil
		err := inst.Install(context.Background(), config.InstallSpec{Method: c.method, Package: c.pkg})
		if err != nil {
			t.Fatal(err)
		}
		if len(run.calls) != 1 || !reflect.DeepEqual(run.calls[0], c.want) {
			t.Fatalf("method %s: got %v, want %v", c.method, run.calls, c.want)
		}
	}
}

func TestInstallPropagatesError(t *testing.T) {
	run := &recordingRunner{err: errors.New("boom")}
	inst := New(run)
	if err := inst.Install(context.Background(), config.InstallSpec{Method: "npm", Package: "x"}); err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected wrapped error, got %v", err)
	}
}

func TestInstallRejectsUnknownMethod(t *testing.T) {
	inst := New(&recordingRunner{})
	if err := inst.Install(context.Background(), config.InstallSpec{Method: "magic"}); err == nil {
		t.Fatal("expected error for unknown method")
	}
}
