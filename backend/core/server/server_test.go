package server

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"seekourney/core/config"
	"seekourney/core/database"
	"seekourney/core/document"
	"seekourney/utils"
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

	_, err = db.Exec(`DROP TABLE collection`)
	panicOnError(err)

	_, err = db.Exec(`DROP TABLE indexer`)
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

// assertBufferEquals checks that the content of two bytes.Buffers are equal,
// failing the test and logging an message if not
func assertBufferEquals(
	test *testing.T,
	expected bytes.Buffer,
	actual bytes.Buffer,
) {
	if !bytes.Equal(expected.Bytes(), actual.Bytes()) {
		test.Error(
			"Buffers do not match, expected:\n",
			expected.String(),
			"\nGot:\n",
			actual.String(),
		)
	}
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

	conf = config.New()

	test.Run(
		"TestHandleAllSingle",
		serverTest(testHandleAllSingle, serverParams),
	)
	test.Run(
		"TestHandleAllMultiple",
		serverTest(testHandleAllMultiple, serverParams),
	)
	test.Run(
		"TestHandleSearchSQLSingle",
		serverTest(testHandleSearchSQLSingle, serverParams),
	)
	test.Run(
		"TestHandleSearchSQLInvalid",
		serverTest(testHandleSearchSQLInvalid, serverParams),
	)
	test.Run(
		"TestHandleSearchSQLMultiple",
		serverTest(testHandleSearchSQLMultiple, serverParams),
	)
	test.Run("TestHandleQuit", serverTest(testHandleQuit, serverParams))

	test.Run("TestHandleDownload", serverTest(testHandleDownload, serverParams))
	err := testDB.Close()
	if err != nil {
		panic(err)
	}
}

func testHandleAllSingle(test *testing.T, serverParams serverFuncParams) {
	var expected bytes.Buffer

	_, err := database.InsertInto(serverParams.db, testIndexer())
	panicOnError(err)

	_, err = database.InsertInto(serverParams.db, testCollection())
	panicOnError(err)

	_, err = database.InsertInto(serverParams.db, testDocument1())
	panicOnError(err)

	jsonData, err := json.Marshal([]document.Document{testDocument1()})
	utils.PanicOnError(err)
	expected.Write(jsonData)
	expected.WriteByte('\n')

	handleAll(serverParams)
	assertBufferEquals(test, expected, buffer)
}

func testHandleAllMultiple(test *testing.T, serverParams serverFuncParams) {
	var expected bytes.Buffer

	_, err := database.InsertInto(serverParams.db, testIndexer())
	panicOnError(err)

	_, err = database.InsertInto(serverParams.db, testCollection())
	panicOnError(err)

	_, err = database.InsertInto(serverParams.db, testDocument1())
	panicOnError(err)

	_, err = database.InsertInto(serverParams.db, testDocument2())
	panicOnError(err)

	jsonData, err := json.Marshal(
		[]document.Document{testDocument1(), testDocument2()},
	)
	utils.PanicOnError(err)
	expected.Write(jsonData)
	expected.WriteByte('\n')

	handleAll(serverParams)
	assertBufferEquals(test, expected, buffer)
}

func testHandleSearchSQLSingle(test *testing.T, serverParams serverFuncParams) {
	var response utils.SearchResponse

	_, err := database.InsertInto(serverParams.db, testIndexer())
	panicOnError(err)

	_, err = database.InsertInto(serverParams.db, testCollection())
	panicOnError(err)

	_, err = database.InsertInto(serverParams.db, testDocument1())
	panicOnError(err)

	handleSearchSQL(serverParams, []string{"key1"})

	err = json.Unmarshal([]byte(buffer.Bytes()), &response)
	panicOnError(err)
	if len(response.Results) != 1 ||
		response.Results[0].Path != testDocument1().Path ||
		response.Results[0].Source != testDocument1().Source {
		test.Error("Recieved incorrect result")
		test.Log(response.Results)
	}
}

func testHandleSearchSQLInvalid(
	test *testing.T,
	serverParams serverFuncParams,
) {
	var response utils.SearchResponse

	_, err := database.InsertInto(serverParams.db, testIndexer())
	panicOnError(err)

	_, err = database.InsertInto(serverParams.db, testCollection())
	panicOnError(err)

	_, err = database.InsertInto(serverParams.db, testDocument1())
	panicOnError(err)

	handleSearchSQL(serverParams, []string{"badkey"})

	err = json.Unmarshal([]byte(buffer.Bytes()), &response)
	panicOnError(err)
	if len(response.Results) != 0 {
		test.Error("Expected no result")
		test.Log(response.Results)
	}
}

func testHandleSearchSQLMultiple(
	test *testing.T,
	serverParams serverFuncParams,
) {
	var response utils.SearchResponse

	_, err := database.InsertInto(serverParams.db, testIndexer())
	panicOnError(err)

	_, err = database.InsertInto(serverParams.db, testCollection())
	panicOnError(err)

	_, err = database.InsertInto(serverParams.db, testDocument1())
	panicOnError(err)

	_, err = database.InsertInto(serverParams.db, testDocument2())
	panicOnError(err)

	// key1 is unique to testDocument1
	handleSearchSQL(serverParams, []string{"key1"})

	err = json.Unmarshal([]byte(buffer.Bytes()), &response)
	panicOnError(err)
	if len(response.Results) != 1 ||
		response.Results[0].Path != testDocument1().Path ||
		response.Results[0].Source != testDocument1().Source {
		test.Error("Recieved incorrect result")
		test.Log(response.Results)
	}
	buffer.Reset()

	// key3 is unique to testDocument2
	handleSearchSQL(serverParams, []string{"key3"})

	err = json.Unmarshal([]byte(buffer.Bytes()), &response)
	panicOnError(err)
	if len(response.Results) != 1 ||
		response.Results[0].Path != testDocument2().Path ||
		response.Results[0].Source != testDocument2().Source {
		test.Error("Recieved incorrect result")
		test.Log(response.Results)
	}
	buffer.Reset()

	// key2 is common among both documents
	handleSearchSQL(serverParams, []string{"key2"})
	err = json.Unmarshal([]byte(buffer.Bytes()), &response)
	panicOnError(err)
	if len(response.Results) != 2 {
		test.Error("Expected two results")
		test.Log(response.Results)
	}
	buffer.Reset()
}

func testHandleDownload(test *testing.T, serverParams serverFuncParams) {
	workingDir, _ := os.Getwd()
	fileDir := filepath.Join(workingDir, "README.md")
	expectedOpenedFile, err := os.ReadFile(fileDir)
	if err != nil {
		test.Error("Error reading file")
	}

	handleDownload(serverParams, []string{fileDir})

	if !bytes.Equal(expectedOpenedFile, buffer.Bytes()) {
		test.Error("Files not equal")
		test.Log(workingDir)
	}
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
