package search

import (
	"seekourney/core/config"
	"seekourney/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseQuery(t *testing.T) {
	config := config.New()
	query :=
		utils.Query("test +at \"hello world\" -the +next \"hihi\" you -joke")
	parsedQuery := parseQuery(config, query)

	// NOTE: Every quote and "-" filter adds one space in the modified quote
	assert.Equal(t, "test at   next  you ", string(parsedQuery.ModifiedQuery))
	assert.Equal(t, 2, len(parsedQuery.PlusWords))
	assert.Equal(t, 2, len(parsedQuery.MinusWords))
	assert.Equal(t, 2, len(parsedQuery.Quotes))

	expectedPlusWords := []string{"at", "next"}
	for wordIndex := range parsedQuery.PlusWords {
		assert.Equal(
			t,
			expectedPlusWords[wordIndex],
			parsedQuery.PlusWords[wordIndex])
	}

	expectedMinusWords := []string{"the", "joke"}
	for wordIndex := range parsedQuery.MinusWords {
		assert.Equal(
			t,
			expectedMinusWords[wordIndex],
			parsedQuery.MinusWords[wordIndex])
	}

	expectedQuotes := []string{"hello world", "hihi"}
	for wordIndex := range parsedQuery.Quotes {
		assert.Equal(
			t,
			expectedQuotes[wordIndex],
			parsedQuery.Quotes[wordIndex])
	}
}

func TestWordsFromQuotes(t *testing.T) {
	quotes := []string{
		"test1 test2  test3  ",
		"hello world!",
		"*^good??error(()",
	}

	expextedWords := []string{
		"test1",
		"test2",
		"test3",
		"hello",
		"world",
		"good",
		"error",
	}

	retrievedWords := wordsFromQuotes(quotes)

	assert.Equal(t, expextedWords, retrievedWords)
}

func TestUpdateStatus(t *testing.T) {
	filterStatus := updateFilterStatus("+")
	assert.Equal(t, _INPLUS_, filterStatus)

	filterStatus = updateFilterStatus("-")
	assert.Equal(t, _INMINUS_, filterStatus)

	filterStatus = updateFilterStatus("\"")
	assert.Equal(t, _INQUOTE_, filterStatus)
}

func TestUpdateModifiedQuery(t *testing.T) {
	config := config.New()
	parsedQuery := utils.ParsedQuery{
		ModifiedQuery: "",
		PlusWords:     make([]string, 0),
		MinusWords:    make([]string, 0),
		Quotes:        make([]string, 0),
	}

	currentByte := " "
	currentFilterString := "hello"
	filterStatus := _INPLUS_

	newStatus, shouldContinue := updateModifiedQuery(
		config,
		&parsedQuery,
		currentByte,
		&currentFilterString,
		filterStatus,
	)

	assert.Equal(t, _NOFILTER_, newStatus)
	assert.Equal(t, true, shouldContinue)
	assert.Equal(t, []string{"hello"}, parsedQuery.PlusWords)
	assert.Equal(t, "", currentFilterString)

	currentByte = " "
	currentFilterString = "bath"
	filterStatus = _INMINUS_

	newStatus, shouldContinue = updateModifiedQuery(
		config,
		&parsedQuery,
		currentByte,
		&currentFilterString,
		filterStatus,
	)

	assert.Equal(t, _NOFILTER_, newStatus)
	assert.Equal(t, true, shouldContinue)
	assert.Equal(t, []string{"bath"}, parsedQuery.MinusWords)
	assert.Equal(t, "", currentFilterString)

	currentByte = "\""
	currentFilterString = "hello world"
	filterStatus = _INQUOTE_

	newStatus, shouldContinue = updateModifiedQuery(
		config,
		&parsedQuery,
		currentByte,
		&currentFilterString,
		filterStatus,
	)

	assert.Equal(t, _NOFILTER_, newStatus)
	assert.Equal(t, false, shouldContinue)
	assert.Equal(t, []string{"hello world"}, parsedQuery.Quotes)
	assert.Equal(t, "", currentFilterString)
}
