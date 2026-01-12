package repository

import (
	"fmt"
	"strings"
	"sync"
	"time"

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

func (r *TaskRepository) Create(task models.Task) models.Task {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = r.nextID
	task.CreatedAt = time.Now()
	r.tasks = append(r.tasks, task)
	r.nextID++
	return task
}

func (r *TaskRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, task := range r.tasks {
		if task.ID == id {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("task not found")
}

func (r *TaskRepository) Update(id int, updates map[string]interface{}) (models.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, task := range r.tasks {
		if task.ID == id {
			if title, ok := updates["title"].(string); ok {
				r.tasks[i].Title = title
			}
			if desc, ok := updates["description"].(string); ok {
				r.tasks[i].Description = desc
			}
			if priority, ok := updates["priority"].(string); ok {
				r.tasks[i].Priority = priority
			}
			if dueDate, ok := updates["due_date"].(string); ok {
				r.tasks[i].DueDate = dueDate
			}
			return r.tasks[i], nil
		}
	}
	return models.Task{}, fmt.Errorf("task not found")
}

func (r *TaskRepository) Complete(id int) (models.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, task := range r.tasks {
		if task.ID == id {
			r.tasks[i].Completed = true
			return r.tasks[i], nil
		}
	}
	return models.Task{}, fmt.Errorf("task not found")
}

func (r *TaskRepository) Search(keyword string) []models.Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make([]models.Task, 0)
	keyword = strings.ToLower(keyword)
	for _, task := range r.tasks {
		if strings.Contains(strings.ToLower(task.Title), keyword) || strings.Contains(strings.ToLower(task.Description), keyword) {
			results = append(results, task)
		}
	}
	return results
}
