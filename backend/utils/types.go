package utils

type Word string
type Path string
type Frequency int
type Score int

type Source int

const (
	// Source is the source of the document
	// SourceLocal is a local file
	SourceLocal Source = iota
	// SourceWeb is a web page
	SourceWeb
)

type FrequencyMap map[Word]Frequency
type ScoreMap map[Path]Score
type ReverseMap map[Word][]Path

type SearchResult struct {
	Path   Path
	Score  Score
	Source Source
}

type SearchResponse struct {
	Query   string
	Results []SearchResult
}
