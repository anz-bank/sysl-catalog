package catalog

import "github.com/anz-bank/sysl/pkg/sysl"

type SourceCoder interface {
	Attr
	GetSourceContext() *sysl.SourceContext
}

type Param interface {
	Typer
	GetName() string
}

type Attr interface {
	GetAttrs() map[string]*sysl.Attribute
}

type Namer interface {
	Attr
	GetName() *sysl.AppName
}

type Typer interface {
	GetType() *sysl.Type
}
