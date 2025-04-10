package server

import (
	"database/sql"
	"os"
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
func testRow(test *testing.T, rows *sql.Rows, expected Page) {
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

func TestQueryAll(test *testing.T) {
	rows := queryAll(testDB)
	testRow(test, rows, page1)
	testRow(test, rows, page2)
	rows.Close()
}

func TestQueryJSONKeysAll(test *testing.T) {
	rows := queryJSONKeysAll(testDB, []string{"key1"})
	testRow(test, rows, page1)
	rows.Close()

	rows = queryJSONKeysAll(testDB, []string{"key3"})
	testRow(test, rows, page2)
	rows.Close()

	rows = queryJSONKeysAll(testDB, []string{"key2"})
	testRow(test, rows, page1)
	testRow(test, rows, page2)
	rows.Close()

	rows = queryJSONKeysAll(testDB, []string{"key1", "key2"})
	testRow(test, rows, page1)
	rows.Close()

	rows = queryJSONKeysAll(testDB, []string{"key1", "key2", "key3"})
	if rows.Next() {
		test.Error("Expected 0 rows")
	}
}

func TestInsertRows(test *testing.T) {
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

func TestMain(m *testing.M) {
	os.Chdir("..")

	go startContainer()

	// TODO, this doesn't seem to stop the container if one of the tests panic
	defer stopContainer()

	testDB = connectToDB()

	m.Run()
}
