package server

import (
	"context"
	"database/sql"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const (
	serverAddress    = ":8080"
	dockerStart      = "./docker-start"
	dockerOutputFile = "./docker.log"
	host             = "localhost"
	port             = 5433
	containerName    = "go-postgres"
	user             = "go-postgres"
	password         = "go-postgres"
	dbname           = "go-postgres"
	emptyJSON        = "{}"
)

// Used to params used by server query handler functions
type serverFuncParams struct {
	server *http.Server
	writer io.Writer
	db     *sql.DB
}

// Starts the docker container using the command defined in dockerStart.
// Blocks until the container is closed
func startDocker() {
	container := exec.Command("/bin/bash", dockerStart)

	outfile, err := os.Create(dockerOutputFile)
	checkIOError(err)
	container.Stdout = outfile
	container.Stderr = outfile

	container.Run()
	outfile.Close()
}

// Stops the docker container, will finish the command started by startDocker()
func stopDocker() {
	err := exec.Command("docker", "stop", "--signal", "SIGTERM",
		containerName).Run()

	if err != nil {
		panic(fmt.Sprintf("Error stopping container: %s\n", err))
	}
}

/*
Runs a http server with a postgres instance within docker container,
can be accessed for example by `curl 'http://localhost:8080/search?q=key1'`
or using the client package: `go run . client <command>`

The server accepts the following paths as commands:

/all - Lists all paths in database, probably won't be used in production but
helpful for tests

/search - Query database, will return all paths containing given keywords.
Keywords are sent using http query under the key 'q'

/add - adds one or several paths to the database, paths are sent using http
query under the key 'p'

/quit - Shuts down the server
*/
func Run(args []string) {
	go startDocker()

	db := connectToDB()

	server := &http.Server{
		Addr: serverAddress,
	}

	queryHandler := func(writer http.ResponseWriter, request *http.Request) {
		serverParams := serverFuncParams{server: server, writer: writer, db: db}

		switch html.EscapeString(request.URL.Path) {
		case "/all":
			handleAll(serverParams)
		case "/search":
			handleSearch(serverParams, request.URL.Query()["q"])
		case "/add":
			handleAdd(serverParams, request.URL.Query()["p"])
		case "/quit":
			handleQuit(serverParams)
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

// Calls recover and writes a message to writer if an SQL function panic'd.
func recoverSQLError(writer io.Writer) {
	if err := recover(); err != nil {
		fmt.Fprintf(writer, "SQL failed: %s\n", err)
	}
}

func handleAll(serverParams serverFuncParams) {
	defer recoverSQLError(serverParams.writer)
	queryAll(serverParams.db, serverParams.writer)
}

func handleSearch(serverParams serverFuncParams, keys []string) {
	defer recoverSQLError(serverParams.writer)
	queryJSONKeysAll(serverParams.db, serverParams.writer, keys)
}

func handleAdd(serverParams serverFuncParams, paths []string) {
	for _, path := range paths {
		_, err := insertRow(serverParams.db, Page{path: path, pathType: pathTypeFile})
		if err != nil {
			fmt.Fprintf(serverParams.writer, "SQL failed: %s\n", err)
		}
	}
}

func handleQuit(serverParams serverFuncParams) {
	fmt.Fprintf(serverParams.writer, "Shutting down\n")

	serverParams.db.Close()
	stopDocker()

	// This needs to be called as a goroutine because the handler needs to return
	// before the server can shutdown
	go serverParams.server.Shutdown(context.Background())
}
