package main

import (
	"context"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"os/exec"

	_ "github.com/lib/pq"
)

type PathType string

const (
	dockerStart            = "docker-start"
	dockerOutput           = "./docker.log"
	host                   = "localhost"
	port                   = 5433
	containerName          = "go-postgres"
	user                   = "go-postgres"
	password               = "go-postgres"
	dbname                 = "go-postgres"
	pathTypeWeb   PathType = "web"
	pathTypeFile  PathType = "file"
	emptyJSON              = "{}"
)

/*
DATABASE SCHEMA
	id: id
	path: string "/some/path"
	path_type: "web" or "text"
	dict: jsonb {"string": number (int)}
*/

// runs an example server, can be accessed for example by
// curl 'http://localhost:8080/search?q=key1'

func main() {
	container := exec.Command("/bin/bash", dockerStart)

	outfile, err := os.Create(dockerOutput)
	checkIOError(err)
	container.Stdout = outfile
	container.Stderr = outfile
	defer outfile.Close()

	go container.Run()

	db := connectToDB()

	server := &http.Server{
		Addr: ":8080",
	}

	queryHandler := func(writer http.ResponseWriter, request *http.Request) {
		switch html.EscapeString(request.URL.Path) {
		case "/quit":
			fmt.Fprintf(writer, "Shutting down\n")
			db.Close()
			err := exec.Command("docker", "stop", "--signal", "SIGTERM", containerName).Run()
			if err != nil {
				fmt.Printf("Error stopping container: %s\n", err)
			}
			go server.Shutdown(context.Background())

		case "/all":
			defer func() {
				if err := recover(); err != nil {
					fmt.Fprintf(writer, "SQL failed: %s\n", err)
				}
			}()
			queryAll(db, writer)

		case "/search":
			// expects query in form "q=term", multiple terms can be given with form:
			// "q=term1&q=term"
			// other query keys are ignored
			queryTerms := request.URL.Query()["q"]
			defer func() {
				if err := recover(); err != nil {
					fmt.Fprintf(writer, "SQL query failed: %s\n", err)
				}
			}()
			queryJSONKeysAll(db, writer, queryTerms)

		case "/add":
			// expects query in form "p=path", multiple paths can be given with form:
			// "p=path1&p=path2"
			// other query keys are ignored
			paths := request.URL.Query()["p"]
			for _, path := range paths {
				defer func() {
					if err := recover(); err != nil {
						fmt.Fprintf(writer, "Insert failed: %s\n", err)
					}
				}()
				insertRow(db, path, pathTypeFile, emptyJSON)
			}

		default:
			fmt.Printf("Error 404 for path: %s\n", request.URL)
			fmt.Fprintf(writer, "Error 404\n")
		}
	}

	http.HandleFunc("/", queryHandler)

	log.Fatal(server.ListenAndServe())
}

func checkIOError(err error) {
	if err != nil {
		panic(err)
	}
}
