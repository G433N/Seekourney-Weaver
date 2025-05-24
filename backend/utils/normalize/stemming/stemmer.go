package stemming

// Source: https://snowballstem.org/algorithms/porter/stemmer.html
// TODO: Explain the algorithm, maybe move from the report

import (
	"seekourney/utils"
	"seekourney/utils/words"
	"strings"
)

type Word = utils.Word

// stem represents a "stem" of a word
type stem struct {
	str string
}

// ruleRes is a type that represents the result of a rule
type ruleRes int

const (
	// __MISS__ is used to indicate that the rule was not matched
	__MISS__ ruleRes = iota
	// __MATCHED__ is used to indicate that the rule was __MATCHED__,
	// but not applied
	__MATCHED__
	// __CHANGED__ is used to indicate that the stem was __CHANGED__
	__CHANGED__
)

// ruleFunc is a type that represents a rule function
type ruleFunc func(stem *stem) ruleRes

// Stem stems alphabetical ascci words and lowercases any other string
//
// It uses the Porter stemming algorithm, as described in the article:
// https://snowballstem.org/algorithms/porter/stemmer.html
func Stem(word Word) Word {

	// Only stem ASCCI words
	for _, b := range []byte(word) {
		if !words.IsASCIIAlpha(b) {
			return utils.Word(strings.ToLower(string(word)))
		}
	}

	stem := wordIntoStemData(word)

	// Step 1a (As described in the article, refrenced at the top of this file)
	rules(stem,
		rule("sses", "ss"),
		rule("ies", "i"),
		rule("ss", "ss"),
		rule("s", ""),
	)

	specalCase := false

	// Step 1b
	if !rules(stem,
		rule("eed", "ee", measureMin(1))) {

		specalCase = rules(stem,
			rule("ing", "", hasVowel),
			rule("ed", "", hasVowel))
	}

	// Step 1b Special case
	if specalCase {
		if !rules(stem,
			rule("at", "ate"),
			rule("bl", "ble"),
			rule("iz", "ize")) {

			if dubbelConsonant(stem) && !endsWith('l', 's', 'z')(stem.str) {
				removeLastChar(stem)
			} else {
				rules(stem,
					rule("", "e", measureIs(1), cvc))
			}
		}
	}

	// Step 1c
	rules(stem,
		rule("y", "i", hasVowel))

	// Step 2
	rules(stem,
		rule("ational", "ate", measureMin(1)),
		rule("tional", "tion", measureMin(1)),
		rule("enci", "ence", measureMin(1)),
		rule("anci", "ance", measureMin(1)),
		rule("izer", "ize", measureMin(1)),
		rule("abli", "able", measureMin(1)),
		rule("alli", "al", measureMin(1)),
		rule("entli", "ent", measureMin(1)),
		rule("eli", "e", measureMin(1)),
		rule("ousli", "ous", measureMin(1)),
		rule("ization", "ize", measureMin(1)),
		rule("ation", "ate", measureMin(1)),
		rule("ator", "ate", measureMin(1)),
		rule("alism", "al", measureMin(1)),
		rule("iveness", "ive", measureMin(1)),
		rule("fulness", "ful", measureMin(1)),
		rule("ousness", "ous", measureMin(1)),
		rule("aliti", "al", measureMin(1)),
		rule("iviti", "ive", measureMin(1)),
		rule("biliti", "ble", measureMin(1)))

	// Step 3
	rules(stem,
		rule("icate", "ic", measureMin(1)),
		rule("ative", "", measureMin(1)),
		rule("alize", "al", measureMin(1)),
		rule("iciti", "ic", measureMin(1)),
		rule("ical", "ic", measureMin(1)),
		rule("ful", "", measureMin(1)),
		rule("ness", "", measureMin(1)))

	// Step 4
	rules(stem,
		rule("al", "", measureMin(2)),
		rule("ance", "", measureMin(2)),
		rule("ence", "", measureMin(2)),
		rule("er", "", measureMin(2)),
		rule("ic", "", measureMin(2)),
		rule("able", "", measureMin(2)),
		rule("ible", "", measureMin(2)),
		rule("ant", "", measureMin(2)),
		rule("ement", "", measureMin(2)),
		rule("ment", "", measureMin(2)),
		rule("ent", "", measureMin(2)),
		rule("ion", "", measureMin(2), endsWith('s', 't')),
		rule("ou", "", measureMin(2)),
		rule("ism", "", measureMin(2)),
		rule("ate", "", measureMin(2)),
		rule("iti", "", measureMin(2)),
		rule("ous", "", measureMin(2)),
		rule("ive", "", measureMin(2)),
		rule("es", "", measureMin(2)))

	// Step 5a
	rules(stem,
		rule("e", "", measureMin(2)),
		rule("e", "", measureIs(1), not(cvc)))

	// Step 5b
	if measureMin(1)(stem.str) &&
		dubbelConsonant(stem) &&
		endsWith('l')(stem.str) {
		removeLastChar(stem)
	}

	// TODO: Map of missed stems example linearl -> linear (from linearly)
	// TODO: Figure out if we want to store the orignal word somewhere

	return Word(stem.str)
}

// wordIntoStemData converts a word into a stem data
func wordIntoStemData(word Word) *stem {
	str := string(word)
	return &stem{
		str: strings.ToLower(str),
	}
}

// rules applies a list of rules to the stem data, returns true if changed
func rules(stem *stem, rules ...ruleFunc) bool {

	for _, rule := range rules {
		res := rule(stem)
		if res == __CHANGED__ {
			return true
		}

		// If the rule was not applied, but it matched, we need to
		// return false, so we can stop applying rules
		if res == __MATCHED__ {
			return false
		}
	}

	return false
}

// rule creates a rule function that applies the given suffix and ending
func rule(suffix string, ending string, conds ...condition) ruleFunc {
	return func(stem *stem) ruleRes {
		return apply(stem, suffix, ending, conds...)
	}
}

// apply tries applies a rule to the stem data
func apply(
	stem *stem,
	suffix string,
	ending string,
	conds ...condition) ruleRes {

	res, ok := strings.CutSuffix(stem.str, suffix)

	if !ok {
		return __MISS__
	}

	for _, cond := range conds {
		if !cond(res) {
			return __MATCHED__
		}
	}

	stem.str = res + ending

	return __CHANGED__
}

// removeLastChar removes the last character from the stem
func removeLastChar(stem *stem) {
	stem.str = stem.str[:len(stem.str)-1]
}
