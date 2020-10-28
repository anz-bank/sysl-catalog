package catalog

import (
	"testing"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/stretchr/testify/require"

	"github.com/anz-bank/sysl/pkg/loader"
	"github.com/anz-bank/sysl/pkg/parse"
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
	t.Parallel()

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
	file := gen.IntegrationPlantuml(m, "", false)
	assert.Equal(t,
		"/svg/~1UDgCpa5Bn30G1U1xViNpr5DXgo8UbhBTjegNBKZ5WqW9RMX2aqoPfAY8_rtMYWTFVSVv7ZEJR8v84cpARxLuQflx-bG_5crTeMog6ccAgi6fQL5N3-t5NtNpris_YaE8akFYhD1cK0XHiQBuCIiHUcaLd7n7TdDrUmsjpAYZ29FnisJfq9ERoIiVyIc0e-odaMdnGqcM67UMMDfdRQ8wA_6WU9MZbVqaW8APtjPHoSO5yk9Bl1JpdBr21dGxxFVQZDgUx-Rv3rskbFsZReSqpT5bug3yi3Zx7G00__yzoceA",
		file)
}
func TestCreateQueryParamDataModelWithPrimitive(t *testing.T) {
	t.Parallel()

	filePath := "../../tests/rest.sysl"
	p, m, err := loadProject(filePath)
	require.NoError(t, err)
	file := p.DataModelParamPlantuml(m.Apps["Bar"], m.Apps["Bar"].Endpoints["GET /address"].RestParams.QueryParam[0])
	assert.Equal(t,
		"http://plantuml.com/plantuml/svg/~1UDfoA2v9B2efpStXKYSQSAchAn05e4eTGqFytLtzN8CSGrnT59pzNLmLT7KLNFmL_Fn355nTF4CKuKg9DfLejt8bvoGM5oie5PSKfQQMA2aa5Yl46oZOs2XekEZa5oLdPAPeAXIN56NcfIlOsIbKSzLoEQJcfO0211000F__-QmtFm00",
		file,
	)
}

func TestCreateQueryParamDataModelWithTypeRef(t *testing.T) {
	t.Parallel()

	filePath := "../../tests/params.sysl"
	p, m, err := loadProject(filePath)
	require.NoError(t, err)
	file := p.DataModelParamPlantuml(m.Apps["App"], m.Apps["App"].Endpoints["GET /testRestQueryParam/{id}"].RestParams.QueryParam[0])
	assert.Equal(t,
		"http://plantuml.com/plantuml/svg/~1UDfoA2v9B2efpStXKYSQSAchAn05e4eTGqFytLtzN8CSGrnT59pzNLmLT7KLNFmL_Fn355nTF4CKuKg9DfLejt8bvoGM5oie5PQc5bK6bnHbvgKhAIGMAyGRADZOA6YuwEGN9MTafcWgE1OKwANbvolOsIbKSsahb6Ha5YjOAHI3TN3LSZcavgM0WWaG003__nwmFcy0",
		file,
	)
}

func TestCreatePathParamDataModelWithPrimitive(t *testing.T) {
	t.Parallel()

	filePath := "../../tests/params.sysl"
	p, m, err := loadProject(filePath)
	require.NoError(t, err)
	file := p.DataModelParamPlantuml(m.Apps["App"], m.Apps["App"].Endpoints["GET /testURLParamPrimitive/{id}"].RestParams.UrlParam[0])
	assert.Equal(t,
		"http://plantuml.com/plantuml/svg/~1UDfoA2v9B2efpStXKYSQSAchAn05e4eTGqFytLtzN8CSGrnT59pzNLmLT7KLNFmL_Fn355nTF4CKuKg9DfLejt8bvoGM5oiePQOeAIGMAyGRADZOA6YuwEGN9MTafcWg59SKPUQbAzZPALHprN8vfEQbW0834000__y4SpLr",
		file,
	)
}

func TestCreatePathParamDataModelWithTypeRef(t *testing.T) {
	t.Parallel()

	filePath := "../../tests/params.sysl"
	p, m, err := loadProject(filePath)
	require.NoError(t, err)
	file := p.DataModelParamPlantuml(m.Apps["App"], m.Apps["App"].Endpoints["GET /testURLParamRef/{id}"].RestParams.UrlParam[0])
	assert.Equal(t,
		"http://plantuml.com/plantuml/svg/~1UDfoA2v9B2efpStXKYSQSAchAn05e4eTGqFytLtzN8CSGrnT59pzNLmLT7KLNFmL_Fn355nTF4CKuKg9DfLejt8bvoGM5oiePQOeAIGMAyGRADZOA6YuwEGN9MTafcWgE1OKwANbvolOsIbKSsahb6Ha5YjOAHIN56NcfNFLSZcavgM0GWSG003__x99Eey0",
		file,
	)
}

func TestCreateParamDataModel(t *testing.T) {
	t.Parallel()

	filePath := "../../tests/datamodel.sysl"
	p, m, err := loadProject(filePath)
	require.NoError(t, err)
	file := p.DataModelParamPlantuml(m.Apps["MobileApp"], m.Apps["MobileApp"].Endpoints["Login"].Param[0])
	assert.Equal(t,
		"http://plantuml.com/plantuml/svg/~1UDgCa47BWa0KHVVlLzpFZL-KqJf4b6QGDWx8j0xHeGPiCjzp5Vtt2ABrdFNXSZabIpVBSXifZORI555yrUfaJQqRtLPMAnoCqiWoA8F6M6Xrj7y_DNer-YlrOyUCn8TfaGGTuxn3dkDVRUvpV_N32lKyzTQn-73PjkwnE1OK1PwqXX-mXmz2BofT63wTtW400F__Kli-gG00",
		file,
	)
}

func TestCreateParamDataModelWithRestParam(t *testing.T) {
	t.Parallel()

	filePath := "../../tests/rest_params.sysl"
	p, m, err := loadProject(filePath)
	require.NoError(t, err)

	cases := []string{
		"http://plantuml.com/plantuml/svg/~1UDfoA2v9B2efpStXKYSQSAchAn05e4eTGqFytLtzN8CSGrnT59pzNLmLT7KLNFmL_Fn355nTF4CKuKg9DfLejt8bvoGM5oie5HHbvcQMP9Qb1YGM9UOgAIGMAyGRADZOA6YuwEGN9MTafcWg59SKPUQbAzZPALHprN8vfEQbW0864000___cfpfo",
		"http://plantuml.com/plantuml/svg/~1UDfoA2v9B2efpStXKYSQSAchAn05e4eTGqFytLtzN8CSGrnT59pzNLmLT7KLNFmL_Fn355nTF4CKuKg9DfLejt8bvoGM5oie9UVd5kIaLYWf91Ohn1iesDWeQBZev1SbPsIcQ2eu5XJeEKCKADZPALHpQIiKPuAu2bOAnIL5cNdfNBLS3gbvAQ200WG00F__YgG_8000",
		"http://plantuml.com/plantuml/svg/~1UDgCaKrBn30GXk_v5Q-vLACfc-koB2rD42yDIFGQIXkqq2-IAXRnlqkn2YuzvBqDyymp0vE5kVBpMz-H93eaIH2L3SsVZBvNfNhCZP8ej5JW75AZr0PAFfYhFpJQ6dqhgRig1D1wxAVEVL1K0VQ0qmdNycxqzMlRt22VfhJu0N0-uvFS8hHhYIF2xDlXXNpzYjwTN-m_czYnFJk_N1Yt6Hp1sDPYR7UJ5M2SWueq5Q2m1vAveLcVz1q00F__lLPaQ000",
	}

	for i, c := range cases {
		file := p.DataModelParamPlantuml(m.Apps["App"], m.Apps["App"].Endpoints["GET /endpoint"].Param[i])
		assert.Equal(t, c, file)
	}
}

func TestCreateReturnDataModelWithSequence(t *testing.T) {
	t.Parallel()

	filePath := "../../tests/return.sysl"
	p, m, err := loadProject(filePath)
	require.NoError(t, err)
	fileStringSequenceRef := p.DataModelReturnPlantuml("App", m.Apps["App"].Endpoints["somefoo"].Stmt[0], m.Apps["App"].Endpoints["somefoo"])
	fileStringSequencePrimitive := p.DataModelReturnPlantuml("App", m.Apps["App"].Endpoints["someprimitivefoo"].Stmt[0], m.Apps["App"].Endpoints["someprimitivefoo"])
	assert.Equal(t,
		"http://plantuml.com/plantuml/svg/~1UDgCaB4AmZ0KHVVt5TSkLRJYBAKqc20k3KYS9RJM1cfJyk8a_hi8ufmp7owNKtEq8JuV8-N9K9uZYPygBaOVLQFEmYY9WvOAHG6fqMW39KzcyLUJLvE_KZjQPcIzznaiuxf3MM8fDpwqW-jM4FEyxRr7LU55QyJ1CVRW6DnqfpVLjwdxvla4003___UGEmO0",
		fileStringSequenceRef,
	)
	assert.Equal(t, retSeqPrimSVG, fileStringSequencePrimitive)
}

func TestCreateTypeDiagramWithRecursiveSequenceTypes(t *testing.T) {
	t.Parallel()

	//TODO: do sequence of aliased and non-aliased primitives, and sequence of self
	filePath := "../../tests/seq_type.sysl"
	p, m, err := loadProject(filePath)
	require.NoError(t, err)
	fileStringSequenceRef := p.DataModelAliasPlantuml("App", "VeryComplex", "VeryComplex", m.Apps["App"].Types["VeryComplex"], true)
	assert.Equal(t,
		"http://plantuml.com/plantuml/svg/~1UDgCaK5Bn30GHk_pApvpgKPJkjUbb2sR85useD5hATcWXMwt9brGfFzTIglieOUysy3ZpS3imb3xuN9gAOc6aWHHB6hvQlIZEgZdqYY9lPOAGa1g7BI1aa_cvb-DhaRVIhQjGm1xS_vxVpxrhVjYg0Eg37cEM_bmzlQZETwXxFlqIa9Hu8Vk4PffzDY2ynVtUN6TTSZjB1MSq_YtOJ7d-cQbRjVAs28ClkdUQQGg0nS2B4jJpb1jQEUwu_IRtm000F__KVnXp000",
		fileStringSequenceRef,
	)

	fileStringSequenceRef = p.DataModelAliasPlantuml("App2", "KindaComplex", "KindaComplex", m.Apps["App2"].Types["KindaComplex"], true)
	assert.Equal(t,
		"http://plantuml.com/plantuml/svg/~1UDgCaazBmp0OX-_v53zxgKPRkXv7AIjB415MqEgrb3He8SjgtWuKyRlBQAVxTyXxItXuyhoG1GsD6xPhmqBlOM48hvdGugKgx-LAFXML55YMGYAA84gioWfIF5HNVwYtrRkLATCG8J2QQBbPzqP_1cW8TO8imxMuspZrthKtAFe-VjNkQuAGv_Xcjw1kTtIOW_Dd_R7LLcbaQeHopjdF_eM97Esp2tProOmvQ72TVViNFNQsROvQqtDR6XzYU-V1bYV59-kDvNK27aK8aZZc0UePgzEZYuGjj6rtl-Ct003__xMAYOC0",
		fileStringSequenceRef,
	)
}

func TestCreateTypeDiagramWithRecursiveNamespacedSequenceTypes(t *testing.T) {
	t.Parallel()

	//TODO: do sequence of aliased and non-aliased primitives, and sequence of self
	filePath := "../../tests/namespaced_seq_type.sysl"
	p, m, err := loadProject(filePath)
	require.NoError(t, err)
	fileStringSequenceRef := p.DataModelAliasPlantuml("First :: App", "VeryComplex", "VeryComplex", m.Apps["First :: App"].Types["VeryComplex"], true)
	assert.Equal(t,
		"http://plantuml.com/plantuml/svg/~1UDgCaS5Bmp0KX-_lh_ZPIpMQqbr3AIjB52yBQFHQIXOqq7KroK2H_UyoDiM3IBxRWu_tFNZBc8QGzjkHocoeB975MUsUZBvJQ_NG6IMnqbA1SqYDjJPGyjvS_AZPMxDl9JiECQ9uTk5ZjTlEetilC4JqDPe6b_9c5-ohtrpXreUO80IwUQv-sMXRVD8reZ-E0GACFRgtlPkiGsDFKiiO7RvJP_EKMVoNiNyb811JyDCB7QYlmJX7KLSLAz0lQEccpV5RNm400F__RwXdH000",
		fileStringSequenceRef,
	)

	fileStringSequenceRef = p.DataModelAliasPlantuml("Second :: App2", "KindaComplex", "KindaComplex", m.Apps["Second :: App2"].Types["KindaComplex"], true)
	assert.Equal(t,
		"http://plantuml.com/plantuml/svg/~1UDgCab5hip0KXk_pAzxF6xzbKwoN8b9fiZ465LQtMv4Oje2DBiV2njX_7qMxjMKRdZlboVEUDsH9G-s6tRPb_knXlezPCw7vGrdSoBMypqifi2g4H11055WN9QIuBVzzeREjitTbggs9uBMQQDLj-rQ_UgW9LOGqnNznrJdrR9eBN1j70v84UT-7ZzgzJJIo3E_i4cJsdmu9ED_ebvssIYDL8vHnpZuEKHCvnYp-Yb4_HI013SxOfs_ZdR5DVR9zVARTQDNPPlB6uV61W_te1ivA_5PzR5L5u440Bmf3xg5Qi2e--H6nfRRD_Wbz1W00__-iqvJ7",
		fileStringSequenceRef,
	)
}

func TestCreateReturnDataModelWithEmpty(t *testing.T) {
	t.Parallel()

	filePath := "../../tests/return.sysl"
	p, m, err := loadProject(filePath)
	require.NoError(t, err)
	fileStringEmpty := p.DataModelReturnPlantuml("App", m.Apps["App"].Endpoints["GET /testReturnNil"].Stmt[0], m.Apps["App"].Endpoints["GET /testReturnNil"])
	assert.Equal(t, "", fileStringEmpty)
}

func TestCreateSequenceMermaid(t *testing.T) {
	t.Parallel()

	filePath := "../../tests/verysimple.sysl"
	p, m, err := loadProject(filePath)
	require.NoError(t, err)
	mermaidString := p.SequenceMermaid("bar", m.Apps["bar"].Endpoints["barendpoint"])
	assert.NotNil(t, mermaidString)
}

func TestCreateIntegrationMermaid(t *testing.T) {
	t.Parallel()

	filePath := "../../tests/verysimple.sysl"
	p, m, err := loadProject(filePath)
	require.NoError(t, err)
	mermaidString := p.IntegrationPlantuml(m, "", true)
	assert.NotNil(t, mermaidString)
}

func TestCreateTypeMermaid(t *testing.T) {
	t.Parallel()

	filePath := "../../tests/datamodel.sysl"
	p, m, err := loadProject(filePath)
	require.NoError(t, err)
	mermaidString := p.DataModelParamPlantuml(m.Apps["MobileApp"], m.Apps["MobileApp"].Endpoints["Login"].Param[0])
	assert.NotNil(t, mermaidString)
}

func TestPackages(t *testing.T) {
	t.Parallel()

	filePath := "../../tests/packages.sysl"
	p, m, err := loadProject(filePath)
	require.NoError(t, err)
	pkgs := p.Packages(m)

	expected := []string{"Org", "Org :: Team :: System", "Org :: Team2", "Project", "Team Two"}
	assert.Equal(t, expected, pkgs)
}

func TestModuleAsPackages(t *testing.T) {
	t.Parallel()

	filePath := "../../tests/packages.sysl"
	p, m, err := loadProject(filePath)
	require.NoError(t, err)
	pkgs := p.ModuleAsPackages(m)

	assert.NotNil(t, pkgs["Org :: Team :: System"].Apps["Org :: Team :: System :: a"])
	assert.Nil(t, pkgs["Org :: Team2 :: System"])
	assert.Nil(t, pkgs["Org :: Team2"].Apps["Org :: Team2 :: System :: c"])
	assert.NotNil(t, pkgs["Team Two"].Apps["Org :: Team2 :: System :: c"])
}

func loadProject(filePath string) (*Generator, *sysl.Module, error) {
	outputDir := "test"
	fs := afero.NewOsFs()
	logger := logrus.New()
	m, _, err := loader.LoadSyslModule("", filePath, fs, logger)
	if err != nil {
		return nil, nil, err
	}
	p := NewProject(filePath, plantumlService, "markdown", logger, m, afero.NewMemMapFs(), outputDir)
	return p, m, nil
}
