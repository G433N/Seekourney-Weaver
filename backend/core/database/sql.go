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

type SQLScan[Self any] interface {

	// This method schould call the Scan method of the sql.Rows
	// and assign the values to the fields of the object
	SQLScan(rows *sql.Rows) (Self, error)
}

type IntoMap[K comparable, V any] interface {
	// NOTE: Might want to extrect this

	IntoKey() K

	IntoValue() V
}

func scan[T SQLScan[T]](rows *sql.Rows) (T, error) {
	var obj T

	return obj.SQLScan(rows)
}

func ScanRows[T SQLScan[T]](riws *sql.Rows) ([]T, error) {

	var objects []T

	for result := range ScanRowsIter[T](riws) {
		if result.Err != nil {
			return nil, result.Err
		}

		objects = append(objects, result.Value)
	}

	return objects, nil
}

func ScanRowsIntoMapRaw[T interface {
	SQLScan[T]
	IntoMap[K, V]
}, K comparable, V any,
	K1 comparable, V1 any](Rows *sql.Rows,
	keyConv func(K) K1,
	valueConv func(V) V1) (map[K1]V1, error) {
	objects := make(map[K1]V1)

	for result := range ScanRowsIter[T](Rows) {
		if result.Err != nil {
			return nil, result.Err
		}

		obj := result.Value
		key := keyConv(obj.IntoKey())
		value := valueConv(obj.IntoValue())

		objects[key] = value
	}

	return objects, nil
}

func ScanRowsIntoMap[T interface {
	SQLScan[T]
	IntoMap[K, V]
}, K comparable, V any](Rows *sql.Rows) (map[K]V, error) {
	return ScanRowsIntoMapRaw[T, K, V, K, V](
		Rows,
		func(k K) K { return k },
		func(v V) V { return v },
	)
}

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

/// Write

type SQLWrite interface {
	SQLGetName() string

	SQLGetFields() []string

	SQLGetValues() []any
}

type objectTemplate string
type valueSubstitution string

type Statment string

/// Insert

func InsertInto(db *sql.DB, object SQLWrite) (sql.Result, error) {
	stmt := InsertIntoStatment(object)

	return db.Exec(string(stmt), object.SQLGetValues()...)
}

func InsertIntoStatment(template SQLWrite) Statment {

	return insertIntoStatment(
		sqlTemplate(template),
		sqlValueSubstition(template),
	)
}

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

func sqlTemplate(template SQLWrite) objectTemplate {
	name := template.SQLGetName()
	fields := template.SQLGetFields()

	return objectTemplate(name + " (" + strings.Join(fields, ",") + ")")
}

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

type SelectStatment string
type SelectQuery string
type SelectFrom string
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

func (s SelectStatment) Queries(query ...string) SelectQuery {
	return SelectQuery(string(s) + " " + strings.Join(query, ", "))
}

func (s SelectStatment) QueryAll() SelectQuery {
	return s.Queries("*")
}

func (s SelectQuery) From(table string) SelectFrom {
	return SelectFrom(string(s) + " " + _FROM_ + " " + table)
}

func (s SelectFrom) Where(condition string) SelectWhere {
	return SelectWhere(string(s) + " " + _WHERE_ + " " + condition)
}
