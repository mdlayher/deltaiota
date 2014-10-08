package models

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"time"

	"code.google.com/p/go.crypto/pbkdf2"
)

// Session represents an application session.
type Session struct {
	ID     uint64 `db:"id" json:"id"`
	UserID uint64 `db:"user_id" json:"userId"`
	Key    string `db:"key" json:"key"`
	Expire uint64 `db:"expire" json:"expire"`
}

// NewSession creates a new session for the specified user ID, which will
// expire at the specified time.
func NewSession(userID uint64, password string, expire time.Time) (*Session, error) {
	// Create new session for input user ID
	s := &Session{
		UserID: userID,
		Expire: uint64(expire.Unix()),
	}

	// Generate salt for use with PBKDF2
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	// Use PBKDF2 to generate a session key based off the user's password
	s.Key = fmt.Sprintf("%x", pbkdf2.Key([]byte(password), salt, 4096, 16, sha1.New))

	// Return generated session
	return s, nil
}

// SQLReadFields returns the correct field order to scan SQL row results into the
// receiving Session struct.
func (u *Session) SQLReadFields() []interface{} {
	return []interface{}{
		&u.ID,
		&u.UserID,
		&u.Key,
		&u.Expire,
	}
}

// SQLWriteFields returns the correct field order for SQL write actions (such as
// insert or update), for the receiving Session struct.
func (u *Session) SQLWriteFields() []interface{} {
	return []interface{}{
		u.UserID,
		u.Key,
		u.Expire,

		// Last argument for WHERE clause
		u.ID,
	}
}
