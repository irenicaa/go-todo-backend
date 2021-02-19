package db

import (
	"database/sql"
	"fmt"

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
		return DB{}, fmt.Errorf("unable to create a pool of DB connections: %v", err)
	}

	if err := pool.Ping(); err != nil {
		return DB{}, fmt.Errorf("unable to ping the DB: %v", err)
	}

	db := DB{pool: pool}
	return db, nil
}

// GetAll ...
func (db DB) GetAll() ([]models.TodoRecord, error) {
	rows, err := db.pool.Query(`SELECT * FROM todo_records`)
	if err != nil {
		return nil, fmt.Errorf("unable to create a cursor: %v", err)
	}
	defer rows.Close()

	var todos []models.TodoRecord
	for rows.Next() {
		var todo models.TodoRecord
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.Order)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal the row: %v", err)
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

// GetSingle ...
func (db DB) GetSingle(id int) (models.TodoRecord, error) {
	var todo models.TodoRecord
	err := db.pool.
		QueryRow(`SELECT * FROM todo_records WHERE id = $1`, id).
		Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.Order)
	return todo, err
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

// Update ...
func (db DB) Update(id int, todo models.TodoRecord) error {
	_, err := db.pool.Exec(
		`UPDATE todo_records
		SET title = $1, completed = $2, "order" = $3
		WHERE id = $4`,
		todo.Title,
		todo.Completed,
		todo.Order,
		id,
	)
	return err
}

// Delete ...
func (db DB) Delete(id int) error {
	_, err := db.pool.Exec(`DELETE FROM todo_records WHERE id = $1`, id)
	return err
}
