package format

import (
	"log"
	"seekourney/utils"
	"strconv"

	"github.com/savioxavier/termlink"
)

// Bold converts text to bold with ANSI escape codes.
func Bold(text string) string {
	return "\033[1m" + text + "\033[0m"
}

// Italic converts text to italic with ANSI escape codes.
func Italic(text string) string {
	return "\033[3m" + text + "\033[0m"
}

// LightBlue converts text to lightblue with ANSI escape codes.
func LightBlue(text string) string {
	return "\033[94m" + text + "\033[0m"
}

// Green converts text to green with ANSI escape codes.
func Green(text string) string {
	return "\033[92m" + text + "\033[0m"
}

// PrintSearchResponse pretty-prints a search response from Core.
func PrintSearchResponse(response utils.SearchResponse) {
	// Perform search using the folder and reverse mapping

	log.Printf(
		"--- Search results for query '%s' ---\n",
		Bold(Italic(string(response.Query))),
	)
	for n, result := range response.Results {
		path := string(result.Path)
		score := float64(result.Score)

		var source string

		switch result.Source {
		case utils.SourceLocal:
			source = "local"
		case utils.SourceWeb:
			source = "web"
		default:
			source = "unknown"
		}

		link := termlink.Link(path, path)
		log.Printf(
			"%d. Path: %s Score: %s, Source: %s\n",
			n,
			LightBlue(Bold(link)),
			Green(strconv.FormatFloat(score, 'f', 2, 64)),
			Bold(source),
		)
	}
}
