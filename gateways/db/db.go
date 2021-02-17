package db

import (
	"database/sql"

	"github.com/irenicaa/go-todo-backend/models"
	_ "github.com/lib/pq"
)

// DB ...
type DB struct {
	pool *sql.DB
}

// OpenDB ...
func OpenDB(dataSourceName string) (DB, error) {
	pool, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return DB{}, err
	}

	db := DB{pool: pool}
	return db, nil
}

// Create ...
func (db DB) Create(todo models.TodoRecord) (id int, err error) {
	err = db.pool.
		QueryRow(
			`INSERT INTO todo_records (title, completed, "order")
			VALUES ($1, $2, $3)
			RETURNING id`,
			todo.Title,
			todo.Completed,
			todo.Order,
		).
		Scan(&id)
	return id, err
}
