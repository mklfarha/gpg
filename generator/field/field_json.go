package field

import (
	"fmt"

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
		stringTemplate.Type = "string"
		stringTemplate.GenFieldType = "RawJSONFieldType"
		stringTemplate.GenRandomValue = "randomvalues.GetRandomRawJSONValue()"
		stringTemplate.JSONRaw = true
		stringTemplate.RepoToMapper = fmt.Sprintf("[]byte(req.%s.%s)", helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(f.Identifier))
		stringTemplate.RepoFromMapper = fmt.Sprintf("string(model.%s)", helpers.ToCamelCase(f.Identifier))
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
	template.RepoToMapper = fmt.Sprintf("%s.%s%sToJSON(req.%s.%s)",
		f.JSONConfig.Identifier,
		helpers.ToCamelCase(fullName),
		slice,
		helpers.ToCamelCase(e.Identifier),
		helpers.ToCamelCase(f.Identifier))
	template.RepoFromMapper = fmt.Sprintf("%s.%s%sFromJSON(model.%s)",
		f.JSONConfig.Identifier,
		helpers.ToCamelCase(fullName),
		slice,
		helpers.ToCamelCase(f.Identifier))

	// graph
	graphOutType := fmt.Sprintf("%s%s", helpers.ToCamelCase(fullName), template.GraphRequired)
	graphInType := fmt.Sprintf("%sInput%s", helpers.ToCamelCase(fullName), template.GraphRequired)
	graphInTypeOptional := fmt.Sprintf("%sInput", helpers.ToCamelCase(fullName))
	if jsonMany {
		graphInType = fmt.Sprintf("[%s]%s", graphInType, template.GraphRequired)
		graphInTypeOptional = fmt.Sprintf("[%s]", graphInTypeOptional)
		graphOutType = fmt.Sprintf("[%s]%s", graphOutType, template.GraphRequired)
	}

	template.GraphOutType = graphOutType
	template.GraphInType = graphInType
	template.GraphInTypeOptional = graphInTypeOptional

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
		template.GraphModelName)
	template.GraphGenFromMapperOptional = fmt.Sprintf("Map%s%sInputOptional(i.%s)",
		helpers.ToCamelCase(fullName),
		slice,
		template.GraphModelName)

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
