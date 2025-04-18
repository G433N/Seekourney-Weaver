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
var buffer bytes.Buffer

func TestServer(test *testing.T) {
	if testing.Short() {
		test.SkipNow()
	}

	os.Chdir("..")

	go startContainer()
	testDB = connectToDB()

	defer testDB.Close()
	defer stopContainer()

	var stop context.CancelFunc
	ctx, stop = context.WithCancel(context.Background())
	defer stop()

	serverParams = serverFuncParams{
		writer: &buffer,
		stop:   stop,
		db:     testDB,
	}

	test.Run("TestHandleAll", safelyTest(testHandleAll))
	test.Run("TestHandleSearch", safelyTest(testHandleSearch))
	test.Run("TestHandleAdd", safelyTest(testHandleAdd))
	test.Run("TestHandleQuit", safelyTest(testHandleQuit))
}

func assertBufferEquals(test *testing.T, expected bytes.Buffer, actual bytes.Buffer) {
	if !bytes.Equal(expected.Bytes(), actual.Bytes()) {
		test.Log("Incorrect output written by server function, expected:")
		test.Log(expected.String())
		test.Log("\n\nGot: ")
		test.Log(actual.String())
		test.Fail()
	}
}

func testHandleAll(test *testing.T) {
	handleAll(serverParams)

	var expected bytes.Buffer
	expected.WriteString(pageString(page1) + "\n")
	expected.WriteString(pageString(page2) + "\n")

	assertBufferEquals(test, expected, buffer)
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
