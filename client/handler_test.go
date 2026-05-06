package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupHandler(t *testing.T) (*Handler, *gorm.DB) {
	t.Helper()
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	AutoMigrate(db)
	return NewHandler(db), db
}

func TestHandlerList(t *testing.T) {
	h, db := setupHandler(t)
	Create(db, &Client{Name: "A", Email: "a@a.com", InboundTag: "t"})
	req := httptest.NewRequest(http.MethodGet, "/clients", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rw.Code)
	}
	var clients []Client
	json.NewDecoder(rw.Body).Decode(&clients)
	if len(clients) != 1 {
		t.Errorf("expected 1 client, got %d", len(clients))
	}
}

func TestHandlerCreate(t *testing.T) {
	h, _ := setupHandler(t)
	body, _ := json.Marshal(Client{Name: "Bob", Email: "bob@b.com", InboundTag: "tag"})
	req := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewReader(body))
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)
	if rw.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rw.Code)
	}
}

func TestHandlerCreateBadJSON(t *testing.T) {
	h, _ := setupHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/clients", bytes.NewBufferString("not-json"))
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rw.Code)
	}
}

func TestHandlerDelete(t *testing.T) {
	h, db := setupHandler(t)
	c := &Client{Name: "Del", Email: "del@d.com", InboundTag: "tag"}
	Create(db, c)
	req := httptest.NewRequest(http.MethodDelete, "/clients/1", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)
	if rw.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rw.Code)
	}
}

func TestHandlerMethodNotAllowed(t *testing.T) {
	h, _ := setupHandler(t)
	req := httptest.NewRequest(http.MethodPut, "/clients", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)
	if rw.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rw.Code)
	}
}

// TestHandlerDeleteNotFound verifies that deleting a non-existent client
// returns 404 rather than silently succeeding with 204.
func TestHandlerDeleteNotFound(t *testing.T) {
	h, _ := setupHandler(t)
	req := httptest.NewRequest(http.MethodDelete, "/clients/999", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)
	if rw.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rw.Code)
	}
}
