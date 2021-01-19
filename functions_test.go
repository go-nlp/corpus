package corpus

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestViterbiSplit(t *testing.T) {
	assert := assert.New(t)
	f, err := os.Open("testdata/corpus_en.txt")
	require.NoError(t, err)

	dict, err := FromTextCorpus(f, nil, func(a string) string { return strings.ToLower(a) })
	require.NoError(t, err)

	s2 := "whiterabbit"
	words := ViterbiSplit(s2, dict)
	assert.Equal([]string{"white", "rabbit"}, words)

	/*
		   // FAILING TEST
			s2 = "curiouserandcuriouser"
			words = ViterbiSplit(s2, dict)
			assert.Equal([]string{"curiouser", "and", "curiouser"}, words)
	*/

	s3 := "thebestwaytoexplainitistodoit"
	words = ViterbiSplit(s3, dict)
	assert.Equal([]string{"the", "best", "way", "to", "explain", "it", "is", "to", "do", "it"}, words)
}
