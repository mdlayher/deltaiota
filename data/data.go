// Package data provides the database abstraction layer and helpers
// for the Phi Mu Alpha Sinfonia - Delta Iota chapter website.
package data

import (
	"database/sql"
	"errors"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

var (
	// ErrMultipleResults is returned when a query should return only zero or a single
	// result, but returns two or more results.
	ErrMultipleResults = errors.New("db: multiple results returned")
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

// Begin starts a transaction on this database instance.
func (db *DB) Begin() (*Tx, error) {
	// Start a transaction on underlying database
	dbtx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}

	// Return wrapped transaction
	return &Tx{
		Tx: dbtx,
	}, nil
}

// withTx creates a new wrapped transaction, invokes an input closure, and
// commits or rolls back the transaction, depending on the result of the
// closure invocation.
func (db *DB) withTx(fn func(tx *Tx) error) error {
	// Start a wrapped transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Invoke input closure, passing in wrapped transaction
	if err := fn(tx); err != nil {
		// Failure, attempt to roll back transaction
		if rErr := tx.Rollback(); rErr != nil {
			return rErr
		}

		// Return error from closure
		return err
	}

	// Attempt to commit transaction
	return tx.Commit()
}

// withPreparedRows creates a new prepared statement with the input SQL query,
// invokes an input closure containing SQL rows, and handles cleanup of rows
// and prepared statements once the closure is complete.
func (db *DB) withPreparedRows(query string, fn func(rows *Rows) error, args ...interface{}) error {
	// Prepare statement using input query
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	// Perform input query, sending arguments from caller
	rows, err := stmt.Query(args...)
	if err != nil {
		return err
	}

	// Invoke input closure with wrapped Rows type, capturing return value for later
	fnErr := fn(&Rows{
		Rows: rows,
	})

	// Close rows
	if err := rows.Close(); err != nil {
		return err
	}

	// Check errors from rows
	if err := rows.Err(); err != nil {
		return err
	}

	// Close prepared statement
	if err := stmt.Close(); err != nil {
		return err
	}

	// Return result of closure
	return fnErr
}

// Tx is a wrapped database transaction, which provides additional methods
// for interacting directly with custom types.
type Tx struct {
	*sql.Tx
}

// Rows is a wrapped set of database rows, which provides additional methods
// for interacting directly with custom types.
type Rows struct {
	*sql.Rows
}
