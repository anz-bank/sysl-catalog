package catalogdiagrams

import (
	"fmt"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/diagrams"
	"github.com/anz-bank/sysl/pkg/sequencediagram"

	"github.com/anz-bank/sysl/pkg/datamodeldiagram"
	"github.com/anz-bank/sysl/pkg/integrationdiagram"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
)

const relationArrow = `}--`
const tupleArrow = `*--`

type DataModelView struct {
	datamodeldiagram.DataModelView
}

type DataModelParam struct {
	datamodeldiagram.DataModelParam
}

func (v *DataModelView) GenerateDataView(dataParam *DataModelParam, appName string, tMap map[string]*sysl.Type) string {
	var isRelation bool
	relationshipMap := map[string]map[string]datamodeldiagram.RelationshipParam{}
	v.StringBuilder.WriteString("@startuml\n")
	if dataParam.Title != "" {
		fmt.Fprintf(v.StringBuilder, "title %s\n", dataParam.Title)
	}
	v.StringBuilder.WriteString(integrationdiagram.PumlHeader)
	//typeMap := map[string]*sysl.Type{}

	ignoredTypes := map[string]struct{}{}
	// TODO: Actually put The appName/project name and the appName in a struct so strings.split and join dont need to be used
	entityNames := []string{}
	//for _, t := range tMap {
	//	RecurseivelyGetTypesHelper(appName, t, m, tMap)
	//}
	for key := range tMap {
		entityNames = append(entityNames, key)
	}
	sort.Strings(entityNames)
	for _, entityName := range entityNames {
		entityType := tMap[entityName]
		if relEntity := entityType.GetRelation(); relEntity != nil {
			isRelation = true
			viewParam := datamodeldiagram.EntityViewParam{
				EntityColor:  `orchid`,
				EntityHeader: `D`,
				EntityName:   entityName,
			}
			v.DrawRelation(viewParam, relEntity, relationshipMap)
		} else if tupEntity := entityType.GetTuple(); tupEntity != nil {
			isRelation = false
			viewParam := datamodeldiagram.EntityViewParam{
				EntityColor:  `orchid`,
				EntityHeader: `D`,
				EntityName:   entityName,
				IgnoredTypes: ignoredTypes,
				Types:        tMap,
			}
			v.DrawTuple(viewParam, tupEntity, relationshipMap)
		} else if pe := entityType.GetPrimitive(); pe != sysl.Type_NO_Primitive && len(strings.TrimSpace(pe.String())) > 0 {
			isRelation = false
			viewParam := datamodeldiagram.EntityViewParam{
				EntityColor:  `orchid`,
				EntityHeader: `D`,
				EntityName:   entityName,
				IgnoredTypes: ignoredTypes,
				Types:        tMap,
			}
			v.DrawPrimitive(viewParam, pe.String(), relationshipMap)
		} else if seq := entityType.GetSequence(); seq != nil {

		} else if syslutil.HasPattern(entityType.Attrs, "empty") {
			if len(strings.Split(entityName, ".")) == 1 {
				entityName = appName + entityName
			}
			v.StringBuilder.WriteString(fmt.Sprintf("class \"%s\" as %s<< (D,orchid) >> {\n}\n", entityName,
				v.UniqueVarForAppName("", entityName)))
		} else if pe := entityType.GetEnum(); pe != nil {
			v.StringBuilder.WriteString(fmt.Sprintf("class \"%s enum\" as %s<< (D,orchid) >> {\n}\n", entityName,
				v.UniqueVarForAppName("", entityName)))
		}
	}
	if isRelation {
		v.DrawRelationship(relationshipMap, relationArrow)
	} else {
		v.DrawRelationship(relationshipMap, tupleArrow)
	}
	v.StringBuilder.WriteString("@enduml\n")
	return v.StringBuilder.String()
}

func RecurseivelyGetTypes(appName string, types map[string]*sysl.Type, m *sysl.Module) map[string]*sysl.Type {
	cummulative := make(map[string]*sysl.Type)
	for _, elem := range types {
		RecurseivelyGetTypesHelper(appName, elem, m, cummulative)
	}
	return cummulative
}

// RecurseivelyGetTypesHelper gets returns a type map of a type and all of its fields recursively.
func RecurseivelyGetTypesHelper(appName string, t *sysl.Type, m *sysl.Module, cummulative map[string]*sysl.Type) map[string]*sysl.Type {
	var typeName string
	if t == nil {
		return nil
	}
	ret := make(map[string]*sysl.Type)
	switch t.Type.(type) {
	case *sysl.Type_Enum_:
		return nil
	case *sysl.Type_Primitive_:
		return nil
	case *sysl.Type_Sequence:
		if path := t.GetSequence().GetTypeRef().GetRef().Path; len(path) > 1 {
			typeName = path[1]
			appName = path[0]
		} else {
			typeName = path[0]
		}

		appName, typeName, t = TypeFromRef(m, appName, t)
		if t != nil {
			ret[appName+"."+typeName] = t.GetSequence()
		}
	case *sysl.Type_TypeRef:
		if path := t.GetTypeRef().GetRef().Path; len(path) > 1 {
			typeName = path[1]
			appName = path[0]
		} else {
			typeName = path[0]
		}

		appName, typeName, t = TypeFromRef(m, appName, t)
		if t != nil {
			ret[appName+"."+typeName] = t

		}
	}
	tuple := t.GetTuple()
	if tuple == nil || tuple.AttrDefs == nil || len(tuple.AttrDefs) == 0 {
		for index, element := range ret {
			cummulative[index] = element
		}
		return ret
	}
	for _, ts := range tuple.AttrDefs {
		var newType *sysl.Type
		appName, typeName, newType = TypeFromRef(m, appName, ts)
		if newType == nil {
			continue
		}
		if _, ok := cummulative[appName+"."+typeName]; ok {
			continue
		}
		ret[appName+"."+typeName] = newType

		for key, v := range RecurseivelyGetTypesHelper(appName, ret[appName+"."+typeName], m, ret) {
			if _, ok := cummulative[key]; ok {
				continue
			}
			switch v.Type.(type) {
			case *sysl.Type_Primitive_:
				continue
			case *sysl.Type_TypeRef:
				typeName = strings.Join(v.GetTypeRef().GetRef().Path, "")
				appName, typeName, newType = TypeFromRef(m, appName, v)
				key = appName + "." + typeName
				if newType != nil {
					ret[appName+"."+typeName] = newType
				}
			case *sysl.Type_Tuple_:
				ret[key] = v
			case *sysl.Type_Enum_:
				ret[key] = v
			}
		}
	}
	for index, element := range ret {
		cummulative[index] = element
	}
	return ret
}

func TypeFromRef(mod *sysl.Module, appName string, t *sysl.Type) (string, string, *sysl.Type) {
	var typeName string
	if t == nil {
		return "", "", nil
	}
	switch t.Type.(type) {
	case *sysl.Type_Primitive_:
		return "", "", nil
	case *sysl.Type_Enum_, *sysl.Type_Tuple_:
		return appName, typeName, t
	case *sysl.Type_Sequence:
		ty := t.GetSequence()
		ref := ty.GetTypeRef().GetRef()
		if ref == nil {
			return "", "", nil // It's most likely a primitive type
		}
		if ref.Appname != nil {
			appName = strings.Join(ref.Appname.Part, "")
		}
		typeName = strings.Join(ref.Path, ".")
		if len(ref.Path) > 1 {
			appName = ref.Path[0]
			typeName = ref.Path[1]
		}
		if appName == "" {
			return "", "", nil
		}
		return appName, typeName, ty

	case *sysl.Type_TypeRef:
		ref := t.GetTypeRef().GetRef()
		if ref.Appname != nil {
			appName = strings.Join(ref.Appname.Part, "")
		}
		typeName = strings.Join(ref.Path, ".")
		if len(ref.Path) > 1 {
			appName = ref.Path[0]
			typeName = ref.Path[1]
		}
		if appName == "" {
			return "", "", nil
		}
		return appName, typeName, mod.Apps[appName].Types[typeName]
	}

	return "", "", nil
}

// GenerateDataModel takes all the types in parentAppName and generates data model diagrams for it
func GenerateDataModel(parentAppName string, t map[string]*sysl.Type) string {
	type datamodelCmd struct {
		diagrams.Plantumlmixin
		cmdutils.CmdContextParamDatagen
	}
	pl := &datamodelCmd{}
	pl.Project = ""
	pl.Direct = true
	pl.ClassFormat = "%(classname)"
	spclass := sequencediagram.ConstructFormatParser("", pl.ClassFormat)
	var stringBuilder strings.Builder
	dataParam := &DataModelParam{}
	v := datamodeldiagram.MakeDataModelView(spclass, dataParam.Mod, &stringBuilder, dataParam.Title, "")
	vNew := &DataModelView{
		DataModelView: *v,
	}
	return vNew.GenerateDataView(dataParam, parentAppName, t)
}
