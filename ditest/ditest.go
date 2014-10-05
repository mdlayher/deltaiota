// Package ditest provides common functionality for testing the Phi Mu Alpha
// Sinfonia - Delta Iota chapter website.
package ditest

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mdlayher/deltaiota/bindata"
	"github.com/mdlayher/deltaiota/data"
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
