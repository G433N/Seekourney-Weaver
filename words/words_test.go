package words

import (
	"seekourney/utils"
	"testing"
)

type word = utils.Word

func TestWordsIter(t *testing.T) {
	s := "Hello World! test gamertag123"
	expected := []word{"Hello", "World", "test", "gamertag123"}

	i := 0
	for w := range WordsIter(s) {
		if w != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], w)
		}
		i++
	}
	if i != len(expected) {
		t.Errorf("Expected %d words, got %d", len(expected), i)
	}
}

func TestWordsIterEmpty(t *testing.T) {
	s := ""

	for range WordsIter(s) {
		t.Errorf("Expected no words, got one")
	}
}

func TestWordsIterSingleWord(t *testing.T) {

	s := "Hello"
	expected := []word{"Hello"}

	i := 0
	for w := range WordsIter(s) {
		if w != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], w)
		}
		i++
	}
	if i != len(expected) {
		t.Errorf("Expected %d words, got %d", len(expected), i)
	}
}

func TestWordsIterPunctionation(t *testing.T) {
	s := "Hello!...."
	expected := []word{"Hello"}

	i := 0
	for w := range WordsIter(s) {
		if w != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], w)
		}
		i++
	}
	if i != len(expected) {
		t.Errorf("Expected %d words, got %d", len(expected), i)
	}
}

func TestWordsNewLine(t *testing.T) {
	s := "Hello\nWorld!"
	expected := []word{"Hello", "World"}

	i := 0
	for w := range WordsIter(s) {
		if w != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], w)
		}
		i++
	}
	if i != len(expected) {
		t.Errorf("Expected %d words, got %d", len(expected), i)
	}
}

func TestWordsIterMultipleSpaces(t *testing.T) {
	s := "   Hello    World!   "
	expected := []word{"Hello", "World"}

	i := 0
	for w := range WordsIter(s) {
		if w != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], w)
		}
		i++
	}
	if i != len(expected) {
		t.Errorf("Expected %d words, got %d", len(expected), i)
	}
}

func TestWordsIterMultipleLines(t *testing.T) {
	s := "Hello\n\nWorld!\n"
	expected := []word{"Hello", "World"}

	i := 0
	for w := range WordsIter(s) {
		if w != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], w)
		}
		i++
	}
	if i != len(expected) {
		t.Errorf("Expected %d words, got %d", len(expected), i)
	}
}

func TestWordsIterNummeric(t *testing.T) {
	s := "str123ing 456"
	expected := []word{"str123ing", "456"}

	i := 0
	for w := range WordsIter(s) {
		if w != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], w)
		}
		i++
	}
	if i != len(expected) {
		t.Errorf("Expected %d words, got %d", len(expected), i)
	}
}

func TestWordsIterUTF8(t *testing.T) {
	s := "Hello 世界"
	expected := []word{"Hello", "世界"}

	i := 0
	for w := range WordsIter(s) {
		if w != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], w)
		}
		i++
	}
	if i != len(expected) {
		t.Errorf("Expected %d words, got %d", len(expected), i)
	}
}

func TestWordsIterParenthesis(t *testing.T) {
	s := "(Hello (World))"
	expected := []word{"Hello", "World"}

	i := 0
	for w := range WordsIter(s) {
		if w != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], w)
		}
		i++
	}
	if i != len(expected) {
		t.Errorf("Expected %d words, got %d", len(expected), i)
	}
}
