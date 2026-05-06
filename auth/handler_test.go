package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/s-ui/auth"
	"github.com/s-ui/database"
)

func setupAuth(t *testing.T) {
	t.Helper()
	auth.SetSecret([]byte("handler-test-secret"))
	if err := database.InitDB(":memory:"); err != nil {
		t.Fatal(err)
	}
	if err := database.CreateUser("admin", "secret"); err != nil {
		t.Fatal(err)
	}
}

func TestLoginSuccess(t *testing.T) {
	setupAuth(t)
	h := auth.NewHandler()
	body, _ := json.Marshal(map[string]string{"username": "admin", "password": "secret"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	h.Login(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal("failed to decode response")
	}
	if resp["token"] == "" {
		t.Error("expected non-empty token")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	setupAuth(t)
	h := auth.NewHandler()
	body, _ := json.Marshal(map[string]string{"username": "admin", "password": "wrong"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	h.Login(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestLoginWrongMethod(t *testing.T) {
	setupAuth(t)
	h := auth.NewHandler()
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rr := httptest.NewRecorder()
	h.Login(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}

func TestLoginBadJSON(t *testing.T) {
	setupAuth(t)
	h := auth.NewHandler()
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte("not-json")))
	rr := httptest.NewRecorder()
	h.Login(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}
