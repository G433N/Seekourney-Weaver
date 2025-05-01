package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"seekourney/tui/format"
	"seekourney/utils"
	"seekourney/utils/timing"
)

// Strings for building URLs in HTTP requests
const (
	link      = "http://localhost:8080"
	search    = "/search?"
	searchKey = "q"
	add       = "/add?"
	addKey    = "p"
	quit      = "/quit"
	all       = "/all"
)

// argumentError prints a usage string and terminates client process.
func argumentError() {
	fmt.Println("usage: client <command> [<args>]")
	fmt.Println()
	fmt.Println("available commands")
	fmt.Println("  all                 request all pages in database")
	fmt.Println("  search [key ...]    request all pages containing keys")
	fmt.Println("  add    [path ...]   add paths to database")
	fmt.Println("  quit                request the server to shutdown")
	os.Exit(1)
}

// init runs before any other code in this package.
// https://golangdocs.com/init-function-in-golang
func init() {
	// Initialize the timing package
	timing.Init(timing.Default())
}

// main takes an array of arguments args as input and calls different functions
// to send a HTTP request to the server
// The first element in args is ignored.
// The second element in args is the command.
// The rest of the elements are for commands that require extra arguments.
// If no valid command is found in the first element, the program terminates.
// args are formatted like commandline arguments ("client", command, args)
func main() {

	args := os.Args

	if len(args) < 2 {
		argumentError()
	}

	switch args[1] {
	case "search":
		searchForTerms(args[2:])
	case "add":
		addPath(args[2:])
	case "all":
		getAll()
	case "quit":
		shutdownServer()
	default:
		argumentError()
	}
}

// checkHTTPError panics on a given error if an error
// occured when sending a request.
func checkHTTPError(err error) {
	if err != nil {
		panic(err)
	}
}

// printResponse prints a HTTP response to stdout.
func printResponse(response *http.Response) {
	bytes, _ := io.ReadAll(response.Body)
	fmt.Print(string(bytes))
}

// searchForTerms requests a search for given terms through Core,
// and prints the results.
// Handler for command /search.
func searchForTerms(terms []string) {
	sw := timing.Measure(timing.Search)
	defer sw.Stop()

	values := url.Values{}
	for _, term := range terms {
		values.Add(searchKey, term)
	}
	response, err := http.Get(link + search + values.Encode())
	checkHTTPError(err)

	bytes, _ := io.ReadAll(response.Body)
	result := utils.SearchResponse{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		log.Println("Response:", string(bytes))
		return
	}

	format.PrintSearchResponse(result)
}

// addPath adds given paths to the database through Core,
// and prints the response.
// Handler for command /add.
func addPath(paths []string) {
	values := url.Values{}
	for _, term := range paths {
		values.Add(addKey, term)
	}
	response, err := http.Get(link + add + values.Encode())
	checkHTTPError(err)
	printResponse(response)
}

// getAll fetches all paths stored in database through Core,
// and prints them.
// Handler for command /all.
func getAll() {
	response, err := http.Get(link + all)
	checkHTTPError(err)
	printResponse(response)
}

// shutdownServer remotely shuts down Core.
// Handler for command /quit.
func shutdownServer() {
	response, err := http.Get(link + quit)
	checkHTTPError(err)
	printResponse(response)
}
