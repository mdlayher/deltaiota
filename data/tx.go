package data

import (
	"database/sql"

	"github.com/mdlayher/deltaiota/data/models"
)

// Tx is a wrapped database transaction, which provides additional methods
// for interacting directly with custom types.
type Tx struct {
	*sql.Tx
}

// SaveUser inserts a new User in the context of the current transaction.
func (tx *Tx) SaveUser(u *models.User) error {
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
