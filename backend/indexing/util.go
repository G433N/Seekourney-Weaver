package indexing

import (
	"context"
	"errors"
	"fmt"
	"html"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"seekourney/utils"
	"strconv"
)

type SourceType int

const (
	FileSource SourceType = iota
	DirSource
	UrlSource
)

func StrToSourceType(str string) (SourceType, error) {
	switch str {
	case "file":
		return FileSource, nil
	case "dir":
		return DirSource, nil
	case "url":
		return UrlSource, nil
	default:
		return 0, errors.New("invalid source type")
	}
}

func SourceTypeToStr(t SourceType) string {
	switch t {
	case FileSource:
		return "file"
	case DirSource:
		return "dir"
	case UrlSource:
		return "url"
	default:
		return "unknown"
	}
}

func SourceTypeFromPath(path utils.Path) (SourceType, error) {

	stat, err := os.Stat(string(path))
	if err != nil {
		return 0, err
	}

	if stat.IsDir() {
		return DirSource, nil
	}

	if stat.Mode().IsRegular() {
		return FileSource, nil
	}

	return 0, errors.New("unknown source type")
}

type IndexerClient struct {
	Port utils.Port
	Name string
}

func NewClient(name string) *IndexerClient {

	client := &IndexerClient{
		Port: utils.Port(0),
		Name: name,
	}

	args := os.Args
	port, err := GetPort(args)
	if err != nil {

		port = utils.MININDEXERPORT
		client.Log("Error getting port: %s Using default %d", err, port)
	}

	client.Port = port

	client.Log("Client initialized")
	return client

}

func (client *IndexerClient) Start(f func(cxt Context, settings Settings)) {

	client.Log("Starting server on port %d", client.Port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	server := &http.Server{
		Addr:        ":" + strconv.Itoa(int(client.Port)),
		BaseContext: func(l net.Listener) context.Context { return ctx },
	}

	queryHandler := func(writer http.ResponseWriter, request *http.Request) {

		switch html.EscapeString(request.URL.Path) {

		case "/index":
			settings, err := client.SettingsFromRequest(request)
			if err != nil {
				log.Println("Error getting settings from request:", err)
			}

			switch settings.Type {
			case FileSource:
				client.Log("Indexing file: %s", settings.Path)
			case DirSource:
				client.Log("Indexing directory: %s", settings.Path)
			case UrlSource:
				client.Log("Indexing URL: %s", settings.Path)
			}

			cxt := Context{
				client: client,
			}

			f(cxt, settings)

		default:
			log.Println("Unknown path:", request.URL)
			client.Log("Unknown path: %s", request.URL)

		}
	}

	http.HandleFunc("/", queryHandler)

	go func() {
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Println("Server encountered an error:", err)
		}
		stop()
	}()

	// Wait until server is finished
	<-ctx.Done()

	log.Println("Shutting down")
	err := server.Shutdown(context.Background())
	if err != nil {
		log.Println("Error while shutting down server: ", err)
	}
}

func (client *IndexerClient) Log(msg string, args ...any) {

	port := strconv.Itoa(int(client.Port))
	name := client.Name

	base := "Indexer " + name + " on port " + port + ": "
	message := fmt.Sprintf(base+msg, args...)

	url := "http://localhost:8080/log?msg=" + url.QueryEscape(message)

	res, err := http.Get(url)

	if err != nil {
		log.Printf("Error sending log: %s", err)
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("Error bad response: %s", res.Status)
		return
	}

}

type Context struct {
	client *IndexerClient
}

func (cxt *Context) Log(msg string, args ...any) {
	cxt.client.Log(msg, args...)
}

func (cxt *Context) Send(doc UnnormalizedDocument) {
	cxt.client.Log("Sending document: %s", doc.Path)
}

type Settings struct {
	Path      utils.Path
	Type      SourceType
	Recursive bool
	Parrallel bool
}

func (client *IndexerClient) SettingsFromRequest(request *http.Request) (Settings, error) {
	path := request.URL.Query().Get("path")
	t := request.URL.Query().Get("type")
	recursive := request.URL.Query().Get("recursive")
	parallel := request.URL.Query().Get("parallel")

	sourceType, err := StrToSourceType(t)

	if err != nil {
		client.Log("Error converting source type: %s", err)
		return Settings{}, err
	}

	recursiveBool, err := strconv.ParseBool(recursive)
	if err != nil {
		client.Log("Error converting recursive: %s", err)
		recursiveBool = false
	}

	parallelBool, err := strconv.ParseBool(parallel)
	if err != nil {
		client.Log("Error converting parallel: %s", err)
		parallelBool = false
	}

	return Settings{
		Path:      utils.Path(path),
		Type:      sourceType,
		Recursive: recursiveBool,
		Parrallel: parallelBool,
	}, nil

}

func (settings *Settings) IntoURL(port utils.Port) (string, error) {

	path := string(settings.Path)
	sourceType := SourceTypeToStr(settings.Type)
	recursive := strconv.FormatBool(settings.Recursive)
	parallel := strconv.FormatBool(settings.Parrallel)

	query := fmt.Sprintf("?path=%s&type=%s&recursive=%s&parallel=%s",
		path, sourceType, recursive, parallel)

	return fmt.Sprintf("http://localhost:%d/index%s", port, query), nil

}

// IntoPort converts an integer to a port.
func IntoPort(integer uint) (utils.Port, bool) {

	if integer < uint(utils.MININDEXERPORT) ||
		integer > uint(utils.MAXINDEXERPORT) {
		return 0, false
	}

	return utils.Port(integer), true
}

// IsValidPort checks if port value is within designated range for indexer API.
func IsValidPort(port utils.Port) bool {
	_, ok := IntoPort(uint(port))
	return ok
}

// GetPort parses the port from command line arguments.
func GetPort(args []string) (utils.Port, error) {

	if len(args) < 2 {
		return 0, errors.New("to few arguments")
	}

	num, err := strconv.ParseInt(args[1], 10, 32)
	if err != nil {
		return 0, err
	}

	port, ok := IntoPort(uint(num))
	if !ok {
		return 0, errors.New("port out of range")
	}

	return utils.Port(port), nil
}
