package repository

import (
	"database/sql"
	"fmt"

	"github.com/vinayb91/task-manager-ai-backend.git/internal/models"
)

type TaskRepository struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

func (r *TaskRepository) List(userID int, status, priority string) ([]models.Task, error) {
	query := "SELECT id, user_id, title, description, priority, due_date, completed, created_at, updated_at FROM tasks WHERE user_id = $1"
	args := []interface{}{userID}
	argCount := 1

	if status == "pending" {
		query += " AND completed = false"
	} else if status == "completed" {
		query += " AND completed = true"
	}

	if priority != "" {
		argCount++
		query += fmt.Sprintf(" AND priority = $%d", argCount)
		args = append(args, priority)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID, &task.UserID, &task.Title, &task.Description,
			&task.Priority, &task.DueDate, &task.Completed,
			&task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskRepository) Create(task models.Task) (*models.Task, error) {
	err := r.DB.QueryRow(
		`INSERT INTO tasks (user_id, title, description, priority, due_date) 
		 VALUES ($1, $2, $3, $4, $5) 
		 RETURNING id, created_at, updated_at`,
		task.UserID, task.Title, task.Description, task.Priority, task.DueDate,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) Delete(userID, taskID int) error {
	result, err := r.DB.Exec("DELETE FROM tasks WHERE id = $1 AND user_id = $2", taskID, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (r *TaskRepository) Update(userID, taskID int, updates map[string]interface{}) (*models.Task, error) {
	query := "UPDATE tasks SET updated_at = CURRENT_TIMESTAMP"
	args := []interface{}{}
	argCount := 0

	if title, ok := updates["title"].(string); ok {
		argCount++
		query += fmt.Sprintf(", title = $%d", argCount)
		args = append(args, title)
	}
	if desc, ok := updates["description"].(string); ok {
		argCount++
		query += fmt.Sprintf(", description = $%d", argCount)
		args = append(args, desc)
	}
	if priority, ok := updates["priority"].(string); ok {
		argCount++
		query += fmt.Sprintf(", priority = $%d", argCount)
		args = append(args, priority)
	}
	if dueDate, ok := updates["due_date"].(string); ok {
		argCount++
		query += fmt.Sprintf(", due_date = $%d", argCount)
		args = append(args, dueDate)
	}

	argCount++
	query += fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, taskID)

	argCount++
	query += fmt.Sprintf(" AND user_id = $%d", argCount)
	args = append(args, userID)

	query += " RETURNING id, user_id, title, description, priority, due_date, completed, created_at, updated_at"

	var task models.Task
	err := r.DB.QueryRow(query, args...).Scan(
		&task.ID, &task.UserID, &task.Title, &task.Description,
		&task.Priority, &task.DueDate, &task.Completed,
		&task.CreatedAt, &task.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) Complete(userID, taskID int) (*models.Task, error) {
	var task models.Task
	err := r.DB.QueryRow(
		`UPDATE tasks SET completed = true, updated_at = CURRENT_TIMESTAMP 
		 WHERE id = $1 AND user_id = $2 
		 RETURNING id, user_id, title, description, priority, due_date, completed, created_at, updated_at`,
		taskID, userID,
	).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Priority,
		&task.DueDate, &task.Completed, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) Search(userID int, keyword string) ([]models.Task, error) {
	query := `SELECT id, user_id, title, description, priority, due_date, completed, created_at, updated_at 
	          FROM tasks 
	          WHERE user_id = $1 AND (title ILIKE $2 OR description ILIKE $2)
	          ORDER BY created_at DESC`

	rows, err := r.DB.Query(query, userID, "%"+keyword+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID, &task.UserID, &task.Title, &task.Description,
			&task.Priority, &task.DueDate, &task.Completed,
			&task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
