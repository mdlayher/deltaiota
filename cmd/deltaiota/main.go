// Command deltaiota serves a HTTP API for the Phi Mu Alpha Sinfonia - Delta
// Iota chapter website.
package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/mdlayher/deltaiota"
	"github.com/stretchr/graceful"
)

const (
	// driver is the database/sql driver used for the database instance
	driver = "sqlite3"
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
	flag.StringVar(&db, "db", "", "DSN for database instance")
	flag.StringVar(&host, "host", ":8080", "HTTP server host")
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "HTTP graceful timeout duration")
}

func main() {
	// Parse all flags
	flag.Parse()

	// Open database connection
	didb := &deltaiota.DB{}
	if err := didb.Open(driver, db); err != nil {
		log.Fatal(err)
	}

	// Start HTTP server using deltaiota handler on specified host
	log.Println("deltaiota: listening:", host)
	if err := graceful.ListenAndServe(&http.Server{
		Addr:    host,
		Handler: deltaiota.NewServeMux(didb),
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
