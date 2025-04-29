package database

import (
	"database/sql"
	"iter"
	"seekourney/utils"
	"strconv"
	"strings"
)

const (
	_INSERT_     = "INSERT"
	_INTO_       = "INTO"
	_VALUES_     = "VALUES"
	_SELECT_     = "SELECT"
	_FROM_       = "FROM"
	_WHERE_      = "WHERE"
	_AS_         = "AS"
	_JSON_VALUE_ = "JSON_VALUE"
)

/// Scan

// SQLScan is an interface that defines a method
// for scanning a SQL row to an object of type Self.
type SQLScan[Self any] interface {

	// SQLScan scans a SQL row into an object of type Self.
	// This method schould call the Scan method of the sql.Rows
	// and assign the values to the fields of the object
	SQLScan(rows *sql.Rows) (Self, error)
}

// scan is a helper function that scans a SQL row into an object of type T.
func scan[T SQLScan[T]](rows *sql.Rows) (T, error) {
	var obj T
	return obj.SQLScan(rows)
}

// ScanRowsIter is a function that takes a sql.Rows object and
// returns an iterator of objects of type T.
// Every row is scanned into an object of type T and yielded to the caller.
func ScanRowsIter[T SQLScan[T]](Rows *sql.Rows) iter.Seq[utils.Result[T]] {

	return func(yield func(utils.Result[T]) bool) {
		for Rows.Next() {
			obj, err := scan[T](Rows)

			result := utils.Result[T]{
				Value: obj,
				Err:   err,
			}

			if !yield(result) {
				break
			}
		}
	}
}

// / Write

type SQLValue = any

// SQLWrite is an interface that defines methods
// for writing SQL rows from a object
type SQLWrite interface {

	// SQLGetName returns the name of the SQL table
	SQLGetName() string

	// SQLGetFields returns the fields of the tables rows
	SQLGetFields() []string

	// SQLGetValues returns the values of a row
	SQLGetValues() []SQLValue
}

// objectTemplate is a type that represents a SQL Object row thing
// TODO: Imporve this
type objectTemplate string

// valueSubstitution is a type that represents a SQL value substitution
type valueSubstitution string

// Statment is a type that represents a SQL statement
type Statment string

/// Insert

// InsertInto executes an INSERT statement into the database
func InsertInto(db *sql.DB, object SQLWrite) (sql.Result, error) {
	stmt := InsertIntoStatment(object)

	return db.Exec(string(stmt), object.SQLGetValues()...)
}

// InsertIntoStatment creates an INSERT statement from a SQLWrite object
func InsertIntoStatment(template SQLWrite) Statment {

	return insertIntoStatment(
		sqlTemplate(template),
		sqlValueSubstition(template),
	)
}

// insertIntoStatment creates an INSERT statement from a template
// and a value substitution
func insertIntoStatment(
	template objectTemplate,
	sub valueSubstitution,
) Statment {
	list := []string{
		_INSERT_,
		_INTO_,
		string(template),
		string(sub),
	}
	return Statment(strings.Join(list, " "))

}

// sqlTemplate creates a SQL template from a Go struct/object
func sqlTemplate(template SQLWrite) objectTemplate {
	name := template.SQLGetName()
	fields := template.SQLGetFields()

	return objectTemplate(name + " (" + strings.Join(fields, ",") + ")")
}

// sqlValueSubstition creates a SQL value substitution from a Go struct/object
func sqlValueSubstition(template SQLWrite) valueSubstitution {
	values := template.SQLGetValues()

	str := _VALUES_ + " ("
	for i := range values {
		str += "$" + strconv.Itoa(i+1)
		if i != len(values)-1 {
			str += ", "
		}
	}
	str += ")"

	return valueSubstitution(str)
}

/// Select

// SelectStatment is a type that represents a SQL SELECT keyword
type SelectStatment string

// SelectQuery is a type that represents a SQL SELECT with a query
type SelectQuery string

// SelectFrom is a type that represents a SQL SELECT with a FROM clause
type SelectFrom string

// SelectWhere is a type that represents a SQL SELECT with a WHERE clause
type SelectWhere string

// / JsonValue creates a JSON_VALUE SELECT statement
// / of the form JSON_VALUE(sqlField, '$.jsonField') AS name
func JsonValue(sqlField string, jsonField string, name string) string {

	s := []string{
		_JSON_VALUE_ + "(",
		sqlField,
		",",
		"'$." + jsonField + "'",
		")",
		_AS_,
		name,
	}

	return strings.Join(s, " ")
}

func Select() SelectStatment {
	return SelectStatment(_SELECT_)
}

// Queries adds a list of queries to the SQL statement
func (s SelectStatment) Queries(query ...string) SelectQuery {
	return SelectQuery(string(s) + " " + strings.Join(query, ", "))
}

// QueryAll adds a wildcard (*) to the SQL statement
func (s SelectStatment) QueryAll() SelectQuery {
	return s.Queries("*")
}

// From adds a FROM clause to the SQL statement
func (s SelectQuery) From(table string) SelectFrom {
	return SelectFrom(string(s) + " " + _FROM_ + " " + table)
}

// Where adds a WHERE clause to the SQL statement
func (s SelectFrom) Where(condition string) SelectWhere {
	return SelectWhere(string(s) + " " + _WHERE_ + " " + condition)
}

// ExecExec executes a SQL statement and returns the result into obj
// The insert function is used to insert the result into obj
func ExecScan[T SQLScan[T], U any](
	db *sql.DB,
	query string,
	obj *U,
	insert func(*U, T),
	args ...any) (resErr error) {

	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}

	defer func() {
		err = rows.Close()
		if resErr != nil {
			resErr = err
		}
	}()

	for row := range ScanRowsIter[T](rows) {
		if row.Err != nil {
			return row.Err
		}

		insert(obj, row.Value)
	}

	return nil
}

type num int

func (n num) SQLScan(rows *sql.Rows) (num, error) {
	var i int
	err := rows.Scan(&i)
	if err != nil {
		return 0, err
	}
	return num(i), nil
}

func RowAmount(db *sql.DB, table string) (int, error) {

	query := Select().Queries("COUNT(*)").From(table)
	var count num

	insert := func(res *num, sqlRes num) {
		*res = sqlRes
	}

	err := ExecScan(db, string(query), &count, insert)

	if err != nil {
		return 0, err
	}

	return int(count), nil
}
