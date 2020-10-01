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

// const appSVG = "http://plantuml.com/plantuml/svg/~1UDfoA2v9B2efpStXKYSQSAchAn05e4eTGqFytLtzN8CSGrnT59pzNLmLT7KLNFmL_Fn355nTF4CKuKg9DfLejt8bvoGM5oieEHOKwANbvoif91Ohn1iesDWeQBZev1SbPsIcQ2hOsIbKSsahb6Ha5YjOAHIN56NcfNFLSZcavgM0WWOG003__moeEQ80"
// const stringSVG = "http://plantuml.com/plantuml/svg/~1UDfoA2v9B2efpStXKYSQSAchAn05e4eTGqFytLtzN8CSGrnT59pzNLmLT7KLNFmL_Fn355nTF4CKuKg9DfLejt8bvoGM5oie1OS4UVgvAoaa5Yl46oZOs2XekEZa5oLdPAPeAXIN56NcfIlOsIbKSzLoEQJcfO0211000F__rXqsVm00"
// const requestSVG = "http://plantuml.com/plantuml/svg/~1UDgCaB4AmZ0KHVVt5TSkLRJYBAMqc51S6YXnbj843TIeUUaa_hi8ufmp7owNKtCSGfnl4-N9K9uZYP_QdBHgPIVxHak1Wn8IHG6Xq2aDAOvwyLUJLvE_qZWDpCZQy1YrvUZyPTlRvsmvPXWOvntA4aknkOVnwimALOKNhU4Czd0-qfjgwystq2S00F__xq4yMG00"
// const retSeqRefSVG = "http://plantuml.com/plantuml/svg/~1UDfoA2v9B2efpStXKYSQSAchAn05e4eTGqFytLtzN8CSGrnT59pzNLmLT7KLNFmL_Fn355nTF4CKuKg9DfLejt8bvoGM5oieEHOKwANbvoif91Ohn1iesDWeQBZev1SbPsIcQ2hOsIbKSsahb6Ha5YjOAHIN56NcfNFLSZcavgM0WWOG003__moeEQ80"
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
	p := NewProject(filePath, plantumlService, "markdown", logger, m, afero.NewMemMapFs(), outputDir)
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
	p := NewProject(filePath, plantumlService, "markdown", logger, m, afero.NewMemMapFs(), outputDir)
	file := p.CreateParamDataModel(m.Apps["App"], m.Apps["App"].Endpoints["GET /testRestQueryParam/{id}"].RestParams.QueryParam[0])
	assert.Equal(t, "App/fooquerystring.svg", file)
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
	p := NewProject(filePath, plantumlService, "markdown", logger, m, afero.NewMemMapFs(), outputDir)
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
	p := NewProject(filePath, plantumlService, "markdown", logger, m, afero.NewMemMapFs(), outputDir)
	file := p.CreateParamDataModel(m.Apps["App"], m.Apps["App"].Endpoints["GET /testURLParamRef/{id}"].RestParams.UrlParam[0])
	assert.Equal(t, "App/fooid.svg", file)
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
	p := NewProject(filePath, plantumlService, "markdown", logger, m, afero.NewMemMapFs(), outputDir)
	file := p.CreateParamDataModel(m.Apps["MobileApp"], m.Apps["MobileApp"].Endpoints["Login"].Param[0])
	assert.Equal(t, "GrpcTesting/requestinput.svg", file)
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
	p := NewProject(filePath, plantumlService, "markdown", logger, m, afero.NewMemMapFs(), outputDir)

	type paramCase struct {
		fileName, svg string
	}
	cases := []paramCase{
		{
			"primitive/stringprimitiveparam.svg",
			"http://plantuml.com/plantuml/svg/~1UDfoA2v9B2efpStXKYSQSAchAn05e4eTGqFytLtzN8CSGrnT59pzNLmLT7KLNFmL_Fn355nTF4CKuKg9DfLejt8bvoGM5oie5HHbvcQMP9Qb1YGM9UOgAIGMAyGRADZOA6YuwEGN9MTafcWg59SKPUQbAzZPALHprN8vfEQbW0864000___cfpfo",
		},
		{
			"App/complexcomplex.svg",
			"http://plantuml.com/plantuml/svg/~1UDfoA2v9B2efpStXKYSQSAchAn05e4eTGqFytLtzN8CSGrnT59pzNLmLT7KLNFmL_Fn355nTF4CKuKg9DfLejt8bvoGM5oie9UVd5kIaLYWf91Ohn1iesDWeQBZev1SbPsIcQ2eu5XJeEKCKADZPALHpQIiKPuAu2bOAnIL5cNdfNBLS3gbvAQ200WG00F__YgG_8000",
		},
		{
			"App/verycomplexvery_complex.svg",
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
	p := NewProject(filePath, plantumlService, "markdown", logger, m, afero.NewMemMapFs(), outputDir)
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

func TestCreateTypeDiagramWithRecursiveSequenceTypes(t *testing.T) {
	//TODO: do sequence of aliased and non-aliased primitives, and sequence of self
	filePath := "../../tests/seq_type.sysl"
	outputDir := "test"
	fs := afero.NewOsFs()
	logger := logrus.New()
	m, _, err := loader.LoadSyslModule("", filePath, fs, logger)
	if err != nil {
		t.Fatal(err)
	}
	p := NewProject(filePath, plantumlService, "markdown", logger, m, afero.NewMemMapFs(), outputDir)
	fileStringSequenceRef := p.CreateTypeDiagram("App", "VeryComplex", m.Apps["App"].Types["VeryComplex"], true)
	assert.Equal(t, "App/verycomplex.svg", fileStringSequenceRef)
	assert.Equal(t,
		"http://plantuml.com/plantuml/svg/~1UDgCaK5Bn30GHk_pApvpgKPJkjUbb2sR85useD5hATcWXMwt9brGfFzTIglieOUysy3ZpS3imb3xuN9gAOc6aWHHB6hvQlIZEgZdqYY9lPOAGa1g7BI1aa_cvb-DhaRVIhQjGm1xS_vxVpxrhVjYg0Eg37cEM_bmzlQZETwXxFlqIa9Hu8Vk4PffzDY2ynVtUN6TTSZjB1MSq_YtOJ7d-cQbRjVAs28ClkdUQQGg0nS2B4jJpb1jQEUwu_IRtm000F__KVnXp000",
		p.FilesToCreate[path.Join(outputDir, fileStringSequenceRef)],
	)

	fileStringSequenceRef = p.CreateTypeDiagram("App2", "KindaComplex", m.Apps["App2"].Types["KindaComplex"], true)
	assert.Equal(t, "App2/kindacomplex.svg", fileStringSequenceRef)
	assert.Equal(t,
		"http://plantuml.com/plantuml/svg/~1UDgCaazBmp0OX-_v53zxgKPRkXv7AIjB415MqEgrb3He8SjgtWuKyRlBQAVxTyXxItXuyhoG1GsD6xPhmqBlOM48hvdGugKgx-LAFXML55YMGYAA84gioWfIF5HNVwYtrRkLATCG8J2QQBbPzqP_1cW8TO8imxMuspZrthKtAFe-VjNkQuAGv_Xcjw1kTtIOW_Dd_R7LLcbaQeHopjdF_eM97Esp2tProOmvQ72TVViNFNQsROvQqtDR6XzYU-V1bYV59-kDvNK27aK8aZZc0UePgzEZYuGjj6rtl-Ct003__xMAYOC0",
		p.FilesToCreate[path.Join(outputDir, fileStringSequenceRef)],
	)
}

func TestCreateTypeDiagramWithRecursiveNamespacedSequenceTypes(t *testing.T) {
	//TODO: do sequence of aliased and non-aliased primitives, and sequence of self
	filePath := "../../tests/namespaced_seq_type.sysl"
	outputDir := "test"
	fs := afero.NewOsFs()
	logger := logrus.New()
	m, _, err := loader.LoadSyslModule("", filePath, fs, logger)
	if err != nil {
		t.Fatal(err)
	}
	p := NewProject(filePath, plantumlService, "markdown", logger, m, afero.NewMemMapFs(), outputDir)
	fileStringSequenceRef := p.CreateTypeDiagram("First :: App", "VeryComplex", m.Apps["First :: App"].Types["VeryComplex"], true)
	assert.Equal(t, "App/verycomplex.svg", fileStringSequenceRef)
	assert.Equal(t,
		"http://plantuml.com/plantuml/svg/~1UDgCaS5Bmp0KX-_lh_ZPIpMQqbr3AIjB52yBQFHQIXOqq7KroK2H_UyoDiM3IBxRWu_tFNZBc8QGzjkHocoeB975MUsUZBvJQ_NG6IMnqbA1SqYDjJPGyjvS_AZPMxDl9JiECQ9uTk5ZjTlEetilC4JqDPe6b_9c5-ohtrpXreUO80IwUQv-sMXRVD8reZ-E0GACFRgtlPkiGsDFKiiO7RvJP_EKMVoNiNyb811JyDCB7QYlmJX7KLSLAz0lQEccpV5RNm400F__RwXdH000",
		p.FilesToCreate[path.Join(outputDir, fileStringSequenceRef)],
	)

	fileStringSequenceRef = p.CreateTypeDiagram("Second :: App2", "KindaComplex", m.Apps["Second :: App2"].Types["KindaComplex"], true)
	assert.Equal(t, "App2/kindacomplex.svg", fileStringSequenceRef)
	assert.Equal(t,
		"http://plantuml.com/plantuml/svg/~1UDgCab5hip0KXk_pAzxF6xzbKwoN8b9fiZ465LQtMv4Oje2DBiV2njX_7qMxjMKRdZlboVEUDsH9G-s6tRPb_knXlezPCw7vGrdSoBMypqifi2g4H11055WN9QIuBVzzeREjitTbggs9uBMQQDLj-rQ_UgW9LOGqnNznrJdrR9eBN1j70v84UT-7ZzgzJJIo3E_i4cJsdmu9ED_ebvssIYDL8vHnpZuEKHCvnYp-Yb4_HI013SxOfs_ZdR5DVR9zVARTQDNPPlB6uV61W_te1ivA_5PzR5L5u440Bmf3xg5Qi2e--H6nfRRD_Wbz1W00__-iqvJ7",
		p.FilesToCreate[path.Join(outputDir, fileStringSequenceRef)],
	)
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
	p := NewProject(filePath, plantumlService, "markdown", logger, m, afero.NewMemMapFs(), outputDir)
	fileStringEmpty := p.CreateReturnDataModel("App", m.Apps["App"].Endpoints["GET /testReturnNil"].Stmt[0], m.Apps["App"].Endpoints["GET /testReturnNil"])
	assert.Equal(t, "", fileStringEmpty)
}

func TestCreateSequenceMermaid(t *testing.T) {
	filePath := "../../tests/verysimple.sysl"
	outputDir := "test"
	fs := afero.NewOsFs()
	logger := logrus.New()
	m, _, err := loader.LoadSyslModule("", filePath, fs, logger)
	if err != nil {
		t.Fatal(err)
	}
	p := NewProject(filePath, plantumlService, "markdown", logger, m, afero.NewMemMapFs(), outputDir)
	mermaidString := p.CreateSequenceDiagramMermaid("bar", m.Apps["bar"].Endpoints["barendpoint"])
	assert.NotNil(t, mermaidString)
}

func TestCreateIntegrationMermaid(t *testing.T) {
	filePath := "../../tests/verysimple.sysl"
	outputDir := "test"
	fs := afero.NewOsFs()
	logger := logrus.New()
	m, _, err := loader.LoadSyslModule("", filePath, fs, logger)
	if err != nil {
		t.Fatal(err)
	}
	p := NewProject(filePath, plantumlService, "markdown", logger, m, afero.NewMemMapFs(), outputDir)
	mermaidString := p.CreateIntegrationDiagram(m, "", true)
	assert.NotNil(t, mermaidString)
}

func TestCreateTypeMermaid(t *testing.T) {
	filePath := "../../tests/datamodel.sysl"
	outputDir := "test"
	fs := afero.NewOsFs()
	logger := logrus.New()
	m, _, err := loader.LoadSyslModule("", filePath, fs, logger)
	if err != nil {
		t.Fatal(err)
	}
	p := NewProject(filePath, plantumlService, "markdown", logger, m, afero.NewMemMapFs(), outputDir)
	mermaidString := p.CreateParamDataModel(m.Apps["MobileApp"], m.Apps["MobileApp"].Endpoints["Login"].Param[0])
	assert.NotNil(t, mermaidString)
}

func TestCreateRedoc(t *testing.T) {
	appName := "myAppName"
	fileName := "/sysl/myfile.yaml"
	sourceContext := &sysl.SourceContext{File: fileName}
	app := &sysl.Application{SourceContext: sourceContext}
	gen := Generator{
		CurrentDir:         "myAppName",
		RedocFilesToCreate: make(map[string]string),
		Redoc:              true,
	}
	link := gen.CreateRedoc(app, appName)
	registeredFile, ok := gen.RedocFilesToCreate["myAppName/myappname.redoc.html"]
	assert.True(t, ok)
	assert.Equal(t, "/sysl/myfile.yaml@master", registeredFile)
	assert.Equal(t, "myappname.redoc.html", link)
}

func TestCreateRedocFlagFalse(t *testing.T) {
	appName := "myAppName"
	fileName := "myfile.yaml"
	sourceContext := &sysl.SourceContext{File: fileName}
	app := &sysl.Application{SourceContext: sourceContext}
	gen := Generator{
		RedocFilesToCreate: make(map[string]string),
		Redoc:              false,
	}
	link := gen.CreateRedoc(app, appName)
	assert.Equal(t, "", link)
}

func TestCreateRedocFromImportRemote(t *testing.T) {
	appName := "myAppName"
	fileName := "github.com/myorg/myrepo/specs/myfile.yaml"
	sourceContext := &sysl.SourceContext{File: fileName}
	app := &sysl.Application{SourceContext: sourceContext}
	gen := Generator{
		RedocFilesToCreate: make(map[string]string),
		Redoc:              true,
	}
	link := gen.CreateRedoc(app, appName)
	assert.Equal(t, "myappname.redoc.html", link)
}
