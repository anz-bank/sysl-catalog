package catalog

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAndWriteRedoc(t *testing.T) {
	t.Parallel()
	fs := afero.NewMemMapFs()
	fileName := "redoc.html"
	err := GenerateAndWriteRedoc(fs, fileName, "")
	assert.NoError(t, err)
	_, err = fs.Open(fileName)
	assert.NoError(t, err)
}
