package field

import (
	"fmt"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/helpers"
)

func IntFieldTemplate(f entity.Field, e entity.Entity) Template {
	pl := pluralize.NewClient()
	graphRequired := ""
	if f.Required {
		graphRequired = "!"
	}
	graphGenToMapper := fmt.Sprintf("int(i.%s)", helpers.ToCamelCase(f.Identifier))
	graphGenFromMapper := fmt.Sprintf("int32(i.%s)", helpers.ToCamelCase(f.Identifier))
	graphGenFromMapperOptional := fmt.Sprintf("IntFromPointer(i.%s)", helpers.ToCamelCase(f.Identifier))
	if !f.Required {
		graphGenToMapper = fmt.Sprintf("IntPointer(i.%s)", helpers.ToCamelCase(f.Identifier))
		graphGenFromMapper = graphGenFromMapperOptional
	}
	return Template{
		Identifier:                 f.Identifier,
		SingularIdentifier:         pl.Singular(f.Identifier),
		Name:                       helpers.ToCamelCase(f.Identifier),
		Type:                       "int32",
		EntityIdentifier:           e.Identifier,
		InternalType:               entity.IntFieldType,
		GenFieldType:               "IntFieldType",
		GenRandomValue:             "randomvalues.GetRandomIntValue()",
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
		GraphInType:                fmt.Sprintf("Int%s", graphRequired),
		GraphInTypeOptional:        "Int",
		GraphOutType:               fmt.Sprintf("Int%s", graphRequired),
		GraphGenType:               "int32",
		GraphGenToMapper:           graphGenToMapper,
		GraphGenFromMapperParam:    f.Identifier,
		GraphGenFromMapper:         graphGenFromMapper,
		GraphGenFromMapperOptional: graphGenFromMapperOptional,
		ProtoType:                  "int64",
		ProtoName:                  helpers.ToSnakeCase(f.Identifier),
		ProtoGenName:               strcase.ToCamel(f.Identifier),
		ProtoToMapper:              fmt.Sprintf("int64(e.%s)", helpers.ToCamelCase(f.Identifier)),
		ProtoFromMapper:            fmt.Sprintf("int32(m.Get%s())", strcase.ToCamel(f.Identifier)),
	}
}
