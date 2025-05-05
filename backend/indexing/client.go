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

		case "/ping":
			_, err := fmt.Fprintf(writer, "%s", string(ResponsePong()))
			client.Log("Responded to ping request from Core")
			utils.PanicOnError(err)
		case "/shutdown":
			_, err := fmt.Fprintf(writer, "%s", string(ResponseExiting()))
			client.Log("Shutdown triggered by Core, shutting down indexer")
			utils.PanicOnError(err)
			os.Exit(0)
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
