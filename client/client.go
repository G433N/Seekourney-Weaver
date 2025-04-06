package main

import (
	"flag"
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

func main() {
	fmt.Println("client starting")
	client := http.Client{}

	var command string
	flag.StringVar(&command, "command", "default", "search, add, all or quit")
	flag.Parse()

	// USAGE: ./client -command=<search, add, all or quit> <arg1 arg2 ... argN>
	// EXAMPLE: ./client -command=search key1 key2
	// EXAMPLE: ./client -command=all
	switch command {
	case "search":
		searchForTerms(client, os.Args[2:])
	case "add":
		addPath(client, os.Args[2:])
	case "all":
		getAll(client)
	case "quit":
		shutdownServer(client)
	default:
		panic("Error: invalid command")
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
