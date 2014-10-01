package data

const (
	// sqlInsertUser is the SQL statement used to insert a new User
	sqlInsertUser = `
		INSERT INTO users (
			"username"
		) VALUES (?);
	`
)

// User represents a user of the application.
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// SaveUser starts a transaction, inserts a new User, and attempts to commit
// the transaction.
func (db *DB) SaveUser(u *User) error {
	return db.withTx(func(tx *Tx) error {
		return tx.SaveUser(u)
	})
}

// SaveUser inserts a new User in the context of the current transaction.
func (tx *Tx) SaveUser(u *User) error {
	// Execute SQL to insert User
	result, err := tx.Tx.Exec(sqlInsertUser)
	if err != nil {
		return err
	}

	// Retrieve generated ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Store generated ID
	u.ID = id
	return nil
}
