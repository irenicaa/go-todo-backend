// +build integration

package db

import (
	"database/sql"
	"flag"
	"fmt"
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

func TestTodoRecord_withGetting(t *testing.T) {
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

func TestTodoRecord_withQuerying(t *testing.T) {
	tests := []struct {
		name          string
		originalTodos []models.TodoRecord
		query         models.Query
		wantTodos     []models.TodoRecord
	}{
		{
			name: "with filtration by the minimal date",
			originalTodos: func() []models.TodoRecord {
				var originalTodos []models.TodoRecord
				for i := 0; i <= 10; i++ {
					originalTodo := models.TodoRecord{
						Date: time.Date(
							2006, time.January, 2+i,
							0, 0, 0, 0,
							time.UTC,
						),
						Title:     "test" + strconv.Itoa(i),
						Completed: true,
						Order:     i,
					}
					originalTodos = append(originalTodos, originalTodo)
				}

				return originalTodos
			}(),
			query: models.Query{
				MinimalDate: models.Date(time.Date(
					2006, time.January, 7,
					0, 0, 0, 0,
					time.UTC,
				)),
			},
			wantTodos: func() []models.TodoRecord {
				var wantTodos []models.TodoRecord
				for i := 10; i >= 5; i-- {
					wantTodo := models.TodoRecord{
						Date: time.Date(
							2006, time.January, 2+i,
							0, 0, 0, 0,
							time.UTC,
						),
						Title:     "test" + strconv.Itoa(i),
						Completed: true,
						Order:     i,
					}
					wantTodos = append(wantTodos, wantTodo)
				}

				return wantTodos
			}(),
		},
		{
			name: "with filtration by the maximal date",
			originalTodos: func() []models.TodoRecord {
				var originalTodos []models.TodoRecord
				for i := 0; i <= 10; i++ {
					originalTodo := models.TodoRecord{
						Date: time.Date(
							2006, time.January, 2+i,
							0, 0, 0, 0,
							time.UTC,
						),
						Title:     "test" + strconv.Itoa(i),
						Completed: true,
						Order:     i,
					}
					originalTodos = append(originalTodos, originalTodo)
				}

				return originalTodos
			}(),
			query: models.Query{
				MaximalDate: models.Date(time.Date(
					2006, time.January, 7,
					0, 0, 0, 0,
					time.UTC,
				)),
			},
			wantTodos: func() []models.TodoRecord {
				var wantTodos []models.TodoRecord
				for i := 5; i >= 0; i-- {
					wantTodo := models.TodoRecord{
						Date: time.Date(
							2006, time.January, 2+i,
							0, 0, 0, 0,
							time.UTC,
						),
						Title:     "test" + strconv.Itoa(i),
						Completed: true,
						Order:     i,
					}
					wantTodos = append(wantTodos, wantTodo)
				}

				return wantTodos
			}(),
		},
		{
			name: "with filtration by the date range",
			originalTodos: func() []models.TodoRecord {
				var originalTodos []models.TodoRecord
				for i := 0; i <= 10; i++ {
					originalTodo := models.TodoRecord{
						Date: time.Date(
							2006, time.January, 2+i,
							0, 0, 0, 0,
							time.UTC,
						),
						Title:     "test" + strconv.Itoa(i),
						Completed: true,
						Order:     i,
					}
					originalTodos = append(originalTodos, originalTodo)
				}

				return originalTodos
			}(),
			query: models.Query{
				MinimalDate: models.Date(time.Date(
					2006, time.January, 5,
					0, 0, 0, 0,
					time.UTC,
				)),
				MaximalDate: models.Date(time.Date(
					2006, time.January, 9,
					0, 0, 0, 0,
					time.UTC,
				)),
			},
			wantTodos: func() []models.TodoRecord {
				var wantTodos []models.TodoRecord
				for i := 7; i >= 3; i-- {
					wantTodo := models.TodoRecord{
						Date: time.Date(
							2006, time.January, 2+i,
							0, 0, 0, 0,
							time.UTC,
						),
						Title:     "test" + strconv.Itoa(i),
						Completed: true,
						Order:     i,
					}
					wantTodos = append(wantTodos, wantTodo)
				}

				return wantTodos
			}(),
		},
		{
			name: "with search by the title fragment",
			originalTodos: func() []models.TodoRecord {
				var originalTodos []models.TodoRecord
				for i := 0; i <= 10; i++ {
					var mark string
					if i%2 == 0 {
						mark = "even"
					} else {
						mark = "odd"
					}

					originalTodo := models.TodoRecord{
						Date: time.Date(
							2006, time.January, 2+i,
							0, 0, 0, 0,
							time.UTC,
						),
						Title:     fmt.Sprintf("test%d (%s)", i, mark),
						Completed: true,
						Order:     i,
					}
					originalTodos = append(originalTodos, originalTodo)
				}

				return originalTodos
			}(),
			query: models.Query{TitleFragment: "even"},
			wantTodos: func() []models.TodoRecord {
				var wantTodos []models.TodoRecord
				for i := 10; i >= 0; i -= 2 {
					wantTodo := models.TodoRecord{
						Date: time.Date(
							2006, time.January, 2+i,
							0, 0, 0, 0,
							time.UTC,
						),
						Title:     fmt.Sprintf("test%d (even)", i),
						Completed: true,
						Order:     i,
					}
					wantTodos = append(wantTodos, wantTodo)
				}

				return wantTodos
			}(),
		},
		{
			name: "with pagination",
			originalTodos: func() []models.TodoRecord {
				var originalTodos []models.TodoRecord
				for i := 0; i <= 10; i++ {
					originalTodo := models.TodoRecord{
						Date: time.Date(
							2006, time.January, 2+i,
							0, 0, 0, 0,
							time.UTC,
						),
						Title:     "test" + strconv.Itoa(i),
						Completed: true,
						Order:     i,
					}
					originalTodos = append(originalTodos, originalTodo)
				}

				return originalTodos
			}(),
			query: models.Query{Pagination: models.Pagination{PageSize: 2, Page: 3}},
			wantTodos: func() []models.TodoRecord {
				var wantTodos []models.TodoRecord
				for i := 6; i >= 5; i-- {
					wantTodo := models.TodoRecord{
						Date: time.Date(
							2006, time.January, 2+i,
							0, 0, 0, 0,
							time.UTC,
						),
						Title:     "test" + strconv.Itoa(i),
						Completed: true,
						Order:     i,
					}
					wantTodos = append(wantTodos, wantTodo)
				}

				return wantTodos
			}(),
		},
		{
			name: "with all parameters",
			originalTodos: func() []models.TodoRecord {
				var originalTodos []models.TodoRecord
				for i := 0; i <= 20; i++ {
					var mark string
					if i%2 == 0 {
						mark = "even"
					} else {
						mark = "odd"
					}

					originalTodo := models.TodoRecord{
						Date: time.Date(
							2006, time.January, 2+i,
							0, 0, 0, 0,
							time.UTC,
						),
						Title:     fmt.Sprintf("test%d (%s)", i, mark),
						Completed: true,
						Order:     i,
					}
					originalTodos = append(originalTodos, originalTodo)
				}

				return originalTodos
			}(),
			query: models.Query{
				MinimalDate: models.Date(time.Date(
					2006, time.January, 5,
					0, 0, 0, 0,
					time.UTC,
				)),
				MaximalDate: models.Date(time.Date(
					2006, time.January, 19,
					0, 0, 0, 0,
					time.UTC,
				)),
				TitleFragment: "even",
				Pagination:    models.Pagination{PageSize: 2, Page: 3},
			},
			wantTodos: func() []models.TodoRecord {
				var wantTodos []models.TodoRecord
				for i := 8; i >= 6; i -= 2 {
					wantTodo := models.TodoRecord{
						Date: time.Date(
							2006, time.January, 2+i,
							0, 0, 0, 0,
							time.UTC,
						),
						Title:     fmt.Sprintf("test%d (even)", i),
						Completed: true,
						Order:     i,
					}
					wantTodos = append(wantTodos, wantTodo)
				}

				return wantTodos
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, err := OpenDB(*dataSourceName)
			require.NoError(t, err)
			db := NewTodoRecord(pool)

			err = db.DeleteAll()
			require.NoError(t, err)

			for _, originalTodo := range tt.originalTodos {
				_, err2 := db.Create(originalTodo)
				require.NoError(t, err2)
			}

			gotTodos, err := db.GetAll(tt.query)
			require.NoError(t, err)
			for index := range gotTodos {
				gotTodos[index].ID = 0
				gotTodos[index].Date = gotTodos[index].Date.In(time.UTC)
			}

			assert.Equal(t, tt.wantTodos, gotTodos)
		})
	}
}

func TestTodoRecord_withModifying(t *testing.T) {
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

func TestTodoRecord_withDeleting(t *testing.T) {
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
