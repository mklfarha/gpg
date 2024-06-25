package field

import (
	"fmt"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/helpers"
)

func OptionsSingleFieldTemplate(f entity.Field, e entity.Entity, prefix *string) Template {
	pl := pluralize.NewClient()
	name := f.Identifier
	if prefix != nil {
		name = fmt.Sprintf("%s_%s", pl.Singular(*prefix), f.Identifier)
	}
	graphRequired := ""
	graphGenToMapper := fmt.Sprintf("i.%s.StringPtr()", helpers.ToCamelCase(f.Identifier))
	if f.Required {
		graphRequired = "!"
		graphGenToMapper = fmt.Sprintf("i.%s.String()", helpers.ToCamelCase(f.Identifier))
	}

	protoType := helpers.ToCamelCase(fmt.Sprintf("%s_%s", e.Identifier, pl.Singular(f.Identifier)))
	return Template{
		Identifier:                 f.Identifier,
		Name:                       pl.Singular(helpers.ToCamelCase(f.Identifier)),
		Type:                       pl.Singular(helpers.ToCamelCase(name)),
		EntityIdentifier:           e.Identifier,
		InternalType:               entity.OptionsSingleFieldType,
		GenFieldType:               "SingleEnumFieldType",
		IsPrimary:                  f.StorageConfig.PrimaryKey,
		Required:                   f.Required,
		Tags:                       helpers.ResolveTags(f),
		Import:                     nil,
		Custom:                     true,
		Generated:                  f.Autogenerated.Type != entity.InvalidAutogeneratedType,
		GeneratedFuncInsert:        resolveGeneratedFuncInsert(e, f),
		GeneratedFuncUpdate:        resolveGeneratedFuncUpdate(e, f),
		Enum:                       true,
		RepoToMapper:               ".ToInt32()",
		RepoFromMapper:             fmt.Sprintf("entity.%s(model.%s)", helpers.ToCamelCase(f.Identifier), helpers.ToCamelCase(f.Identifier)),
		GraphName:                  f.Identifier,
		GraphModelName:             helpers.ToCamelCase(f.Identifier),
		GraphInType:                fmt.Sprintf("String%s", graphRequired),
		GraphInTypeOptional:        "String",
		GraphOutType:               fmt.Sprintf("String%s", graphRequired),
		GraphGenType:               "string",
		GraphGenToMapper:           graphGenToMapper,
		GraphGenFromMapperParam:    fmt.Sprintf("%sentity.%sFromString(%s)", e.Identifier, helpers.ToCamelCase(f.Identifier), f.Identifier),
		GraphGenFromMapper:         fmt.Sprintf("%s.%sFromString(i.%s)", e.Identifier, pl.Singular(helpers.ToCamelCase(name)), helpers.ToCamelCase(f.Identifier)),
		GraphGenFromMapperOptional: fmt.Sprintf("%s.%sFromPointerString(i.%s)", e.Identifier, pl.Singular(helpers.ToCamelCase(name)), helpers.ToCamelCase(f.Identifier)),
		ProtoType:                  protoType,
		ProtoName:                  helpers.ToSnakeCase(f.Identifier),
		ProtoEnumOptions:           helpers.ProtoEnumOptions(protoType, f.OptionValues),
		ProtoToMapper:              fmt.Sprintf("pb.%s(e.%s)", protoType, pl.Singular(helpers.ToCamelCase(f.Identifier))),
		ProtoFromMapper:            fmt.Sprintf("entity.%s(m.Get%s())", pl.Singular(helpers.ToCamelCase(name)), strcase.ToCamel(f.Identifier)),
		ProtoGenName:               strcase.ToCamel(f.Identifier),
	}
}
