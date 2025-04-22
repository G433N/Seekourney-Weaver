package server

import (
	"bytes"
	"context"
	"fmt"
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

	test.Chdir("..")

	go startContainer()
	testDB = connectToDB()

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

	err := testDB.Close()
	if err != nil {
		panic(err)
	}
	stopContainer()
}

func assertBufferEquals(
	test *testing.T,
	label string,
	expected bytes.Buffer,
	actual bytes.Buffer,
) {
	if !bytes.Equal(expected.Bytes(), actual.Bytes()) {
		fmt.Println("Incorrect output written by server function, expected:")
		fmt.Println(expected.String())
		fmt.Println("\n\nGot: ")
		fmt.Println(actual.String())
		test.Error(label)
	}
}

func testHandleAll(test *testing.T) {
	var expected bytes.Buffer
	expected.WriteString(pageString(page1) + "\n")
	expected.WriteString(pageString(page2) + "\n")

	buffer.Reset()
	handleAll(serverParams)
	assertBufferEquals(test, "HandleAll", expected, buffer)
}

func testHandleSearch(test *testing.T) {
	var expected bytes.Buffer

	buffer.Reset()
	expected.WriteString(pageString(page1) + "\n")
	handleSearch(serverParams, []string{"key1"})
	assertBufferEquals(test, "Search key1", expected, buffer)

	expected.Reset()
	buffer.Reset()
	expected.WriteString(pageString(page2) + "\n")
	handleSearch(serverParams, []string{"key3"})
	assertBufferEquals(test, "Search key3", expected, buffer)

	expected.Reset()
	buffer.Reset()
	expected.WriteString(pageString(page1) + "\n")
	expected.WriteString(pageString(page2) + "\n")
	handleSearch(serverParams, []string{"key2"})
	assertBufferEquals(test, "Search key2", expected, buffer)
}

func testHandleAdd(test *testing.T) {
	// Add one
	path1 := "/a/path"
	handleAdd(serverParams, []string{path1})

	page, ok := getRowByPath(testDB, path1)
	if !ok ||
		page.path != path1 ||
		page.pathType != PathTypeFile ||
		page.dict != emptyJSON {
		test.Error("handleAdd one, did not add correct page")
	}

	// Add many
	path2 := "/another/path"
	path3 := "/one/more/path"
	handleAdd(serverParams, []string{path2, path3})
	page, ok = getRowByPath(testDB, path2)
	if !ok ||
		page.path != path2 ||
		page.pathType != PathTypeFile ||
		page.dict != emptyJSON {
		test.Error("handleAdd multiple, did not add correct page (2)")
	}
	page, ok = getRowByPath(testDB, path3)
	if !ok ||
		page.path != path3 ||
		page.pathType != PathTypeFile ||
		page.dict != emptyJSON {
		test.Error("handleAdd multiple, did not add correct page (3)")
	}

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
