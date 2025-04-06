package search

import (
	"indexer/folder"
	"indexer/indexing"
	"indexer/timing"
	"indexer/words"
	"log"
	"sort"
)

type SearchResult struct {
	Path  string
	Value int
}

func Search(f *folder.Folder, rMap map[string][]string, query string) []SearchResult {

	t := timing.Mesure("Search")
	defer t.Stop()

	m := make(map[string]int)

	for word := range words.WordsIter(query) {
		word = indexing.NormalizeWord(word)
		paths, ok := rMap[word]
		if !ok {
			log.Fatalf("Word %s not found in reverse mapping", word)
			continue
		}

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
