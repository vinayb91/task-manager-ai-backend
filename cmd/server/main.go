package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/vinayb91/task-manager-ai-backend.git/internal/handlers"
	"github.com/vinayb91/task-manager-ai-backend.git/internal/repository"
	"github.com/vinayb91/task-manager-ai-backend.git/internal/services"
)

var Version = "1.0.0"

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()
	repo := repository.NewTaskRepository()
	agent := services.NewAgentService(repo)
	handler := handlers.NewHandler(agent, repo)

	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"version": Version,
		})
	}).Methods("GET")

	r.HandleFunc("/api/agent/query", handler.HandleAgentQuery).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/tasks", handler.GetTasks).Methods("GET", "OPTIONS")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s (Version: %s)", port, Version)
	log.Fatal(http.ListenAndServe(":"+port, enableCORS(r)))
}
