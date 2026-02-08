package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		is_admin BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		priority VARCHAR(20) NOT NULL,
		due_date DATE,
		completed BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS app_settings (
		id SERIAL PRIMARY KEY,
		speech_enabled BOOLEAN DEFAULT TRUE,
		gpt_model VARCHAR(50) DEFAULT 'gpt-4o',
		max_tokens INTEGER DEFAULT 1000,
		voice_model VARCHAR(50) DEFAULT 'tts-1',
		voice_enabled BOOLEAN DEFAULT TRUE,
		max_tasks_per_user INTEGER DEFAULT 100,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Insert default settings if not exists
	INSERT INTO app_settings (id, speech_enabled, gpt_model, max_tokens, voice_model, voice_enabled, max_tasks_per_user)
	SELECT 1, TRUE, 'gpt-4o', 1000, 'tts-1', TRUE, 100
	WHERE NOT EXISTS (SELECT 1 FROM app_settings WHERE id = 1);
	`

	_, err = db.Exec(schema)
	return db, err
}
