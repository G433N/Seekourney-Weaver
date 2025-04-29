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
		word = config.Normalizer.Word(word)

		freqMap, err := database.FreqMap(db, word)
		calculateIdf(freqMap, docAmount)

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

			result[path] += utils.Score(doc.CalculateTf(word) * calculateIdf(freqMap, docAmount))
		}
	}

	return topN(scoreMapIntoSearchResult(result), 10)
}

func calculateIdf(freqMap utils.WordFrequencyMap, docAmount int) float64 {

	popularity := float64(len(freqMap))

	return math.Log2(float64(docAmount) / (popularity + 1))
}

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

func topN(
	results []SearchResult,
	n int,
) []SearchResult {
	if len(results) < n {
		return results
	}

	return results[:n]
}
