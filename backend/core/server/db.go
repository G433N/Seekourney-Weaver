package server

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// PathType is used for database enumerable type, can be either 'web' or 'file'
type PathType string

// JSONString is a string in JSON format.
type JSONString string

const (
	// _TYPEWEB_           PathType = "web"
	// _TYPEFILE_          PathType = "file"
	_CONNECTIONRETRIES_ int           = 10
	_RETRYDELAY_        time.Duration = 500 * time.Millisecond
)

// connectToDB attempts to connect to the database,
// and will retry every half second for 5 seconds in case the docker container
// is still starting up.
// Returns database file descriptor ptr if the connection succeeds.
// Panics with error on connection failure.
func connectToDB() *sql.DB {
	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	log.Println("Connecting to database")
	// Waiting animation
	fmt.Print("Waiting for database.")
	for range _CONNECTIONRETRIES_ {
		db, err := sql.Open("postgres", psqlconn)

		if err != nil {
			log.Println("Error opening database connection:", err)
		}

		if err = db.Ping(); err == nil {
			// Need to add a new line to "end" the waiting animation
			fmt.Println("")
			log.Println("Database ready")
			return db
		}
		time.Sleep(_RETRYDELAY_)
		fmt.Print(".")
	}
	fmt.Print("\n")
	panic("Could not connect to database, check docker.log for more info")
}
