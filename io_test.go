package corpus

import (
	"bytes"
	"encoding/gob"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCorpusGob(t *testing.T) {
	buf := new(bytes.Buffer)

	c := New()
	c.Add("Hello")
	c.Add("World")

	helloID, _ := c.Id("Hello")
	worldID, _ := c.Id("World")

	encoder := gob.NewEncoder(buf)
	decoder := gob.NewDecoder(buf)

	if err := encoder.Encode(c); err != nil {
		t.Fatal(err)
	}

	c2 := New()
	if err := decoder.Decode(c2); err != nil {
		t.Fatal(err)
	}

	if hid, ok := c2.Id("Hello"); !ok || (ok && hid != helloID) {
		t.Errorf("\"Hello\" not found after decoding.")
	}

	if wid, ok := c2.Id("World"); !ok || (ok && wid != worldID) {
		t.Errorf("\"World\" not found after decoding.")
	}
}

func TestCorpusToDict(t *testing.T) {
	assert := assert.New(t)
	c, _ := Construct(WithWords([]string{"World", "Hello", "World"}))

	d := ToDict(c)
	c2, err := Construct(FromDict(d))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(c.Words, c2.Words, "Expected Words to be the same")
	assert.Equal(c.Ids, c2.Ids, "Expected IDs to be the same")
	assert.NotEqual(c.Frequencies, c2.Frequencies, "Expected Frequencies to not be the same")
	assert.Equal(c.MaxID, c2.MaxID, "Expected maxID to be the same")
	assert.NotEqual(c.TotalWordFreq, c2.TotalWordFreq, "Expected TotalWordFreq to be different.")
	assert.Equal(c.MaxWordLength_, c2.MaxWordLength_, "Expected MaxWordLength_ to be the same")
}

func TestCorpusToDictWithFreq(t *testing.T) {
	assert := assert.New(t)
	c, _ := Construct(WithWords([]string{"World", "Hello", "World"}))

	d := ToDictWithFreq(c)
	c2, err := Construct(FromDictWithFreq(d))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(c, c2)
}

func TestLoadOneGram(t *testing.T) {
	assert := assert.New(t)
	r := strings.NewReader(sample1Gram)

	c := New()
	err := c.LoadOneGram(r)
	assert.Nil(err)
	assert.Equal(10, c.Size())

	id, ok := c.Id("for")
	if !ok {
		t.Errorf("Expected \"for\" to be in corpus after loading one gram file")
	}
	assert.Equal(int(c.MaxID-1), id)
}

func TestFromTextCorpus(t *testing.T) {
	f, err := os.Open("testdata/corpus_en.txt")
	require.NoError(t, err)

	// set up a dumb english tokenizer for a more "realistic" tokenization than just whitespace.
	pattern := regexp.MustCompile(`'s|'t|'re|'ve|'m|'ll|'d| ?\pL+| ?\pN+| ?[^\s\pL\pN]+|\s+`)
	dumbTokenizer := func(a string) []string {
		strs := pattern.FindAllString(a, -1)
		for i := range strs {
			strs[i] = strings.Trim(strs[i], "\r\n ")
		}
		return strs
	}
	dumbNormalizer := func(a string) string { return strings.ToLower(a) }
	c, err := FromTextCorpus(f, dumbTokenizer, dumbNormalizer)
	require.NoError(t, err)

	aliceID, ok := c.Id("alice")
	assert.True(t, ok)
	assert.Equal(t, 128, aliceID)

	freq := c.IDFreq(aliceID)
	assert.Equal(t, 399, freq)

	// FOR DEBUG PURPOSES
	// g, err := os.OpenFile("testdata/tmp", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	// require.NoError(t, err)
	// for i, w := range c.Words {
	// 	fmt.Fprintf(g, "%v %d\n", w, c.Frequencies[i])
	// }
	// g.Close()
}
