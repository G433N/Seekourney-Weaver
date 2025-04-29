package stemming

import "iter"

// charType represents the type of a character in the context of
// the Porter stemming algorithm. It can be either a vowel or a consonant.
type charType int

const (
	vowel charType = iota
	consonant
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
		return vowel
	}
	return consonant
}

// inverse returns the inverse of the charType
func (t charType) inverse() charType {
	switch t {
	case vowel:
		return consonant
	case consonant:
		return vowel
	}
	panic("unreachable")
}

// charTypeIter iterates over the characters of a string and yields
func charTypeIter(word string) iter.Seq2[charType, charType] {
	bytes := []byte(word)

	return func(yield func(charType, charType) bool) {
		prev := consonant
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
		if charType == vowel {
			return true
		}
	}
	return false
}
