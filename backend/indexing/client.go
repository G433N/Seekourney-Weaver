package indexing

import (
	"context"
	"errors"
	"flag"
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
	Port       utils.Port
	Name       string
	Parallel   bool
	ConfigPath *utils.Path
	channel    chan *UnnormalizedDocument
}

func NewClient(name string) *IndexerClient {
	portFlag := flag.Uint("port", 0,
		"Port to run the indexer on, in the range of 39000-39499")
	parrallelFlag := flag.Bool("par", false, "Run in parallel, Optional")
	configFlag := flag.String("conf", "", "Config file to use, Optional")

	flag.Parse()

	var port utils.Port

	if *portFlag != 0 {
		temp, ok := IntoPort(*portFlag)
		if !ok {
			log.Fatalf("Port %d is out of range", *portFlag)
		}
		port = temp
	} else {
		port = utils.MININDEXERPORT
	}

	parrallel := *parrallelFlag

	var configPath *utils.Path

	if *configFlag != "" {
		temp := utils.Path(*configFlag)
		configPath = &temp
	}

	channel := make(chan *UnnormalizedDocument, 100)

	client := &IndexerClient{
		Port:       port,
		Name:       name,
		Parallel:   parrallel,
		ConfigPath: configPath,
		channel:    channel,
	}

	go func() {
		for doc := range client.channel {
			if doc != nil {
				bytes := ResponseDocs([]UnnormalizedDocument{*doc})
				body := utils.BytesBody(bytes)
				port := utils.Port(8080)
				resp, err := utils.PostRequest(
					body,
					"http://localhost",
					port,
					"push",
					"docs",
				)
				if err != nil {
					client.Log("Error sending document: %s", err)
					return
				}
				log.Println("Response:", resp)
			}
		}
	}()

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
		utils.EnableCORS(&writer)
		switch html.EscapeString(request.URL.Path) {

		case "/index":
			settings, err := client.SettingsFromRequest(request)
			if err != nil {
				client.Log("Error getting settings from request: %s", err)
				log.Println("Error getting settings from request:", err)
			}

			switch settings.Type {
			case utils.FileSource:
				client.Log("Indexing file: %s", settings.Path)
			case utils.DirSource:
				client.Log("Indexing directory: %s", settings.Path)
			case utils.UrlSource:
				client.Log("Indexing URL: %s", settings.Path)
			}

			cxt := NewContext(client)

			f(cxt, settings)

		case "/name":
			_, err := fmt.Fprintf(writer, "%s\n", string(client.Name))
			utils.PanicOnError(err)
		case "/ping":
			_, err := fmt.Fprintf(writer, "%s", string(ResponsePing()))
			client.Log("Responded to ping request from Core")
			utils.PanicOnError(err)
		case "/shutdown":
			_, err := fmt.Fprintf(writer, "%s", string(ResponseExiting()))
			client.Log("Shutdown triggered by Core, shutting down indexer")
			utils.PanicOnError(err)
			// TODO: Shutdown gracefully, like in #85
			// Currently this never sends a response
			os.Exit(0)
		default:
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

	log.Printf(msg, args...)

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
