package corpus

import (
	"sort"
	"sync/atomic"
	"unicode/utf8"

	"github.com/pkg/errors"
	"github.com/xtgo/set"
)

// ConsOpt is a construction option for manual creation of a Corpus
type ConsOpt func(c *Corpus) error

// WithWords creates a corpus from a word list. It may have repeated Words
func WithWords(a []string) ConsOpt {
	f := func(c *Corpus) error {
		s := set.Strings(a)
		c.Words = s
		c.Frequencies = make([]int, len(s))

		ids := make(map[string]int)
		maxID := len(s)

		var totalFreq, maxWL int
		// NOTE: here we're iterating over the set of Words
		for i, w := range s {
			runeCount := utf8.RuneCountInString(w)
			if runeCount > c.MaxWordLength_ {
				maxWL = runeCount
			}

			ids[w] = i
		}

		// NOTE: here we're iterating over the original word list.
		for _, w := range a {
			c.Frequencies[ids[w]]++
			totalFreq++
		}

		c.Ids = ids
		atomic.AddInt64(&c.MaxID, int64(maxID))
		c.TotalWordFreq = totalFreq
		c.MaxWordLength_ = maxWL
		return nil
	}
	return f
}

// WithOrderedWords creates a Corpus with the given word order
func WithOrderedWords(a []string) ConsOpt {
	f := func(c *Corpus) error {
		s := a
		c.Words = s
		c.Frequencies = make([]int, len(s))
		for i := range c.Frequencies {
			c.Frequencies[i] = 1
		}

		ids := make(map[string]int)
		maxID := len(s)
		totalFreq := len(s)
		var maxWL int
		for i, w := range a {
			runeCount := utf8.RuneCountInString(w)
			if runeCount > c.MaxWordLength_ {
				maxWL = runeCount
			}
			ids[w] = i
		}

		c.Ids = ids
		atomic.AddInt64(&c.MaxID, int64(maxID))
		c.TotalWordFreq = totalFreq
		c.MaxWordLength_ = maxWL
		return nil
	}
	return f
}

// WithSize preallocates all the things in Corpus
func WithSize(size int) ConsOpt {
	return func(c *Corpus) error {
		c.Words = make([]string, 0, size)
		c.Frequencies = make([]int, 0, size)
		return nil
	}
}

// FromDict is a construction option to take a map[string]int where the int represents the word ID.
// This is useful for constructing corpuses from foreign sources where the ID mappings are important
func FromDict(d map[string]int) ConsOpt {
	return func(c *Corpus) error {
		var a sortutil
		for k, v := range d {
			a.words = append(a.words, k)
			a.ids = append(a.ids, v)
		}
		sort.Sort(&a)
		c.Ids = make(map[string]int)
		for i, w := range a.words {
			if i != a.ids[i] {
				return errors.Errorf("Unmarshaling error. Expected %dth ID to be %d. Got %d instead. Perhaps something went wrong during sorting? SLYTHERIN IT IS!", i, i, a.ids[i])
			}
			c.Words = append(c.Words, w)
			c.Frequencies = append(c.Frequencies, 1)
			c.Ids[w] = i

			c.TotalWordFreq++
			runeCount := utf8.RuneCountInString(w)
			if runeCount > c.MaxWordLength_ {
				c.MaxWordLength_ = runeCount
			}
		}
		c.MaxID = int64(len(a.words))
		return nil
	}

}

// FromDictWithFreq is like FromDict, but also has a frequency.
func FromDictWithFreq(d map[string]struct{ ID, Freq int }) ConsOpt {
	return func(c *Corpus) error {
		var a sortutil
		for k, v := range d {
			a.words = append(a.words, k)
			a.ids = append(a.ids, v.ID)
			a.freqs = append(a.freqs, v.Freq)
		}
		sort.Sort(&a)
		c.Ids = make(map[string]int)
		for i, w := range a.words {
			if i != a.ids[i] {
				return errors.Errorf("Unmarshaling error. Expected %dth ID to be %d. Got %d instead. Perhaps something went wrong during sorting? SLYTHERIN IT IS!", i, i, a.ids[i])
			}
			c.Words = append(c.Words, w)
			c.Frequencies = append(c.Frequencies, a.freqs[i])
			c.Ids[w] = i

			c.TotalWordFreq += a.freqs[i]
			runeCount := utf8.RuneCountInString(w)
			if runeCount > c.MaxWordLength_ {
				c.MaxWordLength_ = runeCount
			}
		}
		c.MaxID = int64(len(a.words))
		return nil
	}
}
