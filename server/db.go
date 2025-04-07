package server

import (
	"database/sql"
	"fmt"
	"io"
	"time"
)

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

func insertRow(db *sql.DB, path string, pathType PathType, dictJSON string) sql.Result {
	insertStmt := `INSERT INTO "paths"("path", "type", "dict") values($1, $2, $3)`
	result, err := db.Exec(insertStmt, path, pathType, dictJSON)
	checkSQLError(err)
	return result
}

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

func queryJSONKeysAll(db *sql.DB, writer io.Writer, keys []string) {
	if len(keys) == 0 {
		panic(`Must enter at least one key to search`)
	}

	// Really would like a one liner for this but still new to Go
	paramsString := ""
	for i := range len(keys) {
		paramsString += fmt.Sprintf("$%d, ", i+1)
	}
	// Cut off last ", "
	paramsString = paramsString[0 : len(paramsString)-2]

	query := fmt.Sprintf(`SELECT * FROM paths WHERE dict ?& ARRAY[%s]`,
		paramsString)

	// Again this could probably be one lined and inside of Query function call
	keysAny := []any{}
	for _, key := range keys {
		keysAny = append(keysAny, key)
	}

	fmt.Printf("%s (%s)\n", query, keysAny)

	rows, err := db.Query(query, keysAny...)
	checkSQLError(err)
	defer rows.Close()

	writeRows(writer, rows)
}

func queryAll(db *sql.DB, writer io.Writer) {
	query := `SELECT * FROM paths`

	fmt.Printf("%s\n", query)

	rows, err := db.Query(query)
	checkSQLError(err)
	defer rows.Close()

	writeRows(writer, rows)
}

// func queryJSONKeyExists(db *sql.DB, key string) {
// 	query := `SELECT * FROM paths WHERE dict ? $1`

// 	fmt.Printf("%s (%s):", query, key)

// 	rows, err := db.Query(query, key)
// 	checkSQLError(err)
// 	defer rows.Close()

// 	writeRows(os.Stdout, rows)
// }
