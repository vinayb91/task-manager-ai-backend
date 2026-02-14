package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/vinayb91/task-manager-ai-backend.git/internal/database"
	"github.com/vinayb91/task-manager-ai-backend.git/internal/handlers"
	"github.com/vinayb91/task-manager-ai-backend.git/internal/middlewares"
	"github.com/vinayb91/task-manager-ai-backend.git/internal/repository"
	"github.com/vinayb91/task-manager-ai-backend.git/internal/services"
)

var Version = "1.0.0"

func main() {
	// godotenv.Load()
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		log.Fatal("JWT_SECRET environment variable must be set")
	}

	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()
	taskRepo := repository.NewTaskRepository()
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

	r.HandleFunc("/api/agent/query", handler.HandleAgentQuery).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/tasks", handler.GetTasks).Methods("GET", "OPTIONS")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s (Version: %s)", port, Version)
	log.Fatal(http.ListenAndServe(":"+port, middlewares.EnableCORS(r)))
}
