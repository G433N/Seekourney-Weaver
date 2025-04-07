package main

import (
	"fmt"
	"indexer/words"
	"log"
	"os"
	"sort"
	"strings"

	"example.com/client"
	"example.com/server"
)

// Usage for running server or client: `go run . <server | client>`
func main() {
	// check commandline args to run server or client
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "client":
			client.Run(os.Args[1:])
			return
		case "server":
			// right now server does not take any commandline arguments
			server.Run(os.Args[1:])
			return
		}
	}

	content, err := os.ReadFile("text.txt")
	if err != nil {
		log.Fatal(err)
	}

	s := string(content)

	wordList := make(map[string]int)

	for w := range words.WordsIter(s) {

		l := strings.ToLower(w)

		wordList[l]++
	}

	type Pair struct {
		Key   string
		Value int
	}

	pairs := make([]Pair, 0)

	for k, v := range wordList {
		pairs = append(pairs, Pair{k, v})
	}

	sort.Slice(pairs, func(i, j int) bool { return pairs[i].Value > pairs[j].Value })

	for _, p := range pairs {
		fmt.Printf("%s: %d\n", p.Key, p.Value)
	}
}
