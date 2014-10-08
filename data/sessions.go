package data

import "github.com/mdlayher/deltaiota/data/models"

const (
	// sqlInsertSession is the SQL statement used to insert a new Session
	sqlInsertSession = `
		INSERT INTO sessions (
			"user_id"
			, "key"
			, "expire"
		) VALUES (?, ?, ?);
	`

	// sqlDeleteSession is the SQL statement used to delete an existing Session
	sqlDeleteSession = `
		DELETE FROM sessions WHERE id = ?;
	`
)

// InsertSession starts a transaction, inserts a new Session, and attempts to commit
// the transaction.
func (db *DB) InsertSession(u *models.Session) error {
	return db.withTx(func(tx *Tx) error {
		return tx.InsertSession(u)
	})
}

// DeleteSession starts a transaction, deletes the input Session by its ID, and attempts
// to commit the transaction.
func (db *DB) DeleteSession(s *models.Session) error {
	return db.withTx(func(tx *Tx) error {
		return tx.DeleteSession(s)
	})
}

// InsertSession inserts a new Session in the context of the current transaction.
func (tx *Tx) InsertSession(u *models.Session) error {
	// Execute SQL to insert Session
	result, err := tx.Tx.Exec(sqlInsertSession, u.SQLWriteFields()...)
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

// DeleteSession updates the input Session by its ID, in the context of the
// current transaction.
func (tx *Tx) DeleteSession(s *models.Session) error {
	_, err := tx.Tx.Exec(sqlDeleteSession, s.ID)
	return err
}
