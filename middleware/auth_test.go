package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/s-ui/auth"
	"github.com/s-ui/middleware"
)

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func setup() {
	auth.SetSecret([]byte("test-secret-key"))
}

func TestAuthenticateMissingHeader(t *testing.T) {
	setup()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	middleware.Authenticate(http.HandlerFunc(okHandler)).ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestAuthenticateInvalidFormat(t *testing.T) {
	setup()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Token abc")
	rr := httptest.NewRecorder()
	middleware.Authenticate(http.HandlerFunc(okHandler)).ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestAuthenticateValidToken(t *testing.T) {
	setup()
	token, err := auth.GenerateToken("admin", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	middleware.Authenticate(http.HandlerFunc(okHandler)).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestAuthenticateExpiredToken(t *testing.T) {
	setup()
	// Using -time.Minute generates a token already expired by 1 minute.
	// Increasing to -2 minutes gives a slightly wider margin in case of
	// clock skew or slow test environments.
	token, err := auth.GenerateToken("admin", -2*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	middleware.Authenticate(http.HandlerFunc(okHandler)).ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}
