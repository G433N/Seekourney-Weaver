package server

import (
	"bytes"
	"context"
	"os"
	"testing"
)

// Globals used during testing, should always be restored to base state after
// each test
var serverParams serverFuncParams
var ctx context.Context

func TestServer(test *testing.T) {
	os.Chdir("..")

	go startContainer()
	testDB = connectToDB()

	defer testDB.Close()
	defer stopContainer()

	var stop context.CancelFunc
	ctx, stop = context.WithCancel(context.Background())
	defer stop()

	var buffer *bytes.Buffer

	serverParams = serverFuncParams{
		writer: buffer,
		stop:   stop,
		db:     testDB,
	}

	test.Run("TestHandleAll", safelyTest(testHandleAll))
	test.Run("TestHandleSearch", safelyTest(testHandleSearch))
	test.Run("TestHandleAdd", safelyTest(testHandleAdd))
	test.Run("TestHandleQuit", safelyTest(testHandleQuit))
}

func testHandleAll(test *testing.T) {
	// TODO
}

func testHandleSearch(test *testing.T) {
	// TODO

}

func testHandleAdd(test *testing.T) {
	// TODO

}

// Expects context to be done after calling handleQuit
func testHandleQuit(test *testing.T) {
	handleQuit(serverParams)

	if ctx.Err() == nil {
		test.Error("Expected context to be done.")
	}

	// Reset context global
	_, stop := context.WithCancel(context.Background())
	serverParams.stop = stop
}
