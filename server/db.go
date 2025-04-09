package server

import (
	"database/sql"
	"fmt"
	"io"
	"time"
)

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

func insertRow(
	db *sql.DB,
	path string,
	pathType PathType,
	dictJSON string) (sql.Result, error) {
	insertStmt := `INSERT INTO "paths"("path", "type", "dict") values($1, $2, $3)`
	return db.Exec(insertStmt, path, pathType, dictJSON)
}

// Writes the contents of database rows to the given writer
func writeRows(writer io.Writer, rows *sql.Rows) {
	for rows.Next() {
		var id int64
		var path string
		var pathType string
		var dict string

		err := rows.Scan(&id, &path, &pathType, &dict)
		checkSQLError(err)

		fmt.Fprintf(writer, "id: %d\npath: %s\npathType: %s\ndict: %s\n\n",
			id, path, pathType, dict)
	}
}

// Querys the database for rows containing ALL the given keys.
// Writes output to writer
func queryJSONKeysAll(db *sql.DB, writer io.Writer, keys []string) {
	if len(keys) == 0 {
		panic(`Must enter at least one key to search`)
	}

	// Create a string of query params of correct amount in form "$1, $2, ..."
	// Really would like a one liner for this but still new to Go
	paramsString := ""
	for i := range len(keys) {
		paramsString += fmt.Sprintf("$%d, ", i+1)
	}
	// Cut off last ", "
	paramsString = paramsString[0 : len(paramsString)-2]

	query := fmt.Sprintf(`SELECT * FROM paths WHERE dict ?& ARRAY[%s]`,
		paramsString)

	keysAny := make([]any, len(keys))
	for i, key := range keys {
		keysAny[i] = key
	}

	fmt.Printf("%s (%s)\n", query, keysAny)

	rows, err := db.Query(query, keysAny...)
	checkSQLError(err)
	defer rows.Close()

	writeRows(writer, rows)
}

// Querys the database for all rows.
// Writes output to writer
func queryAll(db *sql.DB, writer io.Writer) {
	query := `SELECT * FROM paths`

	fmt.Printf("%s\n", query)

	rows, err := db.Query(query)
	checkSQLError(err)
	defer rows.Close()

	writeRows(writer, rows)
}
