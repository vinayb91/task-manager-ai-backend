package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vinayb91/task-manager-ai-backend.git/internal/models"
	"github.com/vinayb91/task-manager-ai-backend.git/internal/repository"
	"github.com/vinayb91/task-manager-ai-backend.git/internal/services"
)

type Handler struct {
	agent *services.AgentService
	repo  *repository.TaskRepository
}

func NewHandler(agent *services.AgentService, repo *repository.TaskRepository) *Handler {
	return &Handler{agent: agent, repo: repo}
}

func (h *Handler) HandleAgentQuery(w http.ResponseWriter, r *http.Request) {
	var req models.AgentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.agent.RunAgent(req.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tasks := h.repo.List("all", "")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.AgentResponse{
		Response: response,
		Tasks:    tasks,
	})
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.repo.List("all", "")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
