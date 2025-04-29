package indexing

// See indexing_API for documentation.

// ResponseData is "data" value in indexerResponse.
type ResponseData struct {
	Message   string                 `json:"message"`
	Documents []UnnormalizedDocument `json:"documents"`
}

// IndexerResponse is the standard format for responses from indexer.
type IndexerResponse struct {
	Status string       `json:"status"`
	Data   ResponseData `json:"data"`
}
