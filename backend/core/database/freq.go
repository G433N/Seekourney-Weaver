package database

import (
	"database/sql"
	"seekourney/utils"
	"strconv"

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
	minusWords []string,
	quotes []string) (utils.WordFrequencyMap, error) {

	wordStr := string(word)
	json := JsonValue("words", wordStr, "score")

	patternString := "%"

	for quoteIndex := range quotes {
		currentQuote := quotes[quoteIndex]
		patternString += currentQuote + "%"
	}

	query := string(Select().Queries("path", json).
		From("document").
		Where("words ?& $1"))

	nextArgumentNumber := 2

	if len(minusWords) > 0 {
		query += " AND words ?& $" + strconv.Itoa(nextArgumentNumber)
		nextArgumentNumber = 3
	}

	if len(quotes) > 0 {
		query += " AND raw_text LIKE $" +
			strconv.Itoa(nextArgumentNumber)
	}

	requiredWords := append(plusWords, wordStr)

	insert := func(res *utils.WordFrequencyMap, sqlRes sqlResult) {
		(*res)[sqlRes.path] = sqlRes.score
	}

	result := make(utils.WordFrequencyMap)

	var err error
	if len(minusWords) > 0 {
		if len(quotes) > 0 {
			err = ExecScan(
				db,
				string(query),
				&result,
				insert,
				pq.StringArray(requiredWords),
				pq.StringArray(minusWords),
				patternString,
			)
		} else {
			err = ExecScan(
				db,
				string(query),
				&result,
				insert,
				pq.StringArray(requiredWords),
				pq.StringArray(minusWords),
			)
		}
	} else {
		if len(quotes) > 0 {
			err = ExecScan(
				db,
				string(query),
				&result,
				insert,
				pq.StringArray(requiredWords),
				patternString,
			)
		} else {
			err = ExecScan(
				db,
				string(query),
				&result,
				insert,
				pq.StringArray(requiredWords),
			)
		}
	}

	return result, err
}
