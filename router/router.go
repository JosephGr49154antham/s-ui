package router

import (
	"net/http"

	"github.com/s-ui/auth"
	"github.com/s-ui/config"
	"github.com/s-ui/database"
	"github.com/s-ui/middleware"
)

// New builds and returns the application HTTP mux.
func New() http.Handler {
	mux := http.NewServeMux()

	// Public routes
	authHandler := auth.NewHandler()
	mux.HandleFunc("/api/login", authHandler.Login)

	// Protected routes
	cfgHandler := config.NewHandler()
	dbHandler := database.NewHandler()

	protected := http.NewServeMux()
	protected.HandleFunc("/api/config", cfgHandler.ServeHTTP)
	protected.HandleFunc("/api/users", dbHandler.ServeHTTP)

	mux.Handle("/api/config", middleware.Authenticate(protected))
	mux.Handle("/api/users", middleware.Authenticate(protected))

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	return mux
}
