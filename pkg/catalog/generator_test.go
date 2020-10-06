package catalog

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/require"

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

const plantumlService = "http://plantuml.com/plantuml"

type AferoRetriever struct {
	fs afero.Fs
}

func (ar AferoRetriever) Retrieve(resource string) ([]byte, bool, error) {
	f, err := ar.fs.Open(resource)
	if err != nil {
		return nil, false, err
	}
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, false, err
	}
	return bs, false, nil
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
	p := NewProject(filePath, plantumlService, "markdown", logger, m, fs, outputDir)
	p.Run()
	// Assert the right files exist
}

func TestNewProjectWithParser(t *testing.T) {
	filePath := "../../tests/params.sysl"
	outputDir := "test"
	fs := afero.NewMemMapFs()
	m, err := parse.NewParser().Parse(filePath, AferoRetriever{afero.NewOsFs()})
	if err != nil {
		log.Fatal(err)
	}
	p := NewProject(filePath, plantumlService, "markdown", logrus.New(), m, fs, outputDir)
	p.Run()
	// Assert the right files exist
}

func TestGenerateDocsWithRedoc(t *testing.T) {
	filePath := "test/openapi.sysl"
	outputDir := "docs"
	fs := afero.NewMemMapFs()
	logger := logrus.New()
	m, _, err := loader.LoadSyslModule(".", filePath, afero.NewOsFs(), logger)
	if err != nil {
		log.Fatal(err)
	}
	p := NewProject(filePath, plantumlService, "markdown", logger, m, fs, outputDir)
	p.SetOptions(false, "", "")
	p.Retriever = retr{content: map[string]string{"test/simple.yaml": "openapi: \"3.0.0\"\ninfo:\n  version: 1.0.0\n  title: simple\n  license:\n    name: MIT\nservers:\n  - url: http://petstore.swagger.io/v1\n"}}
	p.Run()
	// Assert the right files exist
	testFile := outputDir + "/Simple/simple.redoc.html"
	_, err = fs.Open(testFile)
	assert.NoError(t, err)
	_, err = afero.ReadFile(fs, testFile)
	assert.NoError(t, err)
}

func TestPrettyPackageNmes(t *testing.T) {
	contents := `
AppName:
	@package = "whatever"
whatever:
	@package_alias = "renamed"
`
	outputDir := "docs"
	fs := afero.NewMemMapFs()
	logger := logrus.New()
	m, err := parse.NewParser().ParseString(contents)
	if err != nil {
		log.Fatal(err)
	}
	p := NewProject("", plantumlService, "markdown", logger, m, fs, outputDir)
	p.SetOptions(false, "", "")
	p.Run()
	// Assert the right files exist
	testFile := outputDir + "/renamed/README.md"
	_, err = fs.Open(testFile)
	assert.NoError(t, err)
	_, err = afero.ReadFile(fs, testFile)
	assert.NoError(t, err)
}

func TestNewProjectFromJson(t *testing.T) {

	expectedModule, err := parse.NewParser().Parse("../../tests/rest.sysl", AferoRetriever{afero.NewOsFs()})
	require.Nil(t, err)

	expected := NewProject("", "", "", logrus.New(), expectedModule, afero.NewMemMapFs(), "/")

	file, err := ioutil.ReadFile("../../tests/rest.json")

	require.Nil(t, err)

	actual := NewProjectFromJson("", "", "", logrus.New(), file, afero.NewMemMapFs(), "/")
	require.Nil(t, err)

	require.Equal(t, expected.String(), actual.String())

}

/* String returns a string of all of the non pointer fields; mainly to be used with require.Equal*/
func (p *Generator) String() string {
	return fmt.Sprint(
		p.FilesToCreate,
		p.MermaidFilesToCreate,
		p.RedocFilesToCreate,
		p.SourceFileName,
		p.ProjectTitle,
		p.ImageDest,
		p.Format,
		p.OutputFileName,
		p.PlantumlService,
		p.StartTemplateIndex,
		p.FilterPackage,
		p.CustomTemplate,
		p.LiveReload,
		p.DisableCss,
		p.Fs,
		p.errs,
		p.CurrentDir,
		p.TempDir,
		p.Title,
		p.OutputDir,
		p.Links,
		p.Server,
	)
}
