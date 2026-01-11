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
	return []models.OpenAITool{
		{
			Type: "function",
			Function: models.OpenAIFunctionDef{
				Name:        "create_task",
				Description: "Creates a new task with title, optional description, priority (low/medium/high), and optional due date",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"title": map[string]interface{}{
							"type":        "string",
							"description": "Task title",
						},
						"description": map[string]interface{}{
							"type":        "string",
							"description": "Detailed description",
						},
						"priority": map[string]interface{}{
							"type":        "string",
							"description": "Task priority",
							"enum":        []string{"low", "medium", "high"},
						},
						"due_date": map[string]interface{}{
							"type":        "string",
							"description": "Due date in YYYY-MM-DD format",
						},
					},
					"required": []string{"title", "priority"},
				},
			},
		},
		{
			Type: "function",
			Function: models.OpenAIFunctionDef{
				Name:        "delete_task",
				Description: "Deletes a task by task ID",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"task_id": map[string]interface{}{
							"type":        "number",
							"description": "ID of the task to delete",
						},
					},
					"required": []string{"task_id"},
				},
			},
		},
		{
			Type: "function",
			Function: models.OpenAIFunctionDef{
				Name:        "update_task",
				Description: "Updates task properties like title, description, priority, or due date",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"task_id": map[string]interface{}{
							"type":        "number",
							"description": "ID of the task to update",
						},
						"title": map[string]interface{}{
							"type":        "string",
							"description": "New title",
						},
						"description": map[string]interface{}{
							"type":        "string",
							"description": "New description",
						},
						"priority": map[string]interface{}{
							"type":        "string",
							"description": "New priority",
							"enum":        []string{"low", "medium", "high"},
						},
						"due_date": map[string]interface{}{
							"type":        "string",
							"description": "New due date",
						},
					},
					"required": []string{"task_id"},
				},
			},
		},
		{
			Type: "function",
			Function: models.OpenAIFunctionDef{
				Name:        "list_tasks",
				Description: "Lists all tasks, optionally filtered by status (pending/completed/all) or priority",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"status": map[string]interface{}{
							"type":        "string",
							"description": "Filter by status",
							"enum":        []string{"pending", "completed", "all"},
						},
						"priority": map[string]interface{}{
							"type":        "string",
							"description": "Filter by priority",
							"enum":        []string{"low", "medium", "high"},
						},
					},
				},
			},
		},
	}
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
		finishReason := response.Choices[0].FinishReason

		if finishReason == "tool_calls" && len(message.ToolCalls) > 0 {
			messages = append(messages, message)

			for _, toolCall := range message.ToolCalls {
				log.Printf("Agent using tool: %s", toolCall.Function.Name)

				var input map[string]interface{}
				if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &input); err != nil {
					return "", fmt.Errorf("failed to parse tool arguments: %v", err)
				}

				result, err := s.ExecuteTool(toolCall.Function.Name, input)
				if err != nil {
					result = map[string]interface{}{"error": err.Error()}
				}

				resultJSON, _ := json.Marshal(result)

				messages = append(messages, models.OpenAIMessage{
					Role:       "tool",
					Content:    string(resultJSON),
					ToolCallID: toolCall.ID,
				})
			}
		} else {
			return message.Content, nil
		}
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

func (s *AgentService) ExecuteTool(toolName string, input map[string]interface{}) (interface{}, error) {
	switch toolName {
	case "create_task":
		task := models.Task{
			Title:       input["title"].(string),
			Priority:    input["priority"].(string),
			Description: getStringOrEmpty(input, "description"),
			DueDate:     getStringOrEmpty(input, "due_date"),
		}
		created := s.repo.Create(task)
		return map[string]interface{}{"success": true, "task": created}, nil

	case "delete_task":
		taskID := int(input["task_id"].(float64))
		err := s.repo.Delete(taskID)
		if err != nil {
			return map[string]interface{}{"success": false, "message": err.Error()}, nil
		}
		return map[string]interface{}{"success": true, "message": "Task deleted"}, nil

	case "update_task":
		taskID := int(input["task_id"].(float64))
		delete(input, "task_id")
		task, err := s.repo.Update(taskID, input)
		if err != nil {
			return map[string]interface{}{"success": false, "message": err.Error()}, nil
		}
		return map[string]interface{}{"success": true, "task": task}, nil

	case "list_tasks":
		status := getStringOrEmpty(input, "status")
		priority := getStringOrEmpty(input, "priority")
		tasks := s.repo.List(status, priority)
		return map[string]interface{}{"tasks": tasks, "count": len(tasks)}, nil

	default:
		return nil, fmt.Errorf("unknown tool: %s", toolName)
	}
}

func getStringOrEmpty(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}
