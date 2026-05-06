package config

import (
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	ListenPort  int    `json:"listen_port"`
	ListenAddr  string `json:"listen_addr"`
	DBPath      string `json:"db_path"`
	LogLevel    string `json:"log_level"`
	SecretKey   string `json:"secret_key"`
	SessionTTL  int    `json:"session_ttl"`
	BasePath    string `json:"base_path"`
}

var (
	instance *Config
	once     sync.Once
	mu       sync.RWMutex
)

func DefaultConfig() *Config {
	return &Config{
		ListenPort: 2095,
		ListenAddr: "0.0.0.0",
		DBPath:     "./db/s-ui.db",
		LogLevel:   "info",
		SecretKey:  "s-ui-secret",
		SessionTTL: 86400,
		BasePath:   "/",
	}
}

func GetConfig() *Config {
	once.Do(func() {
		instance = DefaultConfig()
	})
	mu.RLock()
	defer mu.RUnlock()
	return instance
}

func LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		instance = DefaultConfig()
	}
	return json.Unmarshal(data, instance)
}

func SaveToFile(path string) error {
	mu.RLock()
	defer mu.RUnlock()
	data, err := json.MarshalIndent(instance, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
