package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Settings struct {
	Language string `json:"language"` // zh | en
	Region   string `json:"region"`   // auto | cn | global
}

func Default() Settings {
	return Settings{Language: "zh", Region: "auto"}
}

type Store struct {
	mu   sync.Mutex
	path string
}

func New(dir string) *Store {
	return &Store{path: filepath.Join(dir, "settings.json")}
}

func (s *Store) Load() (Settings, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, err := os.ReadFile(s.path)
	if err != nil {
		return Default(), nil
	}
	var st Settings
	if err := json.Unmarshal(b, &st); err != nil {
		_ = os.Rename(s.path, s.path+".bak")
		return Default(), nil
	}
	if st.Language == "" {
		st.Language = "zh"
	}
	if st.Region == "" {
		st.Region = "auto"
	}
	return st, nil
}

func (s *Store) Save(st Settings) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o644)
}
