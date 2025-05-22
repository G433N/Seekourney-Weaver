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
	"seekourney/core/indexAPI"
	"seekourney/core/search"
	"seekourney/indexing"
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
	_ALL_             string = "/all"
	_ALL_INDEXERS_    string = "/all/indexers"
	_ALL_COLLECTIONS_ string = "/all/collections"
	_SEARCH_          string = "/search"
	_QUIT_            string = "/quit"
	_PUSHPATHS_       string = "/push/paths"
	_PUSHDOCS_        string = "/push/docs"
	_INDEX_           string = "/index"
	_PUSHCOLLECTION_  string = "/push/collection"
	_PUSHINDEXER_     string = "/push/indexer"
	_LOG_             string = "/log"
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

	defer func() {
		// TODO: Do something similar in IndexHandler for every indexer

		if recover := recover(); recover != nil {
			// TODO: Do we want to dot this, before starting the container?
			err := exec.Command("docker", "kill", "go-postgres").Run()
			utils.PanicOnError(err)
			log.Fatalf(
				"Error starting container: %s\nPlease, start the server again",
				recover,
			)
		}
	}()

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

	// Indexhandler is used to manage running indexers.
	indexHandler := indexAPI.NewIndexHandler()

	queryHandler := func(writer http.ResponseWriter, request *http.Request) {
		utils.EnableCORS(&writer)
		serverParams := serverFuncParams{writer: writer, db: db}

		switch html.EscapeString(request.URL.Path) {
		// TODO figure out where we can do goroutines
		// we cant use resp writer concurrently nor read resp body
		case _ALL_:
			handleAll(serverParams)
		case _ALL_INDEXERS_:
			handleAllIndexers(serverParams)
		case _ALL_COLLECTIONS_:
			handleAllCollections(serverParams)
		case _SEARCH_:
			handleSearchSQL(serverParams, request.URL.Query()["q"])
		case _PUSHPATHS_:
			handlePushPaths(serverParams, request.URL.Query()["p"])
		case _PUSHDOCS_:
			handlePushDocs(serverParams, request)
		case _INDEX_:
			handleIndex(serverParams)
		case _PUSHCOLLECTION_:
			handlePushCollection(serverParams, request, &indexHandler)
		case _PUSHINDEXER_:
			handlePushIndexer(serverParams, request)
		case _QUIT_:
			handleQuit(serverParams)
		case _LOG_:
			msg := request.URL.Query().Get("msg")
			log.Printf("Log: %s\n", msg)
		default:
			log.Println("Unknown path:", request.URL)
		}
	}

	http.HandleFunc("/", queryHandler)

	go func() {
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Server encountered an error:", err)
		}
		indexHandler.ForceShutdownAll()
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

// respondWithSuccess sends an indexing success response through writer.
func respondWithSuccess(writer io.Writer) {
	_, err := fmt.Fprintf(
		writer,
		"%s",
		string(indexing.ResponseSuccess("handling request to Core")),
	)
	utils.PanicOnError(err)
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

// handleAllIndexers handles an /all/indexers request,
// by querying all indexers in database and writing output to response writer.
func handleAllIndexers(serverParams serverFuncParams) {

	query := database.Select().
		QueryAll().
		From("indexer")

	insert := func(docs *[]indexAPI.IndexerData, ind indexAPI.IndexerData) {
		*docs = append(*docs, ind)
	}

	indexers := make([]indexAPI.IndexerData, 0)

	err := database.ExecScan(serverParams.db, string(query), &indexers, insert)
	if err != nil {
		sendError(serverParams.writer, "SQL failed", err)
		return
	}

	sendJSON(serverParams.writer, indexers)
}

// handleAllCollections handles an /all/collections request,
// by querying all collections in database and writing output to response writer.
func handleAllCollections(serverParams serverFuncParams) {

	query := database.Select().
		QueryAll().
		From("collection")

	insert := func(docs *[]indexAPI.Collection, col indexAPI.Collection) {
		*docs = append(*docs, col)
	}

	collections := make([]indexAPI.Collection, 0)

	err := database.ExecScan(serverParams.db, string(query), &collections, insert)
	if err != nil {
		sendError(serverParams.writer, "SQL failed", err)
		return
	}

	sendJSON(serverParams.writer, collections)
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

	for _, result := range results {
		log.Print("Result: ", result.Path, " Score: ", result.Score)
	}

	response := utils.SearchResponse{
		Query:   query,
		Results: results,
	}
	sendJSON(serverParams.writer, response)
}

// handlePushPaths handles an /push/path request,
// by inserting a row to the database for each given path.
func handlePushPaths(serverParams serverFuncParams, paths []string) {
	panic("Not implemented")
}

// handlePushDocs handles a /push/docs request,
// by normalizing documents send in request and adding them to db.
func handlePushDocs(serverParams serverFuncParams, request *http.Request) {
	respondWithSuccess(serverParams.writer)

	body, err := io.ReadAll(request.Body)
	utils.PanicOnError(err)

	resp := indexing.IndexerResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Print("Main server failed to parse PushDocs request" +
			" from indexer with error: " + err.Error())
		return
	}

	if resp.Status != indexing.STATUSSUCCESSFUL {
		log.Print("indexing request failed (messaged with PushDocs request)" +
			" with message: " + resp.Data.Message)
		return
	}

	if len(resp.Data.Documents) == 0 {
		log.Print("indexer indexed path and produced zero documents " +
			"(pushdocs request)")
		return
	}

	// Create goroutine for normalising and inserting into db,
	// as it might take significant time.
	go func() {
		for _, rawDoc := range resp.Data.Documents {
			normalizedDoc := document.Normalize(rawDoc, conf.Normalizer)

			// TODO fix
			// Error inserting row: pq: duplicate key value violates
			// unique constraint "document_pkey"

			exists, err := document.DocumentExsitsDB(
				serverParams.db,
				normalizedDoc.Path,
			)

			if err != nil {
				log.Printf("Error checking if document exists: %s\n", err)
				continue
			}

			if exists {
				err := normalizedDoc.UpdateDB(serverParams.db)
				if err != nil {
					log.Printf("Error updating document: %s\n", err)
				}
				continue
			}

			_, err = database.InsertInto(serverParams.db, normalizedDoc)

			if err != nil {
				log.Printf("Error inserting row: %s\n", err)
			}

			log.Print("Inserted document with path: ", normalizedDoc.Path)
		}
		log.Print("Handled pushdocs request successfully")
	}()
}

// handleIndex handles an /index request by dispatching an indexing request
// to the appropriate indexer.
func handleIndex(serverParams serverFuncParams) {
	panic("Not implemented")
}

// handlePushIndexer handles a /push/indexer request from frontend client
// by generating a new Indexer, storing it, and starting the indexer.
func handlePushIndexer(
	serverParams serverFuncParams,
	request *http.Request,
) {
	// TODO fail or success response after unmarshall?
	respondWithSuccess(serverParams.writer)

	startupCMD, err := utils.RequestBodyString(request)

	if err != nil {
		log.Print("Main server failed to parse PushIndexer request" +
			" from indexer with error: " + err.Error())
		return
	}

	_, err = indexAPI.RegisterIndexer(serverParams.db, startupCMD)

	if err != nil {
		log.Print(
			"Main server failed to register indexer with error: " + err.Error(),
		)
		return
	}

}

// handlePushCollection handles a /push/collection request from frontend client
// by generating a new Collection, storing it, and indexing its associated path.
func handlePushCollection(
	serverParams serverFuncParams,
	request *http.Request,
	indexers *indexAPI.IndexHandler,
) {
	// TODO fail or success response after unmarshall?
	respondWithSuccess(serverParams.writer)

	body, err := io.ReadAll(request.Body)
	utils.PanicOnError(err)

	// TODO change to correct reponse json format
	unreg := indexAPI.UnregisteredCollection{}
	err = json.Unmarshal(body, &unreg)
	if err != nil {
		log.Print("Main server failed to parse PushDocs request" +
			" from indexer with error: " + err.Error())
		return
	}

	go func() {
		collection, err := indexAPI.RegisterCollection(serverParams.db, unreg)
		if err != nil {
			panic("TODO")
		}

		// Dispatch may startup indexer and add to handler
		// if it is not already running.
		indexers.Mutex.Lock()
		errs := indexers.DispatchFromCollection(serverParams.db, collection)
		indexers.Mutex.Unlock()

		logDispatchErrors(errs)
	}()
}

// handleQuit handles a /quit request by initiating the shutdown process
// by cancelling the server context.
func handleQuit(serverParams serverFuncParams) {
	_, err := fmt.Fprintf(serverParams.writer, "Shutting down\n")
	utils.PanicOnError(err)

	serverParams.stop()
}
