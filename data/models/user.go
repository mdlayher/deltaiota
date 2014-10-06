package models

import (
	"net/mail"

	"code.google.com/p/go.crypto/bcrypt"
)

// User represents a user of the application.
type User struct {
	ID        uint64 `db:"id" json:"id"`
	Username  string `db:"username" json:"username"`
	FirstName string `db:"first_name" json:"firstName"`
	LastName  string `db:"last_name" json:"lastName"`
	Email     string `db:"email" json:"email"`
	Phone     string `db:"phone" json:"phone"`

	password string `db:"password"`
}

// Password returns the value of the unexported password field.
// This method is used for interactions with the database.
func (u *User) Password() string {
	return u.password
}

// SetPassword hashes the input password using bcrypt, storing the password
// within the receiving User struct.
func (u *User) SetPassword(password string) error {
	// Check for empty password
	if password == "" {
		return ErrEmpty
	}

	// Generate password hash using bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	if err != nil {
		return err
	}
	u.password = string(hash)

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
	}
}

// Validate verifies that all fields for the receiving User struct contain
// valid input.
func (u *User) Validate() error {
	// Check for required fields
	if u.Username == "" || u.FirstName == "" || u.LastName == "" || u.Email == "" || u.password == "" {
		return ErrEmpty
	}

	// Perform basic validation of email address
	address, err := mail.ParseAddress(u.Email)
	if err != nil {
		return ErrInvalid
	}
	u.Email = address.Address

	return nil
}
