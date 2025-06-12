package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	_, err := os.Stat("./build/seekourney-weaver")
	if !os.IsNotExist(err) {
		homeDirectory, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory:", err)
		}
		configDir := filepath.Join(homeDirectory, ".config")
		err = exec.Command("cp", "-u", "-r", "./build/seekourney-weaver",
			configDir).
			Run()
		if err != nil {
			fmt.Println("Error executing cp: ", err)
		}
	}
	// check commandline args to run server or client
	server.Run(os.Args[1:])
}
