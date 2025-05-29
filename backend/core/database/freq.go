package database

import (
	"database/sql"
	"log"
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
	j := JsonValue("words", wordStr, "score")

	pattern_string := "%"

	for quoteIndex := range quotes {
		currentQuote := quotes[quoteIndex]
		pattern_string += currentQuote + "%"
	}

	nextArgumentNumber := 2

	q := string(Select().Queries("D.path AS path", j).
		From("document AS D, path_text AS P").
		Where("D.words ?& $1"))

	if len(minusWords) > 0 {
		q += " AND NOT D.words ?& $" + strconv.Itoa(nextArgumentNumber)
		nextArgumentNumber += 1
	}

	if len(quotes) > 0 {
		q += " AND D.path = P.path AND P.plain_text LIKE $" +
			strconv.Itoa(nextArgumentNumber)

		nextArgumentNumber += 1
	}

	log.Println("Query: ", q)

	requiredWords := append(plusWords, wordStr)
	log.Println(requiredWords)
	insert := func(res *utils.WordFrequencyMap, sqlRes sqlResult) {
		(*res)[sqlRes.path] = sqlRes.score
	}

	result := make(utils.WordFrequencyMap)

	var err error
	if len(minusWords) > 0 {
		if len(quotes) > 0 {
			err = ExecScan(
				db,
				string(q),
				&result,
				insert,
				pq.StringArray(requiredWords),
				pq.StringArray(minusWords),
				pattern_string,
			)
		} else {
			err = ExecScan(
				db,
				string(q),
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
				string(q),
				&result,
				insert,
				pq.StringArray(requiredWords),
				pattern_string,
			)
		} else {
			err = ExecScan(
				db,
				string(q),
				&result,
				insert,
				pq.StringArray(requiredWords),
			)
		}
	}

	return result, err
}
