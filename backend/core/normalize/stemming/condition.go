package stemming

// condition decides if a substitution is valid
type condition func(stem string) bool

// not negates the condition
func not(c condition) condition {
	return func(stem string) bool {
		return !c(stem)
	}
}

// hasVowel checks if the stem contains a vowel
func hasVowel(stem string) bool {
	return containsVowel(stem)
}

// measureMin checks if the measure is greater than or equal to the given value
func measureMin(value uint) condition {
	return func(stem string) bool {
		return calcMeasure(stem) >= measure(value)
	}
}

// measureIs checks if the measure is equal to the given value
func measureIs(value uint) condition {
	return func(stem string) bool {
		return calcMeasure(stem) == measure(value)
	}
}

// doubleConsonant checks if the last
// two characters of the stem are the same consonant
func dubbelConsonant(stem *stem) bool {

	str := stem.str

	if len(str) < 2 {
		return false
	}

	last := str[len(str)-2:]

	return !isVowel(last[0]) && last[0] == last[1]
}

// endsWith checks if the stem ends with any of the given characters
func endsWith(chars ...byte) condition {
	return func(stem string) bool {

		l := len(stem)

		if l == 0 {
			return false
		}

		for _, char := range chars {
			if stem[l-1] == char {
				return true
			}
		}

		return false
	}

}

// cvc checks if the stem ends with a consonant-vowel-consonant pattern
// the secnd consonant must not be a w, x or y
func cvc(stem string) bool {

	l := len(stem)

	if l < 3 {
		return false
	}

	last := stem[l-3:]

	return !isVowel(last[0]) && isVowel(last[1]) && !isVowel(last[2]) &&
		last[2] != 'w' && last[2] != 'x' && last[2] != 'y'
}
