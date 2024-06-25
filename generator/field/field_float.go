package field

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/helpers"
)

func FloatFieldTemplate(f entity.Field, e entity.Entity) Template {
	graphRequired := ""
	if f.Required {
		graphRequired = "!"
	}

	graphGenToMapper := fmt.Sprintf("i.%s", helpers.ToCamelCase(f.Identifier))
	graphGenFromMapper := fmt.Sprintf("i.%s", helpers.ToCamelCase(f.Identifier))
	graphGenFromMapperOptional := fmt.Sprintf("FloatFromPointer(i.%s)", helpers.ToCamelCase(f.Identifier))
	graphGenFromMapperParam := f.Identifier
	graphGenType := "float64"
	if !f.Required {
		graphGenToMapper = fmt.Sprintf("&i.%s", helpers.ToCamelCase(f.Identifier))
		graphGenFromMapper = graphGenFromMapperOptional
		graphGenFromMapperParam = fmt.Sprintf("mapper.FloatFromPointer(%s)", f.Identifier)
		graphGenType = "*float64"
	}
	return Template{
		Identifier:                 f.Identifier,
		Name:                       helpers.ToCamelCase(f.Identifier),
		Type:                       "float64",
		EntityIdentifier:           e.Identifier,
		InternalType:               entity.FloatFieldType,
		GenFieldType:               "FloatFieldType",
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
		GraphInType:                fmt.Sprintf("Float%s", graphRequired),
		GraphInTypeOptional:        "Float",
		GraphOutType:               fmt.Sprintf("Float%s", graphRequired),
		GraphGenType:               graphGenType,
		GraphGenToMapper:           graphGenToMapper,
		GraphGenFromMapperParam:    graphGenFromMapperParam,
		GraphGenFromMapper:         graphGenFromMapper,
		GraphGenFromMapperOptional: graphGenFromMapperOptional,
		ProtoType:                  "double",
		ProtoName:                  helpers.ToSnakeCase(f.Identifier),
		ProtoGenName:               strcase.ToCamel(f.Identifier),
	}
}
