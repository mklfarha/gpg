package field

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/helpers"
)

func JSONFieldTemplate(f entity.Field, e entity.Entity, dependant bool) Template {
	optional := ""
	if !f.Required {
		optional = "Optional"
	}

	jsonMany := f.JSONConfig.Type == entity.ManyJSONConfigType
	slice := ""
	if jsonMany {
		slice = "Slice"
	}

	if len(f.JSONConfig.Fields) == 0 && !f.JSONConfig.Reuse {
		stringTemplate := StringFieldTemplate(f, e)
		stringTemplate.Type = "json.RawMessage"
		stringTemplate.GenFieldType = "RawJSONFieldType"
		stringTemplate.GenRandomValue = "randomvalues.GetRandomRawJSONValue()"
		stringTemplate.JSONRaw = true
		stringTemplate.GraphGenToMapper = fmt.Sprintf("StringFromJsonRaw%s(%s)",
			optional,
			stringTemplate.GraphGenToMapper)

		if !f.Required {
			stringTemplate.GraphGenFromMapper = strings.ReplaceAll(stringTemplate.GraphGenFromMapper, "StringFromPointer", "JsonRawFromStringOptional")
			stringTemplate.GraphGenFromMapperOptional = strings.ReplaceAll(stringTemplate.GraphGenFromMapperOptional, "StringFromPointer", "JsonRawFromStringOptional")
		} else {
			stringTemplate.GraphGenFromMapper = fmt.Sprintf("JsonRawFromString(i.%s)", stringTemplate.Name)
			stringTemplate.GraphGenFromMapperOptional = fmt.Sprintf("JsonRawFromStringOptional(i.%s)", stringTemplate.Name)
		}
		stringTemplate.ProtoToMapper = fmt.Sprintf("JSONRawToString(%s)", stringTemplate.ProtoToMapper)
		stringTemplate.ProtoFromMapper = fmt.Sprintf("StringToJSONRaw(%s)", stringTemplate.ProtoFromMapper)
		return stringTemplate
	}

	template := BaseFieldTemplate(f, e)

	fullName := f.JSONConfig.Identifier

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
	template.RepoFromMapper = fmt.Sprintf("%s.%sFromJSON(model.%s)", f.JSONConfig.Identifier, helpers.ToCamelCase(fullName), helpers.ToCamelCase(f.Identifier))

	// graph
	graphOutType := fmt.Sprintf("%s%s", helpers.ToCamelCase(fullName), template.GraphRequired)
	graphInType := fmt.Sprintf("%sInput%s", helpers.ToCamelCase(fullName), template.GraphRequired)
	graphInTypeOptional := fmt.Sprintf("%sInput", helpers.ToCamelCase(fullName))
	if jsonMany {
		graphInType = fmt.Sprintf("[%s]", graphInType)
		graphInTypeOptional = fmt.Sprintf("[%s]", graphInTypeOptional)
		graphOutType = fmt.Sprintf("[%s]", graphOutType)
	}

	template.GraphInType = fmt.Sprintf("%s%s", graphInType, template.GraphRequired)
	template.GraphInTypeOptional = graphInTypeOptional
	template.GraphOutType = fmt.Sprintf("%s%s", graphOutType, template.GraphRequired)
	template.GraphGenType = helpers.ToCamelCase(fullName)
	template.GraphGenToMapper = fmt.Sprintf("Map%s%s%s(i.%s)",
		helpers.ToCamelCase(fullName),
		slice,
		optional,
		helpers.ToCamelCase(f.Identifier),
	)
	template.GraphGenFromMapperParam = ""
	template.GraphGenFromMapper = fmt.Sprintf("Map%s%sInput(i.%s)",
		helpers.ToCamelCase(fullName),
		slice,
		helpers.ToCamelCase(f.Identifier))
	template.GraphGenFromMapperOptional = fmt.Sprintf("Map%s%sInputOptional(i.%s)",
		helpers.ToCamelCase(fullName),
		slice,
		helpers.ToCamelCase(f.Identifier))

	// proto
	template.ProtoType = helpers.ToCamelCase(fullName)
	template.ProtoToMapper = fmt.Sprintf("%s%sToProto(e.%s)",
		helpers.ToCamelCase(fullName),
		slice,
		helpers.ToCamelCase(f.Identifier))
	template.ProtoFromMapper = fmt.Sprintf("%s%sFromProto(m.Get%s())",
		helpers.ToCamelCase(fullName),
		slice,
		strcase.ToCamel(f.Identifier))

	return template
}
