package inbound

import (
	"s-ui/database"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	database.SetDB(db)
	if err := AutoMigrate(); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
}

func TestCreateInbound(t *testing.T) {
	setupTestDB(t)
	inbound := &Inbound{Tag: "vless-in", Protocol: "vless", Port: 443}
	if err := Create(inbound); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if inbound.ID == 0 {
		t.Error("expected inbound ID to be set")
	}
}

func TestCreateInboundMissingTag(t *testing.T) {
	setupTestDB(t)
	err := Create(&Inbound{Protocol: "vmess", Port: 8080})
	if err == nil {
		t.Error("expected error for missing tag")
	}
}

func TestCreateInboundMissingPort(t *testing.T) {
	setupTestDB(t)
	err := Create(&Inbound{Tag: "test", Protocol: "vmess"})
	if err == nil {
		t.Error("expected error for missing port")
	}
}

func TestGetAll(t *testing.T) {
	setupTestDB(t)
	Create(&Inbound{Tag: "in-1", Protocol: "vless", Port: 443})
	Create(&Inbound{Tag: "in-2", Protocol: "vmess", Port: 8080})

	inbounds, err := GetAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(inbounds) != 2 {
		t.Errorf("expected 2 inbounds, got %d", len(inbounds))
	}
}

func TestGetByTag(t *testing.T) {
	setupTestDB(t)
	Create(&Inbound{Tag: "find-me", Protocol: "trojan", Port: 9000})

	inbound, err := GetByTag("find-me")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inbound == nil {
		t.Fatal("expected inbound, got nil")
	}
	if inbound.Protocol != "trojan" {
		t.Errorf("expected protocol trojan, got %s", inbound.Protocol)
	}
}

func TestGetByTagNotFound(t *testing.T) {
	setupTestDB(t)
	inbound, err := GetByTag("ghost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inbound != nil {
		t.Error("expected nil for missing tag")
	}
}

func TestDelete(t *testing.T) {
	setupTestDB(t)
	Create(&Inbound{Tag: "delete-me", Protocol: "vless", Port: 1234})

	if err := Delete("delete-me"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	inbound, _ := GetByTag("delete-me")
	if inbound != nil {
		t.Error("expected inbound to be deleted")
	}
}

func TestDeleteNotFound(t *testing.T) {
	setupTestDB(t)
	if err := Delete("nonexistent"); err == nil {
		t.Error("expected error when deleting nonexistent inbound")
	}
}
