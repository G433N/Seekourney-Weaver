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
			Path:       "test/path/1",
			Source:     utils.SOURCE_LOCAL,
			Words:      utils.FrequencyMap{"first": 1, "second": 2},
			Collection: "test_source_id",
		},
		{
			Path:       "test/path/2",
			Source:     utils.SOURCE_WEB,
			Words:      utils.FrequencyMap{"green": 42, "blue": 5},
			Collection: "test_source_id",
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
	assert.Equal(t, goData.Data.Documents[0].Collection, udocs[0].Collection)

	assert.Equal(t, goData.Data.Documents[1].Path, udocs[1].Path)
	assert.Equal(t, goData.Data.Documents[1].Source, udocs[1].Source)
	assert.True(
		t,
		reflect.DeepEqual(goData.Data.Documents[1].Words, udocs[1].Words),
	)
	assert.Equal(t, goData.Data.Documents[1].Collection, udocs[1].Collection)
}

func TestResponsePathText(t *testing.T) {
	pathTexts := []PathText{
		{
			Path: "test/path/1",
			Text: "test text 1",
		},
		{
			Path: "test/path/2",
			Text: "test text 2",
		},
	}

	jsonData := ResponsePathText(pathTexts)
	goData := IndexerTextResponse{}
	err := json.Unmarshal(jsonData, &goData)

	assert.NoError(t, err)
	assert.Equal(t, goData.Status, STATUSSUCCESSFUL)

	assert.Equal(t, goData.Data.PathTexts[0].Path, pathTexts[0].Path)
	assert.Equal(t, goData.Data.PathTexts[0].Text, pathTexts[0].Text)

	assert.Equal(t, goData.Data.PathTexts[1].Path, pathTexts[1].Path)
	assert.Equal(t, goData.Data.PathTexts[1].Text, pathTexts[1].Text)
}
