package server

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Used for database enumerable type, can be either 'web' or 'file'
type PathType string

// const (
// 	pathTypeWeb  PathType = "web"
// 	pathTypeFile PathType = "file"
// )

type JSONString string

// Attempts to connect to the database, will retry every half second for
// 5 seconds in case the docker container is still starting up.
// Returns a pointer to a database file descriptor if the connection succeeds.
// Terminates with an error if it fails to connect.
func connectToDB() *sql.DB {
	retries := 10

	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	log.Println("Connecting to database")
	// Waiting animation
	fmt.Print("Waiting for database.")
	for range retries {
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
		time.Sleep(500 * time.Millisecond)
		fmt.Print(".")
	}
	fmt.Print("\n")
	panic("Could not connect to database, check docker.log for more info")
}
