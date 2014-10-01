package data

import (
	"database/sql"

	"github.com/mdlayher/deltaiota/data/models"

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

// SaveUser starts a transaction, inserts a new User, and attempts to commit
// the transaction.
func (db *DB) SaveUser(u *models.User) error {
	return db.withTx(func(tx *Tx) error {
		return tx.SaveUser(u)
	})
}

// FetchAllUsers returns a slice of all Users from the database.
func (db *DB) FetchAllUsers() ([]models.User, error) {
	return db.fetchUsers(sqlSelectAllUsers)
}

// fetchUsers returns a slice of Users from the database, based upon an input
// SQL query and arguments
func (db *DB) fetchUsers(query string, args ...interface{}) ([]models.User, error) {
	// Slice of users to return
	var users []models.User

	// Invoke closure with prepared statement and wrapped rows,
	// passing any arguments from the caller
	err := db.withPreparedRows(query, func(rows *Rows) error {
		// Scan rows into a slice of Users
		var err error
		users, err = rows.ScanUsers()

		// Return errors from scanning
		return err
	}, args...)

	// Return any matching users and error
	return users, err
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
