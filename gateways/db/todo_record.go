package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	utilmodels "github.com/irenicaa/go-http-utils/models"
	"github.com/irenicaa/go-todo-backend/v2/models"
)

// TodoRecord ...
type TodoRecord struct {
	pool *sql.DB
}

// NewTodoRecord ...
func NewTodoRecord(pool *sql.DB) TodoRecord {
	return TodoRecord{pool: pool}
}

// GetAll ...
func (db TodoRecord) GetAll(query models.Query) ([]models.TodoRecord, error) {
	sql := "SELECT * FROM todo_records"
	var args []interface{}

	sql += " WHERE TRUE"
	var argNumber int
	if query.MinimalDate != (utilmodels.Date{}) {
		argNumber++
		sql += " AND date >= $" + strconv.Itoa(argNumber)
		args = append(args, time.Time(query.MinimalDate))
	}
	if query.MaximalDate != (utilmodels.Date{}) {
		argNumber++
		sql += " AND date <= $" + strconv.Itoa(argNumber)
		args = append(args, time.Time(query.MaximalDate))
	}
	if query.TitleFragment != "" {
		argNumber++
		sql += " AND lower(title) LIKE $" + strconv.Itoa(argNumber)
		args = append(args, "%"+strings.ToLower(query.TitleFragment)+"%")
	}

	sql += ` ORDER BY "date" DESC, "order", id`
	if query.Pagination != (models.Pagination{}) {
		sql += fmt.Sprintf(
			" OFFSET %d LIMIT %d",
			(query.Pagination.Page-1)*query.Pagination.PageSize,
			query.Pagination.PageSize,
		)
	}

	rows, err := db.pool.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("unable to create a cursor: %v", err)
	}
	defer rows.Close()

	var todos []models.TodoRecord
	for rows.Next() {
		var todo models.TodoRecord
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Completed,
			&todo.Order,
			&todo.Date,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal the row: %v", err)
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

// GetSingle ...
func (db TodoRecord) GetSingle(id int) (models.TodoRecord, error) {
	var todo models.TodoRecord
	err := db.pool.
		QueryRow("SELECT * FROM todo_records WHERE id = $1", id).
		Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.Order, &todo.Date)
	return todo, err
}

// Create ...
func (db TodoRecord) Create(todo models.TodoRecord) (id int, err error) {
	err = db.pool.
		QueryRow(
			`INSERT INTO todo_records (title, completed, "order", "date")
			VALUES ($1, $2, $3, $4)
			RETURNING id`,
			todo.Title,
			todo.Completed,
			todo.Order,
			todo.Date,
		).
		Scan(&id)
	return id, err
}

// Update ...
func (db TodoRecord) Update(id int, todo models.TodoRecord) error {
	_, err := db.pool.Exec(
		`UPDATE todo_records
		SET title = $1, completed = $2, "order" = $3, "date" = $4
		WHERE id = $5`,
		todo.Title,
		todo.Completed,
		todo.Order,
		todo.Date,
		id,
	)
	return err
}

// DeleteAll ...
func (db TodoRecord) DeleteAll() error {
	_, err := db.pool.Exec("DELETE FROM todo_records")
	return err
}

// DeleteSingle ...
func (db TodoRecord) DeleteSingle(id int) error {
	_, err := db.pool.Exec("DELETE FROM todo_records WHERE id = $1", id)
	return err
}
