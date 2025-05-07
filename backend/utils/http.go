package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

func GetRequestBytes(host string, port Port, urlPath ...string) ([]byte, error) {

	url := host + ":" + port.String() + "/" + strings.Join(urlPath, "/")

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New(
			"indexer did not respond to request: " + err.Error(),
		)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("indexer did not respond to request, " +
			"alternatively did not respond with ok statuscode")
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
	return respByte, err
}

func GetRequestJSON[T any](host string, port Port, urlPath ...string) (T, error) {

	var respData T
	respByte, err := GetRequestBytes(host, port, urlPath...)
	if err != nil {
		return respData, err
	}

	err = json.Unmarshal(respByte, &respData)
	return respData, err
}

func GetRequest(host string, port Port, urlPath ...string) (string, error) {
	respByte, err := GetRequestBytes(host, port, urlPath...)
	if err != nil {
		return "", err
	}

	respString := string(respByte)
	return respString, nil
}

// PostRequestBytes sends a POST request to the indexer and returns the response as bytes.
func PostRequestBytes(host string, port Port, urlPath ...string) ([]byte, error) {
	// TODO: Implement PostRequestBytes
	return nil, nil
}

// PostRequestJSON sends a POST request to the indexer and returns the response as a JSON object.
func PostRequestJSON[T any](host string, port Port, urlPath ...string) (T, error) {
	// TODO: Implement PostRequestJSON
	var respData T
	return respData, nil
}

// PostRequest sends a POST request to the indexer and returns the response as a string.
func PostRequest(host string, port Port, urlPath ...string) (string, error) {
	respByte, err := PostRequestBytes(host, port, urlPath...)
	if err != nil {
		return "", err
	}

	respString := string(respByte)
	return respString, nil
}
