package corpus

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViterbiSplit(t *testing.T) {
	assert := assert.New(t)
	dict := GenerateCorpus(mediumSentence())

	s2 := "twoindividuals"
	words := ViterbiSplit(s2, dict)
	assert.Equal([]string{"two", "individuals"}, words)

	s2 = "FederalCourts"
	words = ViterbiSplit(s2, dict)
	assert.Equal([]string{"federal", "courts"}, words)

	s3 := "toreplaceon"
	words = ViterbiSplit(s3, dict)
	assert.Equal([]string{"to", "replace", "on"}, words)
}
