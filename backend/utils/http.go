package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
)

// EnableCORS sets Cross-origin resource sharing on for a ResponseWriter.
func EnableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// RequestBodyBytes reads the request body and returns it as a byte slice.
func RequestBodyBytes(body *http.Request) ([]byte, error) {
	if body == nil {
		return nil, errors.New("body is nil")
	}
	bytes, err := io.ReadAll(body.Body)
	if err != nil {
		return nil, errors.New("could not read body: " + err.Error())
	}

	return bytes, nil
}

// RequestBodyString reads the request body and returns it as a string.
func RequestBodyString(body *http.Request) (string, error) {
	if body == nil {
		return "", errors.New("body is nil")
	}
	bytes, err := RequestBodyBytes(body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// RequestBodyJson unmarshals the request body into a struct.
func RequestBodyJson[T any](
	body *http.Request,
) (T, error) {

	var data T

	b, err := RequestBodyBytes(body)
	if err != nil {
		return data, err
	}

	var buffer bytes.Buffer
	err = json.Indent(&buffer, b, "", "  ")
	if err != nil {
		log.Printf("Error indenting JSON: %v", err)
	} else {
		name := reflect.TypeOf(data)
		log.Printf("%s: %s", name, buffer.String())
	}

	err = json.Unmarshal(b, &data)
	return data, err
}

type HttpBody struct {
	body []byte
}

func EmptyBody() *HttpBody {
	return &HttpBody{
		body: []byte{},
	}
}

func JsonBody[T any](body T) *HttpBody {
	bytes, err := json.Marshal(body)
	PanicOnError(err)
	return &HttpBody{
		body: bytes,
	}
}

func StrBody(body string) *HttpBody {
	bytes := []byte(body)
	return &HttpBody{
		body: bytes,
	}
}

func BytesBody(body []byte) *HttpBody {
	return &HttpBody{
		body: body,
	}
}

func intoReader(body *HttpBody) io.ReadCloser {
	if body == nil {
		return nil
	}
	return io.NopCloser(bytes.NewReader(body.body))
}

func respIntoBytes(resp *http.Response) ([]byte, error) {

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("indexer did not respond to request, " +
			"alternatively did not respond with ok statuscode")
	}
	defer func() {
		err := resp.Body.Close()
		PanicOnError(err)
	}()

	respByte, err := io.ReadAll(resp.Body)
	return respByte, err
}

func GetRequestBytes(
	host string,
	port Port,
	urlPath ...string,
) ([]byte, error) {

	url := host + ":" + port.String() + "/" + strings.Join(urlPath, "/")

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New(
			"indexer did not respond to request: " + err.Error(),
		)
	}
	return respIntoBytes(resp)
}

func GetRequestJSON[T any](
	host string,
	port Port,
	urlPath ...string,
) (T, error) {

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

// PostRequestBytes sends a POST request to the indexer and returns the response
// as bytes.
func PostRequestBytes(
	body *HttpBody,
	host string,
	port Port,
	urlPath ...string,
) ([]byte, error) {
	// TODO: Implement PostRequestBytes

	url := host + ":" + port.String() + "/" + strings.Join(urlPath, "/")
	req, err := http.NewRequest("POST", url, intoReader(body))
	if err != nil {
		return nil, errors.New(
			"indexer did not respond to request: " + err.Error(),
		)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, errors.New(
			"indexer did not respond to request: " + err.Error(),
		)
	}

	return respIntoBytes(resp)
}

// PostRequestJSON sends a POST request to the indexer and returns the response
// as a JSON object.
func PostRequestJSON[T any](
	body *HttpBody,
	host string,
	port Port,
	urlPath ...string,
) (T, error) {
	var respData T

	respByte, err := PostRequestBytes(body, host, port, urlPath...)

	if err != nil {
		return respData, err
	}

	err = json.Unmarshal(respByte, &respData)
	return respData, err
}

// PostRequest sends a POST request to the indexer and returns the response as a
// string.
func PostRequest(
	body *HttpBody,
	host string,
	port Port,
	urlPath ...string,
) (string, error) {
	respByte, err := PostRequestBytes(body, host, port, urlPath...)
	if err != nil {
		return "", err
	}

	respString := string(respByte)
	return respString, nil
}
