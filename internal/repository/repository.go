package repository

import (
	"sync"

	"github.com/vinayb91/task-manager-ai-backend.git/internal/models"
)

// Task Repository (In-Memory)
type TaskRepository struct {
	tasks  []models.Task
	nextID int
	mu     sync.RWMutex
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks:  make([]models.Task, 0),
		nextID: 1,
	}
}

func (r *TaskRepository) List(status string, priority string) []models.Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filtered := make([]models.Task, 0)
	for _, task := range r.tasks {
		if status == "pending" && task.Completed {
			continue
		}
		if status == "completed" && !task.Completed {
			continue
		}
		if priority != "" && task.Priority != priority {
			continue
		}
		filtered = append(filtered, task)
	}
	return filtered
}
