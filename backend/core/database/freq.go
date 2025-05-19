package database

import (
	"database/sql"
	"log"
	"seekourney/utils"
	"strings"

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
	j := JsonValue("words", wordStr, "score")

	pattern := make([]string, 0)
	pattern = append(pattern, "%")

	for quoteIndex := range quotes {
		currentQuote := quotes[quoteIndex]

		pattern = append(pattern, currentQuote+"%")
	}

	// execute <unnamed>: SELECT path, JSON_VALUE( words , '$.linear' ) AS score FROM document WHERE words ?& $1

	// TODO: Use fuctions to generate the query
	q := strings.Join([]string{"SELECT D.path AS path,",
		j,
		"FROM document AS D",
		"WHERE D.words ?& $1",
		// "AND NOT D.words ?& $2",
		// "AND D.path = T.path",
		// "AND T.plain_text LIKE $3",
	}, " ")

	log.Println("Query: ", q)

	requiredWords := []string(append(plusWords, wordStr))

	insert := func(res *utils.WordFrequencyMap, sqlRes sqlResult) {
		(*res)[sqlRes.path] = sqlRes.score
	}

	result := make(utils.WordFrequencyMap)

	err := ExecScan(
		db,
		string(q),
		&result,
		insert,
		pq.StringArray(requiredWords),
		// pq.StringArray(minusWords),
	// pq.StringArray(pattern)
	)
	return result, err
}

// // FreqMap returns a map of paths to frequencies for a given word.
// func FreqMap(db *sql.DB, word utils.Word, _ []string, _ []string, _ []string) (utils.WordFrequencyMap, error) {
//
// 	wordStr := string(word)
//
// 	json := JsonValue("words", wordStr, "score")
// 	q := Select().Queries("path", json).From("document").Where("words ?& $1")
//
// 	w := []string{wordStr}
//
// 	insert := func(res *utils.WordFrequencyMap, sqlRes sqlResult) {
// 		(*res)[sqlRes.path] = sqlRes.score
// 	}
//
// 	result := make(utils.WordFrequencyMap)
// 	err := ExecScan(db, string(q), &result, insert, pq.StringArray(w))
// 	return result, err
// }
