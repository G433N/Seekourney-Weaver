package database

import (
	"database/sql"
	"seekourney/utils"

	"github.com/lib/pq"
)

// sqlResult is a struct that implements the SQLScan interface.
// Its used for WordFrequencyMap
type sqlResult struct {
	path  utils.Path
	score utils.Frequency
}

// SQLScan scans a SQL row into an sqlResult object.
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

// FreqMap returns a map of paths to frequencies for a given word.
func FreqMap(
	db *sql.DB,
	word utils.Word,
	plusWords []string,
	minusWords []string) (utils.WordFrequencyMap, error) {

	wordStr := string(word)

	json := JsonValue("words", wordStr, "score")
	q1 := Select().Queries("path", json).From("document")
	q2 := q1.Where("words ?& $1 AND NOT words ?& $2")

	requiredWords := []string(append(plusWords, wordStr))

	insert := func(res *utils.WordFrequencyMap, sqlRes sqlResult) {
		(*res)[sqlRes.path] = sqlRes.score
	}

	result := make(utils.WordFrequencyMap)
	err := ExecScan(
		db,
		string(q2),
		&result,
		insert,
		pq.StringArray(requiredWords),
		pq.StringArray(minusWords))
	return result, err
}
