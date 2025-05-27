package stemming

// Source: https://snowballstem.org/algorithms/porter/stemmer.html

// measure is a value of a word that describes the number of pairs of
// successive vowels and consonants in the word.
// As described in the article above
// If the word starts on a consonant, skip the first consonant sequence
// and if it ends on a vowel, skip the last vowel sequence
type measure uint

// calcMeasure calculates the measure of a word, as defined in
// the Porter stemming algorithm
func calcMeasure(word string) measure {

	count := 0
	first := true
	firstType := charType(-1)
	lastType := charType(-1)

	for charType, prev := range charTypeIter(word) {

		if first {
			count++
			firstType = charType
			lastType = charType
			first = false
			continue
		}

		if prev != charType {
			count++
			lastType = charType
		}

	}

	if firstType == __CONSONANT__ {
		count--
	}

	if lastType == __VOWEL__ {
		count--
	}

	return measure(count / 2)
}
