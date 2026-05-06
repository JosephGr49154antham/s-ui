package config

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"
)

func TestHandlerGetConfig(t *testing.T) {
	resetSingleton()
	h := NewHandler("")
	req := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	rr := httptest.NewRecorder()
	h.GetConfig(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var cfg Config
	if err := json.NewDecoder(rr.Body).Decode(&cfg); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if cfg.ListenPort != 2095 {
		t.Errorf("expected port 2095, got %d", cfg.ListenPort)
	}
}

func TestHandlerGetConfigWrongMethod(t *testing.T) {
	resetSingleton()
	h := NewHandler("")
	req := httptest.NewRequest(http.MethodPost, "/api/config", nil)
	rr := httptest.NewRecorder()
	h.GetConfig(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}

func TestHandlerUpdateConfig(t *testing.T) {
	resetSingleton()
	instance = DefaultConfig()
	once = sync.Once{}
	GetConfig()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	h := NewHandler(path)
	update := Config{
		ListenPort: 9090,
		ListenAddr: "127.0.0.1",
		DBPath:     "./db/test.db",
		LogLevel:   "warn",
		SecretKey:  "newsecret",
		SessionTTL: 3600,
		BasePath:   "/admin/",
	}
	body, _ := json.Marshal(update)
	req := httptest.NewRequest(http.MethodPost, "/api/config/update", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	h.UpdateConfig(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	cfg := GetConfig()
	if cfg.ListenPort != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.ListenPort)
	}
	if cfg.BasePath != "/admin/" {
		t.Errorf("expected base_path /admin/, got %s", cfg.BasePath)
	}
}

func TestHandlerUpdateConfigBadJSON(t *testing.T) {
	resetSingleton()
	h := NewHandler("")
	req := httptest.NewRequest(http.MethodPost, "/api/config/update", bytes.NewBufferString("not-json"))
	rr := httptest.NewRecorder()
	h.UpdateConfig(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}
