package repository

import (
	"github.com/Aytya/projects-manager-HL/internal/entity"
	"github.com/jmoiron/sqlx"
	"time"
)

type User interface {
	CreateUser(user entity.User) (string, time.Time, error)
	UpdateUser(id string, user entity.User) error
	DeleteUser(id string) error
	GetUserById(id string) (entity.User, error)
	GetUserByName(name string) (entity.User, error)
	GetAllUsers() ([]entity.User, error)
	GetUserByEmail(email string) (entity.User, error)
}

type Task interface {
	CreateTask(task entity.Task) (string, time.Time, error)
	UpdateTask(id string, task entity.Task) error
	DeleteTask(id string) error
	GetTaskById(id string) (entity.Task, error)
	GetAllTasks() ([]entity.Task, error)
	GetTaskByTitle(title string) (entity.Task, error)
	GetTaskByStatus(status string) ([]entity.Task, error)
	GetTaskByPriority(priority string) ([]entity.Task, error)
	GetTasksByUserId(userId, column string) ([]entity.Task, error)
	GetTasksByProjectId(projectId string) ([]entity.Task, error)
}

type Project interface {
	CreateProject(project entity.Project) (string, time.Time, error)
	UpdateProject(id string, project entity.Project) error
	DeleteProject(id string) error
	GetProjectById(id string) (entity.Project, error)
	GetProjectByTitle(title string) (entity.Project, error)
	GetProjectByManagerId(managerId string) ([]entity.Project, error)
	GetAllProjects() ([]entity.Project, error)
}

type Repository struct {
	User
	Task
	Project
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:    NewUserPostgres(db),
		Task:    NewTaskPostgres(db),
		Project: NewProjectPostgres(db),
	}
}
