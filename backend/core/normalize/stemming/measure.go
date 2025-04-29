package stemming

type measure uint

// clacMeasure calculates the measure of a word, as defined in
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

	if firstType == consonant {
		count--
	}

	if lastType == vowel {
		count--
	}

	return measure(count / 2)
}
