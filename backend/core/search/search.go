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
)

// Parses a query for filters and removes anything
// that is part of the "-" and quote filters.
// The return data is the query plus slices with words associated with the filters.
// TODO: This implementation assumes correct syntax. E.g, 'terminal +dog -cat "bad"'
func parseQuery(query utils.Query) utils.ParsedQuery {
	parsedQuery := utils.ParsedQuery{
		ModifiedQuery: query,
		PlusWords:     make([]string, 0),
		MinusWords:    make([]string, 0),
		Quotes:        make([]string, 0)}

	currentFilterWord := ""
	currentQuote := ""
	inPlus := false
	inMinus := false
	inQuote := false

	for byteIndex := range query {
		currentByte := string(query[byteIndex])

		if currentByte == " " && inPlus {
			parsedQuery.PlusWords = append(parsedQuery.PlusWords, currentFilterWord)
			inPlus = false
			currentFilterWord = ""
		}

		if currentByte == " " && inMinus {
			parsedQuery.MinusWords = append(parsedQuery.MinusWords, currentFilterWord)
			inMinus = false
			currentFilterWord = ""
		}

		if currentByte == "\"" && inQuote {
			parsedQuery.Quotes = append(parsedQuery.Quotes, currentQuote)
			inQuote = false
			currentQuote = ""
		}

		if inPlus {
			currentFilterWord += currentByte
			query += utils.Query(currentByte)
		} else if inMinus {
			currentFilterWord += currentByte
		} else if inQuote {
			currentQuote += currentByte
		} else {
			query += utils.Query(currentByte)
		}

		if currentByte == "+" && !inPlus && !inMinus && !inQuote {
			inPlus = true
		}

		if currentByte == "-" && !inPlus && !inMinus && !inQuote {
			inMinus = true
		}

		if currentByte == "\"" {
			inQuote = true
		}
	}

	if inPlus && len(currentFilterWord) > 0 {
		parsedQuery.PlusWords = append(parsedQuery.PlusWords, currentFilterWord)
	}

	if inMinus && len(currentFilterWord) > 0 {
		parsedQuery.MinusWords = append(parsedQuery.MinusWords, currentFilterWord)
	}

	return parsedQuery
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

	parsedQuery := parseQuery(query)

	for word := range words.WordsIter(string(parsedQuery.ModifiedQuery)) {
		word = config.Normalizer.Word(word)

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
