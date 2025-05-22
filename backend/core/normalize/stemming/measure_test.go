package stemming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testMeasure(t *testing.T, word string, expected uint) {
	result := calcMeasure(word)
	assert.Equal(t, measure(expected), result, "%s", word)
}

func TestMeasure0(t *testing.T) {
	testMeasure(t, "tr", 0)
	testMeasure(t, "ee", 0)
	testMeasure(t, "tree", 0)
	testMeasure(t, "y", 0)
	testMeasure(t, "by", 0)
}

func TestMeasure1(t *testing.T) {
	testMeasure(t, "trouble", 1)
	testMeasure(t, "oats", 1)
	testMeasure(t, "trees", 1)
	testMeasure(t, "ivy", 1)
}

func TestMeasure2(t *testing.T) {
	testMeasure(t, "troubles", 2)
	testMeasure(t, "private", 2)
	testMeasure(t, "oaten", 2)
	testMeasure(t, "orrery", 2)
}
