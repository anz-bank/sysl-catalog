package catalog

import (
	"fmt"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
)

// DataModelReturnTable prints out a markdown table for a given statement and endpoint
func (p *Generator) DataModelReturnTable(appName string, stmt *sysl.Statement, endpoint *sysl.Endpoint) string {
	appName, typeName, _, _ := p.ExtractReturnInfo(appName, stmt, endpoint)
	if appName == "" && typeName == "" {
		return ""
	}
	return p.DataModelTable(appName, typeName, "")
}

// DataModelAliasTable prints out a markdown table for a given application and parameter
func (p *Generator) DataModelAliasTable(app *sysl.Application, param Param) string {
	defer func() {
		if err := recover(); err != nil {
			p.Log.Errorf("error creating param data model: %s", err)
		}
	}()
	info, typeName, aliasTypeName, _ := p.ExtractTypeInfo(app, param)
	return p.DataModelTable(info, typeName, aliasTypeName)
}

// DataModelTable prints out a markdown table which describes a type
func (p *Generator) DataModelTable(appName, typeName, aliasName string) string {
	var markdownTable string

	if typeName == "" {
		return ""
	}

	if appName == "primitive" {
		return printPrimitiveTable(typeName, aliasName)
	}

	simpleType, ok := p.Mapper.SimpleTypes[appName+"."+typeName]
	if !ok {
		p.Log.Errorf("Unable to find type: %s.%s with alias %s", appName, typeName, aliasName)
		return ""
	}

	markdownTable += fmt.Sprintf("%s %s \n", simpleType.Type, typeName)
	markdownTable += "| Field name | Type | Description |\n"
	markdownTable += "|----|----|----|\n"

	switch simpleType.Type {
	case "enum":
		var enumKeys []int
		for k := range simpleType.Enum {
			enumKeys = append(enumKeys, int(k))
		}
		sort.Ints(enumKeys)
		for _, k := range enumKeys {
			markdownTable += fmt.Sprintf("| %d | %s | %s |\n", k, simpleType.Enum[int64(k)], "")
		}
	case "tuple", "map", "relation":
		for _, fieldName := range SortedKeys(simpleType.Properties) {
			field := simpleType.Properties[fieldName]
			switch simpleType.Properties[fieldName].Type {
			case "ref":
				reference, ok := p.Mapper.SimpleTypes[field.Reference]
				if !ok {
					p.Log.Errorf("Unable to find type: %s with alias %s", field.Reference, fieldName)
					markdownTable += fmt.Sprintf("| %s | %s | %s |\n", fieldName, field.Reference, field.Description)
				} else {
					markdownTable += fmt.Sprintf("| %s | %s (%s) | %s |\n", fieldName, convertReferenceToLink(field.Reference), reference.Type, field.Description)
				}
			case "list":
				markdownTable += fmt.Sprintf("| %s | sequence of %s | %s |\n", fieldName, convertReferenceToLink(field.Items[0].Reference), field.Description)
			default:
				markdownTable += fmt.Sprintf("| %s | %s | %s |\n", fieldName, field.Type, field.Description)
			}
		}
	case "list":
		for _, field := range simpleType.Items {
			if field.Type == "ref" {
				markdownTable += fmt.Sprintf("| %s | sequence of %s | %s |\n", typeName, convertReferenceToLink(field.Reference), field.Description)
			} else {
				markdownTable += fmt.Sprintf("| %s | sequence of %s | %s |\n", typeName, field.Type, field.Description)
			}
		}
	case "ref":
		markdownTable += fmt.Sprintf("| %s | %s | %s |\n", typeName, convertReferenceToLink(simpleType.Reference), simpleType.Description)
	}

	return markdownTable
}

func printPrimitiveTable(aliasName, typeName string) (markdownTable string) {
	lowerCaseTypeName := strings.ToLower(typeName)
	markdownTable += aliasName + "<:" + lowerCaseTypeName + "\n"
	markdownTable += "| Field name | Type | Description |\n"
	markdownTable += "|----|----|----|\n"
	markdownTable += fmt.Sprintf("| %s | %s | %s |\n", aliasName, lowerCaseTypeName, "")
	return
}

func convertReferenceToLink(reference string) string {
	safeReference := SanitiseOutputName(reference)
	typeName := strings.Split(reference, ".")
	if len(typeName) < 2 {
		return "" // Not a valid reference
	}
	return fmt.Sprintf(`<a href="#%s">%s</a>`, safeReference, typeName[1])
}
