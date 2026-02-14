package repository

import (
	"database/sql"

	"github.com/vinayb91/task-manager-ai-backend.git/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user models.User) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	err = r.DB.QueryRow(
		"INSERT INTO users (email, name, password) VALUES ($1, $2, $3) RETURNING id, created_at",
		user.Email, user.Name, string(hashedPassword),
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	user.Password = ""
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRow(
		"SELECT id, email, name, password, is_admin, created_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.IsAdmin, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRow(
		"SELECT id, email, name, is_admin, created_at FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Email, &user.Name, &user.IsAdmin, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
