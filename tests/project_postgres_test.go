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

func TestCreateProject(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := repository.NewProjectPostgres(db)

	tests := []struct {
		name    string
		mock    func()
		input   entity.Project
		want    string
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now())
				mock.ExpectQuery(`INSERT INTO projects \(title,description,manager\) VALUES \(\$1,\$2,\$3\) RETURNING id, created_at`).WillReturnRows(rows)
			},
			input: entity.Project{
				Title:       "New Project 2",
				Description: "This is a easy project.",
				Manager:     "manager1",
			},
			want:    "1",
			wantErr: false,
		},
		{
			name: "error - empty title",
			mock: func() {
				mock.ExpectQuery(`INSERT INTO "projects"`).WillReturnError(fmt.Errorf("some error"))
			},
			input: entity.Project{
				Title:       "",
				Description: "This is a easy project.",
				Manager:     "manager1",
			},
			wantErr: true,
		},
		{
			name: "error - database error",
			mock: func() {
				mock.ExpectQuery(`INSERT INTO "projects"`).WillReturnError(fmt.Errorf("some error"))
			},
			input: entity.Project{
				Title:       "New Project 2",
				Description: "This is a easy project.",
				Manager:     "manager1",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, _, err := r.CreateProject(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetProjectById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewProjectPostgres(db)
	tests := []struct {
		name    string
		mock    func()
		id      string
		want    entity.Project
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "manager", "created_at"}).AddRow(1, "title1", "description", "manager1", time.Now())
				mock.ExpectQuery("SELECT (.+) FROM projects WHERE id = \\$1").WillReturnRows(rows)
			},
			id: "1",
			want: entity.Project{
				ID:          "1",
				Title:       "title1",
				Description: "description",
				Manager:     "manager1",
				CreatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "error - no project found",
			mock: func() {
				mock.ExpectQuery("SELECT (.+) FROM projects WHERE id = \\$1").WillReturnError(fmt.Errorf("no rows in result set"))
			},
			id:      "1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := r.GetProjectById(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.Title, got.Title)
				assert.Equal(t, tt.want.Description, got.Description)
				assert.Equal(t, tt.want.Manager, got.Manager)
				assert.WithinDuration(t, tt.want.CreatedAt, got.CreatedAt, time.Second)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateProject(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewProjectPostgres(db)
	tests := []struct {
		name    string
		mock    func()
		input   entity.Project
		want    string
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				mock.ExpectExec("UPDATE projects SET title=\\$1, description=\\$2, manager=\\$3 WHERE id = \\$4").
					WithArgs("title2", "description2", "manager2", "1").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input: entity.Project{
				ID:          "1",
				Title:       "title2",
				Description: "description2",
				Manager:     "manager2",
			},
			want:    "Project updated",
			wantErr: false,
		},
		{
			name: "error - update failure",
			mock: func() {
				mock.ExpectExec("UPDATE projects SET title=\\$1, description=\\$2, manager=\\$3 WHERE id = \\$4").
					WithArgs("title2", "description2", "manager2", "1").
					WillReturnError(fmt.Errorf("update failed"))
			},
			input: entity.Project{
				ID:          "1",
				Title:       "title2",
				Description: "description2",
				Manager:     "manager2",
			},
			want:    "Project updated",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.UpdateProject(tt.input.ID, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, "Project updated")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteProjectById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewProjectPostgres(db)
	tests := []struct {
		name    string
		mock    func()
		id      string
		want    string
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				mock.ExpectExec("DELETE FROM projects WHERE id = \\$1").
					WithArgs("1").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			id:      "1",
			want:    "Project deleted",
			wantErr: false,
		},
		{
			name: "error - no project found",
			mock: func() {
				mock.ExpectExec("DELETE FROM projects WHERE id = \\$1").
					WithArgs("1").WillReturnError(fmt.Errorf("no rows in result set"))
			},
			id:      "1",
			wantErr: true,
		},
		{
			name: "error - database error",
			mock: func() {
				mock.ExpectExec("DELETE FROM projects WHERE id = \\$1").
					WithArgs("1").WillReturnError(fmt.Errorf("some error"))
			},
			id:      "1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.DeleteProject(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListProjects(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewProjectPostgres(db)

	staticTime := time.Date(2024, time.July, 29, 15, 45, 24, 0, time.UTC)

	tests := []struct {
		name    string
		mock    func()
		want    []entity.Project
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "manager", "created_at"}).
					AddRow(1, "title1", "description1", "manager1", staticTime).
					AddRow(2, "title2", "description2", "manager2", staticTime)
				mock.ExpectQuery("SELECT (.+) FROM projects").WillReturnRows(rows)
			},
			want: []entity.Project{
				{
					ID:          "1",
					Title:       "title1",
					Description: "description1",
					CreatedAt:   staticTime,
					Manager:     "manager1",
				},
				{
					ID:          "2",
					Title:       "title2",
					Description: "description2",
					CreatedAt:   staticTime,
					Manager:     "manager2",
				},
			},
			wantErr: false,
		},
		{
			name: "error - database error",
			mock: func() {
				mock.ExpectQuery("SELECT (.+) FROM projects").WillReturnError(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := r.GetAllProjects()
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
