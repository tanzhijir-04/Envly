package executor

import (
	"context"
	"errors"
	"testing"

	"github.com/tanzhijir-04/Envly/engine/internal/config"
	"github.com/tanzhijir-04/Envly/engine/internal/store"
)

type fakeInstaller struct {
	installed []string
	err       error
}

func (f *fakeInstaller) Install(_ context.Context, spec config.InstallSpec) error {
	if f.err != nil {
		return f.err
	}
	f.installed = append(f.installed, spec.Package)
	return nil
}

type fakeEnv struct {
	mirrorsRegion string
	err           error
}

func (f *fakeEnv) ApplyMirrors(_ context.Context, region string) error { f.mirrorsRegion = region; return f.err }
func (f *fakeEnv) ApplyPipMirror(_ context.Context, region string) error { return f.err }
func (f *fakeEnv) ApplyProxy(context.Context) error { return f.err }
func (f *fakeEnv) ApplyProfile(context.Context) error { return f.err }
func (f *fakeEnv) ApplyPath(context.Context, string) error { return f.err }

type fakeVerifier struct {
	installed map[string]bool
}

func (f *fakeVerifier) Check(_ context.Context, cmdline string) (string, bool) {
	if f.installed[cmdline] {
		return "9.9.9", true
	}
	return "", false
}

func newReal(t *testing.T, inst InstallRunner, env EnvApplier, ver Verifier) *Real {
	t.Helper()
	return NewReal(inst, env, ver, store.New(t.TempDir()), func() string { return "cn" }, func() string { return "win" })
}

func TestRealInstallsAndRecords(t *testing.T) {
	inst := &fakeInstaller{}
	r := newReal(t, inst, &fakeEnv{}, &fakeVerifier{})
	var events []Progress
	err := r.Run(context.Background(), []string{"nodejs"}, func(p Progress) { events = append(events, p) })
	if err != nil {
		t.Fatal(err)
	}
	if len(inst.installed) != 1 || inst.installed[0] != "OpenJS.NodeJS.LTS" {
		t.Fatalf("installed %v", inst.installed)
	}
	last := events[len(events)-1]
	if last.Status != "success" {
		t.Fatalf("expected success, got %+v", last)
	}
}

func TestRealSkipsAlreadyInstalled(t *testing.T) {
	inst := &fakeInstaller{}
	r := newReal(t, inst, &fakeEnv{}, &fakeVerifier{installed: map[string]bool{"node -v": true}})
	events := runCollect(t, r, []string{"nodejs"})
	if len(inst.installed) != 0 {
		t.Fatalf("should not install, got %v", inst.installed)
	}
	if events[len(events)-1].Status != "skipped" {
		t.Fatalf("expected skipped, got %+v", events[len(events)-1])
	}
}

func TestRealSkipsMirrorsOutsideCN(t *testing.T) {
	r := NewReal(&fakeInstaller{}, &fakeEnv{}, &fakeVerifier{}, store.New(t.TempDir()), func() string { return "global" }, func() string { return "win" })
	events := runCollect(t, r, []string{"env-mirrors"})
	if events[len(events)-1].Status != "skipped" {
		t.Fatalf("expected skipped outside CN, got %+v", events[len(events)-1])
	}
}

func TestRealAppliesMirrorsInCN(t *testing.T) {
	env := &fakeEnv{}
	r := newReal(t, &fakeInstaller{}, env, &fakeVerifier{})
	runCollect(t, r, []string{"env-mirrors"})
	if env.mirrorsRegion != "cn" {
		t.Fatalf("expected cn, got %q", env.mirrorsRegion)
	}
}

func TestRealFailedInstallContinues(t *testing.T) {
	inst := &fakeInstaller{err: errors.New("boom")}
	r := newReal(t, inst, &fakeEnv{}, &fakeVerifier{})
	events := runCollect(t, r, []string{"nodejs", "git"})
	if events[len(events)-1].Status != "failed" {
		t.Fatalf("expected failed, got %+v", events[len(events)-1])
	}
	if len(events) < 2 {
		t.Fatal("expected both tools to emit events")
	}
}

func runCollect(t *testing.T, r *Real, ids []string) []Progress {
	t.Helper()
	var events []Progress
	if err := r.Run(context.Background(), ids, func(p Progress) { events = append(events, p) }); err != nil {
		t.Fatal(err)
	}
	return events
}

func TestRealUnknownToolEmittedFailed(t *testing.T) {
	r := newReal(t, &fakeInstaller{}, &fakeEnv{}, &fakeVerifier{})
	events := runCollect(t, r, []string{"ghost"})
	if events[len(events)-1].Status != "failed" {
		t.Fatalf("expected failed, got %+v", events[len(events)-1])
	}
}
