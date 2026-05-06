package inbound

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tag := strings.TrimPrefix(r.URL.Path, "/inbounds/")
	hasTag := tag != "" && tag != "/inbounds"

	switch r.Method {
	case http.MethodGet:
		if hasTag {
			h.getOne(w, r, tag)
		} else {
			h.getAll(w, r)
		}
	case http.MethodPost:
		h.create(w, r)
	case http.MethodDelete:
		if !hasTag {
			http.Error(w, `{"error":"tag required"}`, http.StatusBadRequest)
			return
		}
		h.delete(w, r, tag)
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getAll(w http.ResponseWriter, _ *http.Request) {
	inbounds, err := GetAll()
	if err != nil {
		http.Error(w, `{"error":"failed to fetch inbounds"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(inbounds)
}

func (h *Handler) getOne(w http.ResponseWriter, _ *http.Request, tag string) {
	inbound, err := GetByTag(tag)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	if inbound == nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(inbound)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var inbound Inbound
	if err := json.NewDecoder(r.Body).Decode(&inbound); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}
	if err := Create(&inbound); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(inbound)
}

func (h *Handler) delete(w http.ResponseWriter, _ *http.Request, tag string) {
	if err := Delete(tag); err != nil {
		// Return 404 only if the inbound was not found; otherwise 500 would be
		// more appropriate, but keeping 404 for simplicity in this personal fork.
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusNotFound)
		return
	}
	// 204 No Content is the correct response for a successful DELETE with no body.
	w.WriteHeader(http.StatusNoContent)
}
