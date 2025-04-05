package main

import (
	"indexer/document"
	"log"
)

func main() {

	d, err := document.DocumentFromFile("text.txt")

	if err != nil {
		log.Fatal(err)
	}

	d.GetWordsSorted()

	// for _, p := range pairs {
	// 	fmt.Printf("%s: %d\n", p.Word, p.Freq)
	// }
}
