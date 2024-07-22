package repository

import (
	"fmt"
	"github.com/Aytya/projects-manager-HL/internal/entity"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type TaskPostgres struct {
	db *sqlx.DB
}

func NewTaskPostgres(db *sqlx.DB) *TaskPostgres {
	return &TaskPostgres{db: db}
}

func (repo *TaskPostgres) CreateTask(task entity.Task) (string, time.Time, error) {
	query := fmt.Sprintf("INSERT INTO %s (title, description, priority, status, assignee, project) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at", tasksTable)
	return Create(repo.db, query, task.Title, task.Description, task.Priority, task.Status, task.Assignee, task.Project)
}

func (repo *TaskPostgres) GetByColumn(column, value string) (entity.Task, error) {
	var task entity.Task
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", tasksTable, column)
	err := repo.db.Get(&task, query, value)
	if err != nil {
		return task, fmt.Errorf("error retrieving task by %s: %v", column, err)
	}

	return task, nil
}

func (repo *TaskPostgres) SelectByColumn(column, value string) ([]entity.Task, error) {
	var tasks []entity.Task
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", tasksTable, column)
	err := repo.db.Select(&tasks, query, value)

	return tasks, err
}

func (repo *TaskPostgres) GetTaskById(id string) (entity.Task, error) {
	return repo.GetByColumn("id", id)
}

func (repo *TaskPostgres) GetTaskByTitle(title string) (entity.Task, error) {
	return repo.GetByColumn("title", title)
}

func (repo TaskPostgres) GetTaskByStatus(status string) ([]entity.Task, error) {
	return repo.SelectByColumn("status", status)
}

func (repo TaskPostgres) GetTaskByPriority(priority string) ([]entity.Task, error) {
	return repo.SelectByColumn("priority", priority)
}

func (repo TaskPostgres) GetTasksByUserId(userId, column string) ([]entity.Task, error) {
	return repo.SelectByColumn(column, userId)
}

func (repo TaskPostgres) GetTasksByProjectId(project string) ([]entity.Task, error) {
	return repo.SelectByColumn("project", project)
}

func (repo TaskPostgres) GetAllTasks() ([]entity.Task, error) {
	var tasks []entity.Task
	query := fmt.Sprintf("SELECT * FROM %s", tasksTable)
	err := repo.db.Select(&tasks, query)
	if err != nil {
		return nil, fmt.Errorf("error retrieving all tasks: %v", err)
	}

	return tasks, nil
}

func (repo TaskPostgres) UpdateTask(id string, task entity.Task) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if task.Title != "" {
		setValues = append(setValues, fmt.Sprintf("title = $%d", argId))
		args = append(args, task.Title)
		argId++
	}

	if task.Description != "" {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argId))
		args = append(args, task.Description)
		argId++
	}

	if task.Priority != "" {
		setValues = append(setValues, fmt.Sprintf("priority = $%d", argId))
		args = append(args, task.Priority)
		argId++
	}

	if task.Status != "" {
		setValues = append(setValues, fmt.Sprintf("status = $%d", argId))
		args = append(args, task.Status)
		argId++
	}

	if task.Assignee != "" {
		setValues = append(setValues, fmt.Sprintf("assignee = $%d", argId))
		args = append(args, task.Assignee)
		argId++
	}

	if task.Project != "" {
		setValues = append(setValues, fmt.Sprintf("project = $%d", argId))
		args = append(args, task.Project)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", tasksTable, setQuery, argId)
	args = append(args, id)

	_, err := repo.db.Exec(query, args...)
	return err
}

func (repo TaskPostgres) DeleteTask(id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tasksTable)
	_, err := repo.db.Exec(query, id)
	return err
}
