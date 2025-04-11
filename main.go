package main

import (
	"os"
	"seekourney/client"
	"seekourney/server"
	"seekourney/timing"
)

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
		case "server":
			// right now server does not take any commandline arguments
			server.Run(os.Args[1:])
		case "search":
			ManualSearch(os.Args[2:])
		}
	} else {
		AutomaticSearch()
	}
}
