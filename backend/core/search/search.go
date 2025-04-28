package search

import (
	"database/sql"
	"log"
	"seekourney/core/config"
	"seekourney/core/database"
	"seekourney/core/folder"
	"seekourney/core/normalize"
	"seekourney/utils"
	"seekourney/utils/timing"
	"seekourney/utils/words"
	"sort"
)

type SearchResult = utils.SearchResult

// / scoreWord takes a folder, a reverse mapping and a word
// It returns a map of document paths and their corresponding score of the word
// Higher score means more relevant document
func scoreWord(
	folder *folder.Folder,
	rm utils.ReverseMap,
	word utils.Word,
) utils.ScoreMap {

	paths, ok := rm[word]
	if !ok {
		log.Printf("Word %s not found in reverse mapping", word)
		return make(utils.ScoreMap)
	}

	result := make(utils.ScoreMap)

	for _, path := range paths {
		if path == "" {
			log.Printf("ERROR: Path is empty\n")
			continue
		}

		doc, ok := folder.GetDoc(path)
		if !ok {
			log.Printf("Document %s not found in folder\n", path)
			continue
		}

		// freq = 0 if not found
		freq := doc.Words[word]
		result[path] += utils.Score(freq)
	}

	return result
}

// search takes a folder, a reverse mapping and a query
// It returns a map of document paths and their corresponding score of the query
// Higher score means more relevant document
func search(
	normalize normalize.Normalizer,
	folder *folder.Folder,
	rm utils.ReverseMap,
	query string,
) utils.ScoreMap {
	result := make(utils.ScoreMap)

	for word := range words.WordsIter(query) {
		word = normalize.Word(word)

		res := scoreWord(folder, rm, word)

		for path, value := range res {
			result[path] += value
		}
	}

	return result
}

// searchParrallel is a parallel version of the search function, currently
// slower
func searchParrallel(
	normalize normalize.Normalizer,
	folder *folder.Folder,
	rm utils.ReverseMap,
	query string,
) utils.ScoreMap {

	// TODO: This is currently slower than the normal search function, I think
	// caching is faster / Marcus
	result := make(utils.ScoreMap)

	channel := make(chan utils.ScoreMap)
	amount := 0

	for word := range words.WordsIter(query) {
		amount++
		go func(word utils.Word) {
			word = normalize.Word(word)
			channel <- scoreWord(folder, rm, word)
		}(word)
	}

	for range amount {
		res := <-channel
		for path, value := range res {
			result[path] += value
		}
	}

	return result
}

// Search performs a search on the folder using the reverse mapping
// It returns a slice of SearchResult sorted by value in descending order, max
// 10 results
func Search(
	config *config.Config,
	f *folder.Folder,
	rm utils.ReverseMap,
	query string,
) []SearchResult {

	// TODO: Support more than 10 results

	t := timing.Measure(timing.Search)
	defer t.Stop()

	var searchResult utils.ScoreMap

	if config.ParrallelSearching {
		searchResult = searchParrallel(config.Normalizer, f, rm, query)
	} else {
		searchResult = search(config.Normalizer, f, rm, query)
	}

	// Convert map to slice of SearchResult
	results := make([]SearchResult, 0, len(searchResult))
	for path, score := range searchResult {
		doc, ok := f.GetDoc(path)

		if !ok {
			log.Fatalf(
				"Document %s not found in folder, this should be imposible\n",
				path,
			)
		}

		results = append(
			results,
			SearchResult{Path: path, Score: score, Source: doc.Source},
		)
	}

	// Sort results by value
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if len(results) < 10 {
		return results
	}

	return results[:10] // Return top 10 results
}

/// New

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
