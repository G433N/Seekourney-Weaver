package indexing

import (
	"seekourney/utils"
	"strings"
)

type Context struct {
	client   *IndexerClient
	metadata *docMetadata
}

func NewContext(client *IndexerClient) Context {
	if client == nil {
		panic("client cannot be nil")
	}

	return Context{
		client:   client,
		metadata: nil,
	}
}

func (cxt *Context) Log(msg string, args ...any) {
	cxt.client.Log(msg, args...)
}

func (cxt *Context) hasMetadata() bool {
	return cxt.metadata != nil
}

func (cxt *Context) StartDoc(path utils.Path, source utils.Source) {
	if cxt.hasMetadata() {
		panic("Document already started")
	}

	cxt.metadata = &docMetadata{
		path:   path,
		source: source,
		text:   make([]string, 0),
	}
}

func (cxt *Context) AddText(text string) {
	if !cxt.hasMetadata() {
		panic("Use StartDoc before AddText")
	}
	cxt.metadata.AddText(text)
}

// Done is called when the document is finished.
// It can take a function to modify the document before sending it.
// This function needs to be thread safe.
// When indexed the document is sent to the server.
func (cxt *Context) Done(f *func(*UnnormalizedDocument)) {
	if !cxt.hasMetadata() {
		panic("Use StartDoc before Done")
	}

	if cxt.client.Parallel {
		go index(cxt.client, *cxt.metadata, f)
	} else {
		index(cxt.client, *cxt.metadata, f)
	}

	cxt.metadata = nil
}

func index(
	client *IndexerClient,
	metadata docMetadata,
	f *func(*UnnormalizedDocument),
) {

	doc, err := metadata.index()
	if err != nil {
		client.Log("Error indexing document: %s", err)
		return
	}

	if f != nil {
		(*f)(doc)
	}
	client.channel <- doc
}

func (cxt *Context) send(doc *UnnormalizedDocument) {
	cxt.Log("Sending document: %s", doc.Path)
}

// TODO: Better name
type docMetadata struct {
	path   utils.Path
	source utils.Source
	text   []string
}

func (metadata *docMetadata) AddText(text string) {
	metadata.text = append(metadata.text, text)
}

// Index takes a Settings struct and returns an UnnormalizedDocument.
// It converts the source type into a source and then creates a document
// from the text.
// Return nil if there was an error converting the source type.
// The text is cleared after the document is created.
func (metadata *docMetadata) index() (*UnnormalizedDocument, error) {

	text := strings.Join(metadata.text, "\n")

	doc := DocFromText(
		metadata.path,
		metadata.source,
		text,
	)

	return &doc, nil

}
