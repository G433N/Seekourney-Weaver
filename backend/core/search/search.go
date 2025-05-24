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

	for word := range words.WordsIter(string(query)) {
		word = config.Normalizer.NormalizeWord(word)

		freqMap, err := database.FreqMap(db, word)
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
