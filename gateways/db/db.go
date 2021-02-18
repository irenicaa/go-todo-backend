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
