package catalog

import (
	"strings"
	"testing"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestServiceMetadata(t *testing.T) {
	t.Parallel()

	m, err := parse.NewParser().ParseString(`
foo:
    @Repo.URL = "1"
    @Owner.Email = "2"
    @Owner.Slack = "3"
    @Server.Prod.URL = "4"
    @Server.UAT.URL = "5"
	@Lifecycle = "6"

bar:
    @Lifecycle = "1"
    @Owner.Slack = "2"
    @Owner.Email = "3"

boo:
	@Lifecycle = "1"
	@Random = "2"
	@Repo.URL = "3"

ree:
	@lifecycle = "1"
	@oWnEr.EMaIL = "2"
    @REPO.URL = "3"
`,
	)
	require.NoError(t, err)
	createRes := func(s []string) string {
		return strings.Join(s, "\n\n") + "\n\n"
	}
	results := map[string]string{
		"foo": createRes([]string{
			"Repo.URL: 1",
			"Owner.Email: 2",
			"Owner.Slack: 3",
			"Server.Prod.URL: 4",
			"Server.UAT.URL: 5",
			"Lifecycle: 6",
		}),
		"bar": createRes([]string{
			"Owner.Email: 3",
			"Owner.Slack: 2",
			"Lifecycle: 1",
		}),
		"boo": createRes([]string{
			"Repo.URL: 3",
			"Lifecycle: 1",
		}),
		"ree": createRes([]string{
			"Repo.URL: 3",
			"Owner.Email: 2",
			"Lifecycle: 1",
		}),
	}

	for app, exp := range results {
		assert.Equal(t, exp, ServiceMetadata(m.GetApps()[app]))
	}
}
