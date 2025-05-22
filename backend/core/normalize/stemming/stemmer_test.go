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

func TestRational(t *testing.T) {
	testStem(t, "rational", "ration")
}

func TestValency(t *testing.T) {
	testStem(t, "valency", "valenc")
}

func TestHesitancy(t *testing.T) {
	testStem(t, "hesitancy", "hesit")
}

func TestConformability(t *testing.T) {
	testStem(t, "conformability", "conform")
}

func TestRadical(t *testing.T) {
	testStem(t, "radical", "radic")
}

func TestFormative(t *testing.T) {
	testStem(t, "formative", "form")
}

func TestFormalize(t *testing.T) {
	testStem(t, "formalize", "formal")
}

func TestElectricity(t *testing.T) {
	testStem(t, "electricity", "electr")
}

func TestElectrical(t *testing.T) {
	testStem(t, "electrical", "electr")
}

func TestHopefulness(t *testing.T) {
	testStem(t, "hopefulness", "hope")
}

func TestGoodness(t *testing.T) {
	testStem(t, "goodness", "good")
}
