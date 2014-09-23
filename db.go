package deltaiota

import (
	"database/sql"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// DB provides the database abstraction layer for the application.
type DB struct {
	*sql.DB
}

// Open opens and initializes a database instance.
func (db *DB) Open(driver string, dsn string) error {
	// Open database
	d, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	db.DB = d

	return nil
}

// Close closes and cleans up a database instance.
func (db *DB) Close() error {
	return db.DB.Close()
}
