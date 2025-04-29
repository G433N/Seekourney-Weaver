package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"seekourney/core/config"
	"seekourney/core/database"
	"seekourney/core/document"
	"seekourney/core/folder"
	"seekourney/core/search"
	"seekourney/indexer/localtext"
	"seekourney/utils"
	"strings"
)

const (
	serverAddress       string = ":8080"
	containerStart      string = "./docker-start"
	containerOutputFile string = "./docker.log"
	host                string = "localhost"
	port                int    = 5433
	containerName       string = "go-postgres"
	user                string = "go-postgres"
	password            string = "go-postgres"
	dbname              string = "go-postgres"
	emptyJSON           string = "{}"
)

// Used to params used by server query handler functions
type serverFuncParams struct {
	server *http.Server
	writer io.Writer
	db     *sql.DB
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

// conf holds the config object for the server, gets initialized in the Run function
var conf *config.Config

// index loads the local file config and creates a folder object
// TODO: This function is temporay
func index() folder.Folder {
	// Load local file config
	localConfig := localtext.Load(conf)

	// TODO: Later when documents comes over the network, we can still use the
	// same code. since it is an iterator
	folder := folder.FromIter(
		conf.Normalizer,
		localConfig.IndexDir("test_data"),
	)

	rm := folder.ReverseMappingLocal()

	files := folder.GetDocAmount()
	words := len(rm)

	log.Printf("Files: %d, Words: %d\n", files, words)

	if files == 0 {
		log.Fatal(
			"No files found, run make downloadTestFiles to download test files",
		)
	}

	return folder
}

// insertFolder inserts all documents in the given folder into the database
func insertFolder(db *sql.DB, folder *folder.Folder) {

	for _, doc := range folder.GetDocs() {

		_, err := database.InsertInto(db, doc)

		if err != nil {
			log.Printf("Error inserting row: %s\n", err)
		}
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

	// Load config

	conf = config.Load()

	go startContainer()

	db := connectToDB()

	loadFromDisc := true

	if len(args) > 1 {
		log.Fatal("Too many arguments")
	} else if len(args) == 1 {
		switch args[0] {
		case "load":
			loadFromDisc = false
		}
	}

	if !loadFromDisc {
		log.Println("Indexing files")
		fold := index()
		insertFolder(db, &fold)
	} else {
		log.Println("Loading from disk")
	}

	server := &http.Server{
		Addr: serverAddress,
	}

	log.Println("Server started at", serverAddress)

	queryHandler := func(writer http.ResponseWriter, request *http.Request) {
		enableCors(&writer)
		serverParams := serverFuncParams{server: server, writer: writer, db: db}

		switch html.EscapeString(request.URL.Path) {
		case "/all":
			handleAll(serverParams)
		case "/search":
			handleSearchSql(serverParams, request.URL.Query()["q"])
		case "/add":
			handleAdd(serverParams, request.URL.Query()["p"])
		case "/quit":
			handleQuit(serverParams)
		}
	}
	http.HandleFunc("/", queryHandler)

	log.Fatal(server.ListenAndServe())
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
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

// sendJSON marshals the given data to JSON and writes it to the writer.
func sendJSON(writer io.Writer, data any) {

	jsonResponse, err := json.Marshal(data)
	if err != nil {
		sendError(writer, "JSON failed", err)
		return
	}

	_, err = fmt.Fprintf(writer, "%s\n", jsonResponse)
	if err != nil {
		sendError(writer, "IO failed", err)
		return
	}
}

// Handles an /all request, queries all rows in database and writes output to
// response writer
func handleAll(serverParams serverFuncParams) {
	defer recoverSQLError(serverParams.writer)
	var doc document.Document
	query := database.Select().Queries(doc.SQLGetFields()...).From("document")

	insert := func(docs *[]document.Document, doc document.Document) {
		*docs = append(*docs, doc)
	}

	docs := make([]document.Document, 0)
	err := database.ExecScan(serverParams.db, string(query), &docs, insert)

	if err != nil {
		sendError(serverParams.writer, "SQL failed", err)
		return
	}

	sendJSON(serverParams.writer, docs)

}

func handleSearchSql(serverParams serverFuncParams, keys []string) {

	defer recoverSQLError(serverParams.writer)

	if len(keys) == 0 {
		sendError(serverParams.writer, "No keys given", nil)
		return
	}

	query := utils.Query(strings.Join(keys, " "))

	results := search.SqlSearch(conf, serverParams.db, query)

	response := utils.SearchResponse{
		Query:   query,
		Results: results,
	}

	sendJSON(serverParams.writer, response)
}

func sendError(writer io.Writer, msg string, err error) {
	_, ioErr := fmt.Fprintf(writer, "%s: %s\n", msg, err)
	checkIOError(ioErr)
}

// Handles an /add request, inserts a row to the database for each path given
func handleAdd(serverParams serverFuncParams, paths []string) {
	// TODO: Albins pr should impl this
	panic("Not implemented")
}

// Handles a /quit request, cleanly shutsdown the database container and server
func handleQuit(serverParams serverFuncParams) {
	_, err := fmt.Fprintf(serverParams.writer, "Shutting down\n")
	checkIOError(err)

	err = serverParams.db.Close()
	checkIOError(err)
	stopContainer()

	// This needs to be called as a goroutine because the handler needs to
	// return
	// before the server can shutdown
	go func() {
		err := serverParams.server.Shutdown(context.Background())
		checkIOError(err)
	}()
}
