package auth

import (
	"testing"
	"time"
)

func TestGenerateAndParseToken(t *testing.T) {
	SetSecret("test-secret")

	token, err := GenerateToken("admin", time.Hour)
	if err != nil {
		t.Fatalf("expected no error generating token, got: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("expected no error parsing token, got: %v", err)
	}
	if claims.Username != "admin" {
		t.Errorf("expected username 'admin', got '%s'", claims.Username)
	}
}

func TestExpiredToken(t *testing.T) {
	SetSecret("test-secret")

	token, err := GenerateToken("admin", -time.Second)
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	_, err = ParseToken(token)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
	if err != ErrExpiredToken {
		t.Errorf("expected ErrExpiredToken, got: %v", err)
	}
}

func TestInvalidToken(t *testing.T) {
	SetSecret("test-secret")

	_, err := ParseToken("this.is.not.a.valid.token")
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
	if err != ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got: %v", err)
	}
}

func TestSecretIsolation(t *testing.T) {
	SetSecret("secret-a")
	tokenA, err := GenerateToken("user", time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	SetSecret("secret-b")
	_, err = ParseToken(tokenA)
	if err == nil {
		t.Fatal("expected error when parsing token signed with different secret")
	}
}
