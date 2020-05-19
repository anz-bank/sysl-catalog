package catalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplace(t *testing.T) {
	txt := "this_is_some_text"
	remove := "_[^_]*?_text"
	assert.Equal(t, "this_is", Remove(txt, remove))
}
