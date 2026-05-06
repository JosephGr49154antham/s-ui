package config

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func resetSingleton() {
	instance = nil
	once = sync.Once{}
}

func TestDefaultConfig(t *testing.T) {
	resetSingleton()
	cfg := DefaultConfig()
	if cfg.ListenPort != 2095 {
		t.Errorf("expected listen_port 2095, got %d", cfg.ListenPort)
	}
	if cfg.ListenAddr != "0.0.0.0" {
		t.Errorf("expected listen_addr 0.0.0.0, got %s", cfg.ListenAddr)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected log_level info, got %s", cfg.LogLevel)
	}
	if cfg.SessionTTL != 86400 {
		t.Errorf("expected session_ttl 86400, got %d", cfg.SessionTTL)
	}
}

func TestGetConfig(t *testing.T) {
	resetSingleton()
	cfg := GetConfig()
	if cfg == nil {
		t.Fatal("GetConfig returned nil")
	}
	if cfg.ListenPort != 2095 {
		t.Errorf("expected default port 2095, got %d", cfg.ListenPort)
	}
}

func TestLoadFromFile(t *testing.T) {
	resetSingleton()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{"listen_port":8080,"log_level":"debug","base_path":"/ui/"}`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	if err := LoadFromFile(path); err != nil {
		t.Fatalf("LoadFromFile error: %v", err)
	}
	cfg := GetConfig()
	if cfg.ListenPort != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.ListenPort)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("expected log_level debug, got %s", cfg.LogLevel)
	}
	if cfg.BasePath != "/ui/" {
		t.Errorf("expected base_path /ui/, got %s", cfg.BasePath)
	}
}

func TestLoadFromFileMissing(t *testing.T) {
	resetSingleton()
	err := LoadFromFile("/nonexistent/path/config.json")
	if err != nil {
		t.Errorf("expected nil error for missing file, got %v", err)
	}
}

func TestSaveToFile(t *testing.T) {
	resetSingleton()
	dir := t.TempDir()
	path := filepath.Join(dir, "config_out.json")
	GetConfig()
	if err := SaveToFile(path); err != nil {
		t.Fatalf("SaveToFile error: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("config file was not created")
	}
}
