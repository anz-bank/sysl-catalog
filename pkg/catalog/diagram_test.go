package catalog

import (
	"path"
	"testing"

	"github.com/anz-bank/sysl/pkg/loader"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

const appSVG = "http://plantuml.com/plantuml/svg/~1UDfoA2v9B2efpStXKYSQSAchAn05e4eTGqFytLtzN8CSGrnT59pzNLmLT7KLNFmL_Fn355nTF4CKuKg9DfLejt8bvoGM5oieEHOKwANbvoif91Ohn1iesDWeQBZev1SbPsIcQ2hOsIbKSsahb6Ha5YjOAHIN56NcfNFLSZcavgM0WWOG003__moeEQ80"
const stringSVG = "http://plantuml.com/plantuml/svg/~1UDfoA2v9B2efpStXKYSQSAchAn05e4eTGqFytLtzN8CSGrnT59pzNLmLT7KLNFmL_Fn355nTF4CKuKg9DfLejt8bvoGM5oie1OS4UVgvAoaa5Yl46oZOs2XekEZa5oLdPAPeAXIN56NcfIlOsIbKSzLoEQJcfO0211000F__rXqsVm00"
const requestSVG = "http://plantuml.com/plantuml/svg/~1UDgCaB4AmZ0KHVVt5TSkLRJYBAMqc51S6YXnbj843TIeUUaa_hi8ufmp7owNKtCSGfnl4-N9K9uZYP_QdBHgPIVxHak1Wn8IHG6Xq2aDAOvwyLUJLvE_qZWDpCZQy1YrvUZyPTlRvsmvPXWOvntA4aknkOVnwimALOKNhU4Czd0-qfjgwystq2S00F__xq4yMG00"
const retSeqRefSVG = "http://plantuml.com/plantuml/svg/~1UDfoA2v9B2efpStXKYSQSAchAn05e4eTGqFytLtzN8CSGrnT59pzNLmLT7KLNFmL_Fn355nTF4CKuKg9DfLejt8bvoGM5oieEHOKwANbvoif91Ohn1iesDWeQBZev1SbPsIcQ2hOsIbKSsahb6Ha5YjOAHIN56NcfNFLSZcavgM0WWOG003__moeEQ80"
const retSeqPrimSVG = "http://plantuml.com/plantuml/svg/~1UDfoA2v9B2efpStXKYSQSAchAn05e4eTGqFytLtzN8CSGrnT59pzNLmLT7KLNFmL_Fn355nTF4CKuKg9DfLejt8bvoGM5oie5PSKPUQbAoaa5Yl46oZOs2XekEZa5oLdPAPeAa3a5Epi5AgvQhaSKlDIG0422000___vlJS_"

func TestCreateIntegrationDiagramPlantuml(t *testing.T) {
	m, err := parse.NewParser().ParseString(`
App1:
	Endpoint1:
		App2 <- Endpoint2
App2:
	Endpoint2:
		...

`)
	assert.Nil(t, err)

	gen := &Generator{RootModule: m, Fs: afero.NewMemMapFs(), FilesToCreate: map[string]string{}, Log: logrus.New()}
	filename := gen.CreateIntegrationDiagramPlantuml(m, "", false)
	contents := gen.FilesToCreate[filename]
	assert.Equal(t,
		"/svg/~1UDgCpa5Bn30G1U1xViNpr5DXgo8UbhBTjegNBKZ5WqW9RMX2aqoPfAY8_rtMYWTFVSVv7ZEJR8v84cpARxLuQflx-bG_5crTeMog6ccAgi6fQL5N3-t5NtNpris_YaE8akFYhD1cK0XHiQBuCIiHUcaLd7n7TdDrUmsjpAYZ29FnisJfq9ERoIiVyIc0e-odaMdnGqcM67UMMDfdRQ8wA_6WU9MZbVqaW8APtjPHoSO5yk9Bl1JpdBr21dGxxFVQZDgUx-Rv3rskbFsZReSqpT5bug3yi3Zx7G00__yzoceA",
		contents)
}
func TestCreateQueryParamDataModelWithPrimitive(t *testing.T) {
	filePath := "../../tests/rest.sysl"
	outputDir := "test"
	fs := afero.NewOsFs()
	logger := logrus.New()
	m, _, err := loader.LoadSyslModule("", filePath, fs, logger)
	if err != nil {
		t.Fatal(err)
	}
	p := NewProject(filePath, plantumlService, "markdown", "", "", logger, m, fs, outputDir)
	file := p.CreateParamDataModel(m.Apps["Bar"], m.Apps["Bar"].Endpoints["GET /address"].RestParams.QueryParam[0])
	assert.Equal(t, "primitive/stringstreet.svg", file)
	assert.Equal(t,
		"http://plantuml.com/plantuml/svg/~1UDfoA2v9B2efpStXKYSQSAchAn05e4eTGqFytLtzN8CSGrnT59pzNLmLT7KLNFmL_Fn355nTF4CKuKg9DfLejt8bvoGM5oie5PSKfQQMA2aa5Yl46oZOs2XekEZa5oLdPAPeAXIN56NcfIlOsIbKSzLoEQJcfO0211000F__-QmtFm00",
		p.FilesToCreate[path.Join(outputDir, file)],
	)
}

func TestCreateQueryParamDataModelWithTypeRef(t *testing.T) {
	filePath := "../../tests/params.sysl"
	outputDir := "test"
	fs := afero.NewOsFs()
	logger := logrus.New()
	m, _, err := loader.LoadSyslModule("", filePath, fs, logger)
	if err != nil {
		t.Fatal(err)
	}
	p := NewProject(filePath, plantumlService, "markdown", "", "", logger, m, fs, outputDir)
	file := p.CreateParamDataModel(m.Apps["App"], m.Apps["App"].Endpoints["GET /testRestQueryParam/{id}"].RestParams.QueryParam[0])
	assert.Equal(t, "App/foo.svg", file)
	assert.Equal(t,
		"http://plantuml.com/plantuml/svg/~1UDfoA2v9B2efpStXKYSQSAchAn05e4eTGqFytLtzN8CSGrnT59pzNLmLT7KLNFmL_Fn355nTF4CKuKg9DfLejt8bvoGM5oie5PQc5bK6bnHbvgKhAIGMAyGRADZOA6YuwEGN9MTafcWgE1OKwANbvolOsIbKSsahb6Ha5YjOAHI3TN3LSZcavgM0WWaG003__nwmFcy0",
		p.FilesToCreate[path.Join(outputDir, file)],
	)
}

func TestCreatePathParamDataModelWithPrimitive(t *testing.T) {
	filePath := "../../tests/params.sysl"
	outputDir := "test"
	fs := afero.NewOsFs()
	logger := logrus.New()
	m, _, err := loader.LoadSyslModule("", filePath, fs, logger)
	if err != nil {
		t.Fatal(err)
	}
	p := NewProject(filePath, plantumlService, "markdown", "", "", logger, m, fs, outputDir)
	file := p.CreateParamDataModel(m.Apps["App"], m.Apps["App"].Endpoints["GET /testURLParamPrimitive/{id}"].RestParams.UrlParam[0])
	assert.Equal(t, "primitive/stringid.svg", file)
	assert.Equal(t,
		"http://plantuml.com/plantuml/svg/~1UDfoA2v9B2efpStXKYSQSAchAn05e4eTGqFytLtzN8CSGrnT59pzNLmLT7KLNFmL_Fn355nTF4CKuKg9DfLejt8bvoGM5oiePQOeAIGMAyGRADZOA6YuwEGN9MTafcWg59SKPUQbAzZPALHprN8vfEQbW0834000__y4SpLr",
		p.FilesToCreate[path.Join(outputDir, file)],
	)
}

func TestCreatePathParamDataModelWithTypeRef(t *testing.T) {
	filePath := "../../tests/params.sysl"
	outputDir := "test"
	fs := afero.NewOsFs()
	logger := logrus.New()
	m, _, err := loader.LoadSyslModule("", filePath, fs, logger)
	if err != nil {
		t.Fatal(err)
	}
	p := NewProject(filePath, plantumlService, "markdown", "", "", logger, m, fs, outputDir)
	file := p.CreateParamDataModel(m.Apps["App"], m.Apps["App"].Endpoints["GET /testURLParamRef/{id}"].RestParams.UrlParam[0])
	assert.Equal(t, "App/foo.svg", file)
	assert.Equal(t,
		"http://plantuml.com/plantuml/svg/~1UDfoA2v9B2efpStXKYSQSAchAn05e4eTGqFytLtzN8CSGrnT59pzNLmLT7KLNFmL_Fn355nTF4CKuKg9DfLejt8bvoGM5oiePQOeAIGMAyGRADZOA6YuwEGN9MTafcWgE1OKwANbvolOsIbKSsahb6Ha5YjOAHIN56NcfNFLSZcavgM0GWSG003__x99Eey0",
		p.FilesToCreate[path.Join(outputDir, file)],
	)
}

func TestCreateParamDataModel(t *testing.T) {
	filePath := "../../tests/datamodel.sysl"
	outputDir := "test"
	fs := afero.NewOsFs()
	logger := logrus.New()
	m, _, err := loader.LoadSyslModule("", filePath, fs, logger)
	if err != nil {
		t.Fatal(err)
	}
	p := NewProject(filePath, plantumlService, "markdown", "", "", logger, m, fs, outputDir)
	file := p.CreateParamDataModel(m.Apps["MobileApp"], m.Apps["MobileApp"].Endpoints["Login"].Param[0])
	assert.Equal(t, "GrpcTesting/request.svg", file)
	assert.Equal(t,
		"http://plantuml.com/plantuml/svg/~1UDgCa47BWa0KHVVlLzpFZL-KqJf4b6QGDWx8j0xHeGPiCjzp5Vtt2ABrdFNXSZabIpVBSXifZORI555yrUfaJQqRtLPMAnoCqiWoA8F6M6Xrj7y_DNer-YlrOyUCn8TfaGGTuxn3dkDVRUvpV_N32lKyzTQn-73PjkwnE1OK1PwqXX-mXmz2BofT63wTtW400F__Kli-gG00",
		p.FilesToCreate[path.Join(outputDir, file)],
	)
}

func TestCreateParamDataModelWithRestParam(t *testing.T) {
	filePath := "../../tests/rest_params.sysl"
	outputDir := "test"
	fs := afero.NewOsFs()
	logger := logrus.New()
	m, _, err := loader.LoadSyslModule("", filePath, fs, logger)
	if err != nil {
		t.Fatal(err)
	}
	p := NewProject(filePath, plantumlService, "markdown", "", "", logger, m, fs, outputDir)

	type paramCase struct {
		fileName, svg string
	}
	cases := []paramCase{
		{
			"primitive/stringprimitiveparam.svg",
			"http://plantuml.com/plantuml/svg/~1UDfoA2v9B2efpStXKYSQSAchAn05e4eTGqFytLtzN8CSGrnT59pzNLmLT7KLNFmL_Fn355nTF4CKuKg9DfLejt8bvoGM5oie5HHbvcQMP9Qb1YGM9UOgAIGMAyGRADZOA6YuwEGN9MTafcWg59SKPUQbAzZPALHprN8vfEQbW0864000___cfpfo",
		},
		{
			"App/complex.svg",
			"http://plantuml.com/plantuml/svg/~1UDfoA2v9B2efpStXKYSQSAchAn05e4eTGqFytLtzN8CSGrnT59pzNLmLT7KLNFmL_Fn355nTF4CKuKg9DfLejt8bvoGM5oie9UVd5kIaLYWf91Ohn1iesDWeQBZev1SbPsIcQ2eu5XJeEKCKADZPALHpQIiKPuAu2bOAnIL5cNdfNBLS3gbvAQ200WG00F__YgG_8000",
		},
		{
			"App/verycomplex.svg",
			"http://plantuml.com/plantuml/svg/~1UDgCaKrBn30GXk_v5Q-vLACfc-koB2rD42yDIFGQIXkqq2-IAXRnlqkn2YuzvBqDyymp0vE5kVBpMz-H93eaIH2L3SsVZBvNfNhCZP8ej5JW75AZr0PAFfYhFpJQ6dqhgRig1D1wxAVEVL1K0VQ0qmdNycxqzMlRt22VfhJu0N0-uvFS8hHhYIF2xDlXXNpzYjwTN-m_czYnFJk_N1Yt6Hp1sDPYR7UJ5M2SWueq5Q2m1vAveLcVz1q00F__lLPaQ000",
		},
	}

	for i, c := range cases {
		file := p.CreateParamDataModel(m.Apps["App"], m.Apps["App"].Endpoints["GET /endpoint"].Param[i])
		assert.Equal(t, c.fileName, file)
		assert.Equal(t, c.svg, p.FilesToCreate[path.Join(outputDir, file)])
	}
}

func TestCreateReturnDataModelWithSequence(t *testing.T) {
	filePath := "../../tests/return.sysl"
	outputDir := "test"
	fs := afero.NewOsFs()
	logger := logrus.New()
	m, _, err := loader.LoadSyslModule("", filePath, fs, logger)
	if err != nil {
		t.Fatal(err)
	}
	p := NewProject(filePath, plantumlService, "markdown", "", "", logger, m, fs, outputDir)
	fileStringSequenceRef := p.CreateReturnDataModel("App", m.Apps["App"].Endpoints["somefoo"].Stmt[0], m.Apps["App"].Endpoints["somefoo"])
	fileStringSequencePrimitive := p.CreateReturnDataModel("App", m.Apps["App"].Endpoints["someprimitivefoo"].Stmt[0], m.Apps["App"].Endpoints["someprimitivefoo"])
	assert.Equal(t, "App/foo.svg", fileStringSequenceRef)
	assert.Equal(t, "App/stringsimple.svg", fileStringSequencePrimitive)
	assert.Equal(t,
		"http://plantuml.com/plantuml/svg/~1UDgCaB4AmZ0KHVVt5TSkLRJYBAKqc20k3KYS9RJM1cfJyk8a_hi8ufmp7owNKtEq8JuV8-N9K9uZYPygBaOVLQFEmYY9WvOAHG6fqMW39KzcyLUJLvE_KZjQPcIzznaiuxf3MM8fDpwqW-jM4FEyxRr7LU55QyJ1CVRW6DnqfpVLjwdxvla4003___UGEmO0",
		p.FilesToCreate[path.Join(outputDir, fileStringSequenceRef)],
	)
	assert.Equal(t, retSeqPrimSVG, p.FilesToCreate[path.Join(outputDir, fileStringSequencePrimitive)])
}

func TestCreateReturnDataModelWithEmpty(t *testing.T) {
	filePath := "../../tests/return.sysl"
	outputDir := "test"
	fs := afero.NewOsFs()
	logger := logrus.New()
	m, _, err := loader.LoadSyslModule("", filePath, fs, logger)
	if err != nil {
		t.Fatal(err)
	}
	p := NewProject(filePath, plantumlService, "markdown", "", "", logger, m, fs, outputDir)
	fileStringEmpty := p.CreateReturnDataModel("App", m.Apps["App"].Endpoints["GET /testReturnNil"].Stmt[0], m.Apps["App"].Endpoints["GET /testReturnNil"])
	assert.Equal(t, "", fileStringEmpty)
}
func TestCreateRedoc(t *testing.T) {
	appName := "myAppName"
	fileName := "myfile.yaml"
	sourceContext := &sysl.SourceContext{File: fileName}
	gen := Generator{
		CurrentDir:         "myAppName",
		RedocFilesToCreate: make(map[string]string),
		Redoc:              true,
	}
	link := gen.CreateRedoc(sourceContext, appName)
	t.Log(gen.RedocFilesToCreate)
	registeredFile, ok := gen.RedocFilesToCreate["myAppName/myappname.redoc.html"]
	assert.True(t, ok)
	assert.Equal(t, "https://raw.githubusercontent.com/anz-bank/sysl-catalog/master/myfile.yaml", registeredFile)
	assert.Equal(t, "myappname.redoc.html", link)
}

func TestCreateRedocFlagFalse(t *testing.T) {
	appName := "myAppName"
	fileName := "myfile.yaml"
	sourceContext := &sysl.SourceContext{File: fileName}
	gen := Generator{
		RedocFilesToCreate: make(map[string]string),
		Redoc:              false,
	}
	link := gen.CreateRedoc(sourceContext, appName)
	assert.Equal(t, "", link)
}

func TestCreateRedocFromImportRemote(t *testing.T) {
	appName := "myAppName"
	fileName := "github.com/myorg/myrepo/specs/myfile.yaml"
	sourceContext := &sysl.SourceContext{File: fileName}
	gen := Generator{
		RedocFilesToCreate: make(map[string]string),
		Redoc:              true,
	}
	link := gen.CreateRedoc(sourceContext, appName)
	assert.Equal(t, "", link)
}
