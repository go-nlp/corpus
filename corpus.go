package corpus

import (
	"sync/atomic"
	"unicode/utf8"

	"github.com/pkg/errors"
)

// Corpus is a data structure holding the relevant metadata and information for a corpus of text.
// It serves as vocabulary with ID for lookup. This is very useful as neural networks rely on the IDs rather than the text themselves
type Corpus struct {
	Words       []string `json:"words"`
	Frequencies []int    `json:"frequencies"`

	Ids map[string]int `json:"ids"`

	// atomic read and write plz
	MaxID          int64 `json:"max_id"`
	TotalWordFreq  int   `json:"total_word_freq"`
	MaxWordLength_ int   `json:"max_word_length"`
}

// New creates a new *Corpus
func New() *Corpus {
	c := &Corpus{
		Words:       make([]string, 0),
		Frequencies: make([]int, 0),
		Ids:         make(map[string]int),
	}

	// add some default Words
	c.Add("") // aka NULL - when there are no Words
	c.Add("-UNKNOWN-")
	c.Add("-ROOT-")
	c.MaxWordLength_ = 0 // specials don't have lengths

	return c
}

// Construct creates a Corpus given the construction options. This allows for more flexibility
func Construct(opts ...ConsOpt) (*Corpus, error) {
	c := new(Corpus)

	// checks
	if c.Words == nil {
		c.Words = make([]string, 0)
	}
	if c.Frequencies == nil {
		c.Frequencies = make([]int, 0)
	}
	if c.Ids == nil {
		c.Ids = make(map[string]int)
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Id returns the ID of a word and whether or not it was found in the corpus
func (c *Corpus) Id(word string) (int, bool) {
	id, ok := c.Ids[word]
	return id, ok
}

// Word returns the word given the ID, and whether or not it was found in the corpus
func (c *Corpus) Word(id int) (string, bool) {
	size := atomic.LoadInt64(&c.MaxID)
	maxid := int(size)

	if id >= maxid {
		return "", false
	}
	return c.Words[id], true
}

// Add adds a word to the corpus and returns its ID. If a word was previously in the corpus, it merely updates the frequency count and returns the ID
func (c *Corpus) Add(word string) int {
	if id, ok := c.Ids[word]; ok {
		c.Frequencies[id]++
		c.TotalWordFreq++
		return id
	}

	id := atomic.AddInt64(&c.MaxID, 1)
	c.Ids[word] = int(id - 1)
	c.Words = append(c.Words, word)
	c.Frequencies = append(c.Frequencies, 1)
	c.TotalWordFreq++

	runeCount := utf8.RuneCountInString(word)
	if runeCount > c.MaxWordLength_ {
		c.MaxWordLength_ = runeCount
	}

	return int(id - 1)
}

// Size returns the size of the corpus.
func (c *Corpus) Size() int {
	size := atomic.LoadInt64(&c.MaxID)
	return int(size)
}

// WordFreq returns the frequency of the word. If the word wasn't in the corpus, it returns 0.
func (c *Corpus) WordFreq(word string) int {
	id, ok := c.Ids[word]
	if !ok {
		return 0
	}

	return c.Frequencies[id]
}

// IDFreq returns the frequency of a word given an ID. If the word isn't in the corpus it returns 0.
func (c *Corpus) IDFreq(id int) int {
	size := atomic.LoadInt64(&c.MaxID)
	maxid := int(size)

	if id >= maxid {
		return 0
	}
	return c.Frequencies[id]
}

// TotalFreq returns the total number of Words ever seen by the corpus. This number includes the count of repeat Words.
func (c *Corpus) TotalFreq() int {
	return c.TotalWordFreq
}

// MaxWordLength returns the length of the longest known word in the corpus.
func (c *Corpus) MaxWordLength() int {
	return c.MaxWordLength_
}

// WordProb returns the probability of a word appearing in the corpus.
func (c *Corpus) WordProb(word string) (float64, bool) {
	id, ok := c.Id(word)
	if !ok {
		return 0, false
	}

	count := c.Frequencies[id]
	return float64(count) / float64(c.TotalWordFreq), true

}

// Merge combines two corpuses. The receiver is the one that is mutated.
func (c *Corpus) Merge(other *Corpus) {
	for i, word := range other.Words {
		freq := other.Frequencies[i]
		if id, ok := c.Ids[word]; ok {
			c.Frequencies[id] += freq
			c.TotalWordFreq += freq
		} else {
			id := c.Add(word)
			c.Frequencies[id] += freq - 1
			c.TotalWordFreq += freq - 1
		}
	}
}

// Replace replaces the content of a word. The old reference remains.
//
// e.g: c.Replace("foo", "bar")
// c.Id("foo") will still return a ID. The ID will be the same as c.Id("bar")
func (c *Corpus) Replace(a, with string) error {
	old, ok := c.Ids[a]
	if !ok {
		return errors.Errorf("Cannot replace %q with %q. %q is not found", a, with, a)
	}
	if _, ok := c.Ids[with]; ok {
		return errors.Errorf("Cannot replace %q with %q. %q exists in the corpus", a, with, with)
	}
	c.Words[old] = with
	c.Ids[with] = old
	return nil

}

// ReplaceWord replaces the word associated with the given ID. The old reference remains.
func (c *Corpus) ReplaceWord(id int, with string) error {
	if id >= len(c.Words) {
		return errors.Errorf("Cannot replace word with ID %d. Out of bounds.", id)
	}
	if _, ok := c.Ids[with]; ok {
		return errors.Errorf("Cannot replace word with ID %d with %q. %q exists in the corpus", id, with, with)
	}
	c.Words[id] = with
	c.Ids[with] = id
	return nil
}
