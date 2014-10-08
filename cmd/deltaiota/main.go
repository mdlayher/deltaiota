// Command deltaiota serves a HTTP API for the Phi Mu Alpha Sinfonia - Delta
// Iota chapter website.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/mdlayher/deltaiota/api"
	"github.com/mdlayher/deltaiota/bindata"
	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/data/models"
	"github.com/mdlayher/deltaiota/ditest"

	"github.com/stretchr/graceful"
)

const (
	// sqlite3 is the name of the sqlite3 driver for the database
	sqlite3 = "sqlite3"

	// sqlite3DBAsset is the name of the bindata asset which stores the sqlite database
	sqlite3DBAsset = "res/sqlite/deltaiota.db"

	// driver is the database/sql driver used for the database instance
	driver = sqlite3
)

// version is the current git hash, injected by the Go linker
var version string

var (
	// db is the DSN used for the database instance
	db string

	// host is the address to which the HTTP server is bound
	host string

	// timeout is the duration the server will wait before forcibly closing
	// ongoing HTTP connections
	timeout time.Duration
)

func init() {
	// Set up flags
	flag.StringVar(&db, "db", "deltaiota.db", "DSN for database instance")
	flag.StringVar(&host, "host", ":1898", "HTTP server host")
	flag.DurationVar(&timeout, "timeout", 5*time.Second, "HTTP graceful timeout duration")
}

func main() {
	// Parse all flags
	flag.Parse()

	// Report information on startup
	log.Println(fmt.Sprintf("deltaiota: starting [pid: %d] [version: %s]", os.Getpid(), version))

	// Determine if database newly created
	var created bool
	var err error

	// If database is sqlite3, perform initial setup
	if driver == sqlite3 {
		// Attempt setup, check if already created
		created, err = sqlite3Setup(db)
		if err != nil {
			log.Fatal(err)
		}

		if !created {
			log.Println("deltaiota: using sqlite3 database:", db)
		}
	}

	// Open database connection
	didb := &data.DB{}
	if err := didb.Open(driver, db); err != nil {
		log.Fatal(err)
	}

	// Perform initial user setup for sqlite3
	if driver == sqlite3 && created {
		log.Println("deltaiota: created sqlite3 database:", db)

		// Generate root user
		root := &models.User{
			Username: "root",
		}

		// Generate a random password
		password := ditest.RandomString(12)
		log.Printf("deltaiota: creating root user [password: %s]", password)
		root.SetPassword(password)

		// Save root user
		if err := didb.InsertUser(root); err != nil {
			log.Fatal(err)
		}
	}

	// Start HTTP server using deltaiota handler on specified host
	log.Println("deltaiota: listening:", host)
	if err := graceful.ListenAndServe(&http.Server{
		Addr:    host,
		Handler: api.NewServeMux(didb),
	}, timeout); err != nil {
		// Ignore error on failed "accept" when closing
		if nErr, ok := err.(*net.OpError); !ok || nErr.Op != "accept" {
			log.Fatal(err)
		}
	}

	log.Println("deltaiota: shutting down")

	// Close database connection
	if err := didb.Close(); err != nil {
		log.Fatal(err)
	}

	log.Println("deltaiota: graceful shutdown complete")
}

// sqlite3Setup performs setup routines specific to a sqlite3 database.
// On success, it returns a boolean indicating if the database was created.
// On failure, it returns an error.
func sqlite3Setup(dsn string) (bool, error) {
	// Check if database already exists at specified location
	dbPath := path.Clean(dsn)
	_, err := os.Stat(dbPath)
	if err == nil {
		// Database exists, skip setup
		return false, nil
	}

	// Any other errors, return now
	if !os.IsNotExist(err) {
		return false, err
	}

	// Retrieve sqlite database asset
	asset, err := bindata.Asset(sqlite3DBAsset)
	if err != nil {
		return false, err
	}

	// Write asset directly to new file
	return true, ioutil.WriteFile(dbPath, asset, 0664)
}
