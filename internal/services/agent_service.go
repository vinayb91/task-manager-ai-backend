package services

import (
	"log"
	"os"

	"github.com/vinayb91/task-manager-ai-backend.git/internal/models"
	"github.com/vinayb91/task-manager-ai-backend.git/internal/repository"
)

// Agent Service with OpenAI
type AgentService struct {
	repo   *repository.TaskRepository
	apiKey string
	apiURL string
	tools  []models.OpenAITool
}

func NewAgentService(repo *repository.TaskRepository) *AgentService {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Println("WARNING: OPENAI_API_KEY not set. Agent will not work.")
	}

	return &AgentService{
		repo:   repo,
		apiKey: apiKey,
		apiURL: os.Getenv("OPENAI_API_URL"),
		tools:  defineOpenAITools(),
	}
}

func defineOpenAITools() []models.OpenAITool {
	return []models.OpenAITool{}
}
