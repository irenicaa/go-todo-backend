package db

import (
	"database/sql"
	"flag"
	"testing"

	"github.com/irenicaa/go-todo-backend/models"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var dataSourceName = flag.String(
	"dataSourceName",
	"postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable",
	"DB connection string",
)

func TestDB_Create(t *testing.T) {
	type args struct {
		todo models.TodoRecord
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				todo: models.TodoRecord{
					Title:     "test",
					Completed: true,
					Order:     42,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := OpenDB(*dataSourceName)
			require.NoError(t, err)

			id, err := db.Create(tt.args.todo)
			require.NoError(t, err)
			tt.args.todo.ID = id

			todo, err := db.GetSingle(id)
			require.NoError(t, err)

			assert.Equal(t, tt.args.todo, todo)
		})
	}
}

func TestDB_Update(t *testing.T) {
	type args struct {
		originalTodo models.TodoRecord
		updatedTodo  models.TodoRecord
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				originalTodo: models.TodoRecord{
					Title:     "test",
					Completed: true,
					Order:     23,
				},
				updatedTodo: models.TodoRecord{
					Title:     "test2",
					Completed: false,
					Order:     42,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := OpenDB(*dataSourceName)
			require.NoError(t, err)

			id, err := db.Create(tt.args.originalTodo)
			require.NoError(t, err)

			err = db.Update(id, tt.args.updatedTodo)
			require.NoError(t, err)
			tt.args.updatedTodo.ID = id

			todo, err := db.GetSingle(id)
			require.NoError(t, err)

			assert.Equal(t, tt.args.updatedTodo, todo)
		})
	}
}

func TestDB_Delete(t *testing.T) {
	type args struct {
		todo models.TodoRecord
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success",
			args: args{
				todo: models.TodoRecord{
					Title:     "test",
					Completed: true,
					Order:     42,
				},
			},
			wantErr: sql.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := OpenDB(*dataSourceName)
			require.NoError(t, err)

			id, err := db.Create(tt.args.todo)
			require.NoError(t, err)

			err = db.Delete(id)
			require.NoError(t, err)

			todo, err := db.GetSingle(id)

			assert.Equal(t, models.TodoRecord{}, todo)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
