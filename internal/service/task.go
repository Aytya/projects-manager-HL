package service

import (
	"github.com/Aytya/projects-manager-HL/internal/entity"
	"github.com/Aytya/projects-manager-HL/internal/repository"
	"time"
)

type TaskService struct {
	repo repository.Task
}

func (t TaskService) CreateTask(task entity.Task) (string, time.Time, error) {
	return t.repo.CreateTask(task)
}

func (t TaskService) UpdateTask(id string, task entity.Task) error {
	return t.repo.UpdateTask(id, task)
}

func (t TaskService) DeleteTask(id string) error {
	return t.repo.DeleteTask(id)
}

func (t TaskService) GetTaskById(id string) (entity.Task, error) {
	return t.repo.GetTaskById(id)
}

func (t TaskService) GetAllTasks() ([]entity.Task, error) {
	return t.repo.GetAllTasks()
}

func (t TaskService) GetTaskByTitle(title string) (entity.Task, error) {
	return t.repo.GetTaskByTitle(title)
}

func (t TaskService) GetTaskByStatus(status string) ([]entity.Task, error) {
	return t.repo.GetTaskByStatus(status)
}

func (t TaskService) GetTaskByPriority(priority string) ([]entity.Task, error) {
	return t.repo.GetTaskByPriority(priority)
}

func (t TaskService) GetTasksByUserId(userId, column string) ([]entity.Task, error) {
	return t.repo.GetTasksByUserId(userId, column)
}

func (t TaskService) GetTasksByProjectId(projectId string) ([]entity.Task, error) {
	return t.repo.GetTasksByProjectId(projectId)
}

func NewTaskService(repo repository.Task) *TaskService {
	return &TaskService{repo: repo}
}
