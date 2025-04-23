package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
)

const (
	serverAddress       string     = ":8080"
	containerStart      string     = "./docker-start"
	containerOutputFile string     = "./docker.log"
	host                string     = "localhost"
	port                int        = 5433
	containerName       string     = "go-postgres"
	user                string     = "go-postgres"
	password            string     = "go-postgres"
	dbname              string     = "go-postgres"
	emptyJSON           JSONString = "{}"
)

// Used to params used by server query handler functions
type serverFuncParams struct {
	writer io.Writer
	db     *sql.DB
	stop   context.CancelFunc
}

// Starts the database container using the command defined in containerStart.
// Blocks until the container is closed
func startContainer() {
	container := exec.Command("/bin/sh", containerStart)

	outfile, err := os.Create(containerOutputFile)
	checkIOError(err)
	container.Stdout = outfile
	container.Stderr = outfile

	err = container.Run()
	checkIOError(err)
	err = outfile.Close()
	checkIOError(err)
}

// Stops the database container, will finish the command started by
// startContainer()
func stopContainer() {
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
	go startContainer()

	db := connectToDB()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	server := &http.Server{
		Addr:        serverAddress,
		BaseContext: func(l net.Listener) context.Context { return ctx },
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		serverParams := serverFuncParams{
			writer: writer,
			db:     db,
			stop:   stop,
		}

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
	})

	go func() {
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Server encountered an error:", err)
		}
		stop()
	}()
	fmt.Println("Server online")

	// Wait until server is finished
	<-ctx.Done()

	fmt.Println("Shutting down")
	err := server.Shutdown(context.Background())
	if err != nil {
		fmt.Println("Error while shutting down server: ", err)
	}
	db.Close()
	stopContainer()
}

func checkIOError(err error) {
	if err != nil {
		panic(err)
	}
}

// Calls recover and writes a message to writer if an SQL function panic'd.
func recoverSQLError(writer io.Writer) {
	if err := recover(); err != nil {
		_, ioErr := fmt.Fprintf(writer, "SQL failed: %s\n", err)
		checkIOError(ioErr)
	}
}

// Handles an /all request, queries all rows in database and writes output to
// response writer
func handleAll(serverParams serverFuncParams) {
	defer recoverSQLError(serverParams.writer)
	rows := queryAll(serverParams.db)
	writeRows(serverParams.writer, rows)
	unsafelyClose(rows)
}

// Handles a /search request, queries database for rows containing ALL keys and
// wrties output to response writer
func handleSearch(serverParams serverFuncParams, keys []string) {
	defer recoverSQLError(serverParams.writer)
	rows := queryJSONKeysAll(serverParams.db, keys)
	writeRows(serverParams.writer, rows)
	unsafelyClose(rows)
}

// Handles an /add request, inserts a row to the database for each path given
func handleAdd(serverParams serverFuncParams, paths []string) {
	for _, path := range paths {
		_, err := insertRow(serverParams.db, Page{
			path:     path,
			pathType: PathTypeFile,
			dict:     emptyJSON,
		})
		if err != nil {
			_, ioErr := fmt.Fprintf(
				serverParams.writer,
				"SQL failed: %s\n",
				err,
			)
			checkIOError(ioErr)
		}
	}
}

// Handles a /quit request initiates the shutdown process by cancelling the
// server context
func handleQuit(serverParams serverFuncParams) {
	_, err := fmt.Fprintf(serverParams.writer, "Shutting down\n")
	checkIOError(err)

	serverParams.stop()
}
