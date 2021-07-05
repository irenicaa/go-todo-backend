// +build integration

package db

import (
	"database/sql"
	"flag"
	"strconv"
	"testing"
	"time"

	"github.com/irenicaa/go-todo-backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var dataSourceName = flag.String(
	"dataSourceName",
	DefaultDataSourceName,
	"DB connection string",
)

func TestDB_withGetting(t *testing.T) {
	pool, err := OpenDB(*dataSourceName)
	require.NoError(t, err)
	db := NewTodoRecord(pool)

	err = db.DeleteAll()
	require.NoError(t, err)

	var createdTodos []models.TodoRecord
	for i := 0; i <= 10; i++ {
		originalTodo := models.TodoRecord{
			Date:      time.Date(2006, time.January, 2+i, 0, 0, 0, 0, time.UTC),
			Title:     "test" + strconv.Itoa(i),
			Completed: true,
			Order:     i,
		}

		id, err2 := db.Create(originalTodo)
		require.NoError(t, err2)
		originalTodo.ID = id

		createdTodos = append(createdTodos, originalTodo)
	}
	// reversing
	for i, j := 0, len(createdTodos)-1; i < j; i, j = i+1, j-1 {
		createdTodos[i], createdTodos[j] = createdTodos[j], createdTodos[i]
	}

	gotTodos, err := db.GetAll(models.Query{})
	require.NoError(t, err)
	for index := range gotTodos {
		gotTodos[index].Date = gotTodos[index].Date.In(time.UTC)
	}

	assert.Equal(t, createdTodos, gotTodos)
}

func TestDB_withModifying(t *testing.T) {
	tests := []struct {
		name         string
		originalTodo models.TodoRecord
		action       func(t *testing.T, db TodoRecord, todoID int)
		wantTodo     models.TodoRecord
	}{
		{
			name: "creation",
			originalTodo: models.TodoRecord{
				Date:      time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC),
				Title:     "test",
				Completed: true,
				Order:     23,
			},
			action: func(t *testing.T, db TodoRecord, todoID int) {},
			wantTodo: models.TodoRecord{
				Date:      time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC),
				Title:     "test",
				Completed: true,
				Order:     23,
			},
		},
		{
			name: "updating",
			originalTodo: models.TodoRecord{
				Date:      time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC),
				Title:     "test",
				Completed: true,
				Order:     23,
			},
			action: func(t *testing.T, db TodoRecord, todoID int) {
				newTodo := models.TodoRecord{
					Date:      time.Date(2006, time.January, 3, 0, 0, 0, 0, time.UTC),
					Title:     "test2",
					Completed: false,
					Order:     42,
				}

				err := db.Update(todoID, newTodo)
				require.NoError(t, err)
			},
			wantTodo: models.TodoRecord{
				Date:      time.Date(2006, time.January, 3, 0, 0, 0, 0, time.UTC),
				Title:     "test2",
				Completed: false,
				Order:     42,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, err := OpenDB(*dataSourceName)
			require.NoError(t, err)
			db := NewTodoRecord(pool)

			id, err := db.Create(tt.originalTodo)
			require.NoError(t, err)

			tt.action(t, db, id)

			gotTodo, err := db.GetSingle(id)
			require.NoError(t, err)
			gotTodo.Date = gotTodo.Date.In(time.UTC)

			tt.wantTodo.ID = id
			assert.Equal(t, tt.wantTodo, gotTodo)
		})
	}
}

func TestDB_DeleteSingle(t *testing.T) {
	pool, err := OpenDB(*dataSourceName)
	require.NoError(t, err)
	db := NewTodoRecord(pool)

	originalTodo := models.TodoRecord{
		Date:      time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC),
		Title:     "test",
		Completed: true,
		Order:     23,
	}

	id, err := db.Create(originalTodo)
	require.NoError(t, err)

	err = db.DeleteSingle(id)
	require.NoError(t, err)

	todo, err := db.GetSingle(id)

	assert.Equal(t, models.TodoRecord{}, todo)
	assert.Equal(t, sql.ErrNoRows, err)
}
