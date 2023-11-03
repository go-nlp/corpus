package corpus

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"io"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// sortutil is a utility struct meant to sort words based on IDs
type sortutil struct {
	words []string
	ids   []int
	freqs []int
}

func (s *sortutil) Len() int           { return len(s.words) }
func (s *sortutil) Less(i, j int) bool { return s.ids[i] < s.ids[j] }
func (s *sortutil) Swap(i, j int) {
	s.words[i], s.words[j] = s.words[j], s.words[i]
	s.ids[i], s.ids[j] = s.ids[j], s.ids[i]
	if len(s.freqs) > 0 {
		s.freqs[i], s.freqs[j] = s.freqs[j], s.freqs[i]
	}
}

// ToDictWithFreq returns a simple marshalable type. Conceptually it's a JSON object with the Words as the keys. The values are a pair - ID and Freq.
func ToDictWithFreq(c *Corpus) map[string]struct{ ID, Freq int } {
	retVal := make(map[string]struct{ ID, Freq int })
	for i, w := range c.Words {
		retVal[w] = struct{ ID, Freq int }{i, c.Frequencies[i]}
	}
	return retVal
}

// ToDict returns a marshalable dict. It returns a copy of the ID mapping.
func ToDict(c *Corpus) map[string]int {
	retVal := make(map[string]int)
	for k, v := range c.Ids {
		retVal[k] = v
	}
	return retVal
}

// GobEncode implements GobEncoder for *Corpus
func (c *Corpus) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(c.Words); err != nil {
		return nil, err
	}

	if err := encoder.Encode(c.Ids); err != nil {
		return nil, err
	}

	if err := encoder.Encode(c.Frequencies); err != nil {
		return nil, err
	}

	if err := encoder.Encode(c.MaxID); err != nil {
		return nil, err
	}

	if err := encoder.Encode(c.TotalWordFreq); err != nil {
		return nil, err
	}

	if err := encoder.Encode(c.MaxWordLength_); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GobDecode implements GobDecoder for *Corpus
func (c *Corpus) GobDecode(buf []byte) error {
	b := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(b)

	if err := decoder.Decode(&c.Words); err != nil {
		return err
	}

	if err := decoder.Decode(&c.Ids); err != nil {
		return err
	}

	if err := decoder.Decode(&c.Frequencies); err != nil {
		return err
	}

	if err := decoder.Decode(&c.MaxID); err != nil {
		return err
	}

	if err := decoder.Decode(&c.TotalWordFreq); err != nil {
		return err
	}

	if err := decoder.Decode(&c.MaxWordLength_); err != nil {
		return err
	}

	return nil
}

// LoadOneGram loads a 1_gram.txt file, which is a tab separated file which lists the frequency counts of Words. Example:
// 		the	23135851162
// 		of	13151942776
// 		and	12997637966
// 		to	12136980858
// 		a	9081174698
// 		in	8469404971
// 		for	5933321709
func (c *Corpus) LoadOneGram(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		splits := strings.Split(line, "\t")

		if len(splits) == 0 {
			break
		}

		word := splits[0] // TODO: normalize
		count, err := strconv.Atoi(splits[1])
		if err != nil {
			return err
		}

		id := c.Add(word)
		c.Frequencies[id] = count
		c.TotalWordFreq--
		c.TotalWordFreq += count

		wc := len([]rune(word))
		if wc > c.MaxWordLength_ {
			c.MaxWordLength_ = wc
		}
	}
	return nil
}

// FromTextCorpus is a utility function to take in a text file, and return a Corpus.
func FromTextCorpus(r io.Reader, tokenizer func(a string) []string, normalizer func(a string) string) (*Corpus, error) {
	if tokenizer == nil {
		tokenizer = func(a string) []string {
			return strings.Split(strings.Trim(a, "\r\n "), " ")
		}
	}
	if normalizer == nil {
		normalizer = func(a string) string { return a }
	}

	var words []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		words = append(words, tokenizer(normalizer(line))...)
	}
	if err := scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "Unable to read from text corpus")
	}

	return Construct(WithWords(words))
}
