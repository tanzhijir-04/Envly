package state

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReturnsDefaultWhenMissing(t *testing.T) {
	store := New(t.TempDir())
	got, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.Language != "zh" || got.Region != "auto" {
		t.Fatalf("unexpected default: %+v", got)
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	store := New(dir)
	want := Settings{Language: "en", Region: "global"}
	if err := store.Save(want); err != nil {
		t.Fatal(err)
	}
	got, err := New(dir).Load()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestLoadCorruptFallsBackToDefaultAndBacksUp(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "settings.json"), []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := New(dir).Load()
	if err != nil {
		t.Fatal(err)
	}
	if got != Default() {
		t.Fatalf("got %+v, want default", got)
	}
	if _, err := os.Stat(filepath.Join(dir, "settings.json.bak")); err != nil {
		t.Fatalf("expected backup file: %v", err)
	}
}
