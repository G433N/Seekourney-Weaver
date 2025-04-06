package main

import (
	"indexer/document"
	"indexer/folder"
	"indexer/indexing"
	"indexer/search"
	"indexer/timing"
	"indexer/utils"
	"log"
	"strconv"
	"strings"

	"github.com/savioxavier/termlink"
)

// TODO: Remove dubble information in fuction names, why is it called document.Documentblbla?

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

func testSearch(c *search.Config, folder *folder.Folder, rm map[string][]string, query string) {

	// Perform search using the folder and reverse mapping
	pairs := c.Search(folder, rm, query)

	log.Printf("--- Search results for query '%s' ---\n", bold(italic(query)))
	for n, result := range pairs {
		path := result.Path
		score := result.Value
		link := termlink.Link(path, path)
		log.Printf("%d. Path: %s Score: %s\n", n, lightBlue(bold(link)), green(strconv.Itoa(score)))
	}
}

func testIndexConfig() *indexing.Config {
	return indexing.NewIndexConfig(strings.ToLower)
}

// TODO: Json config

func testFolderConfig(index *indexing.Config) *folder.FolderConfig {
	documentConfig := document.DocumentConfigFromIndexConfig(index)
	dirConfig := utils.NewWalkDirConfig().SetAllowedExts([]string{".txt", ".md", ".json", ".xml", ".html", "htm", ".xhtml", ".csv"})
	folderConfig := folder.FolderConfigFromConfig(dirConfig, documentConfig)
	return folderConfig
}

func testSearchConfig(index *indexing.Config) *search.Config {
	return search.NewSearchConfig(index)
}

func main() {

	t := timing.Mesure("Main")
	defer t.Stop()

	indexConfig := testIndexConfig()
	folderConfig := testFolderConfig(indexConfig)
	searchConfig := testSearchConfig(indexConfig)

	folder, err := folderConfig.FolderFromDir("test_data")
	if err != nil {
		log.Fatal(err)
	}

	rm := folder.ReverseMappingLocal()

	words := len(rm)

	log.Printf("Words: %d\n", words)

	// TODO: Automated testing
	testSearch(searchConfig, &folder, rm, "Linear Interpolation")
	testSearch(searchConfig, &folder, rm, "Linearly Interpolate")
	testSearch(searchConfig, &folder, rm, "Color")
	testSearch(searchConfig, &folder, rm, "Color Interpolation")
	testSearch(searchConfig, &folder, rm, "Color Interpolation in 3D")

	// for word, paths := range rm {
	// 	log.Printf("Word: %s, Paths: %v\n", word, paths)
	// }
}
