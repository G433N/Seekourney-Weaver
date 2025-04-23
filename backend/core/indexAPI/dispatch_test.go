package indexAPI

import (
	"os/exec"
	"seekourney/utils"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

const (
	_TESTURI_ restEndpoint = restEndpoint(_ENDPOINTPREFIX_ + "39100")
	_STATUS_  string       = "status"
	_DATA_    string       = "data"
)

func TestStartupPingFail(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	gock.New(string(_TESTURI_)).
		Get(_PING_).
		Reply(200).
		JSON(map[string]string{_STATUS_: _STATUSFAILURE_, _DATA_: ""})
	info := IndexerInfo{
		name:             "TestIndexerName",
		cmd:              exec.Command("ls"),
		fileTypesHandled: []utils.FileType{"txt"},
		id:               42,
		endpoint:         _TESTURI_,
	}
	// Wait() is used instead of shutdownIndexerGraceful.
	// This is needed as shutdownIndexerGraceful will force kill if
	// correct JSON is not returned, which won't cleanup resources.
	defer info.cmd.Wait()

	assert.Error(t, startupIndexer(info))
	assert.True(t, gock.IsDone())
}

func TestStartupPingSuccess(t *testing.T) {
	defer gock.Off()
	gock.New(string(_TESTURI_)).
		Get(_PING_).
		Reply(200).
		JSON(map[string]string{_STATUS_: _STATUSSUCCESSFUL_, _DATA_: _PONG_})
	info := IndexerInfo{
		name:             "TestIndexerName",
		cmd:              exec.Command("ls"),
		fileTypesHandled: []utils.FileType{"txt"},
		id:               42,
		endpoint:         _TESTURI_,
	}
	// Wait() is used instead of shutdownIndexerGraceful.
	// This is needed as shutdownIndexerGraceful will force kill if
	// correct JSON is not returned, which won't cleanup resources.
	defer info.cmd.Wait()

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
	// Wait() is used instead of shutdownIndexerGraceful.
	// This is needed as shutdownIndexerGraceful will force kill if
	// correct JSON is not returned, which won't cleanup resources.
	defer info.cmd.Wait()

	assert.Error(t, startupIndexer(info))
	assert.True(t, gock.IsDone())
}

func TestShutdownValidResponse(t *testing.T) {
	defer gock.Off()
	gock.New(string(_TESTURI_)).
		Get(_PING_).
		Reply(200).
		JSON(map[string]string{_STATUS_: _STATUSSUCCESSFUL_, _DATA_: _PONG_})
	info := IndexerInfo{
		name:             "TestIndexerName",
		cmd:              exec.Command("ls"),
		fileTypesHandled: []utils.FileType{"txt"},
		id:               42,
		endpoint:         _TESTURI_,
	}

	startupIndexer(info)

	gock.New(string(_TESTURI_)).
		Get(_SHUTDOWN_).
		Reply(200).
		JSON(map[string]string{_STATUS_: _STATUSSUCCESSFUL_, _DATA_: _EXITING_})

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
		JSON(map[string]string{_STATUS_: _STATUSSUCCESSFUL_, _DATA_: _PONG_})
	info := IndexerInfo{
		name:             "TestIndexerName",
		cmd:              exec.Command("ls"),
		fileTypesHandled: []utils.FileType{"txt"},
		id:               42,
		endpoint:         _TESTURI_,
	}

	startupIndexer(info)

	gock.New(string(_TESTURI_)).
		Get(_SHUTDOWN_).
		Reply(200).
		JSON(map[string]string{_STATUS_: _STATUSFAILURE_, _DATA_: ":("})

	assert.Error(t, shutdownIndexerGraceful(info))
	assert.True(t, gock.IsDone())
}
