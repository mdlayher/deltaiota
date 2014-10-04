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

// Password returns the value of the unexported password field.
// This method is used for interactions with the database.
func (u *User) Password() string {
	return u.password
}

// Salt returns the value of the unexported salt field.
// This method is used for interactions with the database.
func (u *User) Salt() string {
	return u.salt
}

// SetPassword generates a new salt and hashes the input password, storing both fields
// within the receiving User struct.
func (u *User) SetPassword(password string) error {
	// Check for empty password
	if password == "" {
		return ErrInvalid
	}

	// TODO(mdlayher): bcrypt password hash, crypto/rand salt generation
	u.password = password
	u.salt = ""

	return nil
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

// Validate verifies that all fields for the receiving User struct contain
// valid input.
func (u *User) Validate() error {
	// Check for required fields
	if u.Username == "" || u.FirstName == "" || u.LastName == "" || u.password == "" {
		return ErrInvalid
	}

	return nil
}
