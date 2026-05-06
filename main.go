package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/s-ui/s-ui/auth"
	"github.com/s-ui/s-ui/client"
	"github.com/s-ui/s-ui/config"
	"github.com/s-ui/s-ui/database"
	"github.com/s-ui/s-ui/inbound"
	"github.com/s-ui/s-ui/router"
)

func main() {
	// Parse command-line flags
	configPath := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadFromFile(*configPath)
	if err != nil {
		log.Printf("Could not load config from %s, using defaults: %v", *configPath, err)
		cfg = config.DefaultConfig()
	}

	// Initialize the database
	db, err := database.InitDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run auto-migrations for all models
	if err := inbound.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate inbound schema: %v", err)
	}
	if err := client.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate client schema: %v", err)
	}

	// Seed default admin user if none exists
	if err := database.CreateUser(db, cfg.AdminUsername, cfg.AdminPassword); err != nil {
		// User likely already exists; not fatal
		log.Printf("Admin user setup: %v", err)
	}

	// Configure JWT secret
	auth.SetSecret(cfg.JWTSecret)

	// Build the router with all handlers wired up
	r := router.New(db, cfg)

	addr := fmt.Sprintf("%s:%d", cfg.ListenAddr, cfg.Port)
	log.Printf("s-ui starting on %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
