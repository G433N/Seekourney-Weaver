package search

import (
	"database/sql"
	"log"
	"seekourney/core/config"
	"seekourney/core/database"
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

	for word := range words.WordsIter(string(query)) {
		word = config.Normalizer.Word(word)

		freqMap, err := database.FreqMap(db, word)

		if err != nil {
			log.Printf("Error: %s\n", err)
			continue
		}

		for path, freq := range freqMap {
			result[path] += utils.Score(freq)
		}
	}

	return topN(scoreMapIntoSearchResult(result), 10)
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
