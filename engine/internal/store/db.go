// Package store 持久化安装记录与环境操作日志。
package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Record struct {
	ToolID  string `json:"tool_id"`
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
	Method  string `json:"method"`
	Time    string `json:"time"`
}

type EnvOp struct {
	Key      string `json:"key"`
	Before   string `json:"before,omitempty"`
	After    string `json:"after,omitempty"`
	Restored bool   `json:"restored"`
}

type DB struct {
	mu  sync.Mutex
	dir string
}

func New(dir string) *DB {
	return &DB{dir: dir}
}

func (d *DB) AppendRecord(r Record) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	records, err := loadJSON[Record](filepath.Join(d.dir, "records.json"))
	if err != nil {
		return err
	}
	if r.Time == "" {
		r.Time = time.Now().Format(time.RFC3339)
	}
	records = append(records, r)
	return saveJSON(filepath.Join(d.dir, "records.json"), records)
}

func (d *DB) Records() ([]Record, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return loadJSON[Record](filepath.Join(d.dir, "records.json"))
}

func (d *DB) AppendEnvOp(op EnvOp) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	ops, err := loadJSON[EnvOp](filepath.Join(d.dir, "envops.json"))
	if err != nil {
		return err
	}
	ops = append(ops, op)
	return saveJSON(filepath.Join(d.dir, "envops.json"), ops)
}

func (d *DB) EnvOps() ([]EnvOp, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return loadJSON[EnvOp](filepath.Join(d.dir, "envops.json"))
}

func (d *DB) MarkEnvOpRestored(key string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	ops, err := loadJSON[EnvOp](filepath.Join(d.dir, "envops.json"))
	if err != nil {
		return err
	}
	changed := false
	for i := range ops {
		if ops[i].Key == key && !ops[i].Restored {
			ops[i].Restored = true
			changed = true
		}
	}
	if !changed {
		return nil
	}
	return saveJSON(filepath.Join(d.dir, "envops.json"), ops)
}

func loadJSON[T any](path string) ([]T, error) {
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, nil // 损坏文件按空处理，不阻塞主流程
	}
	var out []T
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, nil
	}
	return out, nil
}

func saveJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
