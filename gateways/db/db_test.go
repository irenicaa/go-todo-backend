package db

import (
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
