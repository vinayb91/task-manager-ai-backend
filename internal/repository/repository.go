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
