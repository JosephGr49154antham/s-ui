package client

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	if err := AutoMigrate(db); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestCreateClient(t *testing.T) {
	db := setupTestDB(t)
	c := &Client{Name: "Alice", Email: "alice@example.com", InboundTag: "vless-in"}
	if err := Create(db, c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if c.ID == 0 {
		t.Error("expected non-zero ID after create")
	}
}

func TestCreateClientMissingName(t *testing.T) {
	db := setupTestDB(t)
	err := Create(db, &Client{Email: "x@x.com", InboundTag: "tag"})
	if err == nil {
		t.Error("expected error for missing name")
	}
}

func TestCreateClientMissingEmail(t *testing.T) {
	db := setupTestDB(t)
	err := Create(db, &Client{Name: "Bob", InboundTag: "tag"})
	if err == nil {
		t.Error("expected error for missing email")
	}
}

func TestGetAll(t *testing.T) {
	db := setupTestDB(t)
	Create(db, &Client{Name: "A", Email: "a@a.com", InboundTag: "t1"})
	Create(db, &Client{Name: "B", Email: "b@b.com", InboundTag: "t2"})
	clients, err := GetAll(db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(clients) != 2 {
		t.Errorf("expected 2 clients, got %d", len(clients))
	}
}

func TestGetByEmail(t *testing.T) {
	db := setupTestDB(t)
	Create(db, &Client{Name: "Carol", Email: "carol@example.com", InboundTag: "tag"})
	c, err := GetByEmail(db, "carol@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Name != "Carol" {
		t.Errorf("expected Carol, got %s", c.Name)
	}
}

func TestDelete(t *testing.T) {
	db := setupTestDB(t)
	c := &Client{Name: "Dave", Email: "dave@example.com", InboundTag: "tag"}
	Create(db, c)
	if err := Delete(db, c.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := Delete(db, c.ID); err == nil {
		t.Error("expected error deleting non-existent client")
	}
}
