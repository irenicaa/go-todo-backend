// +build integration

package db

import (
	"database/sql"
	"flag"
	"testing"

	"github.com/irenicaa/go-todo-backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var dataSourceName = flag.String(
	"dataSourceName",
	DefaultDataSourceName,
	"DB connection string",
)

func TestDB_GetAll(t *testing.T) {
	type args struct {
		todos []models.TodoRecord
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				todos: []models.TodoRecord{
					{
						Title:     "test",
						Completed: true,
						Order:     23,
					},
					{
						Title:     "test2",
						Completed: false,
						Order:     42,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, err := OpenDB(*dataSourceName)
			require.NoError(t, err)

			db := NewTodoRecord(pool)
			err = db.DeleteAll()
			require.NoError(t, err)

			for index, todo := range tt.args.todos {
				id, err2 := db.Create(todo)
				require.NoError(t, err2)

				tt.args.todos[index].ID = id
			}

			todos, err := db.GetAll(models.Query{})
			require.NoError(t, err)

			assert.Equal(t, tt.args.todos, todos)
		})
	}
}

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
			pool, err := OpenDB(*dataSourceName)
			require.NoError(t, err)
			db := NewTodoRecord(pool)

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
			pool, err := OpenDB(*dataSourceName)
			require.NoError(t, err)
			db := NewTodoRecord(pool)

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

func TestDB_DeleteAll(t *testing.T) {
	type args struct {
		todos []models.TodoRecord
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				todos: []models.TodoRecord{
					{
						Title:     "test",
						Completed: true,
						Order:     23,
					},
					{
						Title:     "test2",
						Completed: false,
						Order:     42,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, err := OpenDB(*dataSourceName)
			require.NoError(t, err)

			db := NewTodoRecord(pool)
			for _, todo := range tt.args.todos {
				_, err2 := db.Create(todo)
				require.NoError(t, err2)
			}

			err = db.DeleteAll()
			require.NoError(t, err)

			todos, err := db.GetAll(models.Query{})
			require.NoError(t, err)

			assert.Empty(t, todos)
		})
	}
}

func TestDB_DeleteSingle(t *testing.T) {
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
			pool, err := OpenDB(*dataSourceName)
			require.NoError(t, err)
			db := NewTodoRecord(pool)

			id, err := db.Create(tt.args.todo)
			require.NoError(t, err)

			err = db.DeleteSingle(id)
			require.NoError(t, err)

			todo, err := db.GetSingle(id)

			assert.Equal(t, models.TodoRecord{}, todo)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
