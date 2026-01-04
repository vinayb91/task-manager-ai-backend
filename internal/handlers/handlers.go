package handlers

import (
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
