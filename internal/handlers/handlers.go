package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vinayb91/task-manager-ai-backend.git/internal/models"
	"github.com/vinayb91/task-manager-ai-backend.git/internal/repository"
	"github.com/vinayb91/task-manager-ai-backend.git/internal/services"
)

type Handler struct {
	authService  *services.AuthService
	agentService *services.AgentService
	taskRepo     *repository.TaskRepository
}

func NewHandler(authService *services.AuthService, agentService *services.AgentService, taskRepo *repository.TaskRepository) *Handler {
	return &Handler{
		authService:  authService,
		agentService: agentService,
		taskRepo:     taskRepo,
	}
}

func (h *Handler) HandleAgentQuery(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	var req models.AgentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.agentService.RunAgent(userID, req.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tasks, _ := h.taskRepo.List(userID, "all", "")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": response,
		"tasks":    tasks,
	})
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	tasks, err := h.taskRepo.List(userID, "all", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.authService.Register(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.authService.Login(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
