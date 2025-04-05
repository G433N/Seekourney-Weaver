package words

import (
	"iter"
)

const UTF8Prefix = 0b10000000

// isAscii returns true if the byte is an ASCII character.
func isAscii(b byte) bool {
	return b&0b10000000 == 0
}

// isUTF8 returns true if the byte is a UTF-8 character.
func isUTF8(b byte) bool {
	return b&UTF8Prefix != 0
}

// isASCIIAlphaNumeric returns true if the byte is an ASCII alphanumeric character.
func isASCIIAlphaNumeric(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9')
}

// charLen returns the number of bytes in a UTF-8 character.
// Including the first byte.
// If the byte is not a UTF-8 character, it returns 0.
// If the byte is a continuation byte, it returns -1.
func charLen(b byte) int {
	if isAscii(b) {
		return 1
	}

	if b&0b11100000 == 0b11000000 {
		return 2
	}

	if b&0b11110000 == 0b11100000 {
		return 3
	}

	if b&0b11111000 == 0b11110000 {
		return 4
	}

	if b&0b11000000 == 0b10000000 {
		return -1
	}

	return 0
}

func wordSplit(b byte) bool {

	// TODO: Check if UTF-8 chars have any characters we want to splkit words with
	if isUTF8(b) {
		return false
	}

	return !isASCIIAlphaNumeric(b)

}

// / yieldWord yields a word from the byte slice.
// / It takes a yield function, a byte slice, and the start and end indices of the word.
// / It returns true if the iteration should continue, and false if it should stop.
func yieldWord(yield func(string) bool, bytes []byte, start int, end int) bool {
	if start != end {
		word := string(bytes[start:end])
		return yield(word)
	}

	return true
}

func WordsIterBytes(bytes []byte) iter.Seq[string] {
	word_iter := func(yield func(string) bool) {

		start := 0
		end := 0

		i := 0

		for i < len(bytes) {
			b := bytes[i]

			if wordSplit(b) {
				if !yieldWord(yield, bytes, start, end) {
					return
				}

				start = i + 1
				end = i + 1
			}
			// Ship UTF-8 continuation bytes
			len := charLen(b)
			i += len
			end = i
		}

		if !yieldWord(yield, bytes, start, end) {
			return
		}
	}

	return word_iter
}

// WordsIter takes a string and returns an iterator that yields each word in the string.
func WordsIter(s string) iter.Seq[string] {
	bytes := []byte(s)
	return WordsIterBytes(bytes)
}
