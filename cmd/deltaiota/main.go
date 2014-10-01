// Command deltaiota serves a HTTP API for the Phi Mu Alpha Sinfonia - Delta
// Iota chapter website.
package main

import (
	"flag"
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
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "HTTP graceful timeout duration")
}

func main() {
	// Parse all flags
	flag.Parse()

	// If database is sqlite3, perform initial setup
	if driver == sqlite3 {
		if err := sqlite3Setup(db); err != nil {
			log.Fatal(err)
		}
	}

	// Open database connection
	didb := &data.DB{}
	if err := didb.Open(driver, db); err != nil {
		log.Fatal(err)
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

	// Close database connection
	if err := didb.Close(); err != nil {
		log.Fatal(err)
	}

	log.Println("deltaiota: graceful shutdown complete")
}

// sqlite3Setup performs setup routines specific to a sqlite3 database
func sqlite3Setup(dsn string) error {
	// Check if database already exists at specified location
	dbPath := path.Clean(dsn)
	_, err := os.Stat(dbPath)
	if err == nil {
		// Database exists, skip setup
		log.Println("deltaiota: using sqlite3 database:", dbPath)
		return nil
	}

	// Any other errors, return now
	if !os.IsNotExist(err) {
		return err
	}

	// Retrieve sqlite database asset
	asset, err := bindata.Asset(sqlite3DBAsset)
	if err != nil {
		return err
	}

	// Write asset directly to new file
	log.Println("deltaiota: creating sqlite3 database:", dbPath)
	return ioutil.WriteFile(dbPath, asset, 0664)
}
