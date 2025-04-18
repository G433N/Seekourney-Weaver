package server

import (
	"database/sql"
	"os/exec"
	"testing"
)

var testDB *sql.DB

// Adds a deferred func before running the test function to ensure that the
// database container is stopped if the test panics. Also resets the database
// to if the tests executed without panicking
func safelyTest(testFunc func(test *testing.T)) func(*testing.T) {
	return func(test *testing.T) {
		// Stop container if test panicked, otherwise reset database
		defer func() {
			if err := recover(); err != nil {
				stopContainer()
				panic(err)
			} else {
				resetSQL(testDB)
			}
		}()
		testFunc(test)
	}
}

// Resets the state of the database by dropping the table and rerunning initdb
func resetSQL(db *sql.DB) {
	if db == nil {
		return
	}
	_, err := db.Exec(`DROP TABLE page`)
	checkSQLError(err)

	const initDB = "/docker-entrypoint-initdb.d/initdb.sql"

	err = exec.Command(
		"docker", "exec", containerName, "psql", "-U", dbname, "-f", initDB).Run()
	if err != nil {
		panic(err)
	}
}
