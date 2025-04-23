package server

import (
	"database/sql"
	"strconv"
	"strings"
)

const (
	_INSERT_ = "INSERT"
	_INTO_   = "INTO"
	_VALUES_ = "VALUES"
)

type SQLObject interface {
	SQLGetName() string

	SQLGetFields() []string

	SQLGetValues() []any
}

type objectTemplate string
type valueSubstitution string

type Statment string

func InsertInto(db *sql.DB, object SQLObject) (sql.Result, error) {
	stmt := InsertIntoStatment(object)

	return db.Exec(string(stmt), object.SQLGetValues()...)
}

func InsertIntoStatment(template SQLObject) Statment {

	return insertIntoStatment(sqlTemplate(template), sqlValueSubstition(template))
}

func insertIntoStatment(template objectTemplate, sub valueSubstitution) Statment {
	list := []string{
		_INSERT_,
		_INTO_,
		string(template),
		string(sub),
	}
	return Statment(strings.Join(list, " "))
}

func sqlTemplate(template SQLObject) objectTemplate {
	name := template.SQLGetName()
	fields := template.SQLGetFields()

	return objectTemplate(name + " (" + strings.Join(fields, ",") + ")")
}

func sqlValueSubstition(template SQLObject) valueSubstitution {
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
