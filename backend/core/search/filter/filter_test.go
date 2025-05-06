package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testFilter(t *testing.T, string string, expected Filter) {
	result := ParseFilter(string)
	assert.Equal(t, expected, result)
}

func TestIncludes(t *testing.T) {
	testFilter(t, "+foo", Includes("foo"))
}

func TestExcludes(t *testing.T) {
	testFilter(t, "-foo", Excludes("foo"))
}
