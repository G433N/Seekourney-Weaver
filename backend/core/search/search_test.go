package search

import (
	"seekourney/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseQuery(t *testing.T) {
	query := utils.Query("test +at _hello world_ -the +next _hihi_ you -joke")
	parsedQuery := parseQuery(query)

	// NOTE: Every quote and "-" filter adds one space in the modified quote
	assert.Equal(t, "test at   next  you ", string(parsedQuery.ModifiedQuery))

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
