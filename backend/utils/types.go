package utils

// Query is a string containing plain words separated by spaces.
// E.g. "1.24.2 golang documentation".
type Query string

// Word is a plain text word.
type Word string

// Path is a file path or web address.
type Path string

// Frequency is the number of occurrences that a Word appears.
type Frequency int

// Score is a number representing how relevant a given Word is when searching.
type Score float64

// Source denotes the type of source indexed.
// E.g. a local file or a web page.
type Source int

const (
	// Source is the source of the document
	// SourceLocal is a local file
	SourceLocal Source = iota
	// SourceWeb is a web page
	SourceWeb
)

// FrequencyMap gives the frequency of a given word.
type FrequencyMap map[Word]Frequency

// ScoreMap gives the relevance score of a given document
// with respect to a current search query.
type ScoreMap map[Path]Score

// ReverseMap gives all paths to files that contain the given word.
type ReverseMap map[Word][]Path

// WordFrequencyMap gives the frequency of a specific word
// for every document that contains it.
// Used when searching.
type WordFrequencyMap map[Path]Frequency

// SearchResult is information about a single document
// with respect to a current search query.
type SearchResult struct {
	Path   Path
	Score  Score
	Source Source
}

// SearchResponse is the format an HTTP search response
// from Core has after unmarshalling JSON.
type SearchResponse struct {
	Query   Query
	Results []SearchResult
}

// Result is a tuple used when handling database data.
type Result[T any] struct {
	Value T
	Err   error
}

// FileType-s are lower-case letters without dot, e.g. "html" or "md".
type FileType string

// Port is a port number used in http request-responses to-from
// indexers, Core, and DB.
// Value for indexing API must be within range
// [MININDEXERPORT, MAXINDEXERPORT].
type Port uint

// Address including port acting as endpoint for http request.
// E.g. "http://localhost:39010".
type Endpoint string
