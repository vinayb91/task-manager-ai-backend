package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/vinayb91/task-manager-ai-backend.git/internal/database"
	"github.com/vinayb91/task-manager-ai-backend.git/internal/handlers"
	"github.com/vinayb91/task-manager-ai-backend.git/internal/middlewares"
	"github.com/vinayb91/task-manager-ai-backend.git/internal/repository"
	"github.com/vinayb91/task-manager-ai-backend.git/internal/services"
)

var Version = "2.0.0"

func main() {

	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		log.Fatal("JWT_SECRET environment variable must be set")
	}

	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()
	taskRepo := repository.NewTaskRepository(db)
	userRepo := repository.NewUserRepository(db)

	agentService := services.NewAgentService(taskRepo)
	authService := services.NewAuthService(userRepo, jwtSecret)
	handler := handlers.NewHandler(authService, agentService, taskRepo)

	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"version": Version,
		})
	}).Methods("GET")

	// public routes
	r.HandleFunc("/api/auth/register", handler.Register).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/auth/login", handler.Login).Methods("POST", "OPTIONS")

	// Protected routes
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middlewares.AuthMiddleware(authService))
	protected.HandleFunc("/agent/query", handler.HandleAgentQuery).Methods("POST", "OPTIONS")
	protected.HandleFunc("/tasks", handler.GetTasks).Methods("GET", "OPTIONS")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: middlewares.RequestLogger(middlewares.EnableCORS(r)),
	}

	go func() {
		log.Printf("Server starting on port %s (Version: %s)", port, Version)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
