package catalog

import (
	"testing"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/syslwrapper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGenerator_DataModelReturnTableHandlesEmpty(t *testing.T) {
	m, err := parse.NewParser().ParseString(`
App1:
	Endpoint1:
		return ok

`)
	assert.Nil(t, err)

	gen := &Generator{RootModule: m, Fs: afero.NewMemMapFs(), FilesToCreate: map[string]string{}, Log: logrus.New()}
	gen.Mapper = syslwrapper.MakeAppMapper(m)
	gen.Mapper.IndexTypes()
	gen.Mapper.ConvertTypes()
	result := gen.DataModelReturnTable("App1", m.Apps["App1"].Endpoints["Endpoint1"].Stmt[0], m.Apps["App1"].Endpoints["Endpoint1"])
	assert.Empty(t, result)
}
func TestGenerator_DataModelAliasTable(t *testing.T) {
	m, err := parse.NewParser().ParseString(`
App1:
	Endpoint1(something <: myTuple):
		App2 <- Endpoint2

	!type myTuple:
		apple <: apples
		pear <: string
		count <: int

	!alias apples:
		sequence of string


App2:
	Endpoint2:
		...

`)
	assert.Nil(t, err)

	gen := &Generator{RootModule: m, Fs: afero.NewMemMapFs(), FilesToCreate: map[string]string{}, Log: logrus.New()}
	gen.Mapper = syslwrapper.MakeAppMapper(m)
	gen.Mapper.IndexTypes()
	gen.Mapper.ConvertTypes()
	result := gen.DataModelAliasTable(m.Apps["App1"], m.Apps["App1"].Endpoints["Endpoint1"].Param[0])
	t.Log(result)
}

func TestGenerator_DataModelAliasTableRef(t *testing.T) {
	m, err := parse.NewParser().ParseString(`
App1:
	Endpoint1:
		App2 <- Endpoint2
		return ok <: myTuple

	!type myTuple:
		apple <: apples
		pear <: string
		count <: int

	!type apples:
		number <: int
		type <: string


App2:
	Endpoint2:
		...

`)
	assert.Nil(t, err)

	gen := &Generator{RootModule: m, Fs: afero.NewMemMapFs(), FilesToCreate: map[string]string{}, Log: logrus.New()}
	gen.Mapper = syslwrapper.MakeAppMapper(m)
	gen.Mapper.IndexTypes()
	gen.Mapper.ConvertTypes()
	result := gen.DataModelReturnTable("App1", m.Apps["App1"].Endpoints["Endpoint1"].Stmt[1], m.Apps["App1"].Endpoints["Endpoint1"])
	assert.Contains(t, result, "| count | int |")
	assert.Contains(t, result, "apple")
	assert.Contains(t, result, "apples")
}

func TestGenerator_DataModelTableEnum(t *testing.T) {
	m, err := parse.NewParser().ParseString(`
App1:
	!enum OrderStatus:
		created: 1
		placed: 2
		shipped: 3
		delivered: 4

`)
	assert.Nil(t, err)

	gen := &Generator{RootModule: m, Fs: afero.NewMemMapFs(), FilesToCreate: map[string]string{}, Log: logrus.New()}
	gen.Mapper = syslwrapper.MakeAppMapper(m)
	gen.Mapper.IndexTypes()
	gen.Mapper.ConvertTypes()
	result := gen.DataModelTable("App1", "OrderStatus", "")
	t.Log(result)
	assert.Contains(t, result, "enum OrderStatus")
	assert.Contains(t, result, "| 1 | created |")
	assert.Contains(t, result, "| 4 | delivered |")
}

func TestGenerator_DataModelTableSequenceString(t *testing.T) {
	m, err := parse.NewParser().ParseString(`
App1:
	!alias Apples:
		sequence of string
`)
	assert.Nil(t, err)

	gen := &Generator{RootModule: m, Fs: afero.NewMemMapFs(), FilesToCreate: map[string]string{}, Log: logrus.New()}
	gen.Mapper = syslwrapper.MakeAppMapper(m)
	gen.Mapper.IndexTypes()
	gen.Mapper.ConvertTypes()
	result := gen.DataModelTable("App1", "Apples", "")
	assert.Contains(t, result, "Apples")
	assert.Contains(t, result, "sequence of string")
}

func TestGenerator_DataModelTableSequenceTuple(t *testing.T) {
	m, err := parse.NewParser().ParseString(`
App1:
	!alias Apples:
		sequence of Apple
	!type Apple:
		seeds <: int
		variant <: string
`)
	assert.Nil(t, err)

	gen := &Generator{RootModule: m, Fs: afero.NewMemMapFs(), FilesToCreate: map[string]string{}, Log: logrus.New()}
	gen.Mapper = syslwrapper.MakeAppMapper(m)
	gen.Mapper.IndexTypes()
	gen.Mapper.ConvertTypes()
	result := gen.DataModelTable("App1", "Apples", "")
	assert.Contains(t, result, "Apples")
	assert.Contains(t, result, "sequence of")
	assert.Contains(t, result, "App1.Apple")
}

func TestGenerator_DataModelTableRef(t *testing.T) {
	m, err := parse.NewParser().ParseString(`
App1:
	!alias Water:
		h20

	!type h20:
		form <: string
`)
	assert.Nil(t, err)

	gen := &Generator{RootModule: m, Fs: afero.NewMemMapFs(), FilesToCreate: map[string]string{}, Log: logrus.New()}
	gen.Mapper = syslwrapper.MakeAppMapper(m)
	gen.Mapper.IndexTypes()
	gen.Mapper.ConvertTypes()
	result := gen.DataModelTable("App1", "Water", "")
	assert.Contains(t, result, "Water")
	assert.Contains(t, result, "App1.h20")
}
