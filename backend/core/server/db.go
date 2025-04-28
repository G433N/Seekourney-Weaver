package server

import (
	"database/sql"
	"fmt"
	"time"
)

// Used for database enumerable type, can be either 'web' or 'file'
type PathType string

const (
	pathTypeWeb  PathType = "web"
	pathTypeFile PathType = "file"
)

type JSONString string

type Page struct {
	// id       int64
	path     string
	pathType PathType
	// dict     JSONString
}

// Attempts to connect to the database, will retry every half second for
// 5 seconds in case the docker container is still starting up.
// Returns a pointer to a database file descriptor if the connection succeeds.
// Terminates with an error if it fails to connect.
func connectToDB() *sql.DB {
	retries := 10

	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	fmt.Print("Waiting for database.")
	for range retries {
		db, _ := sql.Open("postgres", psqlconn)
		if err := db.Ping(); err == nil {
			fmt.Println("\nDatabase ready")
			return db
		}
		time.Sleep(500 * time.Millisecond)
		fmt.Print(".")
	}
	fmt.Print("\n")
	panic("Could not connect to database, check docker.log for more info")
}

func checkSQLError(err error) {
	if err != nil {
		panic(err)
	}
}
