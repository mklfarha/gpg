package field

import (
	"fmt"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/helpers"
)

func OptionsManyFieldTemplate(f entity.Field, e entity.Entity, dependantEntity *Template) Template {
	name := f.Identifier
	pl := pluralize.NewClient()
	if dependantEntity != nil && dependantEntity.EntityIdentifier != "" {
		name = fmt.Sprintf("%s_%s", pl.Singular(dependantEntity.Identifier), f.Identifier)
	}
	graphRequired := ""
	if f.Required {
		graphRequired = "!"
	}

	protoType := helpers.ToCamelCase(fmt.Sprintf("%s_%s", pl.Singular(e.Identifier), pl.Singular(f.Identifier)))
	if dependantEntity != nil && dependantEntity.EntityIdentifier != "" {
		protoType = helpers.ToCamelCase(fmt.Sprintf("%s_%s", dependantEntity.EntityIdentifier, protoType))
	}

	return Template{
		Identifier:          f.Identifier,
		SingularIdentifier:  pl.Singular(f.Identifier),
		Name:                helpers.ToCamelCase(f.Identifier),
		Type:                pl.Singular(helpers.ToCamelCase(name)),
		EntityIdentifier:    e.Identifier,
		InternalType:        entity.OptionsManyFieldType,
		GenFieldType:        "MultiEnumFieldType",
		IsPrimary:           f.StorageConfig.PrimaryKey,
		Required:            f.Required,
		Tags:                helpers.ResolveTags(f),
		Import:              nil,
		Custom:              true,
		Generated:           f.Autogenerated.Type != entity.InvalidAutogeneratedType,
		GeneratedFuncInsert: resolveGeneratedFuncInsert(e, f),
		GeneratedFuncUpdate: resolveGeneratedFuncUpdate(e, f),
		Enum:                true,
		EnumMany:            true,
		GraphName:           f.Identifier,
		GraphModelName:      helpers.ToCamelCase(f.Identifier),
		GraphInType:         fmt.Sprintf("[String%s]%s", graphRequired, graphRequired),
		GraphInTypeOptional: "[String]",
		GraphOutType:        fmt.Sprintf("[String%s]%s", graphRequired, graphRequired),
		GraphGenType:        "[]string",
		GraphGenToMapper:    fmt.Sprintf("Map%sToModel(i.%s)", pl.Singular(helpers.ToCamelCase(name)), helpers.ToCamelCase(f.Identifier)),
		GraphGenFromMapper:  fmt.Sprintf("Map%sFromModel(i.%s)", pl.Singular(helpers.ToCamelCase(name)), helpers.ToCamelCase(f.Identifier)),
		ProtoType:           protoType,
		ProtoName:           helpers.ToSnakeCase(f.Identifier),
		ProtoEnumOptions:    helpers.ProtoEnumOptions(protoType, f.OptionValues),
		ProtoGenName:        strcase.ToCamel(f.Identifier),
		ProtoToMapper:       fmt.Sprintf("%sSliceToProto(e.%s)", protoType, helpers.ToCamelCase(f.Identifier)),
		ProtoFromMapper:     fmt.Sprintf("%sSliceFromProto(m.Get%s())", protoType, strcase.ToCamel(f.Identifier)),
	}
}
