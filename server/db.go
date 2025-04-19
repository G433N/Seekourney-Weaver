package server

import (
	"database/sql"
	"fmt"
	"io"
	"time"

	"github.com/lib/pq"
)

// Used for database enumerable type, can be either 'web' or 'file'
type PathType string

const (
	pathTypeWeb  PathType = "web"
	pathTypeFile PathType = "file"
)

type JSONString string

type Page struct {
	id       int64
	path     string
	pathType PathType
	dict     JSONString
}

// Prints the values of a page to a string
func pageString(page Page) string {
	return fmt.Sprintf("id: %d, \npath: %s, \npathType: %s, \ndict: %s\n",
		page.id, page.path, page.pathType, page.dict)
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

func unsafelyClose(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		panic(err)
	}
}

func insertRow(db *sql.DB, page Page) (sql.Result, error) {
	insert := `INSERT INTO "page"("path", "type") values($1, $2)`
	fmt.Printf("%s (%s, %s)\n", insert, page.path, page.pathType)
	return db.Exec(insert, page.path, page.pathType)
}

// func insertRowWithJSON(db *sql.DB, page Page) (sql.Result, error) {
// 	insertStmt := `INSERT INTO "page"("path", "type", "dict") values($1, $2,
// $3)`
// 	return db.Exec(insertStmt, page.path, page.pathType, page.dict)
// }

// Writes the contents of database rows to the given writer
func writeRows(writer io.Writer, rows *sql.Rows) {
	for rows.Next() {
		var page Page

		err := rows.Scan(&page.id, &page.path, &page.pathType, &page.dict)
		checkSQLError(err)

		fmt.Fprint(writer, pageString(page), "\n")
	}
}

// Querys the database for rows containing ALL the given keys.
// Writes output to writer
func queryJSONKeysAll(db *sql.DB, writer io.Writer, keys []string) {
	query := `SELECT * FROM page WHERE dict ?& $1`

	if len(keys) == 0 {
		panic(`No keys given`)
	}

	fmt.Printf("%s (%s)\n", query, keys)

	rows, err := db.Query(query, pq.StringArray(keys))
	checkSQLError(err)
	defer func() {
		err = rows.Close()
		checkIOError(err)
	}()

	writeRows(writer, rows)
}

// Querys the database for all rows.
// Writes output to writer
func queryAll(db *sql.DB, writer io.Writer) {
	query := `SELECT * FROM page`

	fmt.Printf("%s\n", query)

	rows, err := db.Query(query)
	checkSQLError(err)
	defer func() {
		err = rows.Close()
		checkIOError(err)
	}()

	writeRows(writer, rows)
}

func getRowByPath(db *sql.DB, path string) (Page, bool) {
	var page Page
	const query = `SELECT * FROM page WHERE path = $1`

	fmt.Printf("%s (%s)\n", query, path)

	rows, err := db.Query(query, path)
	checkSQLError(err)
	defer unsafelyClose(rows)

	if rows.Next() {
		err = rows.Scan(&page.id, &page.path, &page.pathType, &page.dict)
		checkSQLError(err)
		return page, true
	} else {
		return page, false
	}
}
