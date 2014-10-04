package data

import (
	"database/sql"

	"github.com/mdlayher/deltaiota/data/models"
)

const (
	// sqlSelectAllUsers is the SQL statement used to select all Users
	sqlSelectAllUsers = `
		SELECT * FROM users;
	`
	// sqlSelectUserByID is the SQL statement used to select a single user by ID
	sqlSelectUserByID = `
		SELECT * FROM users WHERE id = ?;
	`

	// sqlInsertUser is the SQL statement used to insert a new User
	sqlInsertUser = `
		INSERT INTO users (
			"username"
			, "first_name"
			, "last_name"
			, "email"
			, "phone"
			, "password"
		) VALUES (?, ?, ?, ?, ?, ?);
	`
)

// SaveUser starts a transaction, inserts a new User, and attempts to commit
// the transaction.
func (db *DB) InsertUser(u *models.User) error {
	return db.withTx(func(tx *Tx) error {
		return tx.InsertUser(u)
	})
}

// SelectAllUsers returns a slice of all Users from the database.
func (db *DB) SelectAllUsers() ([]*models.User, error) {
	return db.selectUsers(sqlSelectAllUsers)
}

// SelectUserByID returns a single User by ID from the database.
func (db *DB) SelectUserByID(id uint64) (*models.User, error) {
	// Fetch users with matching ID
	users, err := db.selectUsers(sqlSelectUserByID, id)
	if err != nil {
		return nil, err
	}

	// Verify only 0 or 1 user returned
	if len(users) == 0 {
		return nil, sql.ErrNoRows
	} else if len(users) == 1 {
		return users[0], nil
	}

	// More than one result returned
	return nil, ErrMultipleResults
}

// selectUsers returns a slice of Users from the database, based upon an input
// SQL query and arguments
func (db *DB) selectUsers(query string, args ...interface{}) ([]*models.User, error) {
	// Slice of users to return
	var users []*models.User

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

// InsertUser inserts a new User in the context of the current transaction.
func (tx *Tx) InsertUser(u *models.User) error {
	// Execute SQL to insert User
	result, err := tx.Tx.Exec(sqlInsertUser, u.Username, u.FirstName, u.LastName, u.Email, u.Phone, u.Password())
	if err != nil {
		return err
	}

	// Retrieve generated ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Store generated ID
	u.ID = uint64(id)
	return nil
}

// ScanUsers returns a slice of Users from wrapped rows.
func (r *Rows) ScanUsers() ([]*models.User, error) {
	// Iterate all returned rows
	var users []*models.User
	for r.Rows.Next() {
		// Scan new user into struct, using specified fields
		u := new(models.User)
		if err := r.Rows.Scan(u.SQLFields()...); err != nil {
			return nil, err
		}

		// Discard any nil results
		if u == nil {
			continue
		}

		// Append user to output slice
		users = append(users, u)
	}

	return users, nil
}
