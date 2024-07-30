package repository

import (
	"fmt"
	"github.com/Aytya/projects-manager-HL/internal/entity"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type ProjectPostgres struct {
	db *sqlx.DB
}

func (repo *ProjectPostgres) CreateProject(project entity.Project) (string, time.Time, error) {
	query := fmt.Sprintf("INSERT INTO %s (title,description,manager) VALUES ($1,$2,$3) RETURNING id, created_at", projectsTable)
	return Create(repo.db, query, project.Title, project.Description, project.Manager)
}

func (repo ProjectPostgres) UpdateProject(id string, project entity.Project) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if project.Title != "" {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, project.Title)
		argId++
	}

	if project.Description != "" {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, project.Description)
		argId++
	}

	if project.Manager != "" {
		setValues = append(setValues, fmt.Sprintf("manager=$%d", argId))
		args = append(args, project.Manager)
		argId++
	}

	if project.FinishedAt.Valid {
		setValues = append(setValues, fmt.Sprintf("finished_at=$%d", argId))
		args = append(args, project.FinishedAt)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", projectsTable, setQuery, argId)
	args = append(args, id)

	_, err := repo.db.Exec(query, args...)
	return err
}

func (repo ProjectPostgres) DeleteProject(id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", projectsTable)
	_, err := repo.db.Exec(query, id)
	return err
}

func (repo ProjectPostgres) GetByColumn(column, value string) (entity.Project, error) {
	var project entity.Project
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", projectsTable, column)
	err := repo.db.Get(&project, query, value)
	if err != nil {
		return entity.Project{}, err
	}
	return project, err
}

func (repo ProjectPostgres) GetProjectById(id string) (entity.Project, error) {
	return repo.GetByColumn("id", id)
}

func (repo ProjectPostgres) GetProjectByTitle(title string) (entity.Project, error) {
	return repo.GetByColumn("title", title)
}

func (repo ProjectPostgres) GetProjectByManagerId(managerId string) ([]entity.Project, error) {
	var projects []entity.Project
	query := fmt.Sprintf("SELECT * FROM %s WHERE Manager = $1", projectsTable)
	err := repo.db.Select(&projects, query, managerId)

	return projects, err
}

func (repo ProjectPostgres) GetAllProjects() ([]entity.Project, error) {
	var projects []entity.Project
	query := fmt.Sprintf("SELECT * FROM %s", projectsTable)
	err := repo.db.Select(&projects, query)
	return projects, err
}

func NewProjectPostgres(db *sqlx.DB) *ProjectPostgres {
	return &ProjectPostgres{db: db}
}
