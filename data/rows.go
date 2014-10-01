package data

import (
	"database/sql"

	"github.com/mdlayher/deltaiota/data/models"
)

// Rows is a wrapped set of database rows, which provides additional methods
// for interacting directly with custom types.
type Rows struct {
	*sql.Rows
}

// ScanUsers returns a slice of Users from wrapped rows.
func (r *Rows) ScanUsers() ([]models.User, error) {
	// Iterate all returned rows
	var users []models.User
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

		// Dereference and append user to output slice
		users = append(users, *u)
	}

	return users, nil
}
