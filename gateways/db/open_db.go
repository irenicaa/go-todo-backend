package db

import (
	"database/sql"
	"fmt"
)

// OpenDB ...
func OpenDB(dataSourceName string) (*sql.DB, error) {
	pool, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("unable to create a pool of DB connections: %v", err)
	}

	if err := pool.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping the DB: %v", err)
	}

	return pool, nil
}
