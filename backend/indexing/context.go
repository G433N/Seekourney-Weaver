package indexing

import (
	"seekourney/utils"
	"strings"
)

// Context is a struct that holds the client and provides methods
// to log messages and create documents. It is used in client.Start()
// this struct is here such that we can implment more feature later, without
// breaking older code.
type Context struct {
	client *IndexerClient
}

// NewContext creates a new Context struct.
func NewContext(client *IndexerClient) Context {
	if client == nil {
		panic("client cannot be nil")
	}

	return Context{
		client: client,
	}
}

// Log is a method that logs messages to the server.
func (cxt *Context) Log(msg string, args ...any) {
	cxt.client.Log(msg, args...)
}

// StartDoc creates a new document builder.
func (cxt *Context) StartDoc(
	path utils.Path,
	source utils.Source,
	settings Settings,
) *docBuilder {

	return &docBuilder{
		path:       path,
		source:     source,
		text:       make([]string, 0),
		collection: settings.CollectionID,
		cxt:        cxt,
	}
}

// index makes a docBuilder into a document and indexes it. After that
// the function f is called on the document, to do any modifications.
// The document is then sent to the server.
func index(
	client *IndexerClient,
	docBuilder docBuilder,
	f *func(*UnnormalizedDocument),
) {

	doc, err := docBuilder.index()
	if err != nil {
		client.Log("Error indexing document: %s", err)
		return
	}

	if f != nil {
		(*f)(doc)
	}
	client.channel <- doc
}

// DocBuilder represents a partally completed document
type docBuilder struct {
	path       utils.Path
	source     utils.Source
	collection CollectionID
	text       []string
	cxt        *Context
}

// AddText adds text to the document.
func (doc *docBuilder) AddText(text string) {
	doc.text = append(doc.text, text)
}

// Done is called when the document is finished.
// It can take a function to modify the document before sending it.
// This function needs to be thread safe.
// When indexed the document is sent to the server.
func (doc *docBuilder) Done(f *func(*UnnormalizedDocument)) {

	if doc.cxt.client.Parallel {
		go index(doc.cxt.client, *doc, f)
	} else {
		index(doc.cxt.client, *doc, f)
	}

	*doc = docBuilder{}
}

// Index takes a Settings struct and returns an UnnormalizedDocument.
// It converts the source type into a source and then creates a document
// from the text.
// Return nil if there was an error converting the source type.
// The text is cleared after the document is created.
func (docBuilder *docBuilder) index() (*UnnormalizedDocument, error) {

	text := strings.Join(docBuilder.text, "\n")

	doc := DocFromText(
		docBuilder.path,
		docBuilder.source,
		docBuilder.collection,
		text,
	)

	return &doc, nil

}
