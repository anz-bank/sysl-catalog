package catalog

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildRedoc(t *testing.T) {
	url := "http://petstore.swagger.io/v2/swagger.json"
	expectedTag := fmt.Sprintf("<redoc spec-url='https://cors-anywhere.herokuapp.com/%s'></redoc>", url)
	redoc := BuildRedoc(url)
	assert.True(t, strings.Contains(string(redoc), expectedTag))
}
