package folder

import (
	"reflect" // map equality
	"seekourney/document"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ( // Can't use const here
	testDocAlpha document.Document = document.New("These are", 42)
	testDocBeta  document.Document = document.New("some bogus", 43)
	testDocGamma document.Document = document.New("file paths", 44)
	// testDocDelta   document.Document = document.New("not important", 45)
	// testDocEpsilon document.Document = document.New("for testing", 46)
)

func TestAddRemoveDoc(t *testing.T) {
	docMap := make(DocMap)
	docMap[testDocAlpha.Path] = testDocAlpha
	docMap[testDocBeta.Path] = testDocBeta
	docMap[testDocGamma.Path] = testDocGamma
	expected := New(docMap)

	result := Default()
	result.AddDoc(testDocBeta.Path, testDocBeta)
	result.AddDoc(testDocGamma.Path, testDocGamma)
	result.AddDoc(testDocAlpha.Path, testDocAlpha)

	assert.True(t, reflect.DeepEqual(result, expected))

	delete(expected.docs, testDocBeta.Path)
	removedDoc, firstOK := result.RemoveDoc(testDocBeta.Path)

	assert.True(t, reflect.DeepEqual(result, expected))
	assert.True(t, firstOK)
	assert.Equal(t, removedDoc, testDocBeta)

	_, secondOK := result.RemoveDoc("nonexistent path to file")
	assert.False(t, secondOK)

	_, thirdOK := result.RemoveDoc(testDocBeta.Path)
	assert.False(t, thirdOK)
}
