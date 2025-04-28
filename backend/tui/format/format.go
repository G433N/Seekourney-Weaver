package format

import (
	"log"
	"seekourney/utils"
	"strconv"

	"github.com/savioxavier/termlink"
)

func Bold(text string) string {
	return "\033[1m" + text + "\033[0m"
}

func Italic(text string) string {
	return "\033[3m" + text + "\033[0m"
}

func LightBlue(text string) string {
	return "\033[94m" + text + "\033[0m"
}

func Green(text string) string {
	return "\033[92m" + text + "\033[0m"
}

func PrintSearchResponse(response utils.SearchResponse) {

	// Perform search using the folder and reverse mapping

	log.Printf(
		"--- Search results for query '%s' ---\n",
		Bold(Italic(response.Query)),
	)
	for n, result := range response.Results {
		path := string(result.Path)
		score := int(result.Score)

		var source string

		if result.Source == utils.SourceLocal {
			source = "local"
		} else if result.Source == utils.SourceWeb {
			source = "web"
		} else {
			source = "unknown"
		}

		link := termlink.Link(path, path)
		log.Printf(
			"%d. Path: %s Score: %s, Source: %s\n",
			n,
			LightBlue(Bold(link)),
			Green(strconv.Itoa(score)),
			Bold(source),
		)
	}
}
