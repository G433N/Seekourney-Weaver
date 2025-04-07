package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	link      = "http://localhost:8080"
	search    = "/search?"
	searchKey = "q"
	add       = "/add?"
	addKey    = "p"
	quit      = "/quit"
	all       = "/all"
)

// args are formatted as commandline arguments would be (client <command> [<args>])
func Run(args []string) {
	client := http.Client{}

	switch args[1] {
	case "search":
		searchForTerms(client, args[2:])
	case "add":
		addPath(client, args[2:])
	case "all":
		getAll(client)
	case "quit":
		shutdownServer(client)
	default:
		fmt.Println("usage:\tclient <search | add | all | quit> [<args>]")
		os.Exit(1)
	}
}

func checkHTTPError(err error) {
	if err != nil {
		panic(err)
	}
}

func printResponse(response *http.Response) {
	bytes, _ := io.ReadAll(response.Body)
	fmt.Print(string(bytes))
}

func searchForTerms(client http.Client, terms []string) {
	values := url.Values{}
	for _, term := range terms {
		values.Add(searchKey, term)
	}
	response, err := client.Get(link + search + values.Encode())
	checkHTTPError(err)
	printResponse(response)
}

func addPath(client http.Client, paths []string) {
	values := url.Values{}
	for _, term := range paths {
		values.Add(addKey, term)
	}
	response, err := client.Get(link + add + values.Encode())
	checkHTTPError(err)
	printResponse(response)
}

func getAll(client http.Client) {
	response, err := client.Get(link + all)
	checkHTTPError(err)
	printResponse(response)
}

func shutdownServer(client http.Client) {
	response, err := client.Get(link + quit)
	checkHTTPError(err)
	printResponse(response)
}
