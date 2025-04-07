package search

import (
	"log"
	"seekourney/config"
	"seekourney/folder"
	"seekourney/indexing"
	"seekourney/timing"
	"seekourney/words"
	"sort"
)

type SearchResult struct {
	Path  string
	Value int
}

// / scoreWord takes a folder, a reverse mapping and a word
// It returns a map of document paths and their corresponding score of the word
// Higher score means more relevant document
func scoreWord(f *folder.Folder, rm map[string][]string, word string) map[string]int {

	paths, ok := rm[word]
	if !ok {
		log.Printf("Word %s not found in reverse mapping", word)
		return make(map[string]int)
	}

	m := make(map[string]int)

	for _, path := range paths {
		if path == "" {
			log.Printf("ERROR: Path is empty\n")
			continue
		}

		doc, ok := f.GetDoc(path)
		if !ok {
			log.Printf("Document %s not found in folder\n", path)
			continue
		}

		// freq = 0 if not found
		freq := doc.Words[word]
		m[path] += freq
	}

	return m
}

// search takes a folder, a reverse mapping and a query
// It returns a map of document paths and their corresponding score of the query
// Higher score means more relevant document
func search(funcId config.NormalizeWordID, f *folder.Folder, rm map[string][]string, query string) map[string]int {
	m := make(map[string]int)

	for word := range words.WordsIter(query) {
		word = indexing.NormalizeWord(funcId, word)

		res := scoreWord(f, rm, word)

		for path, value := range res {
			m[path] += value
		}
	}

	return m
}

// searchParrallel is a parallel version of the search function, currently slower
func searchParrallel(funcId config.NormalizeWordID, f *folder.Folder, rm map[string][]string, query string) map[string]int {

	// TODO: This is currently slower than the normal search function, I think caching is faster / Marcus
	m := make(map[string]int)

	channel := make(chan map[string]int)
	amount := 0

	for word := range words.WordsIter(query) {
		amount++
		go func(word string) {
			word = indexing.NormalizeWord(funcId, word)
			channel <- scoreWord(f, rm, word)
		}(word)
	}

	for range amount {
		res := <-channel
		for path, value := range res {
			m[path] += value
		}
	}

	return m
}

// Search performs a search on the folder using the reverse mapping
// It returns a slice of SearchResult sorted by value in descending order, max 10 results
func Search(c *config.Config, f *folder.Folder, rm map[string][]string, query string) []SearchResult {

	// TODO: Support more than 10 results

	t := timing.Mesure(timing.Search)
	defer t.Stop()

	var m map[string]int

	if c.ParrallelSearching {
		m = searchParrallel(c.NormalizeWordFunc, f, rm, query)
	} else {
		m = search(c.NormalizeWordFunc, f, rm, query)
	}

	// Convert map to slice of SearchResult
	results := make([]SearchResult, 0, len(m))
	for path, value := range m {
		results = append(results, SearchResult{Path: path, Value: value})
	}

	// Sort results by value
	sort.Slice(results, func(i, j int) bool {
		return results[i].Value > results[j].Value
	})

	if len(results) < 10 {
		return results
	}

	return results[:10] // Return top 10 results
}
