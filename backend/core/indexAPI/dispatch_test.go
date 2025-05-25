package indexAPI

import (
	"seekourney/indexing"
	"seekourney/utils"
	"seekourney/utils/normalize"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

const (
	_TESTINDEXERID_ IndexerID = "testid"
	// URI for mocking index requests.
	// Workaround as we cant call function when definition const.
	_TESTURI_  string     = "http://localhost:39042"
	_TESTPORT_ utils.Port = 39042
	// Dispatch now sends path with settings in hppt body,
	// so this may be unused in test mocks.
	_TESTPATH_ utils.Path = "home/george/my_cool_text_files"
)

// Post request functions add a slash between all words, but when mocking
// we need to add it manually, this is shorthand for that.
// Only use this in tests.
const _SLASHINDEX_ string = "/index"

// Use NewIndexHandler for making indexing handlers.

// nameTestIndexerData creates IndexerData struct needed when testing dispatch.
// By using a function we avoid any potential data to be modified in other
// tests alternatively let use use less boiler-plate setting up.
func makeTestIndexerData() IndexerData {
	return IndexerData{
		ID:       _TESTINDEXERID_,
		Name:     "The Test Indexer",
		ExecPath: "ls",
		Args:     []string{""},
		Port:     _TESTPORT_,
	}
}

// nameTestIndexerData creates Collection struct needed when testing dispatch.
// By using a function we avoid any potential data to be modified in other
// tests alternatively let use use less boiler-plate setting up.
func makeTestCollection() Collection {
	return Collection{
		UnregisteredCollection: UnregisteredCollection{
			Path:                _TESTPATH_,
			IndexerID:           "testid",
			SourceType:          utils.FILE_SOURCE,
			Recursive:           false,
			RespectLastModified: false,
			Normalfunc:          normalize.TO_LOWER,
		},
		ID: "ID",
	}
}

// TODO change startup and shutdown tests to work with new startup/shutdown

/*
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

	defer waitOnTestCMD(info)

	assert.Error(t, startupIndexer(info))
	assert.True(t, gock.IsDone())
}

func TestStartupPingSuccess(t *testing.T) {
	defer gock.Off()
	gock.New(string(_TESTURI_)).
		Get(_PING_).
		Reply(200).
		JSON(indexing.ResponsePing())

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

	defer waitOnTestCMD(info)

	assert.Error(t, startupIndexer(info))
	assert.True(t, gock.IsDone())
}

func TestShutdownValidResponse(t *testing.T) {
	defer gock.Off()
	gock.New(string(_TESTURI_)).
		Get(_PING_).
		Reply(200).
		JSON(indexing.ResponsePing())


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
		JSON(indexing.ResponsePing())


	assert.NoError(t, startupIndexer(info))
	assert.True(t, gock.IsDone())

	gock.New(string(_TESTURI_)).
		Get(_SHUTDOWN_).
		Reply(200).
		JSON(indexing.ResponseFail("failing to shut down indexer"))

	assert.Error(t, shutdownIndexerGraceful(info))
	assert.True(t, gock.IsDone())
}
*/

func TestNewDispatchErrors(t *testing.T) {
	errs := newDispatchErrors()
	assert.True(t, errs.IndexerWasRunning)
	assert.NoError(t, errs.StartupAttempt)
	assert.NoError(t, errs.DispatchAttempt)
}

// This test will fail due to Dispatch expecting the indexer to be started
func TestDispatchSuccessIsRunning(t *testing.T) {
	defer gock.Off()
	gock.New(string(_TESTURI_)).
		Post(_SLASHINDEX_).
		Reply(200).
		JSON(indexing.ResponseSuccess(""))

	handler := NewIndexHandler()
	errs := handler.Dispatch(makeTestIndexerData(), makeTestCollection())
	assert.True(t, gock.IsDone())

	assert.True(t, errs.IndexerWasRunning)
	assert.NoError(t, errs.StartupAttempt)
	assert.NoError(t, errs.DispatchAttempt)
	// Indexer was not added to running indexers map.
	assert.Equal(t, len(handler.Indexers), 0)
}

func TestDispatchSuccessNotRunning(t *testing.T) {
	defer gock.Off()
	gock.New(string(_TESTURI_)).
		Post(_SLASHINDEX_).
		Reply(500).
		JSON("")
	gock.New(string(_TESTURI_)).
		Get(_PING_).
		Reply(200).
		JSON(indexing.ResponsePing())
	gock.New(string(_TESTURI_)).
		Post(_SLASHINDEX_).
		Reply(200).
		JSON(indexing.ResponseSuccess(""))

	handler := NewIndexHandler()
	errs := handler.Dispatch(makeTestIndexerData(), makeTestCollection())
	assert.True(t, gock.IsDone())

	assert.False(t, errs.IndexerWasRunning)
	assert.NoError(t, errs.StartupAttempt)
	assert.NoError(t, errs.DispatchAttempt)
	// Indexer was added to running indexers map.
	assert.Equal(t, len(handler.Indexers), 1)
}

func TestDispatchStartupFail(t *testing.T) {
	defer gock.Off()
	gock.New(string(_TESTURI_)).
		Post(_SLASHINDEX_).
		Reply(500).
		JSON("")
	gock.New(string(_TESTURI_)).
		Get(_PING_).
		Reply(200).
		JSON(indexing.ResponseFail("failed to startup indexer"))

	handler := NewIndexHandler()
	errs := handler.Dispatch(makeTestIndexerData(), makeTestCollection())
	assert.True(t, gock.IsDone())

	assert.False(t, errs.IndexerWasRunning)
	assert.Error(t, errs.StartupAttempt)
	assert.Error(t, errs.DispatchAttempt)
	// Indexer was not added to running indexers map because startup failed.
	assert.Equal(t, len(handler.Indexers), 0)
}

func TestDispatchIndexFail(t *testing.T) {
	defer gock.Off()
	gock.New(string(_TESTURI_)).
		Post(_SLASHINDEX_).
		Reply(200).
		JSON(indexing.ResponseFail("unable to fulfill indexing request"))

	handler := NewIndexHandler()
	errs := handler.Dispatch(makeTestIndexerData(), makeTestCollection())
	assert.True(t, gock.IsDone())

	assert.True(t, errs.IndexerWasRunning)
	assert.NoError(t, errs.StartupAttempt)
	assert.Error(t, errs.DispatchAttempt)
	// Indexer was not added to running indexers map.
	assert.Equal(t, len(handler.Indexers), 0)
}
