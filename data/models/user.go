package models

// User represents a user of the application.
type User struct {
	ID        uint64 `db:"id" json:"id"`
	Username  string `db:"username" json:"username"`
	FirstName string `db:"first_name" json:"firstName"`
	LastName  string `db:"last_name" json:"lastName"`
	Email     string `db:"email" json:"email"`
	Phone     string `db:"phone" json:"phone"`

	password string `db:"password"`
	salt     string `db:"salt"`
}

// SQLFields returns the correct field order to scan SQL row results into the
// receiving User struct.
func (u *User) SQLFields() []interface{} {
	return []interface{}{
		&u.ID,
		&u.Username,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Phone,

		&u.password,
		&u.salt,
	}
}
