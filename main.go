package main

import (
	"log"
	"seekourney/config"
	"seekourney/folder"
	"seekourney/search"
	"seekourney/timing"
	"strconv"

	"github.com/savioxavier/termlink"
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

func testSearch(c *config.Config, folder *folder.Folder, rm map[string][]string, query string) {

	// Perform search using the folder and reverse mapping
	pairs := search.Search(c, folder, rm, query)

	log.Printf("--- Search results for query '%s' ---\n", bold(italic(query)))
	for n, result := range pairs {
		path := result.Path
		score := result.Value
		link := termlink.Link(path, path)
		log.Printf("%d. Path: %s Score: %s\n", n, lightBlue(bold(link)), green(strconv.Itoa(score)))
	}
}

func main() {

	t := timing.Mesure("Main")
	defer t.Stop()

	// Load config
	config := config.Load()

	folder, err := folder.FromDir(config, "test_data")
	if err != nil {
		log.Fatal(err)
	}

	rm := folder.ReverseMappingLocal()

	words := len(rm)

	log.Printf("Words: %d\n", words)

	// TODO: Automated testing
	testSearch(config, &folder, rm, "Linear Interpolation")
	testSearch(config, &folder, rm, "Linearly Interpolate")
	testSearch(config, &folder, rm, "Color")
	testSearch(config, &folder, rm, "Color Interpolation")
	testSearch(config, &folder, rm, "Color Interpolation in 3D")

	// for word, paths := range rm {
	// 	log.Printf("Word: %s, Paths: %v\n", word, paths)
	// }
}
