package field

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/helpers"
)

func JSONFieldTemplate(f entity.Field, e entity.Entity, dependant bool) Template {
	if len(f.JSONConfig.Fields) == 0 && !f.JSONConfig.Reuse {
		stringTemplate := StringFieldTemplate(f, e)
		stringTemplate.Type = "json.RawMessage"
		stringTemplate.GenFieldType = "RawJSONFieldType"
		stringTemplate.GenRandomValue = "randomvalues.GetRandomRawJSONValue()"
		stringTemplate.JSONRaw = true
		stringTemplate.GraphGenToMapper = fmt.Sprintf("StringFromJsonRaw(%s)", stringTemplate.GraphGenToMapper)
		stringTemplate.GraphGenFromMapper = strings.ReplaceAll(stringTemplate.GraphGenFromMapper, "StringFromPointer", "JsonRawFromString")
		stringTemplate.GraphGenFromMapperOptional = strings.ReplaceAll(stringTemplate.GraphGenFromMapperOptional, "StringFromPointer", "JsonRawFromString")
		stringTemplate.ProtoToMapper = fmt.Sprintf("JSONRawToString(%s)", stringTemplate.ProtoToMapper)
		stringTemplate.ProtoFromMapper = fmt.Sprintf("StringToJSONRaw(%s)", stringTemplate.ProtoFromMapper)
		return stringTemplate
	}

	jsonMany := f.JSONConfig.Type == entity.ManyJSONConfigType

	template := BaseFieldTemplate(f, e)

	fullName := template.SingularIdentifier
	if dependant {
		fullName = f.JSONConfig.Identifier
	}

	//base
	template.Type = helpers.ToCamelCase(fullName)
	template.InternalType = entity.JSONFieldType
	template.JSON = true
	template.JSONMany = jsonMany
	template.JSONIdentifier = f.JSONConfig.Identifier
	template.Custom = true
	genFieldType := "SingleDependantEntityFieldType"
	genRandomValue := fmt.Sprintf("%s.New%sWithRandomValues()", f.JSONConfig.Identifier, template.Type)
	if jsonMany {
		genFieldType = "MultiDependantEntityFieldType"
		genRandomValue = fmt.Sprintf("%s.New%sSliceWithRandomValues(rand.Intn(10))", f.JSONConfig.Identifier, template.Type)
	}
	template.GenFieldType = genFieldType
	template.GenRandomValue = genRandomValue
	template.RepoToMapper = fmt.Sprintf("%s.%sSliceToJSON(req.%s.%s)", f.JSONConfig.Identifier, helpers.ToCamelCase(fullName), helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(f.Identifier))
	template.RepoFromMapper = fmt.Sprintf("%s.%sFromJSON(model.%s)", f.JSONConfig.Identifier, helpers.ToCamelCase(template.SingularIdentifier), helpers.ToCamelCase(f.Identifier))

	// graph
	graphInTypeOptional := fmt.Sprintf("%sInput", helpers.ToCamelCase(fullName))
	graphOutType := helpers.ToCamelCase(fullName)
	if jsonMany {
		graphInTypeOptional = fmt.Sprintf("[%s]", graphInTypeOptional)
		graphOutType = fmt.Sprintf("[%s]", graphOutType)
	}

	template.GraphInType = fmt.Sprintf("%s%s", graphInTypeOptional, template.GraphRequired)
	template.GraphInTypeOptional = graphInTypeOptional
	template.GraphOutType = fmt.Sprintf("%s%s", graphOutType, template.GraphRequired)
	template.GraphGenType = helpers.ToCamelCase(fullName)
	template.GraphGenToMapper = fmt.Sprintf("Map%s(i.%s)", helpers.ToCamelCase(fullName), helpers.ToCamelCase(f.Identifier))
	if jsonMany {
		template.GraphGenToMapper = fmt.Sprintf("Map%sSlice(i.%s)", helpers.ToCamelCase(fullName), helpers.ToCamelCase(f.Identifier))
	}
	template.GraphGenFromMapperParam = ""
	template.GraphGenFromMapper = fmt.Sprintf("Map%sInput(i.%s)", helpers.ToCamelCase(fullName), helpers.ToCamelCase(f.Identifier))
	template.GraphGenFromMapperOptional = fmt.Sprintf("Map%sInput(i.%s)", helpers.ToCamelCase(fullName), helpers.ToCamelCase(f.Identifier))

	// proto
	template.ProtoType = helpers.ToCamelCase(fullName)
	template.ProtoToMapper = fmt.Sprintf("%sSliceToProto(e.%s)", helpers.ToCamelCase(fullName), helpers.ToCamelCase(f.Identifier))
	template.ProtoFromMapper = fmt.Sprintf("%sSliceFromProto(m.Get%s())", helpers.ToCamelCase(fullName), strcase.ToCamel(f.Identifier))

	return template
}
