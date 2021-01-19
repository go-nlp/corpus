package corpus

import (
	"math"
	"strings"
	"unicode/utf8"
)

// ViterbiSplit is a Viterbi algorithm for splitting words given a corpus
func ViterbiSplit(input string, c *Corpus) []string {
	s := strings.ToLower(input)
	probabilities := []float64{1.0}
	lasts := []int{0}

	runes := []int{}
	for i := range s {
		runes = append(runes, i)
	}
	runes = append(runes, len(s)+1)

	for i := range s {
		probs := make([]float64, 0)
		ls := make([]int, 0)

		// m := maxInt(0, i-c.maxWordLength)

		for j, r := range runes {
			if r > i {
				break
			}

			p, ok := c.WordProb(s[r : i+1])
			if !ok {
				// http://stackoverflow.com/questions/195010/how-can-i-split-multiple-joined-words#comment48879458_481773
				p = (math.Log(float64(1)/float64(c.totalFreq)) - float64(c.maxWordLength) - float64(1)) * float64(i-r) // note it should be i-r not j-i as per the SO post
			}
			prob := probabilities[j] * p

			probs = append(probs, prob)
			ls = append(ls, r)
		}

		maxProb := -math.SmallestNonzeroFloat64
		maxK := -1 << 63
		for j, p := range probs {
			if p > maxProb {
				maxProb = p
				maxK = ls[j]
			}
		}
		probabilities = append(probabilities, maxProb)
		lasts = append(lasts, maxK)
	}

	words := make([]string, 0)
	i := utf8.RuneCountInString(s)

	for i > 0 {
		start := lasts[i]
		words = append(words, s[start:i])
		i = start
	}

	// reverse it
	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i]
	}

	return words
}
