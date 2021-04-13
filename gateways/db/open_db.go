package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// DefaultDataSourceName ...
const DefaultDataSourceName = "postgresql://postgres:postgres@localhost:5432" +
	"/postgres?sslmode=disable"

// OpenDB ...
func OpenDB(dataSourceName string) (*sql.DB, error) {
	pool, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("unable to create a pool of DB connections: %v", err)
	}

	return pool, nil
}
