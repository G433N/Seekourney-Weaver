package database

import (
	"database/sql"
	"seekourney/utils"

	"github.com/lib/pq"
)

type sqlResult struct {
	path  utils.Path
	score utils.Frequency
}

func (sqlPath sqlResult) SQLScan(rows *sql.Rows) (sqlResult, error) {
	var path utils.Path
	var score utils.Frequency
	err := rows.Scan(&path, &score)
	if err != nil {
		return sqlResult{}, err
	}
	return sqlResult{
		path:  path,
		score: score,
	}, nil
}

func (sqlResult sqlResult) IntoKey() string {
	return string(sqlResult.path)
}

func (sqlResult sqlResult) IntoValue() utils.Frequency {
	return sqlResult.score
}

// FreqMap returns a map of paths to frequencies for a given word.
func FreqMap(db *sql.DB, word utils.Word) (utils.WordFrequencyMap, error) {

	wordStr := string(word)

	json := JsonValue("words", wordStr, "score")
	q := Select().Queries("path", json).From("document").Where("words ?& $1")

	w := []string{wordStr}

	insert := func(res *utils.WordFrequencyMap, sqlRes sqlResult) {
		(*res)[sqlRes.path] = sqlRes.score
	}

	result := make(utils.WordFrequencyMap)
	err := ExecScan(db, string(q), &result, insert, pq.StringArray(w))
	return result, err
}
