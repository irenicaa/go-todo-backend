package db

import (
	"database/sql"

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
