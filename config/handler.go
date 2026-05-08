package config

import (
	"encoding/json"
	"net/http"
)

// Handler exposes config read/write over HTTP (admin only).
type Handler struct {
	configPath string
}

func NewHandler(configPath string) *Handler {
	return &Handler{configPath: configPath}
}

func (h *Handler) GetConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cfg := GetConfig()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cfg); err != nil {
		http.Error(w, "failed to encode config", http.StatusInternalServerError)
	}
}

func (h *Handler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var incoming Config
	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	mu.Lock()
	*instance = incoming
	mu.Unlock()
	if err := SaveToFile(h.configPath); err != nil {
		http.Error(w, "failed to persist config", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// Use 200 OK with a body rather than 204 so clients can easily check the response.
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// RegisterRoutes mounts the config endpoints under the given prefix.
// Example: prefix="/api/v1" → GET /api/v1/config, POST /api/v1/config/update
func (h *Handler) RegisterRoutes(mux *http.ServeMux, prefix string) {
	mux.HandleFunc(prefix+"/config", h.GetConfig)
	mux.HandleFunc(prefix+"/config/update", h.UpdateConfig)
}
