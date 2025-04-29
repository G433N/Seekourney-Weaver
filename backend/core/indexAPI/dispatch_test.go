package indexAPI

import (
	"os/exec"
	"reflect"
	"seekourney/indexing"
	"seekourney/utils"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

const (
	_TESTURI_ restEndpoint = restEndpoint(_ENDPOINTPREFIX_ + "39100")
)

var testResponseFail IndexerResponse = IndexerResponse{
	Status: _STATUSFAILURE_,
	Data:   ResponseData{Message: "failed to server response"},
}

var testResponsePong IndexerResponse = IndexerResponse{
	Status: _STATUSSUCCESSFUL_,
	Data:   ResponseData{Message: _PONG_},
}

var testResponseExiting IndexerResponse = IndexerResponse{
	Status: _STATUSSUCCESSFUL_,
	Data:   ResponseData{Message: _EXITING_},
}

// waitOnTestCMD is used instead of shutdownIndexerGraceful for some tests.
// This is needed as shutdownIndexerGraceful will force kill if
// correct JSON is not returned, which won't clean up resources.
func waitOnTestCMD(info IndexerInfo) {
	if err := info.cmd.Wait(); err != nil {
		panic("wait on test command failed")
	}
}

func TestStartupPingFail(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	gock.New(string(_TESTURI_)).
		Get(_PING_).
		Reply(200).
		JSON(testResponseFail)
	info := IndexerInfo{
		name:             "TestIndexerName",
		cmd:              exec.Command("ls"),
		fileTypesHandled: []utils.FileType{"txt"},
		id:               42,
		endpoint:         _TESTURI_,
	}
	defer waitOnTestCMD(info)

	assert.Error(t, startupIndexer(info))
	assert.True(t, gock.IsDone())
}

func TestStartupPingSuccess(t *testing.T) {
	defer gock.Off()
	gock.New(string(_TESTURI_)).
		Get(_PING_).
		Reply(200).
		JSON(testResponsePong)
	info := IndexerInfo{
		name:             "TestIndexerName",
		cmd:              exec.Command("ls"),
		fileTypesHandled: []utils.FileType{"txt"},
		id:               42,
		endpoint:         _TESTURI_,
	}
	defer waitOnTestCMD(info)

	assert.NoError(t, startupIndexer(info))
	assert.True(t, gock.IsDone())
}

// Same as TestStartupPingFail but with worse JSON response.
func TestStartupInvalidJSON(t *testing.T) {
	defer gock.Off()
	gock.New(string(_TESTURI_)).
		Get(_PING_).
		Reply(200).
		JSON(map[string]string{"invalid": "JSON data send back"})
	info := IndexerInfo{
		name:             "TestIndexerName",
		cmd:              exec.Command("ls"),
		fileTypesHandled: []utils.FileType{"txt"},
		id:               42,
		endpoint:         _TESTURI_,
	}
	defer waitOnTestCMD(info)

	assert.Error(t, startupIndexer(info))
	assert.True(t, gock.IsDone())
}

func TestShutdownValidResponse(t *testing.T) {
	defer gock.Off()
	gock.New(string(_TESTURI_)).
		Get(_PING_).
		Reply(200).
		JSON(testResponsePong)
	info := IndexerInfo{
		name:             "TestIndexerName",
		cmd:              exec.Command("ls"),
		fileTypesHandled: []utils.FileType{"txt"},
		id:               42,
		endpoint:         _TESTURI_,
	}

	assert.NoError(t, startupIndexer(info))
	assert.True(t, gock.IsDone())

	gock.New(string(_TESTURI_)).
		Get(_SHUTDOWN_).
		Reply(200).
		JSON(testResponseExiting)

	assert.NoError(t, shutdownIndexerGraceful(info))
	assert.True(t, gock.IsDone())
}

func TestShutdownInvalidResponse(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping shutdown test to avoid small resource leaks")
	}

	defer gock.Off()
	gock.New(string(_TESTURI_)).
		Get(_PING_).
		Reply(200).
		JSON(testResponsePong)
	info := IndexerInfo{
		name:             "TestIndexerName",
		cmd:              exec.Command("ls"),
		fileTypesHandled: []utils.FileType{"txt"},
		id:               42,
		endpoint:         _TESTURI_,
	}

	assert.NoError(t, startupIndexer(info))
	assert.True(t, gock.IsDone())

	gock.New(string(_TESTURI_)).
		Get(_SHUTDOWN_).
		Reply(200).
		JSON(testResponseFail)

	assert.Error(t, shutdownIndexerGraceful(info))
	assert.True(t, gock.IsDone())
}

// Test data for mocking indexing requests.
const testIndexFolderPath1 utils.Path = "home/george/my_cool_text_files"

const testIndexFilePath1 utils.Path = "home/george/my_cool_text_files/first.txt"
const testIndexFilePath2 utils.Path = "home/george/my_cool_text_files/other.txt"

var testResponseDoc1 indexing.UnnormalizedDocument = indexing.UnnormalizedDocument{
	Path:   testIndexFilePath1,
	Source: utils.SourceLocal,
	Words: utils.FrequencyMap{
		"blue":   5,
		"black":  2,
		"red":    50,
		"green":  34,
		"orange": 11,
	},
}

var testResponseDoc2 indexing.UnnormalizedDocument = indexing.UnnormalizedDocument{
	Path:   testIndexFilePath2,
	Source: utils.SourceLocal,
	Words: utils.FrequencyMap{
		"wood":  234,
		"steel": 52,
	},
}

var testIndexingResponse1 IndexerResponse = IndexerResponse{
	Status: _STATUSSUCCESSFUL_,
	Data: ResponseData{
		Documents: []indexing.UnnormalizedDocument{testResponseDoc1},
	},
}
var testIndexingResponse2 IndexerResponse = IndexerResponse{
	Status: _STATUSSUCCESSFUL_,
	Data: ResponseData{
		Documents: []indexing.UnnormalizedDocument{testResponseDoc1, testResponseDoc2},
	},
}

func TestRequestIndexingSimple(t *testing.T) {
	defer gock.Off()
	gock.New(string(_TESTURI_)).
		Get(_INDEXFULL_ + "/" + string(testIndexFilePath1)).
		Reply(200).
		JSON(testIndexingResponse1)
	info := IndexerInfo{
		name:             "TestIndexerName",
		cmd:              exec.Command("ls"),
		fileTypesHandled: []utils.FileType{"txt"},
		id:               42,
		endpoint:         _TESTURI_,
	}

	docs, err := requestIndexing(info, testIndexFilePath1)
	assert.NoError(t, err)
	assert.True(t, gock.IsDone())

	assert.Equal(t, len(docs), 1)
	assert.Equal(t, docs[0].Path, testIndexFilePath1)
	assert.Equal(t, docs[0].Source, utils.SourceLocal)
	assert.True(t, reflect.DeepEqual(docs[0].Words, testResponseDoc1.Words))
}

func TestRequestIndexingTwo(t *testing.T) {
	defer gock.Off()
	gock.New(string(_TESTURI_)).
		Get(_INDEXFULL_ + "/" + string(testIndexFolderPath1)).
		Reply(200).
		JSON(testIndexingResponse2)
	info := IndexerInfo{
		name:             "TestIndexerName",
		cmd:              exec.Command("ls"),
		fileTypesHandled: []utils.FileType{"txt"},
		id:               42,
		endpoint:         _TESTURI_,
	}

	docs, err := requestIndexing(info, testIndexFolderPath1)
	assert.NoError(t, err)
	assert.True(t, gock.IsDone())

	assert.Equal(t, len(docs), 2)
	assert.True(
		t,
		docs[0].Path == testIndexFilePath1 ||
			docs[1].Path == testIndexFilePath1,
	)
	assert.True(
		t,
		docs[0].Path == testIndexFilePath2 ||
			docs[1].Path == testIndexFilePath2,
	)
}

func TestRequestIndexingFail(t *testing.T) {
	defer gock.Off()
	gock.New(string(_TESTURI_)).
		Get(_INDEXFULL_ + "/" + string(testIndexFilePath1)).
		Reply(200).
		JSON(testResponseFail)
	info := IndexerInfo{
		name:             "TestIndexerName",
		cmd:              exec.Command("ls"),
		fileTypesHandled: []utils.FileType{"txt"},
		id:               42,
		endpoint:         _TESTURI_,
	}

	_, err := requestIndexing(info, testIndexFilePath1)
	assert.Error(t, err)
	assert.True(t, gock.IsDone())
}

// Same as TestRequestIndexingFail but with worse JSON response.
func TestRequestIndexingInvalidJSON(t *testing.T) {
	defer gock.Off()
	gock.New(string(_TESTURI_)).
		Get(_INDEXFULL_ + "/" + string(testIndexFilePath1)).
		Reply(200).
		JSON(map[string]string{"invalid": "JSON data send back"})
	info := IndexerInfo{
		name:             "TestIndexerName",
		cmd:              exec.Command("ls"),
		fileTypesHandled: []utils.FileType{"txt"},
		id:               42,
		endpoint:         _TESTURI_,
	}

	_, err := requestIndexing(info, testIndexFilePath1)
	assert.Error(t, err)
	assert.True(t, gock.IsDone())
}
