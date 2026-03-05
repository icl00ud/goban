package bot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// SessionStore persists per-chat JWT tokens to disk.
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[int64]string // chatID → JWT token
	path     string
}

func newSessionStore() (*SessionStore, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(home, ".config", "goban", "bot_sessions.json")
	ss := &SessionStore{
		sessions: make(map[int64]string),
		path:     path,
	}
	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, &ss.sessions)
	}
	return ss, nil
}

func (s *SessionStore) Get(chatID int64) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessions[chatID]
}

func (s *SessionStore) Set(chatID int64, token string) {
	s.mu.Lock()
	s.sessions[chatID] = token
	s.mu.Unlock()
	_ = s.save()
}

func (s *SessionStore) Delete(chatID int64) {
	s.mu.Lock()
	delete(s.sessions, chatID)
	s.mu.Unlock()
	_ = s.save()
}

func (s *SessionStore) save() error {
	s.mu.RLock()
	data, err := json.MarshalIndent(s.sessions, "", "  ")
	s.mu.RUnlock()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}
