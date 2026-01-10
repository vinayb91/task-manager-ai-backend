package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

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

func (s *AgentService) RunAgent(userMessage string) (string, error) {
	messages := []models.OpenAIMessage{
		{Role: "system", Content: fmt.Sprintf("You are a helpful assistant. Today is %s.", time.Now().Format("2006-01-02 (Monday)"))},
	}

	messages = append(messages, models.OpenAIMessage{Role: "user", Content: userMessage})

	maxIterations := 10
	for i := 0; i < maxIterations; i++ {
		response, err := s.callOpenAI(messages)
		if err != nil {
			return "", err
		}

		if len(response.Choices) == 0 {
			return "", fmt.Errorf("no response from OpenAI")
		}

		message := response.Choices[0].Message

		return message.Content, nil

	}

	return "", fmt.Errorf("max iterations reached")
}

func (s *AgentService) callOpenAI(messages []models.OpenAIMessage) (*models.OpenAIResponse, error) {
	reqBody := models.OpenAIRequest{
		Model:    os.Getenv("OPENAI_MODEL"),
		Messages: messages,
		Tools:    s.tools,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var openAIResp models.OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return nil, err
	}

	return &openAIResp, nil
}
