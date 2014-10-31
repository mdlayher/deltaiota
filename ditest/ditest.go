// Package ditest provides common functionality for testing the Phi Mu Alpha
// Sinfonia - Delta Iota chapter website.
package ditest

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/mdlayher/deltaiota/bindata"
	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/data/models"
)

// WithTemporaryDB generates a temporary, in-memory copy of the deltaiota sqlite3
// database from bindata SQL schema, invokes an input closure, and destroys the
// in-memory database once the closure returns.
func WithTemporaryDB(fn func(db *data.DB) error) error {
	// Retrieve sqlite3 database schema asset
	asset, err := bindata.Asset("res/sqlite/deltaiota.sql")
	if err != nil {
		return err
	}

	// Open in-memory database
	didb := &data.DB{}
	if err := didb.Open("sqlite3", ":memory:"); err != nil {
		return err
	}

	// Execute schema to build database
	if _, err := didb.Exec(string(asset)); err != nil {
		return err
	}

	// Invoke input closure with database
	fnErr := fn(didb)

	// Close and destroy database
	if err := didb.Close(); err != nil {
		return err
	}

	// Return error from closure
	return fnErr
}

// WithTemporaryDBNew is a temporary scaffolding function which will be used for
// refactoring tests, and will eventually replace WithTemporaryDB.
func WithTemporaryDBNew(t *testing.T, fn func(t *testing.T, db *data.DB)) {
	// Retrieve sqlite3 database schema asset
	asset, err := bindata.Asset("res/sqlite/deltaiota.sql")
	if err != nil {
		t.Fatal(err)
	}

	// Open in-memory database
	didb := &data.DB{}
	if err := didb.Open("sqlite3", ":memory:"); err != nil {
		t.Fatal(err)
	}

	// Execute schema to build database
	if _, err := didb.Exec(string(asset)); err != nil {
		t.Fatal(err)
	}

	// Invoke input closure with test and database
	fn(t, didb)

	// Close and destroy database
	if err := didb.Close(); err != nil {
		t.Fatal(err)
	}
}

// MockUser generates a single User with mock data, used for testing.
// The user is randomly generated, but is not guaranteed to be unique.
func MockUser() *models.User {
	return &models.User{
		Username:  RandomString(10),
		FirstName: RandomString(10),
		LastName:  RandomString(10),
		Email:     fmt.Sprintf("%s@%s.com", RandomString(6), RandomString(6)),
		Password:  RandomString(10),
	}
}

// RandomString generates a random string of length n.
// Adapter from: http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func RandomString(n int) string {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Random letters slice
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	// Generate slice of length n
	str := make([]rune, n)
	for i := range str {
		str[i] = letters[rand.Intn(len(letters))]
	}

	return string(str)
}
