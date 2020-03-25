package templategeneration

import (
	"fmt"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/pkg/syslutil"

	"github.com/anz-bank/sysl/pkg/datamodeldiagram"
	"github.com/anz-bank/sysl/pkg/integrationdiagram"
	"github.com/anz-bank/sysl/pkg/sysl"
)

const relationArrow = `}--`
const tupleArrow = `*--`

type DataModelView struct {
	datamodeldiagram.DataModelView
}

type DataModelParam struct {
	datamodeldiagram.DataModelParam
}

func (v *DataModelView) GenerateDataView(dataParam *DataModelParam, appName string, t *sysl.Type, p Project) string {
	var isRelation bool
	relationshipMap := map[string]map[string]datamodeldiagram.RelationshipParam{}
	v.StringBuilder.WriteString("@startuml\n")
	if dataParam.Title != "" {
		fmt.Fprintf(v.StringBuilder, "title %s\n", dataParam.Title)
	}
	v.StringBuilder.WriteString(integrationdiagram.PumlHeader)

	// sort and iterate over each entity type the selected application
	// *Type_Tuple_ OR *Type_Relation_
	typeMap := map[string]*sysl.Type{}
	ignoredTypes := map[string]struct{}{}
	// typeMap := dataParam.App.GetTypes()
	// TODO: Actually put The appName/project name and the appName in a struct so strings.split and join dont need to be used
	entityNames := []string{}
	typeMap = RecurseivelyGetTypes(appName, t, p)
	for key := range typeMap {
		entityNames = append(entityNames, key)
	}
	sort.Strings(entityNames)
	for _, entityName := range entityNames {
		entityType := typeMap[entityName]
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
			}
			v.DrawTuple(viewParam, tupEntity, relationshipMap)
		} else if pe := entityType.GetPrimitive(); pe != sysl.Type_NO_Primitive && len(strings.TrimSpace(pe.String())) > 0 {
			isRelation = false
			viewParam := datamodeldiagram.EntityViewParam{
				EntityColor:  `orchid`,
				EntityHeader: `D`,
				EntityName:   entityName,
			}
			v.DrawPrimitive(viewParam, pe.String(), relationshipMap)
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

// RecurseivelyGetTypes gets returns a type map of a type and all of its fields recursively.
func RecurseivelyGetTypes(appName string, t *sysl.Type, p Project) map[string]*sysl.Type {
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
	case *sysl.Type_TypeRef:
		typeName = strings.Join(t.GetTypeRef().GetRef().Path, "")
		appName, typeName, t = TypeFromRef(p.Module, appName, t)
		ret[appName+"."+typeName] = t
	}
	tuple := t.GetTuple()
	if tuple == nil || tuple.AttrDefs == nil || len(tuple.AttrDefs) == 0 {
		return nil
	}
	for _, ts := range tuple.AttrDefs {
		var newType *sysl.Type
		appName, typeName, newType = TypeFromRef(p.Module, appName, ts)
		if newType == nil {
			continue
		}
		ret[appName+"."+typeName] = newType
		for key, v := range RecurseivelyGetTypes(appName, ret[appName+"."+typeName], p) {
			switch v.Type.(type) {
			case *sysl.Type_Primitive_:
				continue
			case *sysl.Type_TypeRef:
				typeName = strings.Join(v.GetTypeRef().GetRef().Path, "")
				appName, typeName, newType = TypeFromRef(p.Module, appName, v)
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
		fmt.Println(appName, typeName)
		if appName == "" {
			return "", "", nil
		}
		return appName, typeName, mod.Apps[appName].Types[typeName]
	}
	return "", "", nil
}
