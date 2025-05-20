package words

import (
	"iter"
	"seekourney/utils"
)

const UTF8PREFIX = 0b10000000

// isAscii returns true if the byte is an ASCII character.
func isAscii(char byte) bool {
	return char&0b10000000 == 0
}

// isUTF8 returns true if the byte is a UTF-8 character.
func isUTF8(char byte) bool {
	return char&UTF8PREFIX != 0
}

// isASCIIAlphaNumeric returns true if the byte is an ASCII alphanumeric
// character.
func isASCIIAlphaNumeric(char byte) bool {
	return (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') ||
		(char >= '0' && char <= '9')
}

// charLen returns the number of bytes in a UTF-8 character.
// Including the first byte.
// If the byte is not valid, it returns 0.
// If the byte is a continuation byte, it returns -1.
func charLen(char byte) int {
	if isAscii(char) {
		return 1
	}

	if char&0b11100000 == 0b11000000 {
		return 2
	}

	if char&0b11110000 == 0b11100000 {
		return 3
	}

	if char&0b11111000 == 0b11110000 {
		return 4
	}

	if char&0b11000000 == 0b10000000 {
		return -1
	}

	return 0
}

// wordSplit returns true if the byte is a word split character.
func wordSplit(char byte) bool {

	// TODO: Check if UTF-8 chars have any characters we want to splkit words
	// with
	if isUTF8(char) {
		return false
	}

	return !isASCIIAlphaNumeric(char)

}

// yieldWord yields a word from the byte slice.
// It takes a yield function, a byte slice, and the start and end indices of
// the word.
// It returns true if the iteration should continue,
// and false if it should stop.
func yieldWord(
	yield func(utils.Word) bool,
	bytes []byte,
	start int,
	end int,
) bool {
	if start != end {
		word := string(bytes[start:end])
		return yield(utils.Word(word))
	}

	return true
}

// WordsIterBytes takes a byte slice and returns an iterator that yields each
// word in the slice.
// Limited UTF-8 support.
func WordsIterBytes(bytes []byte) iter.Seq[utils.Word] {
	word_iter := func(yield func(utils.Word) bool) {

		start := 0
		end := 0

		i := 0

		for i < len(bytes) {
			c := bytes[i]

			if wordSplit(c) {
				if !yieldWord(yield, bytes, start, end) {
					return
				}

				start = i + 1
			}
			// Ship UTF-8 continuation bytes
			len := charLen(c)
			i += len
			end = i
		}

		if !yieldWord(yield, bytes, start, end) {
			return
		}
	}

	return word_iter
}

// WordsIter takes a string and returns an iterator that yields each word in the
// string.
// Limited UTF-8 support.
// Internally, it converts the string to a byte slice and calls WordsIterBytes.
// So if you got a byte slice, use that function directly.
func WordsIter(s string) iter.Seq[utils.Word] {
	bytes := []byte(s)
	return WordsIterBytes(bytes)
}
