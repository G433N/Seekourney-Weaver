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
func FreqMap(db *sql.DB, word utils.Word, norm utils.Normalizer) (utils.WordFrequencyMap, error) {

	wordStr := string(word)

	subQuery :=
		"collection_id = (SELECT collection_id FROM collection WHERE normalizer = $2);"

	json := JsonValue("words", wordStr, "score")
	q := Select().Queries("path", json).From("document").Where("words ?& $1 and " + subQuery)

	w := []string{wordStr}

	insert := func(res *utils.WordFrequencyMap, sqlRes sqlResult) {
		(*res)[sqlRes.path] = sqlRes.score
	}

	result := make(utils.WordFrequencyMap)
	err := ExecScan(db, string(q), &result, insert, pq.StringArray(w), norm)
	return result, err
}
