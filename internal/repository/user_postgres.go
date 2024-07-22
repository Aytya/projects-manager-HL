package repository

import (
	"fmt"
	"github.com/Aytya/projects-manager-HL/internal/entity"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (repo *UserPostgres) CreateUser(user entity.User) (string, time.Time, error) {
	query := fmt.Sprintf("INSERT INTO %s (name, email, role) VALUES ($1, $2, $3) RETURNING id, registered_at", usersTable)
	return Create(repo.db, query, user.Name, user.Email, user.Role)
}

func (repo *UserPostgres) GetByColumn(value, column string) (entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", usersTable, column)
	err := repo.db.Get(&user, query, value)

	return user, err
}

func (repo *UserPostgres) GetUserById(id string) (entity.User, error) {
	return repo.GetByColumn(id, "id")
}

func (repo *UserPostgres) GetUserByName(name string) (entity.User, error) {
	return repo.GetByColumn(name, "name")
}

func (repo *UserPostgres) GetUserByEmail(email string) (entity.User, error) {
	return repo.GetByColumn(email, "email")
}

func (repo *UserPostgres) UpdateUser(id string, user entity.User) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if user.Name != "" {
		setValues = append(setValues, fmt.Sprintf("name = $%d", argId))
		args = append(args, user.Name)
		argId++
	}

	if user.Email != "" {
		setValues = append(setValues, fmt.Sprintf("email = $%d", argId))
		args = append(args, user.Email)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", usersTable, setQuery, argId)
	args = append(args, id)

	_, err := repo.db.Exec(query, args...)
	return err
}

func (repo *UserPostgres) DeleteUser(id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", usersTable)
	_, err := repo.db.Exec(query, id)
	return err
}

func (repo *UserPostgres) GetAllUsers() ([]entity.User, error) {
	var users []entity.User
	query := fmt.Sprintf("SELECT * FROM %s", usersTable)
	if err := repo.db.Select(&users, query); err != nil {
		return nil, err
	}
	return users, nil
}
