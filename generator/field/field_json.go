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
	pl := pluralize.NewClient()
	singularIdentifier := pl.Singular(f.Identifier)

	fullName := singularIdentifier
	if dependant {
		fullName = fmt.Sprintf("%s_%s", pl.Singular(e.Identifier), singularIdentifier)
	}
	graphRequired := ""
	if f.Required {
		graphRequired = "!"
	}

	jsonMany := f.JSONConfig.Type == entity.ManyJSONConfigType
	fieldType := helpers.ToCamelCase(singularIdentifier)
	graphInTypeOptional := fmt.Sprintf("%s%sInput", helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(singularIdentifier))
	graphOutType := fmt.Sprintf("%s%s", helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(singularIdentifier))
	protoType := helpers.ToCamelCase(fullName)
	genFieldType := "SingleDependantEntityFieldType"
	genRandomValue := fmt.Sprintf("New%sWithRandomValues()", fieldType)
	if jsonMany {
		graphInTypeOptional = fmt.Sprintf("[%s]", graphInTypeOptional)
		graphOutType = fmt.Sprintf("[%s]", graphOutType)
		genFieldType = "MultiDependantEntityFieldType"
		genRandomValue = fmt.Sprintf("New%sSliceWithRandomValues(rand.Intn(10))", fieldType)
	}

	graphInType := fmt.Sprintf("%s%s", graphInTypeOptional, graphRequired)
	graphOutType = fmt.Sprintf("%s%s", graphOutType, graphRequired)

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

	return Template{
		Identifier:                 f.Identifier,
		SingularIdentifier:         singularIdentifier,
		Name:                       helpers.ToCamelCase(f.Identifier),
		Type:                       fieldType,
		EntityIdentifier:           e.Identifier,
		InternalType:               entity.JSONFieldType,
		GenFieldType:               genFieldType,
		GenRandomValue:             genRandomValue,
		IsPrimary:                  f.StorageConfig.PrimaryKey,
		Required:                   f.Required,
		Tags:                       helpers.ResolveTags(f),
		Import:                     nil,
		JSON:                       true,
		JSONMany:                   jsonMany,
		Custom:                     true,
		Generated:                  f.Autogenerated.Type != entity.InvalidAutogeneratedType,
		GeneratedFuncInsert:        resolveGeneratedFuncInsert(e, f),
		GeneratedFuncUpdate:        resolveGeneratedFuncUpdate(e, f),
		RepoToMapper:               fmt.Sprintf("entity.%sSliceToJSON(req.%s.%s)", helpers.ToCamelCase(fullName), helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(f.Identifier)),
		RepoFromMapper:             fmt.Sprintf("entity.%sFromJSON(model.%s)", helpers.ToCamelCase(singularIdentifier), helpers.ToCamelCase(f.Identifier)),
		GraphName:                  f.Identifier,
		GraphModelName:             helpers.ToCamelCase(f.Identifier),
		GraphInType:                graphInType,
		GraphInTypeOptional:        graphInTypeOptional,
		GraphOutType:               graphOutType,
		GraphGenType:               helpers.ToCamelCase(fullName),
		GraphGenToMapper:           fmt.Sprintf("Map%s%s(i.%s)", helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(fullName), helpers.ToCamelCase(f.Identifier)),
		GraphGenFromMapperParam:    "",
		GraphGenFromMapper:         fmt.Sprintf("Map%s%sInput(i.%s)", helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(fullName), helpers.ToCamelCase(f.Identifier)),
		GraphGenFromMapperOptional: fmt.Sprintf("Map%s%sInput(i.%s)", helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(fullName), helpers.ToCamelCase(f.Identifier)),
		ProtoType:                  protoType,
		ProtoName:                  helpers.ToSnakeCase(f.Identifier),
		ProtoGenName:               strcase.ToCamel(f.Identifier),
		ProtoToMapper:              fmt.Sprintf("%sSliceToProto(e.%s)", helpers.ToCamelCase(f.Identifier), helpers.ToCamelCase(f.Identifier)),
		ProtoFromMapper:            fmt.Sprintf("%sSliceFromProto(m.Get%s())", helpers.ToCamelCase(f.Identifier), strcase.ToCamel(f.Identifier)),
	}
}
