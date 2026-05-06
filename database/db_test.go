package database

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestDB(t *testing.T) string {
	t.Helper()
	resetSingleton()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")
	if err := InitDB(path); err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	return path
}

func TestInitDB(t *testing.T) {
	setupTestDB(t)
	if GetDB() == nil {
		t.Fatal("expected non-nil db after InitDB")
	}
}

func TestInitDBCreatesDirectory(t *testing.T) {
	resetSingleton()
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "nested", "test.db")
	if err := InitDB(path); err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		t.Fatal("expected directories to be created")
	}
}

func TestCreateAndGetUser(t *testing.T) {
	setupTestDB(t)

	u, err := CreateUser("alice", "hashed_pw")
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if u.ID == 0 {
		t.Fatal("expected non-zero ID")
	}

	found, err := GetUserByUsername("alice")
	if err != nil {
		t.Fatalf("GetUserByUsername: %v", err)
	}
	if found == nil || found.Username != "alice" {
		t.Fatal("expected to find user alice")
	}
}

func TestCreateUserEmptyUsername(t *testing.T) {
	setupTestDB(t)
	_, err := CreateUser("", "pw")
	if err == nil {
		t.Fatal("expected error for empty username")
	}
}

func TestGetUserByUsernameNotFound(t *testing.T) {
	setupTestDB(t)
	u, err := GetUserByUsername("nobody")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u != nil {
		t.Fatal("expected nil for missing user")
	}
}

func TestDeleteUser(t *testing.T) {
	setupTestDB(t)
	u, _ := CreateUser("bob", "pw")
	if err := DeleteUser(u.ID); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
	found, _ := GetUserByUsername("bob")
	if found != nil {
		t.Fatal("expected user to be soft-deleted")
	}
}
