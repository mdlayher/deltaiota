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

// Tx is a wrapped database transaction, which provides additional methods
// for interacting directly with custom types.
type Tx struct {
	*sql.Tx
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
		if rErr := tx.Rollback(); err != nil {
			return rErr
		}

		// Return error from closure
		return err
	}

	// Attempt to commit transaction
	return tx.Commit()
}
