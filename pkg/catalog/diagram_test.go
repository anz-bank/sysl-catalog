package catalog

import (
	"log"
	"testing"

	"github.com/anz-bank/sysl/pkg/loader"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func TestNewProjectWithLoadSyslModuleQueryParamPrimitive(t *testing.T) {
	filePath := "../../tests/rest.sysl"
	outputDir := "test"
	fs := afero.NewOsFs()
	logger := logrus.New()
	m, _, err := loader.LoadSyslModule("", filePath, fs, logger)
	if err != nil {
		log.Fatal(err)
	}
	p := NewProject(filePath, plantumlService, "markdown", logger, m, fs, outputDir, false)
	file := p.CreateQueryParamDataModel("Bar", m.Apps["Bar"].Endpoints["GET /address"].RestParams.QueryParam[0])
	t.Log(file)
	t.Log(p.FilesToCreate[file])
}
