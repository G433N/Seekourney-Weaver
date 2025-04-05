package main

import (
	"indexer/document"
	"indexer/timing"
	"log"
)

func main() {

	t := timing.Mesure("Main")

	docs, err := document.DocumentsFromDir("test_data")
	if err != nil {
		log.Fatal(err)
	}

	// for _, doc := range docs {
	// 	doc.DebugPrint()
	// }

	rm := document.ReverseMapping(&docs)

	words := len(rm)

	log.Printf("Words: %d\n", words)
	// for word, paths := range rm {
	// 	log.Printf("Word: %s, Paths: %v\n", word, paths)
	// }

	// for _, p := range pairs {
	// 	fmt.Printf("%s: %d\n", p.Word, p.Freq)
	// }

	t.Stop()
}
