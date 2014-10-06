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

// CopyFrom copies fields from an input User into the receiving User struct.
func (u *User) CopyFrom(user *User) {
	u.Username = user.Username
	u.FirstName = user.FirstName
	u.LastName = user.LastName
	u.Email = user.Email
	u.Phone = user.Phone
	u.password = user.password
}

// SetPassword hashes the input password using bcrypt, storing the password
// within the receiving User struct.
func (u *User) SetPassword(password string) error {
	// Check for empty password
	if password == "" {
		return &EmptyFieldError{
			Field: "password",
		}
	}

	// Generate password hash using bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	if err != nil {
		return err
	}
	u.password = string(hash)

	return nil
}

// SetTestPassword directly stores the input password in the receiving
// User struct.
//
// This method should ONLY be used for testing, to save time on bcrypt hashing
// for test user creation.
func (u *User) SetTestPassword(password string) {
	u.password = "test-password-" + password
	return
}

// SQLReadFields returns the correct field order to scan SQL row results into the
// receiving User struct.
func (u *User) SQLReadFields() []interface{} {
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

// SQLWriteFields returns the correct field order for SQL write actions (such as
// insert or update), for the receiving User struct.
func (u *User) SQLWriteFields() []interface{} {
	return []interface{}{
		u.Username,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,

		u.password,

		// Last argument for WHERE clause
		u.ID,
	}
}

// Validate verifies that all fields for the receiving User struct contain
// valid input.
func (u *User) Validate() error {
	// Check for required fields
	if u.Username == "" {
		return &EmptyFieldError{
			Field: "username",
		}
	}
	if u.FirstName == "" {
		return &EmptyFieldError{
			Field: "firstName",
		}
	}
	if u.LastName == "" {
		return &EmptyFieldError{
			Field: "lastName",
		}
	}
	if u.Email == "" {
		return &EmptyFieldError{
			Field: "email",
		}
	}
	if u.password == "" {
		return &EmptyFieldError{
			Field: "password",
		}
	}

	// Perform basic validation of email address
	address, err := mail.ParseAddress(u.Email)
	if err != nil {
		return &InvalidFieldError{
			Field:   "email",
			Err:     err,
			Details: "could not parse valid email address",
		}
	}
	u.Email = address.Address

	return nil
}
