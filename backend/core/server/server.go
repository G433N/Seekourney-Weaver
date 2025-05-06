package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"seekourney/core/config"
	"seekourney/core/database"
	"seekourney/core/document"
	"seekourney/core/folder"
	"seekourney/core/indexAPI"
	"seekourney/core/search"
	"seekourney/utils"
	"strings"
)

const (
	_SERVERADDRESS_       string     = ":8080"
	_CONTAINERSTART_      string     = "./docker-start"
	_CONTAINEROUTPUTFILE_ string     = "./docker.log"
	_HOST_                string     = "localhost"
	_DBPORT_              utils.Port = 5433
	_CONTAINERNAME_       string     = "go-postgres"
	_USER_                string     = "go-postgres"
	_PASSWORD_            string     = "go-postgres"
	_DBNAME_              string     = "go-postgres"
	_EMPTYJSON_           JSONString = "{}"
)

// HTTP requests.
const (
	_ALL_    string = "/all"
	_SEARCH_ string = "/search"
	_ADD_    string = "/add"
	_QUIT_   string = "/quit"
)

// serverFuncParams is used by server query handler functions.
type serverFuncParams struct {
	writer io.Writer
	db     *sql.DB
	stop   context.CancelFunc
}

// startContainer start the database container using
// the command defined in _CONTAINERSTART_.
// Blocks until the container is closed.
func startContainer() {
	container := exec.Command("/bin/sh", _CONTAINERSTART_)

	outfile, err := os.Create(_CONTAINEROUTPUTFILE_)
	utils.PanicOnError(err)
	container.Stdout = outfile
	container.Stderr = outfile

	err = container.Run()
	utils.PanicOnError(err)
	err = outfile.Close()
	utils.PanicOnError(err)
}

// stopContainer signals the database container to stop,
// and will finish the command started by startContainer().
func stopContainer() {
	err := exec.Command("docker", "stop", "--signal", "SIGTERM",
		_CONTAINERNAME_).Run()

	if err != nil {
		panic(fmt.Sprintf("Error stopping container: %s\n", err))
	}
}

// conf holds the config object for the server.
// Gets initialized in the Run function.
var conf *config.Config

// index loads the local file config and creates a folder object.
// TODO: This function is temporay
// func index() folder.Folder {
// 	// Load local file config
// 	localConfig := localtext.Load(conf)
//
// 	// TODO: Later when documents comes over the network, we can still use the
// 	// same code. since it is an iterator
// 	folder := folder.FromIter(
// 		conf.Normalizer,
// 		localConfig.IndexDir("test_data"),
// 	)
//
// 	rm := folder.ReverseMappingLocal()
//
// 	files := folder.GetDocAmount()
// 	words := len(rm)
//
// 	log.Printf("Files: %d, Words: %d\n", files, words)
//
// 	if files == 0 {
// 		log.Fatal(
// 			"No files found, run make downloadTestFiles to download test files",
// 		)
// 	}
//
// 	return folder
// }

// insertFolder inserts all documents in the given folder into the database.
func insertFolder(db *sql.DB, folder *folder.Folder) {

	for _, doc := range folder.GetDocs() {

		_, err := database.InsertInto(db, doc)

		if err != nil {
			log.Printf("Error inserting row: %s\n", err)
		}
	}
}

func test() {
	// Test index registration
	cmd := "go run indexer/localtext/main.go indexer/localtext/localtext.go"
	_, err := indexAPI.RegisterIndexer(cmd)

	if err != nil {
		log.Printf("Error registering indexer: %s\n", err)
	}
}

/*
Run runs an http server with a postgres instance within docker container.
It can be accessed for example by `curl 'http://localhost:8080/search?q=key1'`
or using the client package: `go run . client <command>`.

The server accepts the following paths as commands:

/all - Lists all paths in database, probably won't be used in production but
helpful for tests.

/search - Query database, will return all paths containing given keywords.
Keywords are sent using http query under the key 'q'.

/add - adds one or several paths to the database, paths are sent using http
query under the key 'p'.

/quit - Shuts down the server.
*/
func Run(args []string) {

	// Load config
	conf = config.Load()

	go startContainer()

	db := connectToDB()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	server := &http.Server{
		Addr:        _SERVERADDRESS_,
		BaseContext: func(l net.Listener) context.Context { return ctx },
	}

	amount, err := database.RowAmount(db, utils.TABLEDOCUMENT)

	if err == nil {
		log.Printf("Row amount: %d\n", amount)
	} else {
		log.Printf("Error getting row amount: %s\n", err)
	}

	log.Println("Server started at", _SERVERADDRESS_)

	queryHandler := func(writer http.ResponseWriter, request *http.Request) {
		enableCORS(&writer)
		serverParams := serverFuncParams{writer: writer, db: db}

		switch html.EscapeString(request.URL.Path) {
		case _ALL_:
			handleAll(serverParams)
		case _SEARCH_:
			handleSearchSQL(serverParams, request.URL.Query()["q"])
		case _ADD_:
			handleAdd(serverParams, request.URL.Query()["p"])
		case _QUIT_:
			handleQuit(serverParams)
		case "/log":
			msg := request.URL.Query().Get("msg")
			log.Printf("Log: %s\n", msg)
		default:
			log.Println("Unknown path:", request.URL)
		}
	}

	http.HandleFunc("/", queryHandler)

	go test()
	go func() {
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Server encountered an error:", err)
		}
		stop()
	}()

	// Wait until server is finished
	<-ctx.Done()

	fmt.Println("Shutting down")
	err = server.Shutdown(context.Background())
	if err != nil {
		fmt.Println("Error while shutting down server: ", err)
	}
	err = db.Close()
	if err != nil {
		fmt.Println("Error while closing database: ", err)
	}

	stopContainer()
}

// enableCORS sets Cross-origin resource sharing on for a ResponseWriter.
func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// recoverSQLError calls recover and writes a message to writer
// if an SQL function panic'd.
func recoverSQLError(writer io.Writer) {
	if err := recover(); err != nil {
		_, ioErr := fmt.Fprintf(writer, "SQL failed: %s\n", err)
		utils.PanicOnError(ioErr)
	}
}

// sendError writes msg and err to the writer.
func sendError(writer io.Writer, msg string, err error) {
	_, ioErr := fmt.Fprintf(writer, "%s: %s\n", msg, err)
	utils.PanicOnError(ioErr)
}

// sendJSON marshals the given data to JSON and writes it to the writer.
func sendJSON(writer io.Writer, data any) {

	jsonData, err := json.Marshal(data)
	if err != nil {
		sendError(writer, "JSON failed", err)
		return
	}

	_, err = fmt.Fprintf(writer, "%s\n", jsonData)
	if err != nil {
		sendError(writer, "IO failed", err)
	}
}

// handleAll handles an /all request,
// by querying all rows in database and writing output to response writer.
func handleAll(serverParams serverFuncParams) {
	defer recoverSQLError(serverParams.writer)

	var doc document.Document
	query := database.Select().
		Queries(doc.SQLGetFields()...).
		From(utils.TABLEDOCUMENT)

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

// handleSearchSQL handles a /search request.
func handleSearchSQL(serverParams serverFuncParams, keys []string) {
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

// handleAdd handles an /add request,
// by inserting a row to the database for each given path.
func handleAdd(serverParams serverFuncParams, paths []string) {
	// TODO: Albins pr should impl this
	panic("Not implemented")
}

// handleQuit handles a /quit request by initiating the shutdown process
// by cancelling the server context.
func handleQuit(serverParams serverFuncParams) {
	_, err := fmt.Fprintf(serverParams.writer, "Shutting down\n")
	utils.PanicOnError(err)

	serverParams.stop()
}
