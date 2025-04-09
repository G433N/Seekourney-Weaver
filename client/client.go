package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
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

// Prints a usage string and terminates
func argumentError() {
	fmt.Println("usage:\tclient <search | add | all | quit> [<args>]")
	os.Exit(1)
}

// Run takes an array of arguments args as input and calls different functions
// to send a HTTP request to the server
// The first element in args is ignored.
// The second element in args is the command.
// The rest of the elements are for commands that require extra arguments.
// If no valid command is found in the first element, the program terminates.
// args are formatted like commandline arguments ("client", command, args)
func Run(args []string) {

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

// Panics on a given error if an error occured sending a request
func checkHTTPError(err error) {
	if err != nil {
		panic(err)
	}
}

// Prints a HTTP response to stdout
func printResponse(response *http.Response) {
	bytes, _ := io.ReadAll(response.Body)
	fmt.Print(string(bytes))
}

func searchForTerms(terms []string) {
	values := url.Values{}
	for _, term := range terms {
		values.Add(searchKey, term)
	}
	response, err := http.Get(link + search + values.Encode())
	checkHTTPError(err)
	printResponse(response)
}

func addPath(paths []string) {
	values := url.Values{}
	for _, term := range paths {
		values.Add(addKey, term)
	}
	response, err := http.Get(link + add + values.Encode())
	checkHTTPError(err)
	printResponse(response)
}

func getAll() {
	response, err := http.Get(link + all)
	checkHTTPError(err)
	printResponse(response)
}

func shutdownServer() {
	response, err := http.Get(link + quit)
	checkHTTPError(err)
	printResponse(response)
}
