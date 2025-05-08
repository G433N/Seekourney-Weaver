package indexing

import (
	"encoding/json"
	"reflect"
	"seekourney/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

const _TESTMESSAGE_ string = "test message"

func TestResponseSuccess(t *testing.T) {
	jsonData := ResponseSuccess(_TESTMESSAGE_)
	goData := IndexerResponse{}
	err := json.Unmarshal(jsonData, &goData)

	assert.NoError(t, err)
	assert.Equal(t, goData.Status, STATUSSUCCESSFUL)
	assert.Equal(t, goData.Data.Message, _TESTMESSAGE_)
}

func TestResponseFail(t *testing.T) {
	jsonData := ResponseFail(_TESTMESSAGE_)
	goData := IndexerResponse{}
	err := json.Unmarshal(jsonData, &goData)

	assert.NoError(t, err)
	assert.Equal(t, goData.Status, STATUSFAILURE)
	assert.Equal(t, goData.Data.Message, _TESTMESSAGE_)
}

func TestResponsePing(t *testing.T) {
	jsonData := ResponsePing()
	goData := IndexerResponse{}
	err := json.Unmarshal(jsonData, &goData)

	assert.NoError(t, err)
	assert.Equal(t, goData.Status, STATUSSUCCESSFUL)
	assert.Equal(t, goData.Data.Message, MESSAGEPONG)
}

func TestResponseExiting(t *testing.T) {
	jsonData := ResponseExiting()
	goData := IndexerResponse{}
	err := json.Unmarshal(jsonData, &goData)

	assert.NoError(t, err)
	assert.Equal(t, goData.Status, STATUSSUCCESSFUL)
	assert.Equal(t, goData.Data.Message, MESSAGEEXITING)
}

func TestResponseDocs(t *testing.T) {
	udocs := []UnnormalizedDocument{
		{
			Path:     "test/path/1",
			Source:   utils.SourceLocal,
			Words:    utils.FrequencyMap{"first": 1, "second": 2},
			SourceID: 99,
		},
		{
			Path:     "test/path/2",
			Source:   utils.SourceWeb,
			Words:    utils.FrequencyMap{"green": 42, "blue": 5},
			SourceID: 99,
		},
	}

	jsonData := ResponseDocs(udocs)
	goData := IndexerResponse{}
	err := json.Unmarshal(jsonData, &goData)

	assert.NoError(t, err)
	assert.Equal(t, goData.Status, STATUSSUCCESSFUL)

	assert.Equal(t, goData.Data.Documents[0].Path, udocs[0].Path)
	assert.Equal(t, goData.Data.Documents[0].Source, udocs[0].Source)
	assert.True(
		t,
		reflect.DeepEqual(goData.Data.Documents[0].Words, udocs[0].Words),
	)
	assert.Equal(t, goData.Data.Documents[0].SourceID, udocs[0].SourceID)

	assert.Equal(t, goData.Data.Documents[1].Path, udocs[1].Path)
	assert.Equal(t, goData.Data.Documents[1].Source, udocs[1].Source)
	assert.True(
		t,
		reflect.DeepEqual(goData.Data.Documents[1].Words, udocs[1].Words),
	)
	assert.Equal(t, goData.Data.Documents[1].SourceID, udocs[1].SourceID)
}
