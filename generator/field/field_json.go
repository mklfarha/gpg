package field

import (
	"fmt"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/helpers"
)

func JSONFieldTemplate(f entity.Field, e entity.Entity, dependant bool) Template {
	if len(f.JSONConfig.Fields) == 0 {
		stringTemplate := StringFieldTemplate(f, e)
		stringTemplate.Type = "json.RawMessage"
		stringTemplate.GenFieldType = "RawJSONFieldType"
		stringTemplate.GenRandomValue = "randomvalues.GetRandomRawJSONValue()"
		stringTemplate.JSONRaw = true
		stringTemplate.GraphGenToMapper = fmt.Sprintf("StringFromJsonRaw(%s)", stringTemplate.GraphGenToMapper)
		stringTemplate.GraphGenFromMapper = strings.ReplaceAll(stringTemplate.GraphGenFromMapper, "StringFromPointer", "JsonRawFromString")
		stringTemplate.GraphGenFromMapperOptional = strings.ReplaceAll(stringTemplate.GraphGenFromMapperOptional, "StringFromPointer", "JsonRawFromString")
		stringTemplate.ProtoToMapper = fmt.Sprintf("mapper.JSONRawToString(%s)", stringTemplate.ProtoToMapper)
		stringTemplate.ProtoFromMapper = fmt.Sprintf("mapper.StringToJSONRaw(%s)", stringTemplate.ProtoFromMapper)
		return stringTemplate
	}

	jsonMany := f.JSONConfig.Type == entity.ManyJSONConfigType

	template := BaseFieldTemplate(f, e)

	fullName := template.SingularIdentifier
	pl := pluralize.NewClient()
	if dependant {
		fullName = fmt.Sprintf("%s_%s", pl.Singular(e.Identifier), template.SingularIdentifier)
	}

	//base
	template.Type = helpers.ToCamelCase(template.SingularIdentifier)
	template.InternalType = entity.JSONFieldType
	template.JSON = true
	template.JSONMany = jsonMany
	template.Custom = true
	genFieldType := "SingleDependantEntityFieldType"
	genRandomValue := fmt.Sprintf("New%sWithRandomValues()", template.Type)
	if jsonMany {
		genFieldType = "MultiDependantEntityFieldType"
		genRandomValue = fmt.Sprintf("New%sSliceWithRandomValues(rand.Intn(10))", template.Type)
	}
	template.GenFieldType = genFieldType
	template.GenRandomValue = genRandomValue
	template.RepoToMapper = fmt.Sprintf("entity.%sSliceToJSON(req.%s.%s)", helpers.ToCamelCase(fullName), helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(f.Identifier))
	template.RepoFromMapper = fmt.Sprintf("entity.%sFromJSON(model.%s)", helpers.ToCamelCase(template.SingularIdentifier), helpers.ToCamelCase(f.Identifier))

	// graph
	graphInTypeOptional := fmt.Sprintf("%s%sInput", helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(template.SingularIdentifier))
	graphOutType := fmt.Sprintf("%s%s", helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(template.SingularIdentifier))
	if jsonMany {
		graphInTypeOptional = fmt.Sprintf("[%s]", graphInTypeOptional)
		graphOutType = fmt.Sprintf("[%s]", graphOutType)
	}

	template.GraphInType = fmt.Sprintf("%s%s", graphInTypeOptional, template.GraphRequired)
	template.GraphInTypeOptional = graphInTypeOptional
	template.GraphOutType = fmt.Sprintf("%s%s", graphOutType, template.GraphRequired)
	template.GraphGenType = helpers.ToCamelCase(fullName)
	template.GraphGenToMapper = fmt.Sprintf("Map%s%s(i.%s)", helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(fullName), helpers.ToCamelCase(f.Identifier))
	template.GraphGenFromMapperParam = ""
	template.GraphGenFromMapper = fmt.Sprintf("Map%s%sInput(i.%s)", helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(fullName), helpers.ToCamelCase(f.Identifier))
	template.GraphGenFromMapperOptional = fmt.Sprintf("Map%s%sInput(i.%s)", helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(fullName), helpers.ToCamelCase(f.Identifier))

	// proto
	template.ProtoType = helpers.ToCamelCase(fullName)
	template.ProtoToMapper = fmt.Sprintf("%sSliceToProto(e.%s)", helpers.ToCamelCase(f.Identifier), helpers.ToCamelCase(f.Identifier))
	template.ProtoFromMapper = fmt.Sprintf("%sSliceFromProto(m.Get%s())", helpers.ToCamelCase(f.Identifier), strcase.ToCamel(f.Identifier))

	return template
}
