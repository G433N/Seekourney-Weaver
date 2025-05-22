package words

import (
	"github.com/stretchr/testify/assert"
	"seekourney/utils"
	"testing"
)

type word = utils.Word

func TestWordsIter(t *testing.T) {
	s := "Hello World! test gamertag123"
	expected := []word{"Hello", "World", "test", "gamertag123"}

	i := 0
	for w := range WordsIter(s) {
		assert.Equal(t, w, expected[i])
		i++
	}

	assert.Equal(t, i, len(expected))
}

func TestWordsIterEmpty(t *testing.T) {
	s := ""

	for range WordsIter(s) {
		assert.FailNow(t, "Expected no words, got one")
	}
}

func TestWordsIterSingleWord(t *testing.T) {

	s := "Hello"
	expected := []word{"Hello"}

	i := 0
	for w := range WordsIter(s) {
		assert.Equal(t, w, expected[i])
		i++
	}

	assert.Equal(t, i, len(expected))
}

func TestWordsIterPunctionation(t *testing.T) {
	s := "Hello!...."
	expected := []word{"Hello"}

	i := 0
	for w := range WordsIter(s) {
		assert.Equal(t, w, expected[i])
		i++
	}

	assert.Equal(t, i, len(expected))
}

func TestWordsNewLine(t *testing.T) {
	s := "Hello\nWorld!"
	expected := []word{"Hello", "World"}

	i := 0
	for w := range WordsIter(s) {
		assert.Equal(t, w, expected[i])
		i++
	}

	assert.Equal(t, i, len(expected))
}

func TestWordsIterMultipleSpaces(t *testing.T) {
	s := "   Hello    World!   "
	expected := []word{"Hello", "World"}

	i := 0
	for w := range WordsIter(s) {
		assert.Equal(t, w, expected[i])
		i++
	}

	assert.Equal(t, i, len(expected))
}

func TestWordsIterMultipleLines(t *testing.T) {
	s := "Hello\n\nWorld!\n"
	expected := []word{"Hello", "World"}

	i := 0
	for w := range WordsIter(s) {
		assert.Equal(t, w, expected[i])
		i++
	}

	assert.Equal(t, i, len(expected))
}

func TestWordsIterNummeric(t *testing.T) {
	s := "str123ing 456"
	expected := []word{"str123ing", "456"}

	i := 0
	for w := range WordsIter(s) {
		assert.Equal(t, w, expected[i])
		i++
	}

	assert.Equal(t, i, len(expected))
}

func TestWordsIterUTF8(t *testing.T) {
	s := "Hello 世界"
	expected := []word{"Hello", "世界"}

	i := 0
	for w := range WordsIter(s) {
		assert.Equal(t, w, expected[i])
		i++
	}

	assert.Equal(t, i, len(expected))
}

func TestWordsIterParenthesis(t *testing.T) {
	s := "(Hello (World))"
	expected := []word{"Hello", "World"}

	i := 0
	for w := range WordsIter(s) {
		assert.Equal(t, w, expected[i])
		i++
	}

	assert.Equal(t, i, len(expected))
}
