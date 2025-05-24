package stemming

import "iter"

// charType represents the type of a character in the context of
// the Porter stemming algorithm. It can be either a vowel or a consonant.
type charType int

const (
	__VOWEL__ charType = iota
	__CONSONANT__
)

// isVowel checks if a byte is a context free vowel
func isVowel(b byte) bool {
	switch b {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	default:
		return false
	}
}

// getCharType returns the charType of a byte
func getCharType(char byte) charType {
	if isVowel(char) {
		return __VOWEL__
	} else {
		return __CONSONANT__
	}
}

// inverse returns the inverse of the charType
func (t charType) inverse() charType {
	switch t {
	case __VOWEL__:
		return __CONSONANT__
	case __CONSONANT__:
		return __VOWEL__
	}
	panic("unreachable")
}

// charTypeIter iterates over the characters of a string and yields
func charTypeIter(word string) iter.Seq2[charType, charType] {
	bytes := []byte(word)

	return func(yield func(charType, charType) bool) {
		prev := __CONSONANT__
		for _, char := range bytes {
			var charType charType

			switch char {
			case 'y':
				charType = prev.inverse()
			default:
				charType = getCharType(char)
			}

			if !yield(charType, prev) {
				return
			}
			prev = charType
		}
	}

}

// containsVowel checks if a string contains a vowel
func containsVowel(str string) bool {

	for charType := range charTypeIter(str) {
		if charType == __VOWEL__ {
			return true
		}
	}
	return false
}
