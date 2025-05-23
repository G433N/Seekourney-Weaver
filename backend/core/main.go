package main

import (
	"os"
	"seekourney/core/server"
	"seekourney/utils/timing"
)

func init() {
	// Initialize the timing package
	timing.Init(timing.Default())
}

// Usage for running server or client: `go run . <server | client>`
func main() {
	t := timing.Measure(timing.Main)
	defer t.Stop()

	// check commandline args to run server or client
	server.Run(os.Args[1:])

}
