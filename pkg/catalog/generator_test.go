package catalog

import (
	"log"
	"testing"

	"github.com/anz-bank/sysl/pkg/loader"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestReplace(t *testing.T) {
	txt := "this_is_some_text"
	remove := "_[^_]*?_text"
	assert.Equal(t, "this_is", Remove(txt, remove))
}

const plantumlService = "http://localhost:8080"

var testFiles = []string{
	"test/App/App/gettestrestqueryparam{id}.svg",
	"test/App/App/gettesturlparamprimitive{id}.svg",
	"test/App/App/gettesturlparamref{id}.svg",
	"test/App/App/foosimple.svg",
	"test/App/App/foo.svg",
	"test/App/App/endpoint.svg",
	"test/App/integration.svg",
	"test/App/README.md",
	"test/App/integration.svg",
	"test/App/primitive/stringsimple.svg",
	"test/README.md",
	"test/integration.svg",
	"test/integrationepa.svg",
}

func TestNewProjectWithLoadSyslModule(t *testing.T) {
	filePath := "../../tests/params.sysl"
	outputDir := "test"
	fs := afero.NewMemMapFs()
	logger := logrus.New()
	m, _, err := loader.LoadSyslModule("", filePath, afero.NewOsFs(), logger)
	if err != nil {
		log.Fatal(err)
	}
	p := NewProject(filePath, plantumlService, "markdown", logger, m, fs, outputDir, false)
	p.Run()
	// Assert the right files exist
	for _, testFile := range testFiles {
		t.Run(testFile, func(t *testing.T) {
			_, err = fs.Open(testFile)
			assert.NoError(t, err)
			_, err := afero.ReadFile(fs, testFile)
			assert.NoError(t, err)
		})
	}
}

func TestNewProjectWithParser(t *testing.T) {
	filePath := "../../tests/params.sysl"
	outputDir := "test"
	fs := afero.NewMemMapFs()
	m, err := parse.NewParser().Parse(filePath, afero.NewOsFs())
	if err != nil {
		log.Fatal(err)
	}
	p := NewProject(filePath, plantumlService, "markdown", logrus.New(), m, fs, outputDir, false)
	p.Run()
	// Assert the right files exist
	for _, testFile := range testFiles {
		t.Run(testFile, func(t *testing.T) {
			_, err = fs.Open(testFile)
			assert.NoError(t, err)
			_, err := afero.ReadFile(fs, testFile)
			assert.NoError(t, err)
		})
	}
}
