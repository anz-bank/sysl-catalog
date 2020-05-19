package catalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplace(t *testing.T) {
	txt := "this_is_some_text"
	remove := "_[^_]*?_text"
	assert.Equal(t, Remove(txt, remove), "this_is_some")
}
