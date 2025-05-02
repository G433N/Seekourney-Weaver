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
	_TESTURI_ utils.Endpoint = utils.Endpoint(_ENDPOINTPREFIX_ + "39100")
)

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
		JSON(indexing.ResponseFail(""))
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
		JSON(indexing.ResponsePong())
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
		JSON(indexing.ResponsePong())
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
		JSON(indexing.ResponseExiting())

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
		JSON(indexing.ResponsePong())
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
		JSON(indexing.ResponseFail("failing to shut down indexer"))

	assert.Error(t, shutdownIndexerGraceful(info))
	assert.True(t, gock.IsDone())
}

// Test data for mocking indexing requests.
const testIndexFolderPath1 utils.Path = "home/george/my_cool_text_files"

const testIndexFilePath1 utils.Path = "home/george/my_cool_text_files/first.txt"
const testIndexFilePath2 utils.Path = "home/george/my_cool_text_files/other.txt"

var testResponseDoc1 UnnormalizedDocument = UnnormalizedDocument{
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

var testResponseDoc2 UnnormalizedDocument = UnnormalizedDocument{
	Path:   testIndexFilePath2,
	Source: utils.SourceLocal,
	Words: utils.FrequencyMap{
		"wood":  234,
		"steel": 52,
	},
}

// TODO TEMP until POST requests exist
var testIndexingResponse1 IndexerResponse = IndexerResponse{
	Status: indexing.STATUSSUCCESSFUL,
	Data: ResponseData{
		Documents: []UnnormalizedDocument{testResponseDoc1},
	},
}
var testIndexingResponse2 IndexerResponse = IndexerResponse{
	Status: indexing.STATUSSUCCESSFUL,
	Data: ResponseData{
		Documents: []UnnormalizedDocument{
			testResponseDoc1,
			testResponseDoc2,
		},
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
		JSON(indexing.ResponseFail("failed to index requested path"))
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

func TestNewIndexErrorsLow(t *testing.T) {
	errs := newIndexErrors(1)
	assert.NoError(t, errs.Startup)
	assert.NoError(t, errs.Shutdown)
	assert.Equal(t, len(errs.Indexing), 1)
	assert.NoError(t, errs.Indexing[0])
}

func TestNewIndexErrorsHigh(t *testing.T) {
	errs := newIndexErrors(42)
	assert.NoError(t, errs.Startup)
	assert.NoError(t, errs.Shutdown)
	assert.Equal(t, len(errs.Indexing), 42)
	for i := range errs.Indexing {
		assert.NoError(t, errs.Indexing[i])
	}
}

func TestIndexOneSuccess(t *testing.T) {
	defer gock.Off()
	gock.New(string(_TESTURI_)).
		Get(_PING_).
		Reply(200).
		JSON(indexing.ResponsePong())
	gock.New(string(_TESTURI_)).
		Get(_INDEXFULL_ + "/" + string(testIndexFilePath1)).
		Reply(200).
		JSON(testIndexingResponse1)
	gock.New(string(_TESTURI_)).
		Get(_SHUTDOWN_).
		Reply(200).
		JSON(indexing.ResponseExiting())
	info := IndexerInfo{
		name:             "TestIndexerName",
		cmd:              exec.Command("ls"),
		fileTypesHandled: []utils.FileType{"txt"},
		id:               42,
		endpoint:         _TESTURI_,
	}

	docs, errStruct := IndexOne(info, testIndexFilePath1)
	assert.True(t, gock.IsDone())

	assert.Equal(t, len(docs), 1) // Path yields 1 Document
	assert.Equal(t, len(errStruct.Indexing), 1)
	assert.NoError(t, errStruct.Startup)
	assert.NoError(t, errStruct.Shutdown)
	assert.NoError(t, errStruct.Indexing[0])
}

func TestIndexOneStartupFail(t *testing.T) {
	defer gock.Off()
	gock.New(string(_TESTURI_)).
		Get(_PING_).
		Reply(200).
		JSON(indexing.ResponseFail("failed to startup indexer"))
	info := IndexerInfo{
		name:             "TestIndexerName",
		cmd:              exec.Command("ls"),
		fileTypesHandled: []utils.FileType{"txt"},
		id:               42,
		endpoint:         _TESTURI_,
	}

	_, errStruct := IndexOne(info, testIndexFilePath1)
	assert.True(t, gock.IsDone())

	assert.Equal(t, len(errStruct.Indexing), 1)
	assert.Error(t, errStruct.Startup)
	assert.Error(t, errStruct.Shutdown)
	assert.Error(t, errStruct.Indexing[0])
}

func TestIndexOneIndexFail(t *testing.T) {
	defer gock.Off()
	gock.New(string(_TESTURI_)).
		Get(_PING_).
		Reply(200).
		JSON(indexing.ResponsePong())
	gock.New(string(_TESTURI_)).
		Get(_INDEXFULL_ + "/" + string(testIndexFilePath1)).
		Reply(200).
		JSON(indexing.ResponseFail("failed to index requested path"))
	gock.New(string(_TESTURI_)).
		Get(_SHUTDOWN_).
		Reply(200).
		JSON(indexing.ResponseExiting())
	info := IndexerInfo{
		name:             "TestIndexerName",
		cmd:              exec.Command("ls"),
		fileTypesHandled: []utils.FileType{"txt"},
		id:               42,
		endpoint:         _TESTURI_,
	}

	_, errStruct := IndexOne(info, testIndexFilePath1)
	assert.True(t, gock.IsDone())

	assert.NoError(t, errStruct.Startup)
	assert.NoError(t, errStruct.Shutdown)
	assert.Error(t, errStruct.Indexing[0])
}

func TestIndexOneShutdownFail(t *testing.T) {
	defer gock.Off()
	gock.New(string(_TESTURI_)).
		Get(_PING_).
		Reply(200).
		JSON(indexing.ResponsePong())
	gock.New(string(_TESTURI_)).
		Get(_INDEXFULL_ + "/" + string(testIndexFilePath1)).
		Reply(200).
		JSON(testIndexingResponse1)
	gock.New(string(_TESTURI_)).
		Get(_SHUTDOWN_).
		Reply(200).
		JSON(indexing.ResponseFail("unable to shut down"))
	info := IndexerInfo{
		name:             "TestIndexerName",
		cmd:              exec.Command("ls"),
		fileTypesHandled: []utils.FileType{"txt"},
		id:               42,
		endpoint:         _TESTURI_,
	}

	docs, errStruct := IndexOne(info, testIndexFilePath1)
	assert.True(t, gock.IsDone())

	assert.Equal(t, len(docs), 1) // Path yields 1 Document
	assert.Equal(t, len(errStruct.Indexing), 1)
	assert.NoError(t, errStruct.Startup)
	assert.Error(t, errStruct.Shutdown)
	assert.NoError(t, errStruct.Indexing[0])
}

func TestIndexManyPartialSuccess(t *testing.T) {
	defer gock.Off()
	gock.New(string(_TESTURI_)).
		Get(_PING_).
		Reply(200).
		JSON(indexing.ResponsePong())
	gock.New(string(_TESTURI_)).
		Get(_INDEXFULL_ + "/" + string(testIndexFilePath1)).
		Reply(200).
		JSON(indexing.ResponseFail("failed to index path"))
	gock.New(string(_TESTURI_)).
		Get(_INDEXFULL_ + "/" + string(testIndexFilePath2)).
		Reply(200).
		JSON(indexing.ResponseFail("failed to index path"))
	gock.New(string(_TESTURI_)).
		Get(_INDEXFULL_ + "/" + string(testIndexFolderPath1)).
		Reply(200).
		JSON(testIndexingResponse2)
	gock.New(string(_TESTURI_)).
		Get(_INDEXFULL_ + "/" + string(testIndexFilePath1)).
		Reply(200).
		JSON(indexing.ResponseFail("failed to index path"))
	gock.New(string(_TESTURI_)).
		Get(_SHUTDOWN_).
		Reply(200).
		JSON(indexing.ResponseExiting())
	info := IndexerInfo{
		name:             "TestIndexerName",
		cmd:              exec.Command("ls"),
		fileTypesHandled: []utils.FileType{"txt"},
		id:               42,
		endpoint:         _TESTURI_,
	}

	paths := []utils.Path{
		testIndexFilePath1,
		testIndexFilePath2,
		testIndexFolderPath1,
		testIndexFilePath1,
	}
	// First fails, second fails, third succeeds, fourth fails.
	manyDocs, errStruct := IndexMany(info, paths)
	assert.True(t, gock.IsDone())

	assert.Equal(t, len(manyDocs), 4)
	assert.Equal(t, len(manyDocs), len(errStruct.Indexing))
	assert.NoError(t, errStruct.Startup)
	assert.NoError(t, errStruct.Shutdown)

	assert.Error(t, errStruct.Indexing[0])
	assert.Error(t, errStruct.Indexing[1])
	assert.NoError(t, errStruct.Indexing[2])
	assert.Error(t, errStruct.Indexing[3])

	// The successful indexing attempt should produce 2 documents.
	assert.Equal(t, len(manyDocs[2]), 2)
}
