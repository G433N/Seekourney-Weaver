package indexing

import "encoding/json"

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

const (
	// Values used in status field in response.
	STATUSSUCCESSFUL string = "success"
	STATUSFAILURE    string = "fail"
	// Values used in message field in response.
	MESSAGEPONG    string = "pong"
	MESSAGEEXITING string = "exiting"
)

// ResponseSuccess creates an indexer response denoting success in JSON format.
// message is optional, empty string is accepted.
func ResponseSuccess(message string) []byte {
	jsonData, err := json.Marshal(IndexerResponse{
		Status: STATUSSUCCESSFUL,
		Data:   ResponseData{Message: message},
	})

	if err != nil {
		panic("indexing ResponseSuccess could not marshal response")
	}

	return jsonData
}

// ResponseSuccess creates an indexer response denoting failure in JSON format.
func ResponseFail(message string) []byte {
	jsonData, err := json.Marshal(IndexerResponse{
		Status: STATUSFAILURE,
		Data:   ResponseData{Message: message},
	})

	if err != nil {
		panic("indexing ResponseFail could not marshal response")
	}

	return jsonData
}

// ResponseSuccess creates an indexer pong response,
// used when starting up indexer, in JSON format.
func ResponsePong() []byte {
	jsonData, err := json.Marshal(IndexerResponse{
		Status: STATUSSUCCESSFUL,
		Data:   ResponseData{Message: MESSAGEPONG},
	})

	if err != nil {
		panic("indexing ResponsePong could not marshal response")
	}

	return jsonData
}

// ResponseSuccess creates an indexer exiting response,
// used when exiting indexer from core, in JSON format.
func ResponseExiting() []byte {
	jsonData, err := json.Marshal(IndexerResponse{
		Status: STATUSSUCCESSFUL,
		Data:   ResponseData{Message: MESSAGEEXITING},
	})

	if err != nil {
		panic("indexing ResponseExiting could not marshal response")
	}

	return jsonData
}
