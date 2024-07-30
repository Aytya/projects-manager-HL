package tests

import (
	"fmt"
	"github.com/Aytya/projects-manager-HL/internal/entity"
	"github.com/Aytya/projects-manager-HL/internal/repository"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewUserPostgres(db)

	tests := []struct {
		name    string
		mock    func()
		input   entity.User
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "registered_at"}).AddRow(1, time.Now())
				mock.ExpectQuery("INSERT INTO users").WillReturnRows(rows)
			},
			input: entity.User{
				Name:  "Test User",
				Email: "test@test.com",
				Role:  "User",
			},
			want:    "1",
			wantErr: false,
		},
		{
			name: "error - empty name",
			mock: func() {
				mock.ExpectQuery("INSERT INTO users").WillReturnError(fmt.Errorf("some error"))
			},
			input: entity.User{
				Name:  "",
				Email: "test@test.com",
				Role:  "",
			},
			wantErr: true,
		},
		{
			name: "error - database error",
			mock: func() {
				mock.ExpectQuery("INSERT INTO users").WithArgs("Test User", "test@test.com", "").WillReturnError(fmt.Errorf("some error"))
			},
			input: entity.User{
				Name:  "Test User",
				Email: "test@test.com",
				Role:  "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, _, err := r.CreateUser(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewUserPostgres(db)
	tests := []struct {
		name    string
		mock    func()
		input   entity.User
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "registered_at"}).AddRow(1, time.Now())
				mock.ExpectQuery("UPDATE users SET (.+) WHERE (.+)").WillReturnRows(rows)
			},
			input: entity.User{
				Name:  "Test User",
				Email: "test@test.com",
				Role:  "User",
			},
			want:    "User updated successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.UpdateUser(tt.input.ID, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}

}

func TestGetUserById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewUserPostgres(db)

	tests := []struct {
		name    string
		mock    func()
		id      string
		want    entity.User
		wantErr bool
	}{
		{
			name: "ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "role"}).AddRow(1, "Test User", "test@test.com", "User")
				mock.ExpectQuery("SELECT \\* FROM users WHERE id = \\$1").WithArgs("1").WillReturnRows(rows)
			},
			id: "1",
			want: entity.User{
				ID:    "1",
				Name:  "Test User",
				Email: "test@test.com",
				Role:  "User",
			},
			wantErr: false,
		},
		{
			name: "error - no user found",
			mock: func() {
				mock.ExpectQuery("SELECT \\* FROM users WHERE id = \\$1").WithArgs("1").WillReturnError(fmt.Errorf("no rows in result set"))
			},
			id:      "1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := r.GetUserById(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewUserPostgres(db)

	tests := []struct {
		name    string
		mock    func()
		id      string
		wantErr bool
	}{
		{
			name: "ok",
			mock: func() {
				mock.ExpectExec("DELETE FROM users WHERE id = \\$1").WithArgs("1").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			id:      "1",
			wantErr: false,
		},
		{
			name: "error - database error",
			mock: func() {
				mock.ExpectExec("DELETE FROM users WHERE id = \\$1").WithArgs("1").WillReturnError(fmt.Errorf("some error"))
			},
			id:      "1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := r.DeleteUser(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewUserPostgres(db)

	tests := []struct {
		name    string
		mock    func()
		want    []entity.User
		wantErr bool
	}{
		{
			name: "ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "role"}).
					AddRow(1, "User1", "user1@test.com", "Admin").
					AddRow(2, "User2", "user2@test.com", "User")
				mock.ExpectQuery("SELECT \\* FROM users").WillReturnRows(rows)
			},
			want: []entity.User{
				{ID: "1", Name: "User1", Email: "user1@test.com", Role: "Admin"},
				{ID: "2", Name: "User2", Email: "user2@test.com", Role: "User"},
			},
			wantErr: false,
		},
		{
			name: "error - database error",
			mock: func() {
				mock.ExpectQuery("SELECT \\* FROM users").WillReturnError(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := r.GetAllUsers()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
