package models

import (
	"errors"
	"net/mail"
	"time"

	"code.google.com/p/go.crypto/bcrypt"
)

var (
	// ErrInvalidPassword is returned when password authentication fails for a
	// specified User.
	ErrInvalidPassword = errors.New("invalid password")
)

// User represents a user of the application.
type User struct {
	ID        uint64 `db:"id" json:"id"`
	Username  string `db:"username" json:"username"`
	FirstName string `db:"first_name" json:"firstName"`
	LastName  string `db:"last_name" json:"lastName"`
	Email     string `db:"email" json:"email"`
	Phone     string `db:"phone" json:"phone"`
	Password  string `db:"password" json:"password,omitempty"`
}

// CopyFrom copies fields from an input User into the receiving User struct.
func (u *User) CopyFrom(user *User) {
	u.Username = user.Username
	u.FirstName = user.FirstName
	u.LastName = user.LastName
	u.Email = user.Email
	u.Phone = user.Phone
	u.Password = user.Password
}

// NewSession generates a new Session for this user.
func (u *User) NewSession(expire time.Time) (*Session, error) {
	return NewSession(u.ID, u.Password, expire)
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
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)

	return nil
}

// TryPassword attempts to verify the input password against the receiving User's
// current password.
func (u *User) TryPassword(password string) error {
	// Attempt to hash password
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	// Check for bcrypt-specific password failure, return more generic failure
	// (other packages should not have to import or know about bcrypt to know
	// the password was incorrect)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrInvalidPassword
	}

	// Return other errors, or no error
	return err
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
		&u.Password,
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
		u.Password,

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
	if u.Password == "" {
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
