package main

import (
	"github.com/savioxavier/termlink"
	"fmt"
	pdftotext "seelourney/pdftotxt"
	"seekourney/words"
	"log"
	"os"
	"seekourney/client"
	"seekourney/config"
	"seekourney/folder"
	"seekourney/indexing/localtext"
	"seekourney/search"
	"seekourney/server"
	"seekourney/timing"
	"seekourney/utils"
	"strconv"
)

// TODO: All this should be moved to client side
func bold(text string) string {
	return "\033[1m" + text + "\033[0m"
}

func italic(text string) string {
	return "\033[3m" + text + "\033[0m"
}

func lightBlue(text string) string {
	return "\033[94m" + text + "\033[0m"
}

func green(text string) string {
	return "\033[92m" + text + "\033[0m"
}

func testSearch(c *config.Config, folder *folder.Folder, rm utils.ReverseMap, query string) {

	// Perform search using the folder and reverse mapping
	pairs := search.Search(c, folder, rm, query)

	log.Printf("--- Search results for query '%s' ---\n", bold(italic(query)))
	for n, result := range pairs {
		path := string(result.Path)
		score := int(result.Value)
		link := termlink.Link(path, path)
		log.Printf("%d. Path: %s Score: %s\n", n, lightBlue(bold(link)), green(strconv.Itoa(score)))
	}
}

func init() {
	// Initialize the timing package
	timing.Init(timing.Default())
}

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


	pdftotext.Run()

	t := timing.Measure(timing.Main)
	defer t.Stop()

	// Load config
	config := config.Load()

	// Load local file config
	localConfig := localtext.Load(config)

	// TODO: Later when documents comes over the network, we can still use the same code. since it is an iterator
	folder := folder.FromIter(config.Normalizer, localConfig.IndexDir("test_data"))

	rm := folder.ReverseMappingLocal()

	queries := []string{
		"Linear Interpolation",
		"Linearly Interpolate",
		"Color",
		"Color Interpolation",
		"Color Interpolation in 3D",
		"macro",
		"neovim",
		"mozilla",
		"curl",
		"math",
	}

	return

	content, err := os.ReadFile("text.txt")
	if err != nil {
		log.Fatal(err)

	// TODO: Automated testing
	for _, query := range queries {
		testSearch(config, &folder, rm, query)
	}

	files := folder.GetDocAmount()
	words := len(rm)

	log.Printf("Files: %d, Words: %d\n", files, words)

	if files == 0 {
		log.Println("No files found, run make downloadTestFiles to download test files")
	}
}
