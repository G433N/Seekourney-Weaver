package server

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"os/exec"
	"testing"

	"seekourney/core/database"
	"seekourney/core/document"
)

// Globally accessible buffer used as mock interface for server handlers
var buffer bytes.Buffer

// Globally accessible context and stop to avoid tons of parameter passing for
// testHandleQuit implementation. Should be reset if used
var ctx context.Context
var stop context.CancelFunc

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// Resets the state of the database by dropping the table and rerunning initdb
func resetSQL(db *sql.DB) {
	if db == nil {
		return
	}
	_, err := db.Exec(`DROP TABLE document`)
	panicOnError(err)

	err = exec.Command(
		"docker",
		"exec",
		_TESTCONTAINERNAME_,
		"psql",
		"-U",
		_DBNAME_,
		"-f",
		"/docker-entrypoint-initdb.d/initdb.sql",
	).Run()
	panicOnError(err)
}

// serverTest defers a func before running the test function to ensure that the
// database container is stopped if the test panics.
// Also resets the state of objects that may have been written to if the test
// did not panic
func serverTest(
	testFunc func(test *testing.T, serverParams serverFuncParams),
	serverParams serverFuncParams,
) func(*testing.T) {
	return func(test *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				stopContainer()
				panic(err)
			} else {
				resetSQL(serverParams.db)
				buffer.Reset()
			}
		}()
		testFunc(test, serverParams)
	}
}

func TestServer(test *testing.T) {
	if testing.Short() {
		test.SkipNow()
	}

	// Navigate back to root directory of codebase
	// Tests seem to run from their own directory, not from where go test is run
	test.Chdir("../..")

	go startContainer()
	defer stopContainer()

	testDB := connectToDB()

	ctx, stop = context.WithCancel(context.Background())
	defer stop()

	serverParams := serverFuncParams{
		writer: &buffer,
		db:     testDB,
	}

	test.Run(
		"TestHandleAllSingle",
		serverTest(testHandleAllSingle, serverParams),
	)
	test.Run(
		"TestHandleAllMultiple",
		serverTest(testHandleAllMultiple, serverParams),
	)
	test.Run("TestHandleSearch", serverTest(testHandleSearch, serverParams))
	test.Run("TestHandleQuit", serverTest(testHandleQuit, serverParams))

	err := testDB.Close()
	if err != nil {
		panic(err)
	}
}

func assertBufferEquals(
	test *testing.T,
	label string,
	expected bytes.Buffer,
	actual bytes.Buffer,
) {
	if !bytes.Equal(expected.Bytes(), actual.Bytes()) {
		test.Log("Buffers do not match, expected:")
		test.Log(expected.String())
		test.Log("Got: ")
		test.Log(actual.String())
		test.Error(label)
	}
}

func testHandleAllSingle(test *testing.T, serverParams serverFuncParams) {
	var expected bytes.Buffer

	database.InsertInto(serverParams.db, testDocument1())

	jsonData, err := json.Marshal([]document.Document{testDocument1()})
	checkIOError(err)
	expected.Write(jsonData)
	expected.WriteByte('\n')

	handleAll(serverParams)
	assertBufferEquals(test, "HandleAllSingle", expected, buffer)
}

func testHandleAllMultiple(test *testing.T, serverParams serverFuncParams) {
	var expected bytes.Buffer

	database.InsertInto(serverParams.db, testDocument1())
	database.InsertInto(serverParams.db, testDocument2())

	jsonData, err := json.Marshal(
		[]document.Document{testDocument1(), testDocument2()},
	)
	checkIOError(err)
	expected.Write(jsonData)
	expected.WriteByte('\n')

	handleAll(serverParams)
	assertBufferEquals(test, "HandleAllMultiple", expected, buffer)
}

func testHandleSearch(test *testing.T, serverParams serverFuncParams) {
	// var expected bytes.Buffer

	// expected.WriteString(pageString(page1) + "\n")
	// handleSearch(serverParams, []string{"key1"})
	// assertBufferEquals(test, "Search key1", expected, buffer)

	// expected.Reset()
	// buffer.Reset()
	// expected.WriteString(pageString(page2) + "\n")
	// handleSearch(serverParams, []string{"key3"})
	// assertBufferEquals(test, "Search key3", expected, buffer)

	// expected.Reset()
	// buffer.Reset()
	// expected.WriteString(pageString(page1) + "\n")
	// expected.WriteString(pageString(page2) + "\n")
	// handleSearch(serverParams, []string{"key2"})
	// assertBufferEquals(test, "Search key2", expected, buffer)
}

// Expects context to be done after calling handleQuit
func testHandleQuit(test *testing.T, serverParams serverFuncParams) {
	handleQuit(serverParams, stop)

	if ctx.Err() == nil {
		test.Error("Expected context to be done.")
	}

	// Reset context globals
	ctx, stop = context.WithCancel(context.Background())
}
