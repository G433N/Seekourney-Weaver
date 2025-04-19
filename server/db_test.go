package server

import (
	"database/sql"
	"os/exec"
	"testing"
)

var testDB *sql.DB

var page1 = Page{
	id:       1,
	path:     "/some/path",
	pathType: "file",
	dict:     `{"key1": 1, "key2": 2}`,
}

var page2 = Page{
	id:       2,
	path:     "/some/other/path",
	pathType: "file",
	dict:     `{"key2": 4, "key3": 6}`,
}

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

func TestDB(test *testing.T) {
	if testing.Short() {
		test.SkipNow()
	}
	test.Chdir("..")

	go startContainer()
	testDB = connectToDB()

	test.Run("TestQueryAll", safelyTest(testQueryAll))
	test.Run("TestQueryJSONKeysAll", safelyTest(testQueryJSONKeysAll))
	test.Run("TestInsertRow", safelyTest(testInsertRow))

	err := testDB.Close()
	if err != nil {
		panic(err)
	}
	stopContainer()
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
		"docker",
		"exec",
		containerName,
		"psql",
		"-U",
		dbname,
		"-f",
		initDB,
	).Run()
	if err != nil {
		panic(err)
	}
}

func pageEquals(a Page, b Page) bool {
	return a.id == b.id &&
		a.path == b.path &&
		a.pathType == b.pathType &&
		a.dict == b.dict
}

func pageEqualsIgnoreId(a Page, b Page) bool {
	return a.path == b.path &&
		a.pathType == b.pathType &&
		a.dict == b.dict
}

// Reads the next row in rows and checks if it matches expected
func checkRow(test *testing.T, rows *sql.Rows, expected Page) {
	var page Page
	if !rows.Next() {
		test.Error("Expected a row")
	}
	err := rows.Scan(&page.id, &page.path, &page.pathType, &page.dict)
	checkSQLError(err)

	if !pageEquals(page, expected) {
		test.Errorf("testRow failed \nexpected: \n%s, \n\ngot: \n%s",
			pageString(expected), pageString(page))
	}
}

func testQueryAll(test *testing.T) {
	rows := queryAll(testDB)
	checkRow(test, rows, page1)
	checkRow(test, rows, page2)
	unsafelyClose(rows)
}

func testQueryJSONKeysAll(test *testing.T) {
	rows := queryJSONKeysAll(testDB, []string{"key1"})
	checkRow(test, rows, page1)
	unsafelyClose(rows)

	rows = queryJSONKeysAll(testDB, []string{"key3"})
	checkRow(test, rows, page2)
	unsafelyClose(rows)

	rows = queryJSONKeysAll(testDB, []string{"key2"})
	checkRow(test, rows, page1)
	checkRow(test, rows, page2)
	unsafelyClose(rows)

	rows = queryJSONKeysAll(testDB, []string{"key1", "key2"})
	checkRow(test, rows, page1)
	unsafelyClose(rows)

	rows = queryJSONKeysAll(testDB, []string{"key1", "key2", "key3"})
	if rows.Next() {
		test.Error("Expected 0 rows")
	}
}

func testInsertRow(test *testing.T) {
	// Test normal
	newPage := Page{
		path:     "/path1",
		pathType: PathTypeFile,
		dict:     `{"uniqueKey": 1}`,
	}
	_, err := insertRow(testDB, newPage)
	checkSQLError(err)
	inserted, ok := getRowByPath(testDB, newPage.path)
	if !ok || !pageEqualsIgnoreId(inserted, newPage) {
		test.Error("Insertion error")
	}

	// Test duplicate path
	res, err := insertRow(testDB, newPage)
	if err == nil {
		test.Error("Expected error", res)
	}

	// Test missing path
	missingPath := Page{
		pathType: PathTypeFile,
		dict:     `{"uniqueKey": 1}`,
	}
	res, err = insertRow(testDB, missingPath)
	if err == nil {
		test.Error("Expected error", res)
	}

	// Test missing pathType
	missingPathType := Page{
		path: "/path2",
		dict: `{"uniqueKey": 1}`,
	}
	_, err = insertRow(testDB, missingPathType)
	if err == nil {
		test.Error("Expected error")
	}

	// Test missing dict
	missingDict := Page{
		path:     "/path3",
		pathType: PathTypeFile,
	}
	missingDictWithJSON := Page{
		path:     missingDict.path,
		pathType: PathTypeFile,
		dict:     emptyJSON,
	}
	_, err = insertRow(testDB, missingDict)
	checkSQLError(err)
	inserted, ok = getRowByPath(testDB, missingDict.path)
	if !ok || !pageEqualsIgnoreId(inserted, missingDictWithJSON) {
		test.Error("Insertion error")
	}
}
