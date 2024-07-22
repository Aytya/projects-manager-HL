package service

import (
	"github.com/Aytya/projects-manager-HL/internal/entity"
	"github.com/Aytya/projects-manager-HL/internal/repository"
	"time"
)

type UserService struct {
	repo repository.User
}

func (u UserService) GetAllUsers() ([]entity.User, error) {
	return u.repo.GetAllUsers()
}

func (u UserService) CreateUser(user entity.User) (string, time.Time, error) {
	return u.repo.CreateUser(user)
}

func (u UserService) UpdateUser(id string, user entity.User) error {
	return u.repo.UpdateUser(id, user)
}

func (u UserService) DeleteUser(id string) error {
	return u.repo.DeleteUser(id)
}

func (u UserService) GetUserById(id string) (entity.User, error) {
	return u.repo.GetUserById(id)
}

func (u UserService) GetUserByName(name string) (entity.User, error) {
	return u.repo.GetUserByName(name)
}

func (u UserService) GetUserByEmail(email string) (entity.User, error) {
	return u.repo.GetUserByEmail(email)
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}
