package service

import (
	"github.com/Aytya/projects-manager-HL/internal/entity"
	"github.com/Aytya/projects-manager-HL/internal/repository"
	"time"
)

type ProjectService struct {
	repo repository.Project
}

func (p ProjectService) GetAllProjects() ([]entity.Project, error) {
	return p.repo.GetAllProjects()
}

func (p ProjectService) CreateProject(project entity.Project) (string, time.Time, error) {
	return p.repo.CreateProject(project)
}

func (p ProjectService) UpdateProject(id string, project entity.Project) error {
	return p.repo.UpdateProject(id, project)
}

func (p ProjectService) DeleteProject(id string) error {
	return p.repo.DeleteProject(id)
}

func (p ProjectService) GetProjectById(id string) (entity.Project, error) {
	return p.repo.GetProjectById(id)
}

func (p ProjectService) GetProjectByTitle(projectTitle string) (entity.Project, error) {
	return p.repo.GetProjectByTitle(projectTitle)
}

func (p ProjectService) GetProjectByManagerId(managerId string) ([]entity.Project, error) {
	return p.repo.GetProjectByManagerId(managerId)
}

func NewProjectService(repo repository.Project) *ProjectService {
	return &ProjectService{repo: repo}
}
