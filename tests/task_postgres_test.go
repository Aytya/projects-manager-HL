package tests

import (
	"database/sql"
	"fmt"
	"github.com/Aytya/projects-manager-HL/internal/entity"
	"github.com/Aytya/projects-manager-HL/internal/repository"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
	"time"
)

func TestCreateTask(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewTaskPostgres(db)

	tests := []struct {
		name    string
		mock    func()
		input   entity.Task
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "registered_at"}).AddRow(1, time.Now())
				mock.ExpectQuery("INSERT INTO tasks").WillReturnRows(rows)
			},
			input: entity.Task{
				Title:       "Write Unit Tests 1",
				Description: "Create unit tests for the authentication module.",
				Priority:    "High",
				Status:      "Not Started",
				Assignee:    "user1",
				Project:     "project1",
			},
			want:    "1",
			wantErr: false,
		},
		{
			name: "error - empty title",
			mock: func() {
				mock.ExpectQuery("INSERT INTO tasks").WithArgs(
					"", "Create unit tests for the authentication module.", "High", "Not Started", "user1", "project1",
				).WillReturnError(fmt.Errorf("ERROR: null value in column \"title\" violates not-null constraint"))
			},
			input: entity.Task{
				Title:       "",
				Description: "Create unit tests for the authentication module.",
				Priority:    "High",
				Status:      "Not Started",
				Assignee:    "user1",
				Project:     "project1",
			},
			wantErr: true,
		},
		{
			name: "error - database error",
			mock: func() {
				mock.ExpectQuery("INSERT INTO tasks").WillReturnError(fmt.Errorf("some db error"))
			},
			input: entity.Task{
				Title:       "Write Unit Tests 1",
				Description: "Create unit tests for the authentication module.",
				Priority:    "High",
				Status:      "Not Started",
				Assignee:    "user1",
				Project:     "project1",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			got, _, err := r.CreateTask(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("CreateTask() error = %v, wantErr %v", err, tc.wantErr)
			} else {
				assert.Equal(t, tc.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewTaskPostgres(db)

	tests := []struct {
		name    string
		mock    func()
		input   entity.Task
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			mock: func() {
				mock.ExpectExec("UPDATE tasks").WithArgs(
					"Updated Title", "Updated Description", "Medium", "In Progress", "user2", "project2", "1",
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input: entity.Task{
				Title:       "Updated Title",
				Description: "Updated Description",
				Priority:    "Medium",
				Status:      "In Progress",
				Assignee:    "user2",
				Project:     "project2",
			},
			want:    "Task updated",
			wantErr: false,
		},
		{
			name: "error - database error",
			mock: func() {
				mock.ExpectExec("UPDATE tasks").WithArgs("Write Unit Tests 2", "Updated description", "Medium", "In Progress", "user2", "project2", "1").
					WillReturnError(fmt.Errorf("db error"))
			},
			input: entity.Task{
				Title:       "Write Unit Tests 2",
				Description: "Updated description",
				Priority:    "Medium",
				Status:      "In Progress",
				Assignee:    "user2",
				Project:     "project2",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.UpdateTask(tt.input.ID, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestGetTaskById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewTaskPostgres(db)

	tests := []struct {
		name    string
		mock    func()
		id      string
		want    entity.Task
		wantErr bool
	}{
		{
			name: "ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "priority", "status", "assignee", "project"}).
					AddRow("1", "Test Title", "Test Description", "High", "Not Started", "user1", "project1")
				mock.ExpectQuery("SELECT \\* FROM tasks WHERE id = \\$1").WithArgs("1").WillReturnRows(rows)
			},
			id: "1",
			want: entity.Task{
				ID:          "1",
				Title:       "Test Title",
				Description: "Test Description",
				Priority:    "High",
				Status:      "Not Started",
				Assignee:    "user1",
				Project:     "project1",
			},
			wantErr: false,
		},
		{
			name: "error - not found",
			mock: func() {
				mock.ExpectQuery("SELECT \\* FROM tasks WHERE id = \\$1").WithArgs("1").WillReturnError(fmt.Errorf("sql: no rows in result set"))
			},
			id:      "1",
			want:    entity.Task{},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			got, err := r.GetTaskById(tc.id)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetTaskById() error = %v, wantErr %v", err, tc.wantErr)
			}
			assert.Equal(t, tc.want, got)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewTaskPostgres(db)

	tests := []struct {
		name    string
		mock    func()
		id      string
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			mock: func() {
				mock.ExpectExec("DELETE FROM tasks WHERE id = ?").WithArgs("1").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			id:      "1",
			want:    "Task deleted",
			wantErr: false,
		},
		{
			name: "error - database error",
			mock: func() {
				mock.ExpectExec("DELETE FROM tasks WHERE id = ?").WithArgs("1").
					WillReturnError(fmt.Errorf("db error"))
			},
			id:      "1",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			err := r.DeleteTask(tc.id)
			if (err != nil) != tc.wantErr {
				t.Errorf("DeleteTask() error = %v, wantErr %v", err, tc.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetAllTasks(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewTaskPostgres(db)

	tests := []struct {
		name    string
		mock    func()
		want    []entity.Task
		wantErr bool
	}{
		{
			name: "ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "priority", "status", "assignee", "project", "created_at", "finished_at"}).
					AddRow(1, "Write Unit Tests 1", "Create unit tests for the authentication module.", "High", "Not Started", "user1", "project1", time.Now(), nil).
					AddRow(2, "Write Unit Tests 2", "Create unit tests for the project management module.", "Medium", "In Progress", "user2", "project2", time.Now(), nil)
				mock.ExpectQuery("SELECT \\* FROM tasks").WillReturnRows(rows)
			},
			want: []entity.Task{
				{
					ID:          "1",
					Title:       "Write Unit Tests 1",
					Description: "Create unit tests for the authentication module.",
					Priority:    "High",
					Status:      "Not Started",
					Assignee:    "user1",
					Project:     "project1",
					// Omit CreatedAt and FinishedAt checks
				},
				{
					ID:          "2",
					Title:       "Write Unit Tests 2",
					Description: "Create unit tests for the project management module.",
					Priority:    "Medium",
					Status:      "In Progress",
					Assignee:    "user2",
					Project:     "project2",
					// Omit CreatedAt and FinishedAt checks
				},
			},
			wantErr: false,
		},
		{
			name: "error - database error",
			mock: func() {
				mock.ExpectQuery("SELECT \\* FROM tasks").WillReturnError(fmt.Errorf("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := r.GetAllTasks()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Compare only the relevant fields
				for i := range tt.want {
					tt.want[i].CreatedAt = time.Time{}
					tt.want[i].FinishedAt = sql.NullTime{}
				}
				for i := range got {
					got[i].CreatedAt = time.Time{}
					got[i].FinishedAt = sql.NullTime{}
				}
				assert.ElementsMatch(t, tt.want, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
