// Package ditest provides common functionality for testing the Phi Mu Alpha
// Sinfonia - Delta Iota chapter website.
package ditest

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/mdlayher/deltaiota/bindata"
	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/data/models"
)

// WithTemporaryDB generates a temporary copy of the embedded sqlite3 database from
// bindata, invokes an input closure, and cleans up the temporary database once
// the closure returns.
func WithTemporaryDB(fn func(db *data.DB)) error {
	// Create temporary working directory in system temporary directory
	tmpDir, err := ioutil.TempDir(os.TempDir(), "ditest-")
	if err != nil {
		return err
	}

	// Retrieve sqlite3 database asset
	asset, err := bindata.Asset("res/sqlite/deltaiota.db")
	if err != nil {
		return err
	}

	// Write sqlite3 database asset to temporary file
	dbPath := filepath.Join(tmpDir, "ditest.db")
	if err := ioutil.WriteFile(dbPath, asset, 0644); err != nil {
		return err
	}

	// Open database connection
	didb := &data.DB{}
	if err := didb.Open("sqlite3", dbPath); err != nil {
		return err
	}

	// Invoke input closure with database
	fn(didb)

	// Close database
	if err := didb.Close(); err != nil {
		return err
	}

	// Remove temporary directory and database
	return os.RemoveAll(tmpDir)
}

// MockUser generates a single User with mock data, used for testing.
// The user is randomly generated, but is not guaranteed to be unique.
func MockUser() *models.User {
	// Generate user
	user := &models.User{
		Username:  RandomString(10),
		FirstName: RandomString(10),
		LastName:  RandomString(10),
		Email:     fmt.Sprintf("%s@%s.com", RandomString(6), RandomString(6)),
	}

	// Generate test password, only used for duration of test
	user.SetTestPassword(RandomString(10))
	return user
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
