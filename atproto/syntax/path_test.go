package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepoPath(t *testing.T) {
	assert := assert.New(t)

	testValid := [][]string{
		{"app.gndr.feed.post/asdf", "app.gndr.feed.post", "asdf"},
	}

	testErr := []string{
		"",
		"/",
		"/app.gndr.feed.post/asdf",
		"/asdf",
		"./app.gndr.feed.post",
		"blob/asdf",
		"app.gndr.feed.post/",
		"app.gndr.feed.post/.",
		"app.gndr.feed.post/!",
	}

	for _, parts := range testValid {
		nsid, rkey, err := ParseRepoPath(parts[0])
		assert.NoError(err)
		assert.Equal(parts[1], nsid.String())
		assert.Equal(parts[2], rkey.String())
	}

	for _, raw := range testErr {
		nsid, rkey, err := ParseRepoPath(raw)
		assert.Error(err)
		assert.Equal("", nsid.String())
		assert.Equal("", rkey.String())
	}
}
