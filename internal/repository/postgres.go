package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

const (
	usersTable    = "users"
	tasksTable    = "tasks"
	projectsTable = "projects"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

func NewPostgresDB(config Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Database, config.Password, config.SSLMode))

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = createUsersTable(db)
	if err != nil {
		return nil, err
	}

	err = createProjectsTable(db)
	if err != nil {
		return nil, err
	}

	err = createTasksTable(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Create(db *sqlx.DB, query string, args ...interface{}) (string, time.Time, error) {
	var id string
	var timestamp time.Time
	row := db.QueryRow(query, args...)
	if err := row.Scan(&id, &timestamp); err != nil {
		return "", time.Time{}, fmt.Errorf("error scanning row: %v", err)
	}
	return id, timestamp, nil
}

func createUsersTable(db *sqlx.DB) error {
	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		registered_at TIMESTAMP DEFAULT NOW(),
		role VARCHAR(50) NOT NULL
	)`, usersTable)

	_, err := db.Exec(query)
	return err
}

func createTasksTable(db *sqlx.DB) error {
	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    priority VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    assignee UUID NOT NULL,
    project UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    finished_at TIMESTAMP NULL,
	FOREIGN KEY (assignee) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (project) REFERENCES projects(id) ON DELETE CASCADE
	    )
    `, tasksTable)

	_, err := db.Exec(query)
	return err
}

func createProjectsTable(db *sqlx.DB) error {
	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		title VARCHAR(255) NOT NULL,
		description VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		finished_at TIMESTAMP NULL,
		manager UUID NOT NULL,
	    FOREIGN KEY (manager) REFERENCES users(id) ON DELETE CASCADE
	)`, projectsTable)

	_, err := db.Exec(query)
	return err
}
