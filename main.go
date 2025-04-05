package main

import (
	"indexer/folder"
	"indexer/search"
	"indexer/timing"
	"log"
)

// func Search() []SearchResult {
//
//
// }

func main() {

	t := timing.Mesure("Main")
	defer t.Stop()

	folder, err := folder.FolderFromDir("test_data")
	if err != nil {
		log.Fatal(err)
	}

	// for _, doc := range docs {
	// 	doc.DebugPrint()
	// }

	rm := folder.ReverseMappingLocal()

	words := len(rm)

	log.Printf("Words: %d\n", words)
	pairs := search.Search(&folder, rm, "Jesus Christ is annoying")

	for _, result := range pairs {
		path := result.Path
		score := result.Value
		log.Printf("Path: %s, Score: %d\n", path, score)
	}
	// for word, paths := range rm {
	// 	log.Printf("Word: %s, Paths: %v\n", word, paths)
	// }

	// for _, p := range pairs {
	// 	fmt.Printf("%s: %d\n", p.Word, p.Freq)
	// }

}
