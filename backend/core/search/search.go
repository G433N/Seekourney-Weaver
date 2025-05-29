package search

import (
	"database/sql"
	"log"
	"math"
	"seekourney/core/config"
	"seekourney/core/database"
	"seekourney/core/document"
	"seekourney/utils"
	"seekourney/utils/words"
	"sort"
	"strings"
)

type status int

const (
	_CONTINUELOOP_ = true
	_STOPLOOP_     = false
	_NOFILTER_     = status(0)
	_INPLUS_       = status(1)
	_INMINUS_      = status(2)
	_INQUOTE_      = status(3)
)

// updateFilterStatus returns a new status
// depending of a given character.
func updateFilterStatus(character string) status {
	if character == "+" {
		return _INPLUS_
	}

	if character == "-" {
		return _INMINUS_
	}

	if character == "\"" {
		return _INQUOTE_
	}

	return _NOFILTER_
}

// updateModifiedQuery adds words to the filter slices
// if certain condiations are satisfied.
func updateModifiedQuery(
	config *config.Config,
	parsedQuery *utils.ParsedQuery,
	currentByte string,
	currentFilterString *string,
	filterStatus status) (status, bool) {

	newStatus := filterStatus

	if (currentByte == " ") &&
		((filterStatus == _INPLUS_) || (filterStatus == _INMINUS_)) {

		normalizedWord := string(config.
			Normalizer.
			NormalizeWord(utils.Word(*currentFilterString)))

		if filterStatus == _INPLUS_ {
			parsedQuery.PlusWords = append(
				parsedQuery.PlusWords,
				normalizedWord,
			)
		} else {
			parsedQuery.MinusWords = append(
				parsedQuery.MinusWords,
				normalizedWord,
			)
		}

		newStatus = _NOFILTER_
		*currentFilterString = ""
	}

	if (currentByte == "\"") && (filterStatus == _INQUOTE_) {
		parsedQuery.Quotes = append(
			parsedQuery.Quotes,
			*currentFilterString,
		)

		newStatus = _NOFILTER_
		*currentFilterString = ""
		return newStatus, _STOPLOOP_
	}

	return newStatus, _CONTINUELOOP_
}

// Parses a query for filters and removes anything
// that is part of the "-" and quote filters.
// The return data is the query plus
// slices with words associated with the filters.
// TODO: This implementation assumes correct syntax.
// E.g, 'terminal +dog -cat "bad"'
func parseQuery(config *config.Config, query utils.Query) utils.ParsedQuery {
	parsedQuery := utils.ParsedQuery{
		ModifiedQuery: "",
		PlusWords:     make([]string, 0),
		MinusWords:    make([]string, 0),
		Quotes:        make([]string, 0)}

	currentFilterString := ""

	filterStatus := _NOFILTER_
	var shouldContinue bool

	for byteIndex := range query {
		currentByte := string(query[byteIndex])

		filterStatus, shouldContinue = updateModifiedQuery(
			config,
			&parsedQuery,
			currentByte,
			&currentFilterString,
			filterStatus,
		)

		if !shouldContinue {
			continue
		}

		if filterStatus == _NOFILTER_ {
			filterStatus = updateFilterStatus(currentByte)

			if filterStatus != _NOFILTER_ {
				continue
			}
		}

		switch filterStatus {
		case _INPLUS_:
			currentFilterString += currentByte
			parsedQuery.ModifiedQuery += utils.Query(currentByte)
		case _INMINUS_:
			currentFilterString += currentByte
		case _INQUOTE_:
			currentFilterString += currentByte
		default:
			parsedQuery.ModifiedQuery += utils.Query(currentByte)
		}
	}

	// When the last word in a query was a filter-word, add it
	if len(currentFilterString) > 0 {
		currentFilterString = string(config.
			Normalizer.
			NormalizeWord(utils.Word(currentFilterString)))

		if filterStatus == _INPLUS_ {
			parsedQuery.PlusWords = append(
				parsedQuery.PlusWords,
				currentFilterString,
			)
		} else if filterStatus == _INMINUS_ {
			parsedQuery.MinusWords = append(
				parsedQuery.MinusWords,
				currentFilterString,
			)
		}
	}

	return parsedQuery
}

// wordsFromQuotes looks in all quotes from a search query
// and returns all words.
func wordsFromQuotes(quotes []string) []string {
	retrievedWords := make([]string, 0)

	for _, quote := range quotes {
		for word := range words.WordsIter(quote) {
			retrievedWords = append(retrievedWords, string(word))
		}
	}

	return retrievedWords
}

type SearchResult = utils.SearchResult

// SqlSearch performs a search in the database using SQL.
func SqlSearch(
	config *config.Config,
	db *sql.DB,
	query utils.Query) []SearchResult {

	result := make(utils.ScoreMap)

	docAmount, err := database.RowAmount(db, "document")
	if err != nil {
		log.Printf("Error: %s\n", err)
		panic(err)
	}

	parsedQuery := parseQuery(config, query)

	stringFromQuotes := strings.Join(
		wordsFromQuotes(parsedQuery.Quotes),
		" ",
	)

	parsedQuery.ModifiedQuery += utils.Query(" " + stringFromQuotes)

	for word := range words.WordsIter(string(parsedQuery.ModifiedQuery)) {
		word = config.Normalizer.NormalizeWord(word)

		freqMap, err := database.FreqMap(
			db,
			word,
			parsedQuery.PlusWords,
			parsedQuery.MinusWords,
			parsedQuery.Quotes)

		idf := calculateIdf(freqMap, docAmount)

		if err != nil {
			log.Printf("Error: %s\n", err)
			continue
		}

		for path := range freqMap {
			doc, err := document.DocumentFromDB(db, path)

			if err != nil {
				log.Printf("Error: %s\n", err)
				continue
			}

			tf := doc.CalculateTf(word)
			result[path] += utils.Score(tf * idf)
		}
	}

	return topN(scoreMapIntoSearchResult(result), 10)
}

// calculateIdf calculates the Inverse Document Frequency (IDF)
// for a WordFrequencyMap.
// See: https://en.wikipedia.org/wiki/Tf%E2%80%93idf#Inverse_document_frequency
func calculateIdf(freqMap utils.WordFrequencyMap, docAmount int) float64 {
	popularity := float64(len(freqMap))
	return math.Log2(float64(docAmount) / (popularity + 1))
}

// scoreMapIntoSearchResult converts a ScoreMap into a slice of SearchResult.
func scoreMapIntoSearchResult(scores utils.ScoreMap) []SearchResult {
	results := make([]SearchResult, 0, len(scores))

	for path, score := range scores {
		results = append(results, SearchResult{Path: path, Score: score})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

// topN returns the top n results from the given results slice.
func topN(results []SearchResult, n int) []SearchResult {
	if len(results) < n {
		return results
	} else {
		return results[:n]
	}
}
