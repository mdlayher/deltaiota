package models

// User represents a user of the application.
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// SQLFields returns the correct field order to scan SQL row results into the
// receiving User struct.
func (u *User) SQLFields() []interface{} {
	return []interface{}{
		&u.ID,
		&u.Username,
	}
}
