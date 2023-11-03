package corpus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorpus(t *testing.T) {
	assert := assert.New(t)
	dict := New()
	assert.Equal(0, dict.WordFreq("hello")) // frequency of a word not in dict ould have to be 0
	assert.Equal(0, dict.IDFreq(3))         // ditto

	id := dict.Add("hello")

	assert.Equal(3, id)
	assert.Equal([]string{"", "-UNKNOWN-", "-ROOT-", "hello"}, dict.Words)
	assert.Equal(map[string]int{"": 0, "-UNKNOWN-": 1, "-ROOT-": 2, "hello": 3}, dict.Ids)
	assert.Equal(4, dict.Size())

	id2, ok := dict.Id("hello")
	if !ok {
		t.Errorf("The ID of null should be  0")
	}
	assert.Equal(id, id2)

	word, ok := dict.Word(3)
	if !ok {
		t.Errorf("Expected word of ID 3 to be found")
	}
	assert.Equal("hello", word)

	dict.Add(word)
	assert.Equal(2, dict.WordFreq(word))
	assert.Equal(2, dict.IDFreq(3))
	assert.Equal(5, dict.TotalFreq())
	assert.Equal(5, dict.MaxWordLength())

	prob, ok := dict.WordProb(word)
	if !ok {
		t.Errorf("Expected a probability")
	}
	assert.Equal(0.4, prob)
	// t.Logf("%q: %v", word, dict.WordProb(word))
}

func TestCorpus_Merge(t *testing.T) {
	assert := assert.New(t)

	dict := New()
	id := dict.Add("hello")
	dict.Frequencies[id] += 4 // freq for "hello" is 5
	dict.TotalWordFreq += 4

	other := New()
	id = other.Add("hello")
	other.Frequencies[id] += 2 // freq for "hello" is 3
	other.TotalWordFreq += 2
	id = other.Add("world")
	other.Frequencies[id] += 1
	other.TotalWordFreq += 1

	dict.Merge(other)

	assert.Equal(8, dict.WordFreq("hello"))
	assert.Equal(2, dict.WordFreq("world"))
}

func TestCorpus_Replace(t *testing.T) {
	dict := New()
	dict.Add("Hello")
	if err := dict.Replace("Hello", "Bye"); err != nil {
		t.Fatal("Replacement caused an error")
	}

	helloID, ok := dict.Id("Hello")
	assert.True(t, ok, "Hello should have an ID")
	byeID, ok := dict.Id("Bye")
	assert.True(t, ok, "Bye should have an ID")
	assert.Equal(t, helloID, byeID)

	// do it a second time and you will get an errorr
	if err := dict.Replace("Hello", "Bye"); err == nil {
		t.Errorf("Expected an error when replacing a word with a known ID")
	}

	if err := dict.Replace("Foo", "bar"); err == nil {
		t.Errorf("Expected an error when replacing an unknown word")
	}

}

func TestCorpus_ReplaceWord(t *testing.T) {
	dict := New()
	helloID := dict.Add("Hello")
	if err := dict.ReplaceWord(helloID, "Bye"); err != nil {
		t.Fatal("Replacement caused an error")
	}

	helloID, ok := dict.Id("Hello")
	assert.True(t, ok, "Hello should have an ID")
	byeID, ok := dict.Id("Bye")
	assert.True(t, ok, "Bye should have an ID")
	assert.Equal(t, helloID, byeID)

	// do it a second time and you will get an errorr
	if err := dict.ReplaceWord(helloID, "Bye"); err == nil {
		t.Errorf("Expected an error when replacing a word with a known ID")
	}

	if err := dict.ReplaceWord(100, "bar"); err == nil {
		t.Errorf("Expected an error when replacing an unknown word")
	}

}
