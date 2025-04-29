package stemming

import (
	"seekourney/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testStem(t *testing.T, word string, expected string) {
	result := Stem(utils.Word(word))
	assert.Equal(t, utils.Word(expected), result, word)
}

func TestElement(t *testing.T) {
	testStem(t, "element", "element")
}

func TestLosing(t *testing.T) {
	testStem(t, "losing", "lose")
}

func TestFizzing(t *testing.T) {
	testStem(t, "fizzing", "fizz")
}

func TestAllowance(t *testing.T) {
	testStem(t, "allowance", "allow")
}

func TestTanned(t *testing.T) {
	testStem(t, "tanned", "tan")
}

func TestControlling(t *testing.T) {
	testStem(t, "controlling", "control")
}

func TestFailing(t *testing.T) {
	testStem(t, "failing", "fail")
}

func TestFiling(t *testing.T) {
	testStem(t, "filing", "file")
}
