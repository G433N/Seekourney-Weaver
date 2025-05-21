package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"seekourney/indexing"
	"seekourney/tui/format"
	"seekourney/utils"
	"seekourney/utils/timing"
)

// Strings for building URLs in HTTP requests
const (
	_COREENDPOINT_ utils.Endpoint = "http://localhost:8080"
	_HOST_         string         = "http://localhost"
	_PORT_         utils.Port     = 8080
	_SEARCH_       string         = "/search?"
	_PUSHPATHS_    string         = "/push/paths?"
	_PUSHDOCS_     string         = "/push/docs"
	_QUIT_         string         = "/quit"
	_ALL_          string         = "/all"
	_SEARCHKEY_    string         = "q"
	_ADDKEY_       string         = "p"
)

// argumentError prints a usage string and terminates client process.
func argumentError() {
	fmt.Println("usage: client <command> [<args>]")
	fmt.Println()
	fmt.Println("available commands")
	fmt.Println("  all                  request all pages in database")
	fmt.Println("  search    [key ...]  request all pages containing keys")
	fmt.Println("  pushpaths [path ...] add paths to database")
	fmt.Println("  pushdocs             add 2 test documents to database")
	fmt.Println("  index     [path ...] test indexing of paths")
	fmt.Println("  quit                 request the server to shutdown")
	os.Exit(1)
}

// init runs before any other code in this package.
// https://golangdocs.com/init-function-in-golang
func init() {
	// Initialize the timing package
	timing.Init(timing.Default())
}

// main takes an array of arguments args as input and calls different functions
// to send a HTTP request to the server
// The first element in args is ignored.
// The second element in args is the command.
// The rest of the elements are for commands that require extra arguments.
// If no valid command is found in the first element, the program terminates.
// args are formatted like commandline arguments ("client", command, args)
func main() {

	var num float32 = 0.0
	log.Println("Num:", num)

	args := os.Args

	if len(args) < 2 {
		argumentError()
	}

	switch args[1] {
	case "search":
		searchForTerms(args[2:])
	case "pushpaths":
		pushPaths(args[2:])
	case "pushdocs":
		pushDocs()
	case "all":
		getAll()
	case "index":
		index(args[2:])
	case "test":
		test()
	case "quit":
		shutdownServer()
	default:
		argumentError()
	}
}

// printResponse prints a HTTP response to stdout.
func printResponse(response *http.Response) {
	bytes, _ := io.ReadAll(response.Body)
	fmt.Print(string(bytes))
}

// searchForTerms requests a search for given terms through Core,
// and prints the results.
// Handler for command /search.
func searchForTerms(terms []string) {
	sw := timing.Measure(timing.Search)
	defer sw.Stop()

	values := url.Values{}
	for _, term := range terms {
		values.Add(_SEARCHKEY_, term)
	}
	resp, err := http.Get(
		string(_COREENDPOINT_) + _SEARCH_ + values.Encode(),
	)
	utils.PanicOnError(err)

	bytes, _ := io.ReadAll(resp.Body)
	result := utils.SearchResponse{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		log.Println("Response:", string(bytes))
		return
	}

	format.PrintSearchResponse(result)
}

// pushPaths adds given paths to the database through Core,
// and prints the response.
// Handler for command /add.
func pushPaths(paths []string) {
	values := url.Values{}
	for _, term := range paths {
		values.Add(_ADDKEY_, term)
	}
	resp, err := http.Get(
		string(_COREENDPOINT_) + _PUSHPATHS_ + values.Encode(),
	)
	utils.PanicOnError(err)
	printResponse(resp)
}

type udoc = indexing.UnnormalizedDocument

// pushDocs acts as an indexer sending unnormalized documents to core
// to add to database.
// 2 test documents are sent.
func pushDocs() {
	testdocs := []udoc{
		{
			Path:   "a/test/document/path",
			Source: 0,
			Words:  utils.FrequencyMap{"good": 42, "bad": 11},
		},
		{
			Path:   "yet/another/path",
			Source: 0,
			Words:  utils.FrequencyMap{"yep": 100},
		}}

	bodyReader := bytes.NewReader(indexing.ResponseDocs(testdocs))
	req, err := http.NewRequest(
		http.MethodPost,
		string(_COREENDPOINT_)+_PUSHDOCS_,
		bodyReader,
	)
	utils.PanicOnError(err)

	resp, err := http.DefaultClient.Do(req)
	utils.PanicOnError(err)
	printResponse(resp)
}

// getAll fetches all paths stored in database through Core,
// and prints them.
// Handler for command /all.
func getAll() {
	resp, err := http.Get(string(_COREENDPOINT_) + _ALL_)
	utils.PanicOnError(err)
	printResponse(resp)
}

// shutdownServer remotely shuts down Core.
// Handler for command /quit.
func shutdownServer() {
	resp, err := http.Get(string(_COREENDPOINT_) + _QUIT_)
	utils.PanicOnError(err)
	printResponse(resp)
}

func index(paths []string) {

	d := "/home/oxygen/Projects/Seekourney-Weaver/backend/test_data/docs.gl/todo.md"
	paths = append(paths, d)
	log.Println("Indexing paths:", paths)

	for _, path := range paths {

		Type, err := indexing.SourceTypeFromPath(utils.Path(path))

		if err != nil {
			log.Println("Error getting source type:", err)
			return
		}

		settings := indexing.Settings{
			Path:         utils.Path(path),
			Type:         Type,
			CollectionID: "",
			Recursive:    true,
			Parrallel:    false,
		}

		test, err := json.MarshalIndent(settings, "", "  ")
		utils.PanicOnError(err)
		log.Printf("Settings: %s", string(test))

		body := utils.JsonBody(settings)
		log.Println("Request body:", *body)
		port := utils.MININDEXERPORT
		_, err = utils.PostRequest(body, _HOST_, port, "index")

		if err != nil {
			log.Println("Error sending request:", err)
			return
		}

		log.Printf("Sent %s successfully\n", path)
	}
}

type extractID struct {
	ID utils.IndexerID
}

func allIndexers() utils.IndexerID {
	res, err := utils.GetRequestBytes(_HOST_, _PORT_, "all", "indexers")

	if err != nil {
		log.Println("Error sending request:", err)
		return ""
	}

	var bytes bytes.Buffer
	err = json.Indent(&bytes, res, "", "  ")
	if err != nil {
		log.Println("Error indenting JSON:", err)
		return ""
	}

	log.Println("Response:", bytes.String())

	var indexers []extractID
	err = json.Unmarshal(res, &indexers)

	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return ""
	}

	return indexers[0].ID
}

func allCollections() {
	res, err := utils.GetRequestBytes(_HOST_, _PORT_, "all", "collections")

	if err != nil {
		log.Println("Error sending request:", err)
		return
	}

	var bytes bytes.Buffer
	err = json.Indent(&bytes, res, "", "  ")
	if err != nil {
		log.Println("Error indenting JSON:", err)
		return
	}

	log.Println("Response:", bytes.String())
}

func test() {

	body := utils.StrBody("go run indexer/localtext/main.go indexer/localtext/localtext.go")
	_, err := utils.PostRequest(body, _HOST_, _PORT_, "push", "indexer")
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}

	id := allIndexers()

	col := utils.UnregisteredCollection{
		Path:                "/home/carbon/Projects/go_indexer/backend/test_data/docs.gl/todo.md",
		IndexerID:           id,
		SourceType:          utils.FileSource,
		Recursive:           true,
		RespectLastModified: false,
		Normalfunc:          utils.Stemming,
	}

	body = utils.JsonBody(col)
	_, err = utils.PostRequest(body, _HOST_, _PORT_, "push", "collection")
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}

	allCollections()

}
