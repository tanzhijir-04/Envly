package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAppendAndReadRecords(t *testing.T) {
	db := New(t.TempDir())
	if err := db.AppendRecord(Record{ToolID: "git", Name: "Git", Version: "2.49.0", Method: "winget"}); err != nil {
		t.Fatal(err)
	}
	records, err := db.Records()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || records[0].ToolID != "git" {
		t.Fatalf("unexpected records: %+v", records)
	}
}

func TestRecordsEmptyWhenFileMissing(t *testing.T) {
	records, err := New(t.TempDir()).Records()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 0 {
		t.Fatalf("expected empty, got %+v", records)
	}
}

func TestEnvOpMarkRestored(t *testing.T) {
	dir := t.TempDir()
	db := New(dir)
	if err := db.AppendEnvOp(EnvOp{Key: "npm_registry", Before: "a", After: "b"}); err != nil {
		t.Fatal(err)
	}
	if err := db.MarkEnvOpRestored("npm_registry"); err != nil {
		t.Fatal(err)
	}
	ops, err := db.EnvOps()
	if err != nil {
		t.Fatal(err)
	}
	if len(ops) != 1 || !ops[0].Restored {
		t.Fatalf("expected restored env op, got %+v", ops)
	}
	// 第二次标记同名 key 不应新增
	if err := db.MarkEnvOpRestored("npm_registry"); err != nil {
		t.Fatal(err)
	}
	ops, _ = db.EnvOps()
	if len(ops) != 1 {
		t.Fatalf("expected still 1 env op, got %d", len(ops))
	}
}

func TestCorruptRecordsFileFallsBackEmpty(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "records.json"), []byte("{bad"), 0o644); err != nil {
		t.Fatal(err)
	}
	records, err := New(dir).Records()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 0 {
		t.Fatalf("expected empty, got %+v", records)
	}
}
