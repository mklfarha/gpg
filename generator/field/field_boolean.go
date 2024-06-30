package field

import (
	"fmt"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/helpers"
)

func BooleanFieldTemplate(f entity.Field, e entity.Entity) Template {
	pl := pluralize.NewClient()
	graphRequired := ""
	if f.Required {
		graphRequired = "!"
	}
	graphGenToMapper := fmt.Sprintf("&i.%s", helpers.ToCamelCase(f.Identifier))
	if f.Required {
		graphGenToMapper = fmt.Sprintf("i.%s", helpers.ToCamelCase(f.Identifier))
	}
	graphGenFromMapper := fmt.Sprintf("i.%s", helpers.ToCamelCase(f.Identifier))
	graphGenFromMapperOptional := fmt.Sprintf("BoolFromPointer(i.%s)", helpers.ToCamelCase(f.Identifier))

	return Template{
		Identifier:                 f.Identifier,
		SingularIdentifier:         pl.Singular(f.Identifier),
		Name:                       helpers.ToCamelCase(f.Identifier),
		Type:                       "bool",
		EntityIdentifier:           e.Identifier,
		InternalType:               entity.BooleanFieldType,
		GenFieldType:               "BooleanFieldType",
		GenRandomValue:             "randomvalues.GetRandomBoolValue()",
		IsPrimary:                  f.StorageConfig.PrimaryKey,
		Required:                   f.Required,
		Tags:                       helpers.ResolveTags(f),
		Import:                     nil,
		Custom:                     false,
		Generated:                  f.Autogenerated.Type != entity.InvalidAutogeneratedType,
		GeneratedFuncInsert:        resolveGeneratedFuncInsert(e, f),
		GeneratedFuncUpdate:        resolveGeneratedFuncUpdate(e, f),
		RepoToMapper:               "",
		RepoFromMapper:             fmt.Sprintf("model.%s", helpers.ToCamelCase(f.Identifier)),
		GraphName:                  f.Identifier,
		GraphModelName:             helpers.ToCamelCase(f.Identifier),
		GraphInType:                fmt.Sprintf("Boolean%s", graphRequired),
		GraphInTypeOptional:        "Boolean",
		GraphOutType:               fmt.Sprintf("Boolean%s", graphRequired),
		GraphGenType:               "bool",
		GraphGenToMapper:           graphGenToMapper,
		GraphGenFromMapperParam:    f.Identifier,
		GraphGenFromMapper:         graphGenFromMapper,
		GraphGenFromMapperOptional: graphGenFromMapperOptional,
		ProtoType:                  "bool",
		ProtoName:                  helpers.ToSnakeCase(f.Identifier),
		ProtoToMapper:              fmt.Sprintf("e.%s", helpers.ToCamelCase(f.Identifier)),
		ProtoFromMapper:            fmt.Sprintf("m.Get%s()", strcase.ToCamel(f.Identifier)),
		ProtoGenName:               strcase.ToCamel(f.Identifier),
	}
}
